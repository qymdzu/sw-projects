"""
全局配置 — 翠花集群管理面板 (sw-dashboard)

配置优先级: 环境变量 > .env 文件 > 默认值
"""

import os
import sys
import warnings
from pathlib import Path
from typing import Dict, List, Optional


class AppConfig:
    """应用全局配置单例"""

    # 服务器
    HOST: str = "0.0.0.0"
    PORT: int = 8900

    # Token 认证
    TOKEN: str = ""

    # 版本
    VERSION: str = "0.1.0"

    # 三库路径（优先从环境变量读取）
    _AGENTS_WORKSPACE_PATH: str = "~/gitee-software/sw-agents-workspace"
    _SKILLS_LIBRARY_PATH: str = "~/gitee-software/sw-skills-library"
    _PROJECTS_PATH: str = "~/gitee-software/sw-projects"
    _HERMES_HOME: str = "/home/ubuntu/.hermes/profiles/software-dev"
    _SKILLS_PATH: str = "~/.hermes/skills"
    _MEMORY_PATH: str = "~/.hermes/profiles/software-dev/MEMORY.md"

    # 已展开的路径缓存
    base_paths: Dict[str, str] = {}
    skills_path: str = ""
    memory_path: str = ""

    # dashboard 自身运行日志
    # 修复 2026-06-05：之前路径 /home/ubuntu/.hermes/profiles/software-dev/var/log
    # 不存在且 dashboard 没自动建，导致 Logs 页永远显示"日志目录不存在"。
    # 改为 hermes home 下的 log/ 子目录（user-owned 目录 dashboard 可写），
    # 启动时自动建。同时把 INFO 日志镜像到文件供 Logs 页 tail。
    LOG_DIR: str = "/home/ubuntu/.hermes/profiles/software-dev/log"

    # 文件限制
    MAX_FILE_SIZE: int = 10 * 1024 * 1024  # 10MB 读取上限
    MAX_EDIT_SIZE: int = 1 * 1024 * 1024   # 1MB 编辑上限
    MAX_TAIL_SCAN: int = 2 * 1024 * 1024   # 2MB tail 扫描上限

    # 可编辑后缀白名单
    EDITABLE_EXTENSIONS: List[str] = [
        ".py", ".js", ".ts", ".yaml", ".yml",
        ".md", ".json", ".env", ".toml", ".sh",
        ".txt", ".cfg", ".conf", ".ini", ".csv",
        ".xml", ".html", ".css", ".scss", ".vue",
        ".svelte", ".tsx", ".jsx",
    ]

    # 日志目录

    _initialized: bool = False

    @classmethod
    def init_config(cls) -> None:
        """初始化配置 — 读取环境变量并展开路径"""
        if cls._initialized:
            return

        # --- 读取 Token ---
        token = os.environ.get("HERMES_DASHBOARD_TOKEN", "")
        if not token:
            # 尝试从 ~/.hermes/.env 读取
            env_path = Path("~/.hermes/.env").expanduser()
            if env_path.exists():
                try:
                    for line in env_path.read_text().splitlines():
                        line = line.strip()
                        if line.startswith("HERMES_DASHBOARD_TOKEN="):
                            token = line.split("=", 1)[1].strip().strip("\"'")
                            break
                except Exception:
                    pass
        if not token:
            # 尝试从 ~/.hermes/profiles/software-dev/.env 读取
            env_path = Path("~/.hermes/profiles/software-dev/.env").expanduser()
            if env_path.exists():
                try:
                    for line in env_path.read_text().splitlines():
                        line = line.strip()
                        if line.startswith("HERMES_DASHBOARD_TOKEN="):
                            token = line.split("=", 1)[1].strip().strip("\"'")
                            break
                except Exception:
                    pass
        if not token:
            warnings.warn(
                "⚠ HERMES_DASHBOARD_TOKEN 未设置，使用默认开发 Token！"
            )
            token = "dev-token-123456"
        cls.TOKEN = token

        # --- 读取环境变量路径覆盖 ---
        _real_user = os.environ.get("USER", "ubuntu")
        _base_default = f"/home/{_real_user}/gitee-software"

        cls._AGENTS_WORKSPACE_PATH = os.environ.get(
            "AGENTS_WORKSPACE_PATH", f"{_base_default}/sw-agents-workspace"
        )
        cls._SKILLS_LIBRARY_PATH = os.environ.get(
            "SKILLS_LIBRARY_PATH", f"{_base_default}/sw-skills-library"
        )
        cls._PROJECTS_PATH = os.environ.get(
            "PROJECTS_PATH", f"{_base_default}/sw-projects"
        )
        cls._HERMES_HOME = os.environ.get(
            "HERMES_HOME", f"/home/{_real_user}/.hermes/profiles/software-dev"
        )
        cls._SKILLS_PATH = os.environ.get(
            "SKILLS_PATH", f"/home/{_real_user}/.hermes/skills"
        )
        cls._MEMORY_PATH = os.environ.get(
            "MEMORY_PATH", f"/home/{_real_user}/.hermes/profiles/software-dev/MEMORY.md"
        )

        # --- 路径展开 ---
        # 注意：hermes_home 是绝对路径，不用 expanduser
        hermes_home = str(Path(cls._HERMES_HOME).resolve())

        cls.base_paths = {
            "workspace": str(Path(cls._AGENTS_WORKSPACE_PATH).expanduser().resolve()),
            "skills-library": str(Path(cls._SKILLS_LIBRARY_PATH).expanduser().resolve()),
            "projects": str(Path(cls._PROJECTS_PATH).expanduser().resolve()),
            "hermes": hermes_home,
        }

        # 技能路径：默认从 HERMES_HOME 推断（profile 内 skills 目录）
        # 修复 2026-06-05：原来硬编码 ~/.hermes/skills 是 default profile 的，
        # 当前 software-dev profile 跑，导致 skills 浏览整棵树为空。
        # 优先 SKILLS_PATH 环境变量 → 否则 {HERMES_HOME}/skills。
        _skills_default = str(Path(hermes_home) / "skills")
        skills_env = os.environ.get("SKILLS_PATH", _skills_default)
        cls.skills_path = str(Path(skills_env).expanduser().resolve())

        # 记忆路径：翠花的实际存储是 {HERMES_HOME}/memories/MEMORY.md
        # 修复 2026-06-05：之前默认 {HERMES_HOME}/MEMORY.md，但实际文件在 memories/ 子目录。
        _memory_default = str(Path(hermes_home) / "memories" / "MEMORY.md")
        memory_env = os.environ.get("MEMORY_PATH", _memory_default)
        cls.memory_path = str(Path(memory_env).expanduser().resolve())

        cls._initialized = True

    @classmethod
    def get_config_file_path(cls, config_name: str) -> str:
        """获取配置文件的完整路径"""
        hermes_home = cls.base_paths.get("hermes", "")
        config_map = {
            "config.yaml": os.path.join(hermes_home, "config.yaml"),
            ".env": os.path.join(hermes_home, ".env"),
            "CLAUDE.md": os.path.join(hermes_home, "CLAUDE.md"),
            "SOUL.md": os.path.join(hermes_home, "SOUL.md"),
            "MEMORY.md": AppConfig.memory_path,  # 2026-06-05 新增：实际在 memories/ 子目录
        }
        return config_map.get(config_name, "")

    @classmethod
    def get_safe_paths(cls) -> List[str]:
        """返回所有安全路径白名单"""
        return list(cls.base_paths.values()) + [
            cls.skills_path,
            cls.memory_path,
            cls.LOG_DIR,
        ]


# 初始化（模块加载时自动执行）
AppConfig.init_config()
