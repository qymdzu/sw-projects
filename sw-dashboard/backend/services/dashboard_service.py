"""
仪表盘聚合服务 — Gateway 状态、会话统计、Pipeline 进度、Cron 状态
"""

import asyncio
import json
import os
from datetime import datetime, timezone, timedelta
from pathlib import Path
from typing import List, Optional

import yaml

from config import AppConfig
from models.schemas import (
    CronJobStatus,
    DashboardData,
    GatewayStatus,
    PipelineProgress,
    SessionStats,
    StageInfo,
)
from utils.cmd_utils import run_cmd, run_crontab, run_systemctl
from services import config_service


async def get_gateway_status() -> GatewayStatus:
    """获取 Gateway 服务状态"""
    try:
        is_active = await run_systemctl("is-active", "hermes-gateway")
        status = is_active.get("status", "unknown")

        if status == "active":
            show_info = await run_systemctl("show", "hermes-gateway")
            return GatewayStatus(
                status="active",
                pid=show_info.get("pid"),
                memory_mb=show_info.get("memory_mb"),
                cpu_percent=show_info.get("cpu_percent"),
            )
        return GatewayStatus(status=status)
    except Exception:
        return GatewayStatus(status="unknown")


async def get_session_stats() -> SessionStats:
    """获取会话统计信息"""
    session_file = Path(AppConfig.base_paths.get("hermes", "~/.hermes")) / ".." / ".." / "hermes-gateway" / "var" / "stats" / "sessions.json"
    session_file = session_file.resolve()

    if not session_file.exists():
        return SessionStats()

    try:
        content = session_file.read_text(encoding="utf-8")
        data = json.loads(content)
        return SessionStats(
            total_sessions=data.get("total_sessions", 0),
            total_messages=data.get("total_messages", 0),
            active_sessions=data.get("active_sessions", 0),
            avg_duration_sec=data.get("avg_duration_sec", 0.0),
        )
    except (json.JSONDecodeError, OSError, PermissionError):
        return SessionStats()


async def get_pipeline_status() -> Optional[PipelineProgress]:
    """获取 Pipeline 进度"""
    project_file = Path(AppConfig.base_paths.get("projects", "")) / ".project-meta.yaml"
    project_file = project_file.resolve()

    if not project_file.exists():
        return None

    try:
        content = project_file.read_text(encoding="utf-8")
        data = yaml.safe_load(content)
        if not data or not isinstance(data, dict):
            return None

        stages_data = data.get("stages", [])
        stages = []
        for s in stages_data:
            stages.append(StageInfo(
                id=s.get("id", 0),
                name=s.get("name", ""),
                status=s.get("status", "pending"),
                progress=s.get("progress", 0.0),
            ))

        total_stages = len(stages)
        current_stage = None
        for i, s in enumerate(stages):
            if s.status != "completed":
                current_stage = i
                break

        overall_progress = (
            sum(s.progress for s in stages) / total_stages
            if total_stages > 0
            else 0.0
        )

        return PipelineProgress(
            current_stage=current_stage,
            total_stages=total_stages,
            stages=stages,
            overall_progress=round(overall_progress, 1),
        )
    except (yaml.YAMLError, OSError, PermissionError, KeyError):
        return None


async def get_cron_status() -> List[CronJobStatus]:
    """获取 Cron 任务状态（仪表盘用）"""
    result = await run_crontab(action="list")
    if not result.get("success", False):
        return []

    jobs_data = result.get("jobs", [])
    cron_statuses = []

    for i, job in enumerate(jobs_data):
        name = job.get("comment") or job.get("command", "").split("/")[-1] or f"Cron Job {i + 1}"
        cron_statuses.append(CronJobStatus(
            name=name,
            schedule=job.get("schedule", ""),
            enabled=job.get("enabled", True),
        ))

    return cron_statuses


async def aggregate_dashboard() -> DashboardData:
    """聚合所有仪表盘数据"""
    # 并行获取所有数据
    gateway_task = get_gateway_status()
    sessions_task = get_session_stats()
    pipeline_task = get_pipeline_status()
    cron_task = get_cron_status()

    results = await asyncio.gather(
        gateway_task,
        sessions_task,
        pipeline_task,
        cron_task,
        return_exceptions=True,
    )

    # 提取结果，异常时使用默认值
    gateway = results[0] if not isinstance(results[0], Exception) else GatewayStatus()
    sessions = results[1] if not isinstance(results[1], Exception) else SessionStats()
    pipeline = results[2] if not isinstance(results[2], Exception) else None
    cron_jobs = results[3] if not isinstance(results[3], Exception) else []

    # 服务器当前时间
    now = datetime.now(tz=timezone(timedelta(hours=8)))
    server_time = now.strftime("%Y-%m-%dT%H:%M:%S+08:00")

    return DashboardData(
        gateway=gateway,
        sessions=sessions,
        pipeline=pipeline,
        cron_jobs=cron_jobs,
        server_time=server_time,
        version=AppConfig.VERSION,
    )
