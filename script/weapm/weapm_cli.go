package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// ==================== 命令行参数 ====================

type CommandLineArgs struct {
	ConfigPath  string
	Env         string
	BaseURL     string
	Username    string
	Password    string
	Timeout     int
	Quiet       bool
	Command     string
	ClusterName string
	Detail      bool
	Search      bool
	SubsysID    string
	Check       string
	Limit       int
	Address     string
	Role        string
	CpuLimit    string
	MemLimit    string
	Topic       string
	BucketNames string
	BackendDomain string
	StorageDomain string
	Status      string
}

func parseArgs() *CommandLineArgs {
	args := &CommandLineArgs{}

	// 全局参数
	flag.StringVar(&args.ConfigPath, "config", "", "配置文件路径")
	flag.StringVar(&args.ConfigPath, "c", "", "配置文件路径 (简写)")
	flag.StringVar(&args.Env, "env", "", "环境名称 (dev/prod)")
	flag.StringVar(&args.Env, "e", "", "环境名称 (简写)")
	flag.StringVar(&args.BaseURL, "base-url", "", "API 基础 URL")
	flag.StringVar(&args.Username, "username", "", "用户名")
	flag.StringVar(&args.Password, "password", "", "密码")
	flag.IntVar(&args.Timeout, "timeout", 30, "请求超时时间(秒)")
	flag.BoolVar(&args.Quiet, "quiet", false, "静默模式,不输出日志")
	flag.BoolVar(&args.Quiet, "q", false, "静默模式 (简写)")

	// 集群管理参数
	flag.StringVar(&args.ClusterName, "cluster-name", "", "集群名称")
	flag.StringVar(&args.ClusterName, "n", "", "集群名称 (简写)")
	flag.BoolVar(&args.Detail, "detail", false, "显示详细信息")
	flag.BoolVar(&args.Detail, "d", false, "显示详细信息 (简写)")

	// 子系统参数
	flag.BoolVar(&args.Search, "search", false, "搜索子系统")
	flag.BoolVar(&args.Search, "s", false, "搜索子系统 (简写)")
	flag.StringVar(&args.SubsysID, "subsys-id", "", "子系统ID")
	flag.StringVar(&args.Check, "check", "", "检查子系统是否存在")
	flag.IntVar(&args.Limit, "limit", 20, "返回结果数量限制")
	flag.IntVar(&args.Limit, "l", 20, "返回结果数量限制 (简写)")

	// 节点管理参数
	flag.StringVar(&args.Address, "address", "", "节点IP地址")
	flag.StringVar(&args.Role, "role", "", "节点角色")
	flag.StringVar(&args.CpuLimit, "cpulimit", "", "CPU限制")
	flag.StringVar(&args.MemLimit, "memlimit", "", "内存限制")
	flag.StringVar(&args.Topic, "topic", "", "Topic")
	flag.StringVar(&args.BucketNames, "bucketnames", "", "存储桶名称")
	flag.StringVar(&args.BackendDomain, "backenddomain", "", "后端域")
	flag.StringVar(&args.StorageDomain, "storagedomain", "", "存储域")
	flag.StringVar(&args.Status, "status", "", "状态")

	flag.Parse()

	// 获取命令 (第一个非标志参数)
	if len(flag.Args()) > 0 {
		args.Command = flag.Args()[0]
	}

	return args
}

// ==================== 命令处理函数 ====================

func cmdDashboard(client *Client) error {
	ctx := context.Background()
	dashboard, err := client.GetDashboard(ctx)
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(dashboard, "", "  ")
	fmt.Println(string(result))
	return nil
}

func cmdClusters(client *Client, args *CommandLineArgs) error {
	ctx := context.Background()

	if args.Detail {
		if args.ClusterName == "" {
			return fmt.Errorf("使用 --detail 时必须指定 --cluster-name")
		}

		result, err := client.GetClusterDetail(ctx, args.ClusterName)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	} else {
		clusters, err := client.GetClusters(ctx)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(clusters, "", "  ")
		fmt.Println(string(output))
	}

	return nil
}

func cmdSubsystems(client *Client, args *CommandLineArgs) error {
	ctx := context.Background()

	var result interface{}
	var err error

	if args.Search {
		result, err = client.SearchSubsystems(ctx, &SearchSubsystemsRequest{
			SubsysID: &args.SubsysID,
			Limit:    args.Limit,
		})
	} else if args.Check != "" {
		result, err = client.CheckSubsystemExists(ctx, args.Check)
	} else if args.Detail != "" {
		result, err = client.GetSubsystemDetail(ctx, args.Detail)
	} else {
		result, err = client.GetSubsystems(ctx)
	}

	if err != nil {
		return err
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
	return nil
}

func cmdAddNode(client *Client, args *CommandLineArgs) error {
	ctx := context.Background()

	node := &AddClusterNodeRequest{
		Address:       args.Address,
		Role:          args.Role,
		CpuLimit:      args.CpuLimit,
		MemLimit:      args.MemLimit,
		Topic:         args.Topic,
		BucketNames:   args.BucketNames,
		BackendDomain: args.BackendDomain,
		StorageDomain: args.StorageDomain,
		Status:        args.Status,
	}

	err := client.AddClusterNode(ctx, args.ClusterName, node)
	if err != nil {
		return err
	}

	fmt.Println(`{"code": 0, "message": "节点添加成功"}`)
	return nil
}

func cmdDeleteNode(client *Client, args *CommandLineArgs) error {
	ctx := context.Background()

	// 从 args 中获取 IP
	ip := ""
	if len(flag.Args()) > 1 {
		ip = flag.Args()[1]
	}

	if ip == "" {
		return fmt.Errorf("请指定节点IP地址")
	}

	err := client.DeleteClusterNode(ctx, ip)
	if err != nil {
		return err
	}

	fmt.Println(`{"code": 0, "message": "节点删除成功"}`)
	return nil
}

// ==================== 主函数 ====================

func main() {
	args := parseArgs()

	// 如果没有指定命令,显示帮助
	if args.Command == "" {
		fmt.Println("WEAPM-LOGSERVER API 客户端命令行工具")
		fmt.Println("\n使用方法:")
		fmt.Println("  weapm_cli <命令> [参数]")
		fmt.Println("\n可用命令:")
		fmt.Println("  dashboard    获取数据大盘信息")
		fmt.Println("  clusters     集群管理")
		fmt.Println("  subsystems   子系统管理")
		fmt.Println("  add-node     添加集群节点")
		fmt.Println("  delete-node  删除集群节点")
		fmt.Println("\n示例:")
		fmt.Println("  ./weapm_cli dashboard")
		fmt.Println("  ./weapm_cli clusters")
		fmt.Println("  ./weapm_cli clusters --detail --cluster-name LOG001")
		fmt.Println("  ./weapm_cli subsystems")
		fmt.Println("  ./weapm_cli subsystems --search --subsys-id SYS001")
		fmt.Println("  ./weapm_cli add-node --cluster-name LOG008 --address 127.0.0.2 --role write")
		fmt.Println("\n使用 --help 查看详细帮助")
		os.Exit(0)
	}

	// 配置日志
	if args.Quiet {
		log.SetOutput(os.NewFile(0, os.DevNull))
	}

	// 加载配置
	var config *Config
	var err error

	if args.ConfigPath != "" || args.Env != "" {
		config, err = LoadConfigFromYAML(args.ConfigPath, args.Env)
	} else if args.BaseURL != "" {
		// 使用命令行参数创建配置
		config = DefaultConfig(args.BaseURL)
		if args.Username != "" {
			config.Username = args.Username
		}
		if args.Password != "" {
			config.Password = args.Password
		}
		if args.Timeout != 30 {
			config.Timeout = time.Duration(args.Timeout) * time.Second
		}
	} else {
		// 默认使用配置文件
		config, err = LoadConfigFromYAML("", "")
	}

	if err != nil {
		log.Fatalf("⚠️  %v\n请先创建配置文件 config.yaml,参考 config.yaml.example", err)
	}

	// 创建客户端
	client := NewClient(config)

	// 执行命令
	var cmdErr error
	switch args.Command {
	case "dashboard":
		cmdErr = cmdDashboard(client)
	case "clusters":
		cmdErr = cmdClusters(client, args)
	case "subsystems":
		cmdErr = cmdSubsystems(client, args)
	case "add-node":
		cmdErr = cmdAddNode(client, args)
	case "delete-node":
		cmdErr = cmdDeleteNode(client, args)
	default:
		log.Fatalf("❌ 未知命令: %s", args.Command)
	}

	if cmdErr != nil {
		log.Fatalf("❌ 错误: %v", cmdErr)
	}
}
