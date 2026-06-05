"""
记忆浏览服务 — 翠花 MEMORY.md 读取、解析、分段
"""

import re
from pathlib import Path
from typing import List, Optional

from config import AppConfig
from models.schemas import FileContent
from utils.path_utils import is_within_base


class MemoryParseError(Exception):
    pass


def read_memory() -> FileContent:
    """
    读取翠花 MEMORY.md 文件内容。

    Returns:
        FileContent 对象

    Raises:
        FileNotFoundError: 记忆文件不存在
        PermissionError: 无权限读取
        MemoryParseError: 解析失败
    """
    memory_path = AppConfig.memory_path

    if not Path(memory_path).exists():
        raise FileNotFoundError(f"记忆文件不存在: {memory_path}")

    try:
        content = Path(memory_path).read_text(encoding="utf-8")
    except PermissionError:
        raise PermissionError(f"无权限读取: {memory_path}")

    stat = Path(memory_path).stat()
    return FileContent(
        path=memory_path,
        name="MEMORY.md",
        content=content,
        size=stat.st_size,
        mtime=stat.st_mtime,
        language="markdown",
        editable=True,  # 记忆文件可编辑
    )


def parse_memory_sections() -> List[dict]:
    """
    把 MEMORY.md 解析成段落结构，供面板渲染。

    MEMORY.md 格式示例：
        ## 用户偏好（2026-06-05）
        内容...

        ## 技术配置
        内容...

    Returns:
        段落列表，每项含 {title, date, content, offset, length}
    """
    memory_path = AppConfig.memory_path
    try:
        content = Path(memory_path).read_text(encoding="utf-8")
    except FileNotFoundError:
        return []

    sections = []
    # 按 ## 标题分段
    pattern = re.compile(r'^(#{1,3})\s+(.+)$', re.MULTILINE)
    matches = list(pattern.finditer(content))

    for i, match in enumerate(matches):
        level = len(match.group(1))  # 1=#, 2=##, 3=###
        title = match.group(2).strip()
        start = match.end()
        end = matches[i + 1].start() if i + 1 < len(matches) else len(content)

        section_content = content[start:end].strip()

        # 提日期（如果有）
        date_match = re.search(r'\((\d{4}-\d{2}-\d{2})\)', title)
        date = date_match.group(1) if date_match else ""

        sections.append({
            "level": level,
            "title": title,
            "date": date,
            "content": section_content,
            "offset": start,
            "length": end - start,
        })

    return sections


def search_memory(keyword: str) -> List[dict]:
    """
    在 MEMORY.md 中搜索关键词。

    Args:
        keyword: 搜索关键词

    Returns:
        匹配段落列表，每项含 {title, snippet, line_num}
    """
    memory_path = AppConfig.memory_path
    try:
        lines = Path(memory_path).read_text(encoding="utf-8").splitlines()
    except FileNotFoundError:
        return []

    keyword_lower = keyword.lower()
    results = []
    current_title = ""
    current_level = 0

    for i, line in enumerate(lines):
        # 跟踪当前标题
        m = re.match(r'^(#{1,3})\s+(.+)$', line.strip())
        if m:
            current_level = len(m.group(1))
            current_title = m.group(2).strip()
        elif keyword_lower in line.lower():
            # 去除标题行本身
            if not re.match(r'^(#{1,3})\s+', line.strip()):
                snippet = line.strip()
                if len(snippet) > 200:
                    idx = snippet.lower().find(keyword_lower)
                    start = max(0, idx - 80)
                    end = min(len(snippet), idx + len(keyword) + 80)
                    snippet = ("..." if start > 0 else "") + snippet[start:end] + ("..." if end < len(snippet) else "")
                results.append({
                    "title": current_title,
                    "snippet": snippet,
                    "line": i + 1,
                })

    return results
