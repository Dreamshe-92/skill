# WEAPM-LOGSERVER API 客户端

基于 WEAPM-LOGSERVER REST API 的 Python 和 Golang 客户端实现。

## 📋 目录结构

```
script/weapm/
├── swagger.md               # API 文档
├── weapm_client.py          # Python 客户端
├── weapm_client.go          # Golang 客户端
├── config.yaml.example      # 配置文件示例
├── config.yaml              # 实际配置文件 (需自行创建)
└── README.md               # 使用说明
```

## ⚙️ 配置文件 (推荐)

### 1. 创建配置文件

复制配置文件示例:
```bash
cd script/weapm/
cp config.yaml.example config.yaml
```

### 2. 编辑配置文件

根据实际情况修改 `config.yaml`:

```yaml
# 开发/测试环境配置
dev:
  base_url: "http://localhost:8080"
  username: "weapmUser"
  password: "Weapm@123admin"
  timeout: 30
  description: "开发测试环境"

# 生产环境配置
prod:
  base_url: "https://weapm.example.com"
  username: "weapm_admin"
  password: "prod_password_here"
  timeout: 60
  description: "生产环境"

# 默认使用的环境 (dev | prod)
# 修改此值来切换环境
active_env: "dev"
```

### 3. 切换环境

只需修改 `active_env` 字段:
- 测试环境: `active_env: "dev"`
- 生产环境: `active_env: "prod"`

## 🚀 快速开始

### Python 客户端

#### 环境要求

- Python 3.7+
- requests 库
- PyYAML 库

#### 安装依赖

```bash
pip install requests pyyaml
```

#### 基本使用

```python
from weapm_client import WeapmClient, WeapmConfig

# 方式 1: 从配置文件加载 (推荐)
config = WeapmConfig.from_yaml()  # 使用 config.yaml 中的 active_env
# 或指定环境
config = WeapmConfig.from_yaml(env="dev")
config = WeapmConfig.from_yaml(env="prod")

# 方式 2: 手动创建配置
config = WeapmConfig(
    base_url="http://localhost:8080",
    username="weapmUser",
    password="Weapm@123admin",
    timeout=30
)

# 创建客户端
with WeapmClient(config) as client:
    # 获取数据大盘信息
    dashboard = client.get_dashboard()
    print(f"子系统数量: {dashboard['result']['subsystemCount']}")

    # 获取所有集群
    clusters = client.get_clusters()
    for cluster in clusters:
        print(f"集群: {cluster['clustername']}")
```

#### 运行示例

```bash
# 直接运行脚本查看完整示例
python weapm_client.py
```

### Golang 客户端

#### 环境要求

- Go 1.16+

#### 安装依赖

```bash
go get gopkg.in/yaml.v3
```

#### 基本使用

```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // 方式 1: 从配置文件加载 (推荐)
    config, err := LoadConfigFromYAML("", "")  // 使用 config.yaml 中的 active_env
    // 或指定环境
    config, err := LoadConfigFromYAML("", "dev")
    config, err := LoadConfigFromYAML("", "prod")

    if err != nil {
        fmt.Printf("加载配置失败: %v\n", err)
        return
    }

    // 方式 2: 手动创建配置
    // config := DefaultConfig("http://localhost:8080")

    // 创建客户端
    client := NewClient(config)

    // 创建上下文
    ctx := context.Background()

    // 获取数据大盘信息
    dashboard, err := client.GetDashboard(ctx)
    if err != nil {
        fmt.Printf("获取数据大盘失败: %v\n", err)
        return
    }
    fmt.Printf("子系统数量: %d\n", dashboard.SubsystemCount)
}
```

#### 运行示例

```bash
# 直接运行脚本查看完整示例
go run weapm_client.go
```

## 📚 API 功能说明

### 数据大盘

- `get_dashboard()` / `GetDashboard()`: 获取数据大盘信息,包括子系统数、集群数、流量数据等

### 集群管理

- `get_clusters()` / `GetClusters()`: 获取所有集群信息
- `get_cluster_detail(cluster_name)` / `GetClusterDetail()`: 获取指定集群的详细信息
- `add_cluster_node(cluster_name, node_data)` / `AddClusterNode()`: 向集群添加节点
- `delete_cluster_node(ip)` / `DeleteClusterNode()`: 从集群删除节点
- `get_cluster_subsystems(cluster_name)` / `GetClusterSubsystems()`: 获取集群纳管的子系统

### 子系统运维

- `check_subsystem_exists(subsystem_id)` / `CheckSubsystemExists()`: 检查子系统是否存在
- `add_subsystem(...)` / `AddSubsystem()`: 新增子系统接入
- `adjust_subsystem_cluster(...)` / `AdjustSubsystemCluster()`: 调整子系统归属集群
- `adjust_subsystem_status(subsystem_id, status)` / `AdjustSubsystemStatus()`: 调整子系统状态
- `enable_subsystem(subsystem_id)` / `EnableSubsystem()`: 启用子系统
- `get_subsystem_detail(subsystem_id)` / `GetSubsystemDetail()`: 获取子系统详情
- `get_subsystems()` / `GetSubsystems()`: 获取所有子系统信息
- `search_subsystems(...)` / `SearchSubsystems()`: 根据条件搜索子系统

## 🔐 认证配置

### 使用配置文件 (推荐)

所有认证信息都在 `config.yaml` 中配置:

```yaml
dev:
  base_url: "http://localhost:8080"
  username: "weapmUser"      # 认证用户名
  password: "Weapm@123admin"  # 认证密码
  timeout: 30
```

### Python 自定义认证

```python
# 方式 1: 使用配置文件 (推荐)
config = WeapmConfig.from_yaml(env="dev")

# 方式 2: 手动创建配置
config = WeapmConfig(
    base_url="http://localhost:8080",
    username="your_username",  # 自定义用户名
    password="your_password",  # 自定义密码
    timeout=30                 # 请求超时时间(秒)
)
```

### Golang 自定义认证

```go
// 方式 1: 使用配置文件 (推荐)
config, err := LoadConfigFromYAML("", "dev")

// 方式 2: 手动创建配置
config := DefaultConfig("http://localhost:8080")
config.Username = "your_username"  // 自定义用户名
config.Password = "your_password"  // 自定义密码
```

## 💡 使用场景

### 场景 1: 监控数据大盘

```python
# Python
with WeapmClient(config) as client:
    dashboard = client.get_dashboard()
    # 处理大盘数据...
```

```go
// Golang
dashboard, err := client.GetDashboard(ctx)
// 处理大盘数据...
```

### 场景 2: 批量操作子系统

```python
# Python
subsystems = client.get_subsystems()
for subsystem in subsystems:
    if subsystem['state'] == 'disabled':
        client.enable_subsystem(subsystem['subsys_id'])
```

```go
// Golang
subsystems, _ := client.GetSubsystems(ctx)
for _, subsystem := range subsystems {
    if subsystem.State == "disabled" {
        client.EnableSubsystem(ctx, subsystem.SubsysID)
    }
}
```

### 场景 3: 集群节点管理

```python
# Python
# 添加节点 (最小化参数)
client.add_cluster_node(
    cluster_name="LOG008",
    address="127.0.0.2",
    role="write",
    cpulimit="8",
    memlimit="16"
)

# 添加节点 (完整参数)
client.add_cluster_node(
    cluster_name="LOG008",
    address="127.0.0.3",
    role="master",
    cpulimit="16",
    memlimit="32",
    topic="log_topic_008",
    bucketnames="log_bucket_008",
    backenddomain="backend.example.com",
    storagedomain="storage.example.com",
    status="active"
)

# 删除节点
client.delete_cluster_node("192.168.1.100")
```

```go
// Golang
// 添加节点 (最小化参数)
minimalNode := &AddClusterNodeRequest{
    Address:  "127.0.0.2",
    Role:     "write",
    CpuLimit: "8",
    MemLimit: "16",
}
client.AddClusterNode(ctx, "LOG008", minimalNode)

// 添加节点 (完整参数)
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
client.AddClusterNode(ctx, "LOG008", fullNode)

// 删除节点
client.DeleteClusterNode(ctx, "192.168.1.100")
```

## ⚠️ 错误处理

### Python

```python
try:
    dashboard = client.get_dashboard()
except requests.RequestException as e:
    print(f"API 调用失败: {str(e)}")
```

### Golang

```go
dashboard, err := client.GetDashboard(ctx)
if err != nil {
    fmt.Printf("获取数据大盘失败: %v\n", err)
    return
}
```

## 📝 完整 API 文档

详细的 API 文档请参考 [swagger.md](swagger.md)

## 🔧 故障排查

### 问题 1: 配置文件不存在

**错误信息**:
```
⚠️  配置文件不存在: /path/to/config.yaml
```

**解决方案**:
1. 复制示例配置文件: `cp config.yaml.example config.yaml`
2. 根据实际情况修改配置

### 问题 2: 连接超时

**错误信息**:
```
请求失败: HTTPConnectionPool(host='localhost', port=8080): Max retries exceeded
```

**解决方案**:
1. 检查网络连接
2. 确认 API 服务地址正确
3. 在 config.yaml 中增加超时时间:
   ```yaml
   dev:
     timeout: 60  # 增加到 60 秒
   ```

### 问题 3: 认证失败

**错误信息**:
```
API错误 (code 401): Unauthorized
```

**解决方案**:
1. 确认 config.yaml 中的用户名密码正确
2. 检查 API 服务认证配置
3. 验证 Basic Auth 凭据

### 问题 4: 环境配置不存在

**错误信息**:
```
⚠️  配置错误: 环境配置不存在: test, 可用环境: ['dev', 'prod']
```

**解决方案**:
1. 检查环境名称拼写
2. 确认使用正确的环境: dev 或 prod
3. 或在 config.yaml 中添加新环境配置

### 问题 5: 依赖包缺失

**错误信息** (Python):
```
ModuleNotFoundError: No module named 'yaml'
```

**解决方案**:
```bash
pip install pyyaml
```

**错误信息** (Golang):
```
cannot find package "gopkg.in/yaml.v3"
```

**解决方案**:
```bash
go get gopkg.in/yaml.v3
```

### 问题 6: 返回错误码

**解决方案**:
1. 检查返回的 `code` 和 `message` 字段
2. 参考 [swagger.md](swagger.md) 中的错误响应格式
3. 确认请求参数正确

## 📄 许可证

MIT License

## 👤 作者

WEAPM Team
