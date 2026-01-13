# WEAPM-LOGSERVER 命令行工具使用指南

本文档说明如何使用 WEAPM-LOGSERVER API 客户端的命令行工具。

## 📋 目录

- [Python 命令行工具](#python-命令行工具)
- [Golang 命令行工具](#golang-命令行工具)
- [命令参考](#命令参考)
- [使用示例](#使用示例)

---

## Python 命令行工具

### 基本用法

```bash
python weapm_cli.py <命令> [参数]
```

### 全局参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--config` | `-c` | 配置文件路径 |
| `--env` | `-e` | 环境名称 (dev/prod) |
| `--base-url` | | API 基础 URL |
| `--username` | | 用户名 |
| `--password` | | 密码 |
| `--timeout` | | 请求超时时间(秒) |
| `--quiet` | `-q` | 静默模式 |

### 示例

```bash
# 使用配置文件
python weapm_cli.py dashboard

# 指定环境
python weapm_cli.py --env prod dashboard

# 自定义 API 地址
python weapm_cli.py --base-url http://192.168.1.100:8080 clusters

# 静默模式
python weapm_cli.py -q dashboard
```

---

## Golang 命令行工具

### 基本用法

```bash
# 编译后使用
go build weapm_cli.go
./weapm_cli <命令> [参数]

# 或直接运行
go run weapm_cli.go <命令> [参数]
```

### 全局参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--config` | `-c` | 配置文件路径 |
| `--env` | `-e` | 环境名称 (dev/prod) |
| `--base-url` | | API 基础 URL |
| `--username` | | 用户名 |
| `--password` | | 密码 |
| `--timeout` | | 请求超时时间(秒) |
| `--quiet` | `-q` | 静默模式 |

### 示例

```bash
# 使用配置文件
./weapm_cli dashboard

# 指定环境
./weapm_cli --env prod clusters

# 自定义 API 地址
./weapm_cli --base-url http://192.168.1.100:8080 clusters

# 静默模式
./weapm_cli -q dashboard
```

---

## 命令参考

### 1. dashboard - 数据大盘

获取系统概览信息。

```bash
# Python
python weapm_cli.py dashboard

# Golang
./weapm_cli dashboard
```

**输出示例:**
```json
{
  "code": 0,
  "message": "success",
  "result": {
    "subsystemCount": 120,
    "clusterNum": 5,
    "clusterTrafficData": [...],
    "topSubsystems": [...]
  }
}
```

---

### 2. clusters - 集群管理

管理 WEAPM 集群。

#### 2.1 获取所有集群

```bash
# Python
python weapm_cli.py clusters

# Golang
./weapm_cli clusters
```

#### 2.2 获取集群详情

```bash
# Python
python weapm_cli.py clusters --detail --cluster-name LOG001

# Golang
./weapm_cli clusters --detail --cluster-name LOG001
```

**参数:**
- `--detail` / `-d` - 显示详细信息
- `--cluster-name` / `-n` - 集群名称

---

### 3. subsystems - 子系统管理

管理子系统信息。

#### 3.1 获取所有子系统

```bash
# Python
python weapm_cli.py subsystems

# Golang
./weapm_cli subsystems
```

#### 3.2 搜索子系统

```bash
# Python
python weapm_cli.py subsystems --search --subsys-id SYS001

# Golang
./weapm_cli subsystems --search --subsys-id SYS001
```

#### 3.3 检查子系统是否存在

```bash
# Python
python weapm_cli.py subsystems --check SYS001

# Golang
./weapm_cli subsystems --check SYS001
```

#### 3.4 获取子系统详情

```bash
# Python
python weapm_cli.py subsystems --detail SYS001

# Golang
./weapm_cli subsystems --detail SYS001
```

**参数:**
- `--search` / `-s` - 搜索模式
- `--subsys-id` - 子系统ID
- `--check` / `-c` - 检查是否存在
- `--detail` / `-d` - 显示详细信息
- `--limit` / `-l` - 返回结果数量限制 (默认: 20)

---

### 4. add-node - 添加集群节点

向集群添加新节点。

```bash
# Python
python weapm_cli.py add-node \
  --cluster-name LOG008 \
  --address 127.0.0.2 \
  --role write \
  --cpulimit 8 \
  --memlimit 16

# Golang
./weapm_cli add-node \
  --cluster-name LOG008 \
  --address 127.0.0.2 \
  --role write \
  --cpulimit 8 \
  --memlimit 16
```

**参数:**
- `--cluster-name` (必填) - 集群名称
- `--address` (必填) - 节点IP地址
- `--role` (必填) - 节点角色
- `--cpulimit` (可选) - CPU限制
- `--memlimit` (可选) - 内存限制
- `--topic` (可选) - Topic
- `--bucketnames` (可选) - 存储桶名称
- `--backenddomain` (可选) - 后端域
- `--storagedomain` (可选) - 存储域
- `--status` (可选) - 状态

**完整参数示例:**

```bash
python weapm_cli.py add-node \
  --cluster-name LOG008 \
  --address 127.0.0.3 \
  --role master \
  --cpulimit 16 \
  --memlimit 32 \
  --topic log_topic \
  --bucketnames log_bucket \
  --backenddomain backend.example.com \
  --storagedomain storage.example.com \
  --status active
```

---

### 5. delete-node - 删除集群节点

从集群删除节点。

```bash
# Python
python weapm_cli.py delete-node --ip 192.168.1.100

# Golang
./weapm_cli delete-node --ip 192.168.1.100
```

**参数:**
- `--ip` (必填) - 节点IP地址

---

## 使用示例

### 场景 1: 快速查看系统状态

```bash
# Python
python weapm_cli.py dashboard

# Golang
./weapm_cli dashboard
```

### 场景 2: 批量查询集群信息

```bash
# 查看所有集群
python weapm_cli.py clusters

# 查看特定集群详情
python weapm_cli.py clusters --detail --cluster-name LOG001
```

### 场景 3: 节点维护

```bash
# 添加节点
python weapm_cli.py add-node --cluster-name LOG008 --address 192.168.1.50 --role write

# 删除节点
python weapm_cli.py delete-node --ip 192.168.1.50
```

### 场景 4: 子系统管理

```bash
# 搜索特定子系统
python weapm_cli.py subsystems --search --subsys-id SYS001 --limit 10

# 查看子系统详情
python weapm_cli.py subsystems --detail SYS001

# 检查子系统是否存在
python weapm_cli.py subsystems --check SYS001
```

### 场景 5: 使用不同环境

```bash
# 开发环境
python weapm_cli.py --env dev dashboard

# 生产环境
python weapm_cli.py --env prod dashboard
```

### 场景 6: 自定义连接信息

```bash
# 使用自定义 API 地址
python weapm_cli.py --base-url http://192.168.1.100:8080 dashboard

# 自定义认证信息
python weapm_cli.py --base-url http://192.168.1.100:8080 \
  --username admin \
  --password secret123 \
  dashboard
```

---

## 输出格式

所有命令输出 JSON 格式数据:

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "result": { ... }
}
```

### 错误响应

```json
{
  "code": 1,
  "message": "错误描述"
}
```

---

## 退出码

| 退出码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1 | 错误 |

---

## 配置文件

命令行工具支持使用配置文件,避免重复输入参数。

### 创建配置文件

```bash
cd script/weapm/
cp config.yaml.example config.yaml
```

### 使用配置文件

```bash
# 使用默认配置文件 (config.yaml)
python weapm_cli.py dashboard

# 指定配置文件路径
python weapm_cli.py --config /path/to/config.yaml dashboard

# 指定环境
python weapm_cli.py --config /path/to/config.yaml --env dev dashboard
```

---

## 故障排查

### 问题 1: 命令未找到

**错误信息:**
```
bash: python: command not found
```

**解决方案:**
```bash
# 使用 python3
python3 weapm_cli.py dashboard

# 或添加执行权限后直接运行
chmod +x weapm_cli.py
./weapm_cli.py dashboard
```

### 问题 2: 权限拒绝

**错误信息:**
```
Permission denied
```

**解决方案:**
```bash
chmod +x weapm_cli.py
chmod +x weapm_cli.go
```

### 问题 3: 配置文件不存在

**错误信息:**
```
⚠️  配置文件不存在: /path/to/config.yaml
```

**解决方案:**
```bash
# 创建配置文件
cp config.yaml.example config.yaml

# 或使用命令行参数指定配置
python weapm_cli.py --base-url http://localhost:8080 dashboard
```

### 问题 4: 连接超时

**解决方案:**
```bash
# 增加超时时间
python weapm_cli.py --timeout 60 dashboard
```

---

## 高级用法

### 脚本自动化

```bash
#!/bin/bash
# monitor.sh - 监控脚本

# 检查数据大盘
echo "=== 检查数据大盘 ==="
python weapm_cli.py dashboard

# 检查所有集群
echo "=== 检查所有集群 ==="
python weapm_cli.py clusters

# 检查所有子系统
echo "=== 检查所有子系统 ==="
python weapm_cli.py subsystems
```

### JSON 处理

结合 `jq` 工具处理 JSON 输出:

```bash
# 提取子系统数量
python weapm_cli.py dashboard | jq '.result.subsystemCount'

# 提取所有集群名称
python weapm_cli.py clusters | jq '.[].clustername'

# 过滤特定集群
python weapm_cli.py clusters | jq '.[] | select(.clustername == "LOG001")'
```

### 定时任务

使用 cron 定时执行监控:

```cron
# 每5分钟检查一次数据大盘
*/5 * * * * /path/to/weapm_cli.py dashboard > /var/log/weapm.log 2>&1

# 每小时检查一次集群状态
0 * * * * /path/to/weapm_cli.py clusters >> /var/log/weapm_clusters.log
```

---

## 更多帮助

查看完整帮助信息:

```bash
# Python
python weapm_cli.py --help

# Golang
./weapm_cli --help
```

查看特定命令帮助:

```bash
# Python
python weapm_cli.py clusters --help
python weapm_cli.py add-node --help

# Golang
./weapm_cli clusters --help
./weapm_cli add-node --help
```
