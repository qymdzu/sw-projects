"""
路径安全工具 — 路径校验、仓库路径解析

防止路径穿越攻击，确保所有文件操作限定在白名单路径内。
"""

import os
from pathlib import Path
from typing import Dict, List, Optional, Union

from config import AppConfig


def safe_resolve(base_path: str, user_path: str) -> Optional[Path]:
    """
    安全解析路径：确保 user_path 在 base_path 白名单内。

    Args:
        base_path: 白名单基路径
        user_path: 用户传入的路径字符串

    Returns:
        安全则返回 Path 对象，否则返回 None
    """
    if not base_path or not user_path:
        return None
    try:
        base = Path(base_path).resolve()
        target = Path(user_path).resolve()
        target.relative_to(base)
        return target
    except (ValueError, PermissionError, OSError, RuntimeError):
        return None


def is_within_base(target_path: str, base_path: str) -> bool:
    """判断路径是否在基路径白名单内"""
    return safe_resolve(base_path, target_path) is not None


def is_within_any_safe_path(target_path: str) -> Optional[str]:
    """
    判断路径是否在任意一个安全路径白名单内。

    Returns:
        匹配到的白名单基路径，或 None
    """
    safe_paths = AppConfig.get_safe_paths()
    target = Path(target_path).resolve()
    for sp in safe_paths:
        try:
            target.relative_to(Path(sp).resolve())
            return sp
        except (ValueError, OSError):
            continue
    return None


def get_repo_path(repo_name: str) -> Optional[Union[str, List[str]]]:
    """
    获取仓库对应的基路径。

    Args:
        repo_name: "workspace" | "skills-library" | "projects" | "hermes" | "all"

    Returns:
        单仓库: 路径字符串
        "all": 三个工作路径的列表
        无效: None
    """
    base_paths = AppConfig.base_paths
    if repo_name == "all":
        return [
            base_paths.get("workspace", ""),
            base_paths.get("skills-library", ""),
            base_paths.get("projects", ""),
        ]
    return base_paths.get(repo_name)


def expand_user_path(path_str: str) -> str:
    """展开路径中的 ~ 为用户目录"""
    return str(Path(path_str).expanduser().resolve())
