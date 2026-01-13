# WEAPM-LOGSERVER API Documentation

## 概述

本文档描述了 WEAPM-LOGSERVER 的 REST API 接口。WEAPM-LOGSERVER 是一个基于 Java 的日志管理系统，负责将应用程序日志导入到 WEAPM-LOSTORE 存储中。

## API 基础路径

所有 API 端点都以 `/operation` 为前缀。

## 认证

API 使用 Basic Auth 进行认证。
username:weapmUser
password:Weapm@123admin

## 接口列表

### 数据大盘

#### 获取数据大盘信息
- **URL**: `GET /operation/dashboard`
- **描述**: 获取数据大盘页面信息，包括接入子系统数、部门数、所有集群整体流量图、集群状态图、topK子系统数据等
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": {
      "subsystemCount": 0,
      "clusterNum": 0,
      "clusterTrafficData": [
        {
          "clusterName": "string",
          "trafficBytes": 0,
          "timestamp": "string"
        }
      ],
      "topSubsystems": [
        {
          "department": "string",
          "subsys_name": "string",
          "business_owner": "string",
          "subsystem_owner": "string",
          "subsys_id": "string",
          "cluster_name": "string",
          "total_log_mb": 0
        }
      ],
      "clusterLogCounts": [
        {
          "clustername": "string",
          "total_log_gb": 0,
          "capacity": 0
        }
      ]
    }
  }
  ```

### 集群管理

#### 获取所有集群信息
- **URL**: `GET /operation/clusters`
- **描述**: 获取所有集群信息
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": [
      {
        "clustername": "string",
        "isdefault": 0,
        "topic": "string",
        "bucketnames": "string",
        "backenddomain": "string",
        "storagedomain": "string"
      }
    ]
  }
  ```

#### 获取集群详情
- **URL**: `GET /operation/clusters/{clusterName}`
- **描述**: 获取指定集群的详细信息
- **参数**:
  - `clusterName` (path): 集群名称
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": {
      "clusterInfo": {
        "clustername": "string",
        "isdefault": 0,
        "topic": "string",
        "bucketnames": "string",
        "backenddomain": "string",
        "storagedomain": "string"
      },
      "nodeGroups": [
        {
          "role": "string",
          "nodes": [
            {
              "address": "string",
              "clustername": "string",
              "role": "string",
              "topic": "string",
              "bucketnames": "string",
              "backenddomain": "string",
              "storagedomain": "string",
              "isdefault": true,
              "status": "string",
              "cpulimit": "string",
              "memlimit": "string",
              "createtime": "2025-12-26T00:00:00.000Z",
              "updateime": "2025-12-26T00:00:00.000Z"
            }
          ]
        }
      ],
      "managedSubSystems": [
        {
          "clustername": "string",
          "subsystemid": "string",
          "subsys_name": "string",
          "subsystem_owner": "string",
          "business_owner": "string",
          "devdept": "string",
          "traffic": 0,
          "status": "string",
          "createtime": "2025-12-26T00:00:00.000Z",
          "updatetime": "2025-12-26T00:00:00.000Z"
        }
      ],
      "reportData": {
        "peakTraffic": 0,
        "peakTime": "string",
        "totalSubSystems": 0,
        "topicBacklog": 0
      }
    }
  }
  ```

#### 向集群添加节点
- **URL**: `POST /operation/clusters/{clusterName}/nodes`
- **描述**: 向指定集群添加新节点
- **参数**:
  - `clusterName` (path): 集群名称
- **请求体**:
  ```json
  {
    "address": "string",
    "clustername": "string",
    "role": "string",
    "topic": "string",
    "bucketnames": "string",
    "backenddomain": "string",
    "storagedomain": "string",
    "isdefault": true,
    "status": "string",
    "cpulimit": "string",
    "memlimit": "string",
    "createtime": "2025-12-26T00:00:00.000Z",
    "updateime": "2025-12-26T00:00:00.000Z"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "message": "节点添加成功"
  }
  ```

#### 从集群删除节点
- **URL**: `DELETE /operation/clusters/nodes/{ip}`
- **描述**: 从集群删除指定IP的节点
- **参数**:
  - `ip` (path): 节点IP地址
- **响应**:
  ```json
  {
    "code": 0,
    "message": "节点删除成功"
  }
  ```

#### 获取集群纳管子系统信息
- **URL**: `GET /operation/cluster/{clusterName}/subsystems`
- **描述**: 获取指定集群纳管的子系统信息
- **参数**:
  - `clusterName` (path): 集群名称
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": [
      {
        "clustername": "string",
        "subsystemid": "string",
        "subsys_name": "string",
        "subsystem_owner": "string",
        "business_owner": "string",
        "devdept": "string",
        "traffic": 0,
        "status": "string",
        "createtime": "2025-12-26T00:00:00.000Z",
        "updatetime": "2025-12-26T00:00:00.000Z"
      }
    ]
  }
  ```

### 子系统运维

#### 检查子系统是否存在
- **URL**: `GET /operation/subsystem/exists/{subsystemId}`
- **描述**: 查询指定子系统是否已接入
- **参数**:
  - `subsystemId` (path): 子系统ID
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": {
      "subsystemId": "string",
      "exists": true,
      "subsystemName": "string",
      "clusterName": "string"
    }
  }
  ```

#### 新增子系统接入
- **URL**: `POST /operation/subsystem`
- **描述**: 新增子系统接入
- **请求体**:
  ```json
  {
    "subSystemId": "string",
    "logImportValue": "string",
    "logImportFiles": "string",
    "traffic": 0,
    "cluster": "string"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "message": "子系统接入成功"
  }
  ```

#### 调整子系统归属集群
- **URL**: `POST /operation/subsystem/{subsystemId}`
- **描述**: 调整子系统归属的小集群
- **参数**:
  - `subsystemId` (path): 子系统ID
  - `targetClusterName` (path): 目标集群名称
  - `logImportValue` (path): 关键字
  - `logImportFiles` (path): 文件列表
  - `traffic` (path): 流量
- **响应**:
  ```json
  {
    "code": 0,
    "message": "子系统集群调整成功"
  }
  ```

#### 调整子系统状态
- **URL**: `POST /operation/subsystem/{subsystemId}/status/{status}`
- **描述**: 临时调整子系统接入
- **参数**:
  - `subsystemId` (path): 子系统ID
  - `status` (path): 状态，disable/enable
- **响应**:
  ```json
  {
    "code": 0,
    "message": "子系统已禁用"
  }
  ```

#### 启用子系统
- **URL**: `PUT /operation/subsystem/{subsystemId}/enable`
- **描述**: 启用子系统接入
- **参数**:
  - `subsystemId` (path): 子系统ID
- **响应**:
  ```json
  {
    "code": 0,
    "message": "子系统已启用"
  }
  ```

#### 获取子系统详情
- **URL**: `GET /operation/subsystem/{subsystemId}`
- **描述**: 查询指定子系统的详细信息
- **参数**:
  - `subsystemId` (path): 子系统ID
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": {
      "subsystemInfo": {
        "id": 0,
        "subsys_id": "string",
        "subsys_name": "string",
        "subsys_chtname": "string",
        "subsys_updtime": "string",
        "devdept": "string",
        "business_owner": "string",
        "subsystem_owner": "string",
        "system_name": "string",
        "state": "string",
        "important_level": "string",
        "create_topic": "string"
      },
      "collected": true,
      "scanFileWhitelist": [
        "string"
      ],
      "expectedTraffic": 0,
      "actualTraffic": 0,
      "keywordFilters": [
        "string"
      ],
      "clusterName": "string",
      "instances": [
        {
          "additionalProp1": [
            "string"
          ]
        }
      ]
    }
  }
  ```

#### 获取所有子系统信息
- **URL**: `GET /operation/subsystems`
- **描述**: 获取系统中所有子系统的信息
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": [
      {
        "id": 0,
        "subsys_id": "string",
        "subsys_name": "string",
        "subsys_chtname": "string",
        "subsys_updtime": "string",
        "devdept": "string",
        "business_owner": "string",
        "subsystem_owner": "string",
        "system_name": "string",
        "state": "string",
        "important_level": "string",
        "create_topic": "string"
      }
    ]
  }
  ```

#### 根据条件搜索子系统
- **URL**: `GET /operation/subsystems/search`
- **描述**: 根据筛选条件获取子系统信息，默认返回前20个匹配的结果
- **参数**:
  - `subsysId` (query, optional): 子系统ID
  - `limit` (query, optional, default: 20): 返回结果数量限制
- **响应**:
  ```json
  {
    "code": 0,
    "message": "success",
    "result": [
      {
        "id": 0,
        "subsys_id": "string",
        "subsys_name": "string",
        "subsys_chtname": "string",
        "subsys_updtime": "string",
        "devdept": "string",
        "business_owner": "string",
        "subsystem_owner": "string",
        "system_name": "string",
        "state": "string",
        "important_level": "string",
        "create_topic": "string"
      }
    ]
  }
  ```

## 数据模型

### DashboardResult
| 字段名 | 类型 | 描述 |
|--------|------|------|
| subsystemCount | integer | 接入子系统数 |
| clusterNum | integer | 集群部门数 |
| clusterTrafficData | array | 瞬时集群流量数据 |
| topSubsystems | array | TopK子系统数据 |
| clusterLogCounts | array | 当天集群流量统计数据 |

### LogClusterInfo
| 字段名 | 类型 | 描述 |
|--------|------|------|
| clustername | string | 唯一集群名称，例如 LOG001, LOG002 |
| isdefault | integer | 默认集群，1为默认，0为子集群 |
| topic | string | wemq topic，与集群不同 |
| bucketnames | string | 存储桶名称 |
| backenddomain | string | logstore后端组件域 |
| storagedomain | string | logstore读写存储域 |

### SubSystem
| 字段名 | 类型 | 描述 |
|--------|------|------|
| id | integer |  |
| subsys_id | string | 子系统ID |
| subsys_name | string | 子系统名称 |
| subsys_chtname | string | 子系统中文名称 |
| subsys_updtime | string | 更新日期 |
| devdept | string | 开发部门 |
| business_owner | string | 业务负责人 |
| subsystem_owner | string | 子系统负责人 |
| system_name | string | 系统名称 |
| state | string | 状态 |
| important_level | string | 重要等级 |
| create_topic | string |  |

### AddSubsystemRequest
| 字段名 | 类型 | 描述 |
|--------|------|------|
| subSystemId | string |  |
| logImportValue | string |  |
| logImportFiles | string |  |
| traffic | integer |  |
| cluster | string |  |

### SubsystemExistsResult
| 字段名 | 类型 | 描述 |
|--------|------|------|
| subsystemId | string | 子系统ID |
| exists | boolean | 是否存在 |
| subsystemName | string | 子系统名称（如果存在） |
| clusterName | string | 归属集群（如果存在） |

### SubsystemDetailResult
| 字段名 | 类型 | 描述 |
|--------|------|------|
| subsystemInfo | SubSystem | 子系统基本信息 |
| collected | boolean | 是否已采集 |
| scanFileWhitelist | array | 扫描文件白名单 |
| expectedTraffic | integer | 预登记流量 |
| actualTraffic | integer | 实际流量 |
| keywordFilters | array | 关键字过滤规则 |
| clusterName | string |  |
| instances | array | 采集的实例信息，包括文件，ip，dcn |

### ClusterDetailResult
| 字段名 | 类型 | 描述 |
|--------|------|------|
| clusterInfo | LogClusterInfo | 集群基本信息 |
| nodeGroups | array | 节点列表（按角色分组） |
| managedSubSystems | array | 集群纳管子系统信息 |
| reportData | ClusterReportData | 集群报表数据 |

### LogStoreInstance
| 字段名 | 类型 | 描述 |
|--------|------|------|
| address | string |  |
| clustername | string |  |
| role | string |  |
| topic | string | 为子集群添加的topic |
| bucketnames | string | 存储桶名称 |
| backenddomain | string | logstore后端组件域 |
| storagedomain | string | logstore读写存储域 |
| isdefault | boolean | 是否为默认集群，不与子系统分组 |
| status | string |  |
| cpulimit | string |  |
| memlimit | string |  |
| createtime | string |  |
| updateime | string |  |

### LogSubClusterSubSystem
| 字段名 | 类型 | 描述 |
|--------|------|------|
| clustername | string |  |
| subsystemid | string |  |
| subsys_name | string |  |
| subsystem_owner | string |  |
| business_owner | string |  |
| devdept | string |  |
| traffic | integer |  |
| status | string |  |
| createtime | string |  |
| updatetime | string |  |

### ClusterLogCount
| 字段名 | 类型 | 描述 |
|--------|------|------|
| clustername | string |  |
| total_log_gb | integer |  |
| capacity | integer |  |

### SubsystemLogDetailCount
| 字段名 | 类型 | 描述 |
|--------|------|------|
| department | string |  |
| subsys_name | string |  |
| business_owner | string |  |
| subsystem_owner | string |  |
| subsys_id | string |  |
| cluster_name | string |  |
| total_log_mb | integer |  |

## 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "code": 1,
  "message": "错误描述信息"
}
```