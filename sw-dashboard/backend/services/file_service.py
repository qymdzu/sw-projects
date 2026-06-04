"""
文件系统服务 — 目录树构建、文件读写、备份、语言检测
"""

import os
import shutil
import time
from pathlib import Path
from typing import List, Optional

import aiofiles

from config import AppConfig
from models.schemas import FileContent, FileNode, FileWriteResponse
from utils.path_utils import is_within_any_safe_path, safe_resolve


class FileTooLargeError(Exception):
    """文件过大异常"""
    pass


class PathTraversalError(Exception):
    """路径穿越异常"""
    pass


class WriteError(Exception):
    """写入失败异常"""
    pass


async def build_tree(
    base_path: str,
    sub_path: str = "/",
    max_depth: int = 5,
    current_depth: int = 0,
    check_skill_md: bool = False,
) -> List[FileNode]:
    """
    递归构建目录树。

    Args:
        base_path: 仓库基路径
        sub_path: 子路径
        max_depth: 最大递归深度
        current_depth: 当前递归深度
        check_skill_md: 是否检测 SKILL.md

    Returns:
        FileNode 列表
    """
    if current_depth > max_depth:
        return []

    # 组合完整路径
    if sub_path.startswith("/"):
        sub_path = sub_path[1:]
    full_path = os.path.join(base_path, sub_path)
    resolved = safe_resolve(base_path, full_path)

    if resolved is None:
        return []

    if not resolved.exists():
        return []

    try:
        entries = sorted(
            resolved.iterdir(),
            key=lambda x: (not x.is_dir(), x.name.lower()),
        )
    except (PermissionError, OSError):
        return []

    nodes = []
    for entry in entries:
        # 跳过隐藏文件
        if entry.name.startswith("."):
            continue

        rel_path = os.path.join(sub_path, entry.name) if sub_path != "/" else entry.name

        if entry.is_dir():
            children = await build_tree(
                base_path=base_path,
                sub_path=rel_path,
                max_depth=max_depth,
                current_depth=current_depth + 1,
                check_skill_md=check_skill_md,
            )
            has_md = False
            if check_skill_md:
                has_md = (entry / "SKILL.md").exists()
            nodes.append(FileNode(
                name=entry.name,
                path=str(entry),
                type="directory",
                isLeaf=False,
                children=children,
                has_skill_md=has_md,
            ))
        else:
            try:
                stat = entry.stat()
                size = stat.st_size
                mtime = stat.st_mtime
            except OSError:
                size = 0
                mtime = 0.0
            nodes.append(FileNode(
                name=entry.name,
                path=str(entry),
                type="file",
                size=size,
                mtime=mtime,
                isLeaf=True,
            ))

    return nodes


async def read_file(file_path: str) -> FileContent:
    """
    读取文件内容。

    Args:
        file_path: 文件绝对路径

    Returns:
        FileContent 对象

    Raises:
        PathTraversalError: 路径穿越
        FileNotFoundError: 文件不存在
        FileTooLargeError: 文件过大
    """
    matched_base = is_within_any_safe_path(file_path)
    if matched_base is None:
        raise PathTraversalError(f"路径不在安全白名单内: {file_path}")

    file_path_obj = Path(file_path).resolve()
    if not file_path_obj.exists():
        raise FileNotFoundError(f"文件不存在: {file_path}")

    stat = file_path_obj.stat()
    size = stat.st_size

    if size > AppConfig.MAX_FILE_SIZE:
        raise FileTooLargeError(
            f"文件过大 ({size} bytes)，超过上限 {AppConfig.MAX_FILE_SIZE} bytes"
        )

    try:
        async with aiofiles.open(file_path_obj, mode="r", encoding="utf-8", errors="replace") as f:
            content = await f.read()
    except PermissionError:
        raise PathTraversalError(f"无权限读取文件: {file_path}")
    except UnicodeDecodeError:
        content = ""

    language = detect_language(str(file_path_obj))
    editable = is_editable(str(file_path_obj))

    return FileContent(
        path=str(file_path_obj),
        name=file_path_obj.name,
        content=content,
        size=size,
        mtime=stat.st_mtime,
        language=language,
        editable=editable,
    )


async def write_with_backup(
    file_path: str,
    content: str,
    create_backup: bool = True,
) -> FileWriteResponse:
    """
    写入文件，可选择先创建 .bak 备份。

    Args:
        file_path: 文件绝对路径
        content: 新内容
        create_backup: 是否创建备份

    Returns:
        FileWriteResponse

    Raises:
        PathTraversalError: 路径穿越
        WriteError: 写入失败
    """
    matched_base = is_within_any_safe_path(file_path)
    if matched_base is None:
        raise PathTraversalError(f"路径不在安全白名单内: {file_path}")

    file_path_obj = Path(file_path).resolve()
    parent = file_path_obj.parent

    # 确保父目录存在
    parent.mkdir(parents=True, exist_ok=True)

    # 检查是否可编辑
    if not is_editable(str(file_path_obj)):
        raise PathTraversalError(f"不支持编辑该文件类型: {file_path}")

    # 检查编辑大小
    content_bytes = content.encode("utf-8")
    if len(content_bytes) > AppConfig.MAX_EDIT_SIZE:
        raise WriteError(
            f"内容过大 ({len(content_bytes)} bytes)，超过编辑上限 {AppConfig.MAX_EDIT_SIZE} bytes"
        )

    backup_path = None
    if create_backup and file_path_obj.exists():
        backup_path = str(file_path_obj) + ".bak"
        try:
            shutil.copy2(str(file_path_obj), backup_path)
        except OSError as e:
            raise WriteError(f"备份失败: {e}")

    try:
        async with aiofiles.open(file_path_obj, mode="w", encoding="utf-8") as f:
            await f.write(content)
    except OSError as e:
        raise WriteError(f"写入失败: {e}")

    try:
        new_stat = file_path_obj.stat()
        new_size = new_stat.st_size
    except OSError:
        new_size = len(content_bytes)

    return FileWriteResponse(
        path=str(file_path_obj),
        size=new_size,
        backup_path=backup_path,
        timestamp=time.time(),
    )


def detect_language(file_path: str) -> str:
    """根据文件扩展名检测语言类型"""
    ext = Path(file_path).suffix.lower()
    ext_map = {
        ".py": "python",
        ".js": "javascript",
        ".jsx": "javascript",
        ".ts": "typescript",
        ".tsx": "typescript",
        ".yaml": "yaml",
        ".yml": "yaml",
        ".md": "markdown",
        ".json": "json",
        ".env": "env",
        ".toml": "toml",
        ".sh": "shell",
        ".log": "text",
        ".txt": "text",
        ".cfg": "ini",
        ".conf": "ini",
        ".ini": "ini",
        ".csv": "csv",
        ".xml": "xml",
        ".html": "html",
        ".css": "css",
        ".scss": "scss",
        ".vue": "vue",
        ".svelte": "svelte",
        ".pyc": "text",
    }
    return ext_map.get(ext, "text")


def is_editable(file_path: str) -> bool:
    """判断文件是否可编辑（后缀在白名单内）"""
    ext = Path(file_path).suffix.lower()
    return ext in AppConfig.EDITABLE_EXTENSIONS
