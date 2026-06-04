"""
依赖注入 — Token 认证、路径依赖
"""

from typing import Dict

from fastapi import Header, HTTPException, Request

from config import AppConfig


async def verify_token(x_token: str = Header(None, alias="X-Token")) -> str:
    """
    Token 认证依赖。

    从请求头 X-Token 中获取 Token 并与配置对比。

    Raises:
        HTTPException(401): Token 缺失或无效
    """
    if not x_token:
        raise HTTPException(
            status_code=401,
            detail="Missing X-Token header",
        )
    if x_token != AppConfig.TOKEN:
        raise HTTPException(
            status_code=401,
            detail="Invalid token",
        )
    return x_token


def get_base_paths() -> Dict[str, str]:
    """获取三库路径字典"""
    return dict(AppConfig.base_paths)


def get_log_path() -> str:
    """获取日志目录路径"""
    return AppConfig.log_path


def get_skills_path() -> str:
    """获取技能目录路径"""
    return AppConfig.skills_path
