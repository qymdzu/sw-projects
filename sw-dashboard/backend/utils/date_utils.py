"""
日期时间工具 — 时间戳格式化、运行时长解析
"""

from datetime import datetime, timezone, timedelta
from typing import Optional, Union


def format_timestamp(
    timestamp: Optional[Union[float, str]],
    fmt: str = "ISO8601",
) -> str:
    """
    格式化时间戳为 ISO 8601 字符串。

    Args:
        timestamp: 浮点数秒数 / ISO 字符串 / None
        fmt: 输出格式（目前仅支持 ISO8601）

    Returns:
        格式化后的时间字符串，无效输入返回空字符串
    """
    if timestamp is None:
        return ""

    try:
        if isinstance(timestamp, (int, float)):
            dt = datetime.fromtimestamp(timestamp, tz=timezone(timedelta(hours=8)))
        elif isinstance(timestamp, str):
            # Try parsing ISO format
            try:
                dt = datetime.fromisoformat(timestamp)
            except ValueError:
                # Try parsing common formats
                for fmt_str in [
                    "%Y-%m-%dT%H:%M:%S",
                    "%Y-%m-%d %H:%M:%S",
                    "%Y/%m/%d %H:%M:%S",
                ]:
                    try:
                        dt = datetime.strptime(timestamp, fmt_str)
                        dt = dt.replace(tzinfo=timezone(timedelta(hours=8)))
                        break
                    except ValueError:
                        continue
                else:
                    return ""
        else:
            return ""

        return dt.strftime("%Y-%m-%dT%H:%M:%S+08:00")
    except (ValueError, OSError, TypeError):
        return ""


def parse_uptime(seconds: Union[float, str, None]) -> str:
    """
    将秒数解析为人类可读的运行时长。

    Args:
        seconds: 秒数

    Returns:
        "X 天, X 小时, X 分钟" 格式
    """
    if seconds is None:
        return "未知"

    try:
        secs = float(seconds)
    except (TypeError, ValueError):
        return "未知"

    if secs < 0:
        return "未知"

    days = int(secs // 86400)
    hours = int((secs % 86400) // 3600)
    minutes = int((secs % 3600) // 60)

    parts = []
    if days > 0:
        parts.append(f"{days} 天")
    if hours > 0:
        parts.append(f"{hours} 小时")
    if minutes > 0 or (days == 0 and hours == 0):
        parts.append(f"{minutes} 分钟")

    return ", ".join(parts) if parts else "0 分钟"
