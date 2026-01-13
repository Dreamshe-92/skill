#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
WEAPM-LOGSERVER API 客户端命令行工具

基于 WEAPM-LOGSERVER REST API 的 Python 命令行客户端
支持数据大盘、集群管理、子系统运维等功能
"""

import argparse
import sys
import json
import logging
from pathlib import Path

# 导入主客户端模块
from weapm_client import WeapmClient, WeapmConfig


# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


# ==================== 命令处理函数 ====================

def cmd_dashboard(args, client):
    """获取数据大盘信息"""
    dashboard = client.get_dashboard()
    result = dashboard.get('result', {})
    print(json.dumps(result, indent=2, ensure_ascii=False))


def cmd_clusters(args, client):
    """获取集群列表"""
    if args.detail:
        if not args.cluster_name:
            print("错误: 使用 --detail 时必须指定 --cluster-name")
            sys.exit(1)
        result = client.get_cluster_detail(args.cluster_name)
        print(json.dumps(result.get('result', {}), indent=2, ensure_ascii=False))
    else:
        clusters = client.get_clusters()
        print(json.dumps(clusters, indent=2, ensure_ascii=False))


def cmd_subsystems(args, client):
    """获取子系统列表"""
    if args.search:
        result = client.search_subsystems(subsys_id=args.subsys_id, limit=args.limit)
        print(json.dumps(result, indent=2, ensure_ascii=False))
    elif args.check:
        result = client.check_subsystem_exists(args.check)
        print(json.dumps(result.get('result', {}), indent=2, ensure_ascii=False))
    elif args.detail:
        result = client.get_subsystem_detail(args.detail)
        print(json.dumps(result.get('result', {}), indent=2, ensure_ascii=False))
    else:
        result = client.get_subsystems()
        print(json.dumps(result, indent=2, ensure_ascii=False))


def cmd_add_node(args, client):
    """添加集群节点"""
    result = client.add_cluster_node(
        cluster_name=args.cluster_name,
        address=args.address,
        role=args.role,
        cpulimit=args.cpulimit,
        memlimit=args.memlimit,
        topic=args.topic,
        bucketnames=args.bucketnames,
        backenddomain=args.backenddomain,
        storagedomain=args.storagedomain,
        status=args.status
    )
    print(json.dumps(result, indent=2, ensure_ascii=False))


def cmd_delete_node(args, client):
    """删除集群节点"""
    result = client.delete_cluster_node(args.ip)
    print(json.dumps(result, indent=2, ensure_ascii=False))


# ==================== 参数解析 ====================

def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description='WEAPM-LOGSERVER API 客户端命令行工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 获取数据大盘信息
  python weapm_cli.py dashboard

  # 获取所有集群
  python weapm_cli.py clusters

  # 获取集群详情
  python weapm_cli.py clusters --detail --cluster-name LOG001

  # 获取所有子系统
  python weapm_cli.py subsystems

  # 搜索子系统
  python weapm_cli.py subsystems --search --subsys-id SYS001

  # 检查子系统是否存在
  python weapm_cli.py subsystems --check SYS001

  # 获取子系统详情
  python weapm_cli.py subsystems --detail SYS001

  # 添加集群节点 (最小参数)
  python weapm_cli.py add-node --cluster-name LOG008 --address 127.0.0.2 --role write --cpulimit 8 --memlimit 16

  # 添加集群节点 (完整参数)
  python weapm_cli.py add-node --cluster-name LOG008 --address 127.0.0.3 --role master \\
    --cpulimit 16 --memlimit 32 --topic log_topic --bucketnames log_bucket

  # 删除集群节点
  python weapm_cli.py delete-node --ip 192.168.1.100

  # 使用自定义 API 地址
  python weapm_cli.py --base-url http://192.168.1.100:8080 dashboard

  # 指定环境
  python weapm_cli.py --env prod clusters
        """
    )

    # 全局参数
    parser.add_argument('--config', '-c', help='配置文件路径')
    parser.add_argument('--env', '-e', choices=['dev', 'prod'], help='环境名称 (dev/prod)')
    parser.add_argument('--base-url', help='API 基础 URL')
    parser.add_argument('--username', help='用户名')
    parser.add_argument('--password', help='密码')
    parser.add_argument('--timeout', type=int, default=30, help='请求超时时间(秒)')
    parser.add_argument('--quiet', '-q', action='store_true', help='静默模式,不输出日志')
    parser.add_argument('--output', '-o', choices=['json', 'pretty'], default='pretty',
                       help='输出格式 (默认: pretty)')

    subparsers = parser.add_subparsers(dest='command', help='可用命令')

    # dashboard 命令
    dashboard_parser = subparsers.add_parser('dashboard', help='获取数据大盘信息')

    # clusters 命令
    clusters_parser = subparsers.add_parser('clusters', help='集群管理')
    clusters_parser.add_argument('--detail', '-d', action='store_true', help='显示集群详情')
    clusters_parser.add_argument('--cluster-name', '-n', help='集群名称')

    # subsystems 命令
    subsystems_parser = subparsers.add_parser('subsystems', help='子系统管理')
    subsystems_parser.add_argument('--search', '-s', action='store_true', help='搜索子系统')
    subsystems_parser.add_argument('--subsys-id', help='子系统ID (用于搜索)')
    subsystems_parser.add_argument('--check', '-c', help='检查子系统是否存在')
    subsystems_parser.add_argument('--detail', '-d', help='获取子系统详情')
    subsystems_parser.add_argument('--limit', '-l', type=int, default=20, help='返回结果数量限制')

    # add-node 命令
    add_node_parser = subparsers.add_parser('add-node', help='添加集群节点')
    add_node_parser.add_argument('--cluster-name', required=True, help='集群名称')
    add_node_parser.add_argument('--address', required=True, help='节点IP地址')
    add_node_parser.add_argument('--role', required=True, help='节点角色')
    add_node_parser.add_argument('--cpulimit', help='CPU限制')
    add_node_parser.add_argument('--memlimit', help='内存限制')
    add_node_parser.add_argument('--topic', help='Topic')
    add_node_parser.add_argument('--bucketnames', help='存储桶名称')
    add_node_parser.add_argument('--backenddomain', help='后端域')
    add_node_parser.add_argument('--storagedomain', help='存储域')
    add_node_parser.add_argument('--status', help='状态')

    # delete-node 命令
    delete_node_parser = subparsers.add_parser('delete-node', help='删除集群节点')
    delete_node_parser.add_argument('--ip', required=True, help='节点IP地址')

    return parser.parse_args()


# ==================== 主函数 ====================

def main():
    """主函数 - 命令行入口"""
    args = parse_args()

    # 如果没有指定命令,显示帮助
    if not args.command:
        print("请指定命令。使用 --help 查看帮助信息。")
        sys.exit(1)

    # 配置日志级别
    if args.quiet:
        logging.getLogger().setLevel(logging.ERROR)

    # 加载配置
    try:
        if args.config or args.env:
            config = WeapmConfig.from_yaml(args.config, args.env)
        elif args.base_url:
            # 使用命令行参数创建配置
            config = WeapmConfig(
                base_url=args.base_url,
                username=args.username or "weapmUser",
                password=args.password or "Weapm@123admin",
                timeout=args.timeout
            )
        else:
            # 默认使用配置文件
            config = WeapmConfig.from_yaml()
    except FileNotFoundError as e:
        print(f"⚠️  {e}")
        print("请先创建配置文件 config.yaml,参考 config.yaml.example")
        sys.exit(1)
    except ValueError as e:
        print(f"⚠️  配置错误: {e}")
        sys.exit(1)

    # 创建客户端
    with WeapmClient(config) as client:
        try:
            # 根据命令执行相应操作
            if args.command == 'dashboard':
                cmd_dashboard(args, client)
            elif args.command == 'clusters':
                cmd_clusters(args, client)
            elif args.command == 'subsystems':
                cmd_subsystems(args, client)
            elif args.command == 'add-node':
                cmd_add_node(args, client)
            elif args.command == 'delete-node':
                cmd_delete_node(args, client)
            else:
                print(f"未知命令: {args.command}")
                sys.exit(1)

        except Exception as e:
            print(f"❌ 错误: {str(e)}")
            sys.exit(1)


if __name__ == "__main__":
    main()
