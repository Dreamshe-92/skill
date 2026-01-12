#!/usr/bin/env python3
"""
获取本周的打卡记录文件
从当前目录下的 daily 目录中获取本周的打卡记录文件
"""

import os
import glob
from datetime import datetime, timedelta


def get_current_week_files(daily_dir="daily"):
    """
    获取本周的打卡记录文件

    Args:
        daily_dir: 打卡记录目录路径

    Returns:
        list: 本周的文件路径列表（按日期排序）
    """
    # 获取当前日期
    today = datetime.now()

    # 计算本周一的日期（周一作为一周的开始）
    monday = today - timedelta(days=today.weekday())

    # 计算本周日的日期
    sunday = monday + timedelta(days=6)

    # 获取所有文件
    pattern = os.path.join(daily_dir, "*")
    all_files = glob.glob(pattern)

    weekly_files = []

    for file_path in all_files:
        # 跳过目录
        if os.path.isdir(file_path):
            continue

        # 从文件名提取日期（格式：yyyymmdd）
        filename = os.path.basename(file_path)
        date_str = filename

        try:
            file_date = datetime.strptime(date_str, '%Y%m%d')

            # 检查文件日期是否在本周范围内
            if monday <= file_date <= sunday:
                weekly_files.append(file_path)
        except ValueError:
            # 如果文件名不符合日期格式，跳过
            continue

    # 按日期排序
    weekly_files.sort()

    return weekly_files


def read_file_content(file_path):
    """
    读取文件内容

    Args:
        file_path: 文件路径

    Returns:
        str: 文件内容
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()
    except Exception as e:
        return f"读取文件失败: {str(e)}"


if __name__ == "__main__":
    # 获取本周文件
    # 使用脚本所在目录的父目录的 daily 子目录
    script_dir = os.path.dirname(os.path.abspath(__file__))
    skill_dir = os.path.dirname(script_dir)
    daily_dir = os.path.join(skill_dir, "..", "..", "daily")

    print(f"脚本目录: {script_dir}")
    print(f"查找目录: {os.path.abspath(daily_dir)}")

    weekly_files = get_current_week_files(daily_dir)

    if not weekly_files:
        print("本周暂无打卡记录文件")
    else:
        print(f"\n找到 {len(weekly_files)} 个本周打卡记录文件:")
        print("=" * 60)

        for file_path in weekly_files:
            filename = os.path.basename(file_path)
            print(f"\n文件: {filename}")
            print("-" * 60)
            content = read_file_content(file_path)
            print(content)
