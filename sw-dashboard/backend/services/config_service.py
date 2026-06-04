"""
配置管理服务 — 配置文件读写、.env 掩码、Cron 管理
"""

import os
import re
import time
from pathlib import Path
from typing import List, Optional

import aiofiles
import yaml

from config import AppConfig
from models.schemas import (
    ConfigFile,
    ConfigListItem,
    CronJob,
    CronToggleResponse,
    FileWriteResponse,
)
from services.file_service import write_with_backup
from utils.cmd_utils import run_crontab
from utils.path_utils import is_within_any_safe_path


class ConfigError(Exception):
    """配置异常"""
    pass


# 敏感 Key 正则（用于 .env 掩码）
SENSITIVE_KEYS = re.compile(
    r"(API_KEY|TOKEN|SECRET|PASSWORD|PRIVATE_KEY|ACCESS_KEY|"
    r"SECRET_KEY|API_SECRET|APP_SECRET|DB_PASSWORD|"
    r"REDIS_PASSWORD|MYSQL_PASSWORD|POSTGRES_PASSWORD)"
)


async def read_config(config_name: str) -> ConfigFile:
    """
    读取配置文件。

    Args:
        config_name: "config.yaml" | ".env" | "CLAUDE.md" | "SOUL.md"

    Returns:
        ConfigFile 对象
    """
    if config_name not in ["config.yaml", ".env", "CLAUDE.md", "SOUL.md"]:
        raise ConfigError(f"不支持的配置文件: {config_name}")

    config_path = AppConfig.get_config_file_path(config_name)
    if not config_path:
        raise ConfigError(f"配置 {config_name} 的路径未定义")

    resolved = Path(config_path).resolve()
    file_exists = resolved.exists()

    if config_name == ".env":
        return await read_env(masked=True)

    content = ""
    language = "text"
    last_modified = None

    if file_exists:
        try:
            async with aiofiles.open(resolved, mode="r", encoding="utf-8", errors="replace") as f:
                content = await f.read()
            last_modified = resolved.stat().st_mtime
        except PermissionError:
            raise ConfigError(f"无权限读取: {config_path}")
        except Exception:
            pass

    # 设置语言
    if config_name == "config.yaml":
        language = "yaml"
    elif config_name in ("CLAUDE.md", "SOUL.md"):
        language = "markdown"

    return ConfigFile(
        name=config_name,
        path=str(resolved),
        content=content,
        language=language,
        editable=True,
        last_modified=last_modified,
    )


async def read_env(masked: bool = True) -> ConfigFile:
    """
    读取 .env 文件，支持掩码敏感信息。

    Args:
        masked: 是否掩码 API Key 等敏感信息

    Returns:
        ConfigFile
    """
    config_path = AppConfig.get_config_file_path(".env")
    if not config_path:
        raise ConfigError(".env 路径未定义")

    resolved = Path(config_path).resolve()
    content = ""
    last_modified = None

    if resolved.exists():
        try:
            async with aiofiles.open(resolved, mode="r", encoding="utf-8", errors="replace") as f:
                raw_content = await f.read()
            last_modified = resolved.stat().st_mtime
            if masked:
                content = mask_secrets(raw_content)
            else:
                content = raw_content
        except PermissionError:
            raise ConfigError(f"无权限读取: {config_path}")
        except Exception:
            pass

    return ConfigFile(
        name=".env",
        path=str(resolved),
        content=content,
        language="env",
        editable=True,
        last_modified=last_modified,
    )


async def write_config(
    config_name: str,
    content: str,
    validate: bool = True,
    create_backup: bool = True,
) -> FileWriteResponse:
    """
    保存配置文件。

    Args:
        config_name: 配置文件名
        content: 新内容
        validate: 是否校验格式
        create_backup: 是否创建备份

    Returns:
        FileWriteResponse
    """
    if config_name not in ["config.yaml", ".env", "CLAUDE.md", "SOUL.md"]:
        raise ConfigError(f"不支持的配置文件: {config_name}")

    config_path = AppConfig.get_config_file_path(config_name)
    if not config_path:
        raise ConfigError(f"配置 {config_name} 的路径未定义")

    if validate:
        if config_name == "config.yaml":
            try:
                yaml.safe_load(content)
            except yaml.YAMLError as e:
                raise ConfigError(f"YAML 格式错误: {e}")
        elif config_name == ".env":
            _validate_env_content(content)

    return await write_with_backup(
        file_path=config_path,
        content=content,
        create_backup=create_backup,
    )


def mask_secrets(content: str) -> str:
    """
    掩码 .env 文件中的敏感信息。

    规则：
    - 保留前4位 + "****" + 保留后4位
    - 值长度 < 8 时全部替换为 "****"
    """
    if not content:
        return ""

    lines = content.split("\n")
    masked_lines = []

    for line in lines:
        if "=" not in line:
            masked_lines.append(line)
            continue

        key, _, value = line.partition("=")
        key_stripped = key.strip()

        if SENSITIVE_KEYS.search(key_stripped):
            value_stripped = value.strip().strip("\"'")
            if len(value_stripped) >= 8:
                masked = value_stripped[:4] + "****" + value_stripped[-4:]
            elif value_stripped:
                masked = "****"
            else:
                masked = value_stripped
            masked_lines.append(f"{key}={masked}")
        else:
            masked_lines.append(line)

    return "\n".join(masked_lines)


def _validate_env_content(content: str) -> None:
    """校验 .env 文件格式"""
    for i, line in enumerate(content.split("\n"), start=1):
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if "=" not in stripped:
            raise ConfigError(f".env 格式错误，第 {i} 行缺少 '=': {line}")
        key = stripped.split("=", 1)[0].strip()
        if not key:
            raise ConfigError(f".env 格式错误，第 {i} 行 Key 为空")


async def list_configs() -> List[ConfigListItem]:
    """获取配置文件列表（含 last_modified）"""
    configs = []
    for name in ["config.yaml", ".env", "CLAUDE.md", "SOUL.md"]:
        config_path = AppConfig.get_config_file_path(name)
        language = "yaml" if name == "config.yaml" else ("markdown" if name.endswith(".md") else "env")
        last_modified = None
        if config_path:
            p = Path(config_path)
            if p.exists():
                try:
                    last_modified = p.stat().st_mtime
                except OSError:
                    pass
        configs.append(ConfigListItem(
            name=name,
            language=language,
            editable=True,
            last_modified=last_modified,
        ))
    return configs


async def list_crontab() -> List[CronJob]:
    """获取 Cron 任务列表"""
    result = await run_crontab(action="list")
    if not result.get("success", False):
        return []
    jobs_data = result.get("jobs", [])
    return [CronJob(**job) for job in jobs_data]


async def toggle_cron(command: str, enable: bool) -> CronToggleResponse:
    """启停 Cron 任务"""
    result = await run_crontab(action="toggle", command=command, enable=enable)
    return CronToggleResponse(
        success=result.get("success", False),
        message=result.get("message", ""),
    )
