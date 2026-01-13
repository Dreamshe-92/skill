package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// 配置日志
var logger = log.New(os.Stdout, "WEAPM: ", log.LstdFlags|log.Lshortfile)

// ==================== 配置和客户端 ====================

// EnvConfig 环境配置
type EnvConfig struct {
	BaseURL           string  `yaml:"base_url"`
	Username          string  `yaml:"username"`
	Password          string  `yaml:"password"`
	Timeout           int     `yaml:"timeout"`
	MaxRetries        int     `yaml:"max_retries"`
	RetryBackoff      float64 `yaml:"retry_backoff_factor"`
	PoolConnections   int     `yaml:"pool_connections"`
	EnableLogging     bool    `yaml:"enable_logging"`
	Description       string  `yaml:"description"`
}

// ConfigFile 配置文件结构
type ConfigFile struct {
	Dev       EnvConfig `yaml:"dev"`
	Prod      EnvConfig `yaml:"prod"`
	ActiveEnv string    `yaml:"active_env"`
}

// Config WEAPM API 配置
type Config struct {
	BaseURL       string
	Timeout       time.Duration
	Username      string
	Password      string
	MaxRetries    int
	RetryBackoff  time.Duration
	EnableLogging bool
}

// LoadConfigFromYAML 从 YAML 文件加载配置
func LoadConfigFromYAML(configPath string, env string) (*Config, error) {
	// 默认配置文件路径
	if configPath == "" {
		execDir, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("获取可执行文件路径失败: %w", err)
		}
		configPath = filepath.Join(filepath.Dir(execDir), "config.yaml")
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML
	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 确定使用的环境
	if env == "" {
		env = configFile.ActiveEnv
	}
	if env == "" {
		env = "dev"
	}

	// 获取环境配置
	var envConfig EnvConfig
	switch env {
	case "dev":
		envConfig = configFile.Dev
	case "prod":
		envConfig = configFile.Prod
	default:
		return nil, fmt.Errorf("不支持的环境: %s, 可用环境: dev, prod", env)
	}

	// 验证必要字段
	if envConfig.BaseURL == "" {
		return nil, fmt.Errorf("环境 %s 缺少必要字段: base_url", env)
	}

	// 设置默认值
	if envConfig.Username == "" {
		envConfig.Username = "weapmUser"
	}
	if envConfig.Password == "" {
		envConfig.Password = "Weapm@123admin"
	}
	if envConfig.Timeout == 0 {
		envConfig.Timeout = 30
	}
	if envConfig.MaxRetries == 0 {
		envConfig.MaxRetries = 3
	}
	if envConfig.RetryBackoff == 0 {
		envConfig.RetryBackoff = 0.5
	}

	desc := envConfig.Description
	if desc == "" {
		desc = env
	}
	fmt.Printf("✅ 加载配置: %s (%s)\n", desc, env)

	return &Config{
		BaseURL:       envConfig.BaseURL,
		Timeout:       time.Duration(envConfig.Timeout) * time.Second,
		Username:      envConfig.Username,
		Password:      envConfig.Password,
		MaxRetries:    envConfig.MaxRetries,
		RetryBackoff:  time.Duration(envConfig.RetryBackoff * float64(time.Second)),
		EnableLogging: envConfig.EnableLogging,
	}, nil
}

// DefaultConfig 返回默认配置 (备用方案)
func DefaultConfig(baseURL string) *Config {
	return &Config{
		BaseURL:       baseURL,
		Timeout:       30 * time.Second,
		Username:      "weapmUser",
		Password:      "Weapm@123admin",
		MaxRetries:    3,
		RetryBackoff:  500 * time.Millisecond,
		EnableLogging: true,
	}
}

// Client WEAPM-LOGSERVER API 客户端
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient 创建新的客户端实例
func NewClient(config *Config) *Client {
	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &loggingRoundTripper{
				logger:   logger,
				next:     http.DefaultTransport,
				enable:   config.EnableLogging,
				baseURL:  config.BaseURL,
			},
		},
	}
	logger.Printf("WEAPM 客户端初始化成功: %s", config.BaseURL)
	return client
}

// loggingRoundTripper 日志记录的 HTTP Transport
type loggingRoundTripper struct {
	logger  *log.Logger
	next    http.RoundTripper
	enable  bool
	baseURL string
}

func (t *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	if t.enable {
		t.logger.Printf("发送请求: %s %s", req.Method, req.URL.String())
	}

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		if t.enable {
			t.logger.Printf("请求失败: %s %s - 错误: %v", req.Method, req.URL.String(), err)
		}
		return nil, err
	}

	if t.enable {
		duration := time.Since(start)
		t.logger.Printf(
			"收到响应: %s %s - 状态码: %d, 耗时: %.2fs",
			req.Method,
			req.URL.String(),
			resp.StatusCode,
			duration.Seconds(),
		)
	}

	return resp, nil
}

// ==================== 数据模型 ====================

// DashboardResult 数据大盘结果
type DashboardResult struct {
	SubsystemCount      int                 `json:"subsystemCount"`
	ClusterNum          int                 `json:"clusterNum"`
	ClusterTrafficData  []ClusterTrafficData `json:"clusterTrafficData"`
	TopSubsystems       []SubsystemLogDetail `json:"topSubsystems"`
	ClusterLogCounts    []ClusterLogCount   `json:"clusterLogCounts"`
}

// ClusterTrafficData 集群流量数据
type ClusterTrafficData struct {
	ClusterName  string `json:"clusterName"`
	TrafficBytes int64  `json:"trafficBytes"`
	Timestamp    string `json:"timestamp"`
}

// SubsystemLogDetail 子系统日志详情
type SubsystemLogDetail struct {
	Department        string `json:"department"`
	SubsysName        string `json:"subsys_name"`
	BusinessOwner     string `json:"business_owner"`
	SubsystemOwner    string `json:"subsystem_owner"`
	SubsysID          string `json:"subsys_id"`
	ClusterName       string `json:"cluster_name"`
	TotalLogMb        int64  `json:"total_log_mb"`
}

// ClusterLogCount 集群日志统计
type ClusterLogCount struct {
	ClusterName string `json:"clustername"`
	TotalLogGb  int    `json:"total_log_gb"`
	Capacity    int    `json:"capacity"`
}

// LogClusterInfo 集群信息
type LogClusterInfo struct {
	ClusterName   string `json:"clustername"`
	IsDefault     int    `json:"isdefault"`
	Topic         string `json:"topic"`
	BucketNames   string `json:"bucketnames"`
	BackendDomain string `json:"backenddomain"`
	StorageDomain string `json:"storagedomain"`
}

// LogStoreInstance 日志存储实例
type LogStoreInstance struct {
	Address       string `json:"address"`
	ClusterName   string `json:"clustername"`
	Role          string `json:"role"`
	Topic         string `json:"topic"`
	BucketNames   string `json:"bucketnames"`
	BackendDomain string `json:"backenddomain"`
	StorageDomain string `json:"storagedomain"`
	IsDefault     bool   `json:"isdefault"`
	Status        string `json:"status"`
	CpuLimit      string `json:"cpulimit"`
	MemLimit      string `json:"memlimit"`
	CreateTime    string `json:"createtime"`
	UpdateTime    string `json:"updateime"`
}

// ClusterDetailResult 集群详情结果
type ClusterDetailResult struct {
	ClusterInfo       LogClusterInfo       `json:"clusterInfo"`
	NodeGroups        []NodeGroup          `json:"nodeGroups"`
	ManagedSubSystems []LogSubClusterSubSystem `json:"managedSubSystems"`
	ReportData        ClusterReportData    `json:"reportData"`
}

// NodeGroup 节点组
type NodeGroup struct {
	Role  string              `json:"role"`
	Nodes []LogStoreInstance  `json:"nodes"`
}

// ClusterReportData 集群报表数据
type ClusterReportData struct {
	PeakTraffic      int64  `json:"peakTraffic"`
	PeakTime         string `json:"peakTime"`
	TotalSubSystems  int    `json:"totalSubSystems"`
	TopicBacklog     int64  `json:"topicBacklog"`
}

// LogSubClusterSubSystem 集群子系统
type LogSubClusterSubSystem struct {
	ClusterName     string `json:"clustername"`
	SubsystemID     string `json:"subsystemid"`
	SubsysName      string `json:"subsys_name"`
	SubsystemOwner  string `json:"subsystem_owner"`
	BusinessOwner   string `json:"business_owner"`
	DevDept         string `json:"devdept"`
	Traffic         int64  `json:"traffic"`
	Status          string `json:"status"`
	CreateTime      string `json:"createtime"`
	UpdateTime      string `json:"updatetime"`
}

// SubSystem 子系统信息
type SubSystem struct {
	ID               int    `json:"id"`
	SubsysID         string `json:"subsys_id"`
	SubsysName       string `json:"subsys_name"`
	SubsysChtname    string `json:"subsys_chtname"`
	SubsysUpdtime    string `json:"subsys_updtime"`
	DevDept          string `json:"devdept"`
	BusinessOwner    string `json:"business_owner"`
	SubsystemOwner   string `json:"subsystem_owner"`
	SystemName       string `json:"system_name"`
	State            string `json:"state"`
	ImportantLevel   string `json:"important_level"`
	CreateTopic      string `json:"create_topic"`
}

// SubsystemExistsResult 子系统存在性检查结果
type SubsystemExistsResult struct {
	SubsystemID   string `json:"subsystemId"`
	Exists        bool   `json:"exists"`
	SubsystemName string `json:"subsystemName"`
	ClusterName   string `json:"clusterName"`
}

// SubsystemDetailResult 子系统详情结果
type SubsystemDetailResult struct {
	SubsystemInfo    SubSystem `json:"subsystemInfo"`
	Collected        bool      `json:"collected"`
	ScanFileWhitelist []string `json:"scanFileWhitelist"`
	ExpectedTraffic  int       `json:"expectedTraffic"`
	ActualTraffic    int       `json:"actualTraffic"`
	KeywordFilters   []string  `json:"keywordFilters"`
	ClusterName      string    `json:"clusterName"`
	Instances        []map[string][]string `json:"instances"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

// ==================== HTTP 请求方法 ====================

// doRequest 执行HTTP请求 (带重试机制)
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body []byte) (*APIResponse, error) {
	var lastErr error

	// 重试逻辑
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// 计算退避时间
			backoff := time.Duration(float64(attempt) * c.config.RetryBackoff.Seconds() * float64(time.Second))
			logger.Printf("第 %d 次重试,退避时间: %.2fs", attempt, backoff.Seconds())
			time.Sleep(backoff)
		}

		// 构建完整URL
		fullURL := c.config.BaseURL + endpoint

		// 创建请求
		var req *http.Request
		var err error

		if body != nil {
			req, err = http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(body))
			if err != nil {
				return nil, fmt.Errorf("创建请求失败: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
			if err != nil {
				return nil, fmt.Errorf("创建请求失败: %w", err)
			}
		}

		// 设置Basic Auth
		req.SetBasicAuth(c.config.Username, c.config.Password)

		// 发送请求
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			logger.Printf("请求失败 (尝试 %d/%d): %v", attempt+1, c.config.MaxRetries+1, err)
			continue
		}

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			logger.Printf("读取响应失败 (尝试 %d/%d): %v", attempt+1, c.config.MaxRetries+1, err)
			continue
		}

		// 检查HTTP状态码
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("服务器错误: %d - %s", resp.StatusCode, string(respBody))
			logger.Printf("服务器错误 (尝试 %d/%d): %d", attempt+1, c.config.MaxRetries+1, resp.StatusCode)
			continue // 服务器错误,重试
		}

		if resp.StatusCode >= 400 {
			// 客户端错误,不重试
			return nil, fmt.Errorf("客户端错误: %d - %s", resp.StatusCode, string(respBody))
		}

		// 解析响应
		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w", err, string(respBody))
		}

		// 检查业务错误码
		if apiResp.Code != 0 {
			return &apiResp, fmt.Errorf("API错误 (code %d): %s", apiResp.Code, apiResp.Message)
		}

		// 成功
		if attempt > 0 {
			logger.Printf("请求成功 (重试 %d 次后)", attempt)
		}
		return &apiResp, nil
	}

	return nil, fmt.Errorf("请求失败,已重试 %d 次: %w", c.config.MaxRetries, lastErr)
}

// ==================== 数据大盘 API ====================

// GetDashboard 获取数据大盘信息
func (c *Client) GetDashboard(ctx context.Context) (*DashboardResult, error) {
	resp, err := c.doRequest(ctx, "GET", "/operation/dashboard", nil)
	if err != nil {
		return nil, err
	}

	var result DashboardResult
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ==================== 集群管理 API ====================

// GetClusters 获取所有集群信息
func (c *Client) GetClusters(ctx context.Context) ([]LogClusterInfo, error) {
	resp, err := c.doRequest(ctx, "GET", "/operation/clusters", nil)
	if err != nil {
		return nil, err
	}

	var clusters []LogClusterInfo
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &clusters); err != nil {
		return nil, err
	}

	return clusters, nil
}

// GetClusterDetail 获取指定集群的详细信息
func (c *Client) GetClusterDetail(ctx context.Context, clusterName string) (*ClusterDetailResult, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/operation/clusters/%s", clusterName), nil)
	if err != nil {
		return nil, err
	}

	var result ClusterDetailResult
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddClusterNodeRequest 向集群添加节点请求参数
type AddClusterNodeRequest struct {
	Address        string `json:"address"`         // 必填: 节点IP地址
	ClusterName    string `json:"clustername"`     // 必填: 集群名称
	Role           string `json:"role"`            // 必填: 节点角色
	CpuLimit       string `json:"cpulimit,omitempty"`        // 可选: CPU限制
	MemLimit       string `json:"memlimit,omitempty"`        // 可选: 内存限制
	Topic          string `json:"topic,omitempty"`           // 可选: Topic
	BucketNames    string `json:"bucketnames,omitempty"`     // 可选: 存储桶名称
	BackendDomain  string `json:"backenddomain,omitempty"`   // 可选: 后端域
	StorageDomain  string `json:"storagedomain,omitempty"`   // 可选: 存储域
	IsDefault      bool   `json:"isdefault,omitempty"`       // 可选: 是否默认
	Status         string `json:"status,omitempty"`          // 可选: 状态
	CreateTime     string `json:"createtime,omitempty"`      // 可选: 创建时间
	UpdateTime     string `json:"updateime,omitempty"`       // 可选: 更新时间
}

// AddClusterNode 向集群添加节点 (简化版,支持部分参数)
func (c *Client) AddClusterNode(ctx context.Context, clusterName string, req *AddClusterNodeRequest) error {
	// 设置集群名称
	req.ClusterName = clusterName

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化节点数据失败: %w", err)
	}

	_, err = c.doRequest(ctx, "POST", fmt.Sprintf("/operation/clusters/%s/nodes", clusterName), body)
	return err
}

// DeleteClusterNode 从集群删除节点
func (c *Client) DeleteClusterNode(ctx context.Context, ip string) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/operation/clusters/nodes/%s", ip), nil)
	return err
}

// GetClusterSubsystems 获取集群纳管的子系统信息
func (c *Client) GetClusterSubsystems(ctx context.Context, clusterName string) ([]LogSubClusterSubSystem, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/operation/cluster/%s/subsystems", clusterName), nil)
	if err != nil {
		return nil, err
	}

	var subsystems []LogSubClusterSubSystem
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &subsystems); err != nil {
		return nil, err
	}

	return subsystems, nil
}

// ==================== 子系统运维 API ====================

// CheckSubsystemExists 检查子系统是否存在
func (c *Client) CheckSubsystemExists(ctx context.Context, subsystemID string) (*SubsystemExistsResult, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/operation/subsystem/exists/%s", subsystemID), nil)
	if err != nil {
		return nil, err
	}

	var result SubsystemExistsResult
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddSubsystemRequest 新增子系统请求
type AddSubsystemRequest struct {
	SubSystemID    string `json:"subSystemId"`
	LogImportValue string `json:"logImportValue"`
	LogImportFiles string `json:"logImportFiles"`
	Traffic        int    `json:"traffic"`
	Cluster        string `json:"cluster"`
}

// AddSubsystem 新增子系统接入
func (c *Client) AddSubsystem(ctx context.Context, req *AddSubsystemRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	_, err = c.doRequest(ctx, "POST", "/operation/subsystem", body)
	return err
}

// AdjustSubsystemCluster 调整子系统归属集群
func (c *Client) AdjustSubsystemCluster(ctx context.Context, subsystemID, targetClusterName, logImportValue, logImportFiles string, traffic int) error {
	params := url.Values{}
	params.Set("targetClusterName", targetClusterName)
	params.Set("logImportValue", logImportValue)
	params.Set("logImportFiles", logImportFiles)
	params.Set("traffic", strconv.Itoa(traffic))

	endpoint := fmt.Sprintf("/operation/subsystem/%s?%s", subsystemID, params.Encode())
	_, err := c.doRequest(ctx, "POST", endpoint, nil)
	return err
}

// AdjustSubsystemStatus 调整子系统状态
func (c *Client) AdjustSubsystemStatus(ctx context.Context, subsystemID, status string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/operation/subsystem/%s/status/%s", subsystemID, status), nil)
	return err
}

// EnableSubsystem 启用子系统
func (c *Client) EnableSubsystem(ctx context.Context, subsystemID string) error {
	_, err := c.doRequest(ctx, "PUT", fmt.Sprintf("/operation/subsystem/%s/enable", subsystemID), nil)
	return err
}

// GetSubsystemDetail 获取子系统详情
func (c *Client) GetSubsystemDetail(ctx context.Context, subsystemID string) (*SubsystemDetailResult, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/operation/subsystem/%s", subsystemID), nil)
	if err != nil {
		return nil, err
	}

	var result SubsystemDetailResult
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSubsystems 获取所有子系统信息
func (c *Client) GetSubsystems(ctx context.Context) ([]SubSystem, error) {
	resp, err := c.doRequest(ctx, "GET", "/operation/subsystems", nil)
	if err != nil {
		return nil, err
	}

	var subsystems []SubSystem
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &subsystems); err != nil {
		return nil, err
	}

	return subsystems, nil
}

// SearchSubsystemsRequest 搜索子系统请求参数
type SearchSubsystemsRequest struct {
	SubsysID *string
	Limit    int
}

// SearchSubsystems 根据条件搜索子系统
func (c *Client) SearchSubsystems(ctx context.Context, req *SearchSubsystemsRequest) ([]SubSystem, error) {
	params := url.Values{}
	if req.SubsysID != nil {
		params.Set("subsysId", *req.SubsysID)
	}
	if req.Limit != 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	} else {
		params.Set("limit", "20")
	}

	endpoint := "/operation/subsystems/search"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var subsystems []SubSystem
	if err := json.Unmarshal(resp.Result.(*json.RawMessage), &subsystems); err != nil {
		return nil, err
	}

	return subsystems, nil
}

// ==================== 主函数示例 ====================

func main() {
	// 方式 1: 从配置文件加载 (推荐)
	config, err := LoadConfigFromYAML("", "")
	if err != nil {
		fmt.Printf("⚠️  %v\n", err)
		fmt.Println("请先创建配置文件 config.yaml,参考 config.yaml.example")
		return
	}

	// 或者指定环境
	// config, err := LoadConfigFromYAML("", "dev")
	// config, err := LoadConfigFromYAML("", "prod")

	// 或者指定配置文件路径
	// config, err := LoadConfigFromYAML("/path/to/config.yaml", "prod")
	// if err != nil {
	//     fmt.Printf("加载配置失败: %v\n", err)
	//     return
	// }

	// 方式 2: 手动创建配置 (备用方案)
	// config := DefaultConfig("http://localhost:8080")

	// 创建客户端
	client := NewClient(config)

	// 创建上下文
	ctx := context.Background()

	// 示例 1: 获取数据大盘信息
	fmt.Println("========================================")
	fmt.Println("1. 获取数据大盘信息")
	dashboard, err := client.GetDashboard(ctx)
	if err != nil {
		fmt.Printf("获取数据大盘失败: %v\n", err)
	} else {
		fmt.Printf("子系统数量: %d\n", dashboard.SubsystemCount)
		fmt.Printf("集群数量: %d\n", dashboard.ClusterNum)
	}

	// 示例 2: 获取所有集群
	fmt.Println("\n========================================")
	fmt.Println("2. 获取所有集群")
	clusters, err := client.GetClusters(ctx)
	if err != nil {
		fmt.Printf("获取集群列表失败: %v\n", err)
	} else {
		for _, cluster := range clusters {
			fmt.Printf("集群名称: %s, 默认: %d\n", cluster.ClusterName, cluster.IsDefault)
		}
	}

	// 示例 3: 获取集群详情
	if len(clusters) > 0 {
		clusterName := clusters[0].ClusterName
		fmt.Printf("\n========================================\n")
		fmt.Printf("3. 获取集群详情: %s\n", clusterName)
		clusterDetail, err := client.GetClusterDetail(ctx, clusterName)
		if err != nil {
			fmt.Printf("获取集群详情失败: %v\n", err)
		} else {
			clusterJSON, _ := json.MarshalIndent(clusterDetail.ClusterInfo, "", "  ")
			fmt.Printf("集群信息: %s\n", string(clusterJSON))
		}
	}

	// 示例 4: 获取所有子系统
	fmt.Println("\n========================================")
	fmt.Println("4. 获取所有子系统")
	subsystems, err := client.GetSubsystems(ctx)
	if err != nil {
		fmt.Printf("获取子系统列表失败: %v\n", err)
	} else {
		fmt.Printf("子系统总数: %d\n", len(subsystems))
		limit := 5
		if len(subsystems) < limit {
			limit = len(subsystems)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("子系统ID: %s, 名称: %s\n", subsystems[i].SubsysID, subsystems[i].SubsysName)
		}
	}

	// 示例 5: 搜索子系统
	fmt.Println("\n========================================")
	fmt.Println("5. 搜索子系统")
	searchResults, err := client.SearchSubsystems(ctx, &SearchSubsystemsRequest{Limit: 10})
	if err != nil {
		fmt.Printf("搜索子系统失败: %v\n", err)
	} else {
		fmt.Printf("搜索到 %d 个子系统\n", len(searchResults))
	}

	// 示例 6: 检查子系统是否存在
	if len(subsystems) > 0 {
		subsystemID := subsystems[0].SubsysID
		fmt.Printf("\n========================================\n")
		fmt.Printf("6. 检查子系统是否存在: %s\n", subsystemID)
		existsResult, err := client.CheckSubsystemExists(ctx, subsystemID)
		if err != nil {
			fmt.Printf("检查子系统存在性失败: %v\n", err)
		} else {
			fmt.Printf("存在: %v\n", existsResult.Exists)
		}
	}

	// 示例 7: 向集群添加节点 (最小化参数)
	fmt.Println("\n========================================")
	fmt.Println("7. 向集群添加节点 (最小化参数)")
	minimalNode := &AddClusterNodeRequest{
		Address:  "127.0.0.2",
		Role:     "write",
		CpuLimit: "8",
		MemLimit: "16",
	}
	err = client.AddClusterNode(ctx, "LOG008", minimalNode)
	if err != nil {
		fmt.Printf("添加节点失败: %v\n", err)
	} else {
		fmt.Printf("节点添加成功\n")
	}

	// 示例 8: 向集群添加节点 (完整参数)
	fmt.Println("\n========================================")
	fmt.Println("8. 向集群添加节点 (完整参数)")
	fullNode := &AddClusterNodeRequest{
		Address:       "127.0.0.3",
		Role:          "master",
		CpuLimit:      "16",
		MemLimit:      "32",
		Topic:         "log_topic_008",
		BucketNames:   "log_bucket_008",
		BackendDomain: "backend.example.com",
		StorageDomain: "storage.example.com",
		Status:        "active",
	}
	err = client.AddClusterNode(ctx, "LOG008", fullNode)
	if err != nil {
		fmt.Printf("添加节点失败: %v\n", err)
	} else {
		fmt.Printf("节点添加成功\n")
	}

	fmt.Println("\n========================================")
	fmt.Println("示例执行完成")
}
