#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
WEAPM-LOGSERVER API 客户端

基于 WEAPM-LOGSERVER REST API 的 Python 客户端实现
支持数据大盘、集群管理、子系统运维等功能
"""

import argparse
import requests
from typing import Optional, Dict, List, Any
from dataclasses import dataclass
import json
import os
import sys
import yaml
import logging
import time
from pathlib import Path
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry


# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@dataclass
class WeapmConfig:
    """WEAPM API 配置"""
    base_url: str
    username: str
    password: str
    timeout: int = 30
    max_retries: int = 3
    retry_backoff_factor: float = 0.5
    pool_connections: int = 10
    pool_maxsize: int = 10
    enable_logging: bool = True

    @classmethod
    def from_dict(cls, config_dict: Dict[str, Any]) -> 'WeapmConfig':
        """从字典创建配置对象"""
        return cls(
            base_url=config_dict.get('base_url', ''),
            username=config_dict.get('username', 'weapmUser'),
            password=config_dict.get('password', 'Weapm@123admin'),
            timeout=config_dict.get('timeout', 30),
            max_retries=config_dict.get('max_retries', 3),
            retry_backoff_factor=config_dict.get('retry_backoff_factor', 0.5),
            pool_connections=config_dict.get('pool_connections', 10),
            pool_maxsize=config_dict.get('pool_maxsize', 10),
            enable_logging=config_dict.get('enable_logging', True)
        )

    @classmethod
    def from_yaml(cls, config_path: str = None, env: str = None) -> 'WeapmConfig':
        """
        从 YAML 配置文件加载配置

        Args:
            config_path: 配置文件路径,默认为 script/weapm/config.yaml
            env: 环境名称 (dev/prod),默认使用配置文件中的 active_env

        Returns:
            WeapmConfig 配置对象

        Raises:
            FileNotFoundError: 配置文件不存在
            ValueError: 环境配置不存在或格式错误
        """
        # 默认配置文件路径
        if config_path is None:
            script_dir = Path(__file__).parent
            config_path = script_dir / "config.yaml"

        # 检查配置文件是否存在
        if not os.path.exists(config_path):
            raise FileNotFoundError(f"配置文件不存在: {config_path}")

        # 读取配置文件
        with open(config_path, 'r', encoding='utf-8') as f:
            config_data = yaml.safe_load(f)

        # 确定使用的环境
        if env is None:
            env = config_data.get('active_env', 'dev')

        # 获取环境配置
        if env not in config_data:
            raise ValueError(f"环境配置不存在: {env}, 可用环境: {list(config_data.keys())}")

        env_config = config_data[env]

        # 验证必要字段
        if 'base_url' not in env_config:
            raise ValueError(f"环境 {env} 缺少必要字段: base_url")

        print(f"✅ 加载配置: {env_config.get('description', env)} ({env})")
        return cls.from_dict(env_config)


class WeapmClient:
    """WEAPM-LOGSERVER API 客户端"""

    def __init__(self, config: WeapmConfig):
        """
        初始化客户端

        Args:
            config: WEAPM 配置对象
        """
        self.config = config
        self.session = requests.Session()
        self.session.auth = (config.username, config.password)
        self.base_url = config.base_url.rstrip('/')
        self.timeout = config.timeout

        # 配置重试策略
        retry_strategy = Retry(
            total=config.max_retries,
            backoff_factor=config.retry_backoff_factor,
            status_forcelist=[429, 500, 502, 503, 504],
            method_whitelist=["HEAD", "GET", "OPTIONS", "POST", "PUT", "DELETE"]
        )

        # 配置连接池
        adapter = HTTPAdapter(
            max_retries=retry_strategy,
            pool_connections=config.pool_connections,
            pool_maxsize=config.pool_maxsize
        )
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)

        logger.info(f"WEAPM 客户端初始化成功: {self.base_url}")

    def _request(self, method: str, endpoint: str, **kwargs) -> Dict[str, Any]:
        """
        发送 HTTP 请求

        Args:
            method: HTTP 方法 (GET, POST, PUT, DELETE)
            endpoint: API 端点路径
            **kwargs: 其他请求参数

        Returns:
            API 响应数据

        Raises:
            requests.RequestException: 请求失败
        """
        url = f"{self.base_url}{endpoint}"
        kwargs.setdefault('timeout', self.timeout)

        # 记录请求信息
        if self.config.enable_logging:
            logger.info(f"发送请求: {method} {url}")

        start_time = time.time()

        try:
            response = self.session.request(method, url, **kwargs)
            elapsed_time = time.time() - start_time

            # 记录响应信息
            if self.config.enable_logging:
                logger.info(
                    f"收到响应: {method} {url} - "
                    f"状态码: {response.status_code}, "
                    f"耗时: {elapsed_time:.2f}s"
                )

            response.raise_for_status()

            # 验证响应内容
            try:
                data = response.json()
            except json.JSONDecodeError as e:
                logger.error(f"JSON 解析失败: {str(e)}")
                raise ValueError(f"无效的 JSON 响应: {response.text[:200]}")

            # 检查业务错误码
            if isinstance(data, dict) and 'code' in data:
                if data['code'] != 0:
                    error_msg = data.get('message', '未知错误')
                    logger.error(f"API 业务错误: code={data['code']}, message={error_msg}")
                    raise requests.HTTPError(
                        f"API 错误 (code {data['code']}): {error_msg}"
                    )

            return data

        except requests.Timeout as e:
            logger.error(f"请求超时: {method} {url} (超时时间: {self.timeout}s)")
            raise
        except requests.ConnectionError as e:
            logger.error(f"连接错误: {method} {url} - {str(e)}")
            raise
        except requests.HTTPError as e:
            logger.error(f"HTTP 错误: {method} {url} - {str(e)}")
            raise
        except requests.RequestException as e:
            logger.error(f"请求失败: {method} {url} - {str(e)}")
            if hasattr(e, 'response') and e.response is not None:
                logger.error(f"响应内容: {e.response.text[:500]}")
            raise

    # ==================== 数据大盘 ====================

    def get_dashboard(self) -> Dict[str, Any]:
        """
        获取数据大盘信息

        Returns:
            包含子系统数、集群数、流量数据、TopK子系统等信息的字典
        """
        return self._request('GET', '/operation/dashboard')

    # ==================== 集群管理 ====================

    def get_clusters(self) -> List[Dict[str, Any]]:
        """
        获取所有集群信息

        Returns:
            集群信息列表
        """
        response = self._request('GET', '/operation/clusters')
        return response.get('result', [])

    def get_cluster_detail(self, cluster_name: str) -> Dict[str, Any]:
        """
        获取指定集群的详细信息

        Args:
            cluster_name: 集群名称

        Returns:
            包含集群信息、节点列表、子系统、报表数据的字典
        """
        return self._request('GET', f'/operation/clusters/{cluster_name}')

    def add_cluster_node(self, cluster_name: str, address: str, role: str,
                         cpulimit: Optional[str] = None, memlimit: Optional[str] = None,
                         **kwargs) -> Dict[str, Any]:
        """
        向集群添加节点

        Args:
            cluster_name: 集群名称
            address: 节点IP地址
            role: 节点角色
            cpulimit: CPU限制 (可选)
            memlimit: 内存限制 (可选)
            **kwargs: 其他可选字段 (topic, bucketnames, backenddomain, storagedomain,
                         isdefault, status, createtime, updateime等)

        Returns:
            操作结果

        Example:
            # 最小化参数
            client.add_cluster_node("LOG008", "127.0.0.2", "write", cpulimit="8", memlimit="16")

            # 完整参数
            client.add_cluster_node(
                cluster_name="LOG008",
                address="127.0.0.2",
                role="write",
                cpulimit="8",
                memlimit="16",
                topic="log_topic",
                bucketnames="log_bucket",
                backenddomain="backend.example.com",
                storagedomain="storage.example.com",
                status="active"
            )
        """
        # 构建基础请求数据
        data = {
            "address": address,
            "clustername": cluster_name,
            "role": role
        }

        # 添加可选参数
        if cpulimit is not None:
            data["cpulimit"] = cpulimit
        if memlimit is not None:
            data["memlimit"] = memlimit

        # 添加其他可选字段
        data.update(kwargs)

        return self._request('POST', f'/operation/clusters/{cluster_name}/nodes', json=data)

    def delete_cluster_node(self, ip: str) -> Dict[str, Any]:
        """
        从集群删除节点

        Args:
            ip: 节点IP地址

        Returns:
            操作结果
        """
        return self._request('DELETE', f'/operation/clusters/nodes/{ip}')

    def get_cluster_subsystems(self, cluster_name: str) -> List[Dict[str, Any]]:
        """
        获取集群纳管的子系统信息

        Args:
            cluster_name: 集群名称

        Returns:
            子系统信息列表
        """
        response = self._request('GET', f'/operation/cluster/{cluster_name}/subsystems')
        return response.get('result', [])

    # ==================== 子系统运维 ====================

    def check_subsystem_exists(self, subsystem_id: str) -> Dict[str, Any]:
        """
        检查子系统是否存在

        Args:
            subsystem_id: 子系统ID

        Returns:
            包含存在状态、子系统名称、集群名称的字典
        """
        return self._request('GET', f'/operation/subsystem/exists/{subsystem_id}')

    def add_subsystem(self, sub_system_id: str, log_import_value: str,
                     log_import_files: str, traffic: int, cluster: str) -> Dict[str, Any]:
        """
        新增子系统接入

        Args:
            sub_system_id: 子系统ID
            log_import_value: 日志导入关键字
            log_import_files: 日志导入文件列表
            traffic: 流量
            cluster: 集群名称

        Returns:
            操作结果
        """
        data = {
            "subSystemId": sub_system_id,
            "logImportValue": log_import_value,
            "logImportFiles": log_import_files,
            "traffic": traffic,
            "cluster": cluster
        }
        return self._request('POST', '/operation/subsystem', json=data)

    def adjust_subsystem_cluster(self, subsystem_id: str, target_cluster_name: str,
                                log_import_value: str, log_import_files: str,
                                traffic: int) -> Dict[str, Any]:
        """
        调整子系统归属集群

        Args:
            subsystem_id: 子系统ID
            target_cluster_name: 目标集群名称
            log_import_value: 关键字
            log_import_files: 文件列表
            traffic: 流量

        Returns:
            操作结果
        """
        endpoint = (f'/operation/subsystem/{subsystem_id}?targetClusterName={target_cluster_name}'
                   f'&logImportValue={log_import_value}&logImportFiles={log_import_files}'
                   f'&traffic={traffic}')
        return self._request('POST', endpoint)

    def adjust_subsystem_status(self, subsystem_id: str, status: str) -> Dict[str, Any]:
        """
        调整子系统状态 (enable/disable)

        Args:
            subsystem_id: 子系统ID
            status: 状态 (enable/disable)

        Returns:
            操作结果
        """
        return self._request('POST', f'/operation/subsystem/{subsystem_id}/status/{status}')

    def enable_subsystem(self, subsystem_id: str) -> Dict[str, Any]:
        """
        启用子系统

        Args:
            subsystem_id: 子系统ID

        Returns:
            操作结果
        """
        return self._request('PUT', f'/operation/subsystem/{subsystem_id}/enable')

    def get_subsystem_detail(self, subsystem_id: str) -> Dict[str, Any]:
        """
        获取子系统详情

        Args:
            subsystem_id: 子系统ID

        Returns:
            包含子系统信息、采集状态、流量、实例等详细信息的字典
        """
        return self._request('GET', f'/operation/subsystem/{subsystem_id}')

    def get_subsystems(self) -> List[Dict[str, Any]]:
        """
        获取所有子系统信息

        Returns:
            子系统信息列表
        """
        response = self._request('GET', '/operation/subsystems')
        return response.get('result', [])

    def search_subsystems(self, subsys_id: Optional[str] = None, limit: int = 20) -> List[Dict[str, Any]]:
        """
        根据条件搜索子系统

        Args:
            subsys_id: 子系统ID (可选)
            limit: 返回结果数量限制 (默认20)

        Returns:
            匹配的子系统列表
        """
        params = {}
        if subsys_id:
            params['subsysId'] = subsys_id
        if limit != 20:
            params['limit'] = limit

        response = self._request('GET', '/operation/subsystems/search', params=params)
        return response.get('result', [])

    # ==================== 辅助方法 ====================

    def close(self):
        """关闭会话"""
        self.session.close()

    def __enter__(self):
        """支持上下文管理器"""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """退出上下文管理器"""
        self.close()


# ==================== 使用示例 ====================

def main():
    """主函数 - 演示 API 使用"""

    # 方式 1: 从配置文件加载 (推荐)
    try:
        # 自动使用 config.yaml 中的 active_env 配置
        config = WeapmConfig.from_yaml()

        # 或者指定环境
        # config = WeapmConfig.from_yaml(env="dev")
        # config = WeapmConfig.from_yaml(env="prod")

        # 或者指定配置文件路径
        # config = WeapmConfig.from_yaml("/path/to/config.yaml", env="prod")
    except FileNotFoundError as e:
        print(f"⚠️  {e}")
        print("请先创建配置文件 config.yaml,参考 config.yaml.example")
        return
    except ValueError as e:
        print(f"⚠️  配置错误: {e}")
        return

    # 方式 2: 手动创建配置 (备用方案)
    # config = WeapmConfig(
    #     base_url="http://localhost:8080",
    #     username="weapmUser",
    #     password="Weapm@123admin",
    #     timeout=30
    # )

    # 创建客户端
    with WeapmClient(config) as client:
        try:
            # 示例 1: 获取数据大盘信息
            print("=" * 50)
            print("1. 获取数据大盘信息")
            dashboard = client.get_dashboard()
            print(f"子系统数量: {dashboard.get('result', {}).get('subsystemCount')}")
            print(f"集群数量: {dashboard.get('result', {}).get('clusterNum')}")

            # 示例 2: 获取所有集群
            print("\n" + "=" * 50)
            print("2. 获取所有集群")
            clusters = client.get_clusters()
            for cluster in clusters:
                print(f"集群名称: {cluster.get('clustername')}, 默认: {cluster.get('isdefault')}")

            # 示例 3: 获取集群详情
            if clusters:
                cluster_name = clusters[0].get('clustername')
                print(f"\n" + "=" * 50)
                print(f"3. 获取集群详情: {cluster_name}")
                cluster_detail = client.get_cluster_detail(cluster_name)
                print(f"集群信息: {json.dumps(cluster_detail.get('result', {}).get('clusterInfo', {}), indent=2, ensure_ascii=False)}")

            # 示例 4: 获取所有子系统
            print("\n" + "=" * 50)
            print("4. 获取所有子系统")
            subsystems = client.get_subsystems()
            print(f"子系统总数: {len(subsystems)}")
            for subsystem in subsystems[:5]:  # 只显示前5个
                print(f"子系统ID: {subsystem.get('subsys_id')}, 名称: {subsystem.get('subsys_name')}")

            # 示例 5: 搜索子系统
            print("\n" + "=" * 50)
            print("5. 搜索子系统")
            search_results = client.search_subsystems(limit=10)
            print(f"搜索到 {len(search_results)} 个子系统")

            # 示例 6: 检查子系统是否存在
            if subsystems:
                subsystem_id = subsystems[0].get('subsys_id')
                print(f"\n" + "=" * 50)
                print(f"6. 检查子系统是否存在: {subsystem_id}")
                exists_result = client.check_subsystem_exists(subsystem_id)
                print(f"存在: {exists_result.get('result', {}).get('exists')}")

            # 示例 7: 向集群添加节点 (最小化参数)
            print("\n" + "=" * 50)
            print("7. 向集群添加节点 (最小化参数)")
            try:
                result = client.add_cluster_node(
                    cluster_name="LOG008",
                    address="127.0.0.2",
                    role="write",
                    cpulimit="8",
                    memlimit="16"
                )
                print(f"节点添加结果: {result.get('message')}")
            except Exception as e:
                print(f"添加节点失败: {str(e)}")

            # 示例 8: 向集群添加节点 (完整参数)
            print("\n" + "=" * 50)
            print("8. 向集群添加节点 (完整参数)")
            try:
                result = client.add_cluster_node(
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
                print(f"节点添加结果: {result.get('message')}")
            except Exception as e:
                print(f"添加节点失败: {str(e)}")

        except requests.RequestException as e:
            print(f"API 调用失败: {str(e)}")
        except Exception as e:
            print(f"发生错误: {str(e)}")


if __name__ == "__main__":
    main()
