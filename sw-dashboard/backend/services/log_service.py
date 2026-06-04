"""
日志服务 — 日志文件列表、tail 读取、搜索、归档
"""

import os
import time
from pathlib import Path
from typing import List, Optional

import aiofiles

from config import AppConfig
from models.schemas import LogArchive, LogContent, LogFile, LogFileListResponse, LogSearchResult


class LogError(Exception):
    """日志异常"""
    pass


async def list_log_files(log_dir: str) -> LogFileListResponse:
    """
    列出日志目录下的所有日志文件。

    Args:
        log_dir: 日志目录路径

    Returns:
        LogFileListResponse
    """
    log_path = Path(log_dir).resolve()
    files = []

    if not log_path.exists():
        return LogFileListResponse(
            files=[],
            log_path=str(log_path),
            default_file="hermes.log",
        )

    try:
        for entry in sorted(log_path.iterdir(), key=lambda x: x.name):
            if entry.is_file() and (entry.suffix == ".log" or entry.suffix == ".gz"):
                try:
                    stat = entry.stat()
                    files.append(LogFile(
                        name=entry.name,
                        path=str(entry),
                        size=stat.st_size,
                        mtime=stat.st_mtime,
                    ))
                except OSError:
                    files.append(LogFile(
                        name=entry.name,
                        path=str(entry),
                    ))
        # Sort by mtime descending
        files.sort(key=lambda f: f.mtime or 0, reverse=True)
    except (PermissionError, OSError):
        pass

    return LogFileListResponse(
        files=files,
        log_path=str(log_path),
        default_file="hermes.log",
    )


async def tail_log(file_path: str, num_lines: int = 100) -> LogContent:
    """
    读取日志文件尾部内容（反向读取优化）。

    Args:
        file_path: 日志文件绝对路径
        num_lines: 返回行数（默认 100，限制 ≤2000）

    Returns:
        LogContent
    """
    num_lines = min(max(num_lines, 1), 2000)
    file_path_obj = Path(file_path).resolve()

    if not file_path_obj.exists():
        raise FileNotFoundError(f"日志文件不存在: {file_path}")

    total_size = file_path_obj.stat().st_size
    truncated = False

    # 文件过大策略：只读取尾部 2MB
    if total_size > AppConfig.MAX_TAIL_SCAN:
        read_size = AppConfig.MAX_TAIL_SCAN
        truncated = True
    else:
        read_size = total_size

    async with aiofiles.open(file_path_obj, mode="r", encoding="utf-8", errors="replace") as f:
        if total_size > 0:
            # Seek to near-end
            seek_pos = max(0, total_size - read_size)
            await f.seek(seek_pos)
            # Read to end
            content = await f.read()
        else:
            content = ""

    # Split and take last N lines
    lines = content.split("\n")
    tail_lines = lines[-num_lines:] if len(lines) > num_lines else lines
    result_content = "\n".join(tail_lines)

    return LogContent(
        file=file_path_obj.name,
        content=result_content,
        lines=len(tail_lines),
        total_size=total_size,
        truncated=truncated or (len(lines) > num_lines),
    )


async def search_logs(
    keyword: str,
    file_path: str,
    max_lines: int = 50,
    case_sensitive: bool = False,
    context_lines: int = 0,
) -> LogSearchResult:
    """
    在日志文件中搜索关键字。

    Args:
        keyword: 搜索关键字
        file_path: 日志文件路径
        max_lines: 最大匹配行数（限制 ≤500）
        case_sensitive: 是否大小写敏感
        context_lines: 上下文行数

    Returns:
        LogSearchResult
    """
    max_lines = min(max(max_lines, 1), 500)
    context_lines = max(context_lines, 0)
    file_path_obj = Path(file_path).resolve()

    if not file_path_obj.exists():
        raise FileNotFoundError(f"日志文件不存在: {file_path}")

    start_time = time.time()
    matches = []
    total_matches = 0
    truncated = False

    async with aiofiles.open(file_path_obj, mode="r", encoding="utf-8", errors="replace") as f:
        all_lines = await f.readlines()

    line_count = len(all_lines)
    matched_indices = set()

    for i, line in enumerate(all_lines):
        line_stripped = line.rstrip("\n").rstrip("\r")
        if not keyword:
            continue
        if case_sensitive:
            match = keyword in line_stripped
        else:
            match = keyword.lower() in line_stripped.lower()
        if match:
            total_matches += 1
            if len(matches) >= max_lines:
                truncated = True
                break
            matched_indices.add(i)
            # Add context lines
            start = max(0, i - context_lines)
            end = min(line_count - 1, i + context_lines)
            for j in range(start, end + 1):
                if j not in matched_indices:
                    matched_indices.add(j)
                    context_flag = "  " if j != i else "> "
                    matches.append(f"{context_flag}{all_lines[j].rstrip()}")

    time_cost = (time.time() - start_time) * 1000

    return LogSearchResult(
        keyword=keyword,
        file=file_path_obj.name,
        total_matches=total_matches,
        matches=matches,
        truncated=truncated,
        time_cost_ms=round(time_cost, 2),
    )


async def list_archives(log_dir: str) -> List[LogArchive]:
    """
    列出历史归档文件。

    Args:
        log_dir: 日志目录路径

    Returns:
        LogArchive 列表
    """
    log_path = Path(log_dir).resolve()

    if not log_path.exists():
        return []

    archives = []
    try:
        for entry in sorted(log_path.iterdir(), key=lambda x: x.name):
            if not entry.is_file():
                continue
            # Match archive patterns: hermes.log.1, hermes.log.1.gz, etc.
            name = entry.name
            if ".log." in name or name.endswith(".gz"):
                try:
                    stat = entry.stat()
                    archives.append(LogArchive(
                        name=name,
                        path=str(entry),
                        size=stat.st_size,
                        mtime=stat.st_mtime,
                    ))
                except OSError:
                    archives.append(LogArchive(
                        name=name,
                        path=str(entry),
                    ))
        archives.sort(key=lambda a: a.mtime or 0, reverse=True)
    except (PermissionError, OSError):
        pass

    return archives
