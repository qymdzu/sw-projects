"""
数据模型 Schema — 翠花集群管理面板 (sw-dashboard)

所有 Pydantic 请求/响应模型
"""

from datetime import datetime
from typing import Any, Dict, Generic, List, Optional, TypeVar

from pydantic import BaseModel, Field

T = TypeVar("T")


# ========== 通用 ==========

class ApiResponse(BaseModel, Generic[T]):
    """通用 API 响应包装"""
    success: bool = True
    data: Optional[T] = None
    message: str = ""


# ========== 健康检查 & 认证 ==========

class HealthResponse(BaseModel):
    """健康检查响应"""
    status: str = "ok"
    version: str = "0.1.0"
    uptime: Optional[float] = None
    python_version: str = ""
    dependencies: Dict[str, str] = {}


class AuthRequest(BaseModel):
    """认证请求"""
    token: str = Field(..., min_length=1)


class AuthResponse(BaseModel):
    """认证响应"""
    success: bool = True
    token: str = ""


# ========== 文件系统 ==========

class FileNode(BaseModel):
    """文件树节点"""
    name: str
    path: str
    type: str  # "file" | "directory"
    size: Optional[int] = None
    mtime: Optional[float] = None
    isLeaf: bool = True
    children: Optional[List["FileNode"]] = None
    has_skill_md: bool = False


class FileContent(BaseModel):
    """文件内容"""
    path: str
    name: str
    content: str
    size: int
    mtime: float
    language: str = "text"
    editable: bool = False


class FileWriteRequest(BaseModel):
    """文件写入请求"""
    path: str = Field(..., min_length=1)
    content: str = ""
    create_backup: bool = True


class FileWriteResponse(BaseModel):
    """文件写入响应"""
    path: str
    size: int = 0
    backup_path: Optional[str] = None
    timestamp: float = 0.0


# ========== 配置管理 ==========

class ConfigFile(BaseModel):
    """配置文件内容"""
    name: str
    path: str
    content: str = ""
    language: str = "text"
    editable: bool = True
    last_modified: Optional[float] = None


class ConfigListItem(BaseModel):
    """配置列表项"""
    name: str
    language: str
    editable: bool = True
    last_modified: Optional[float] = None


class ConfigListResponse(BaseModel):
    """配置列表响应"""
    configs: List[ConfigListItem] = []


class ConfigWriteRequest(BaseModel):
    """配置写入请求"""
    content: str = ""
    create_backup: bool = True
    do_validate: bool = True


class CronJob(BaseModel):
    """Cron 任务"""
    raw: str = ""
    schedule: str = ""
    command: str = ""
    comment: str = ""
    enabled: bool = True
    last_run: Optional[str] = None
    last_output: Optional[str] = None


class CronToggleRequest(BaseModel):
    """Cron 启停请求"""
    command: str = Field(..., min_length=1)
    enable: bool = True


class CronToggleResponse(BaseModel):
    """Cron 启停响应"""
    success: bool = True
    message: str = ""


# ========== 日志 ==========

class LogFile(BaseModel):
    """日志文件信息"""
    name: str
    path: str
    size: int = 0
    mtime: Optional[float] = None


class LogFileListResponse(BaseModel):
    """日志文件列表响应"""
    files: List[LogFile] = []
    log_path: str = ""
    default_file: str = "hermes.log"


class LogContent(BaseModel):
    """日志尾部内容"""
    file: str = ""
    content: str = ""
    lines: int = 0
    total_size: int = 0
    truncated: bool = False


class LogSearchResult(BaseModel):
    """日志搜索结果"""
    keyword: str = ""
    file: str = ""
    total_matches: int = 0
    matches: List[str] = []
    truncated: bool = False
    time_cost_ms: float = 0.0


class LogArchive(BaseModel):
    """日志归档文件"""
    name: str
    path: str
    size: int = 0
    mtime: Optional[float] = None


# ========== 仪表盘 ==========

class GatewayStatus(BaseModel):
    """Gateway 服务状态"""
    status: str = "unknown"  # active | inactive | failed | unknown
    uptime: Optional[str] = None
    pid: Optional[int] = None
    memory_mb: Optional[float] = None
    cpu_percent: Optional[float] = None


class SessionStats(BaseModel):
    """会话统计"""
    total_sessions: int = 0
    total_messages: int = 0
    active_sessions: int = 0
    avg_duration_sec: float = 0.0


class StageInfo(BaseModel):
    """Pipeline 阶段信息"""
    id: int = 0
    name: str = ""
    status: str = "pending"  # pending | running | completed | failed
    progress: float = 0.0


class PipelineProgress(BaseModel):
    """Pipeline 进度"""
    current_stage: Optional[int] = None
    total_stages: int = 0
    stages: List[StageInfo] = []
    overall_progress: float = 0.0


class CronJobStatus(BaseModel):
    """Cron 任务状态（仪表盘用）"""
    name: str = ""
    schedule: str = ""
    last_run: Optional[str] = None
    last_status: Optional[str] = None
    enabled: bool = True


class DashboardData(BaseModel):
    """仪表盘聚合数据"""
    gateway: GatewayStatus = GatewayStatus()
    sessions: SessionStats = SessionStats()
    pipeline: Optional[PipelineProgress] = None
    cron_jobs: List[CronJobStatus] = []
    server_time: str = ""
    version: str = "0.1.0"


# ========== 技能 ==========

class SkillNode(BaseModel):
    """技能树节点"""
    name: str
    path: str
    type: str  # "file" | "directory"
    size: Optional[int] = None
    mtime: Optional[float] = None
    isLeaf: bool = True
    children: Optional[List["SkillNode"]] = None
    has_skill_md: bool = False


class SkillFileContent(BaseModel):
    """技能文件内容（只读）"""
    path: str
    name: str
    content: str
    size: int
    mtime: float
    language: str = "text"
