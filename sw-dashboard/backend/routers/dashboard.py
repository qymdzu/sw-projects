"""
仪表盘路由 — 健康检查、Token 认证、仪表盘聚合
"""

import platform
import sys
import time
from datetime import datetime, timezone, timedelta

from fastapi import APIRouter, Depends, HTTPException, Request

from config import AppConfig
from dependencies import verify_token
from models.schemas import ApiResponse, AuthRequest, AuthResponse, HealthResponse
from services import dashboard_service

router = APIRouter(tags=["dashboard"])


@router.get("/api/health", response_model=HealthResponse)
async def health_check():
    """健康检查（免认证）"""
    # Python 版本
    python_ver = f"{sys.version_info.major}.{sys.version_info.minor}.{sys.version_info.micro}"

    # 依赖版本
    dependencies = {
        "fastapi": "unknown",
        "uvicorn": "unknown",
        "pyyaml": "unknown",
        "aiofiles": "unknown",
        "httpx": "unknown",
    }
    try:
        import fastapi
        dependencies["fastapi"] = fastapi.__version__
    except (ImportError, AttributeError):
        pass
    try:
        import uvicorn
        dependencies["uvicorn"] = uvicorn.__version__
    except (ImportError, AttributeError):
        pass
    try:
        import yaml
        dependencies["pyyaml"] = yaml.__version__
    except (ImportError, AttributeError):
        pass
    try:
        import aiofiles
        dependencies["aiofiles"] = aiofiles.__version__
    except (ImportError, AttributeError):
        pass
    try:
        import httpx
        dependencies["httpx"] = httpx.__version__
    except (ImportError, AttributeError):
        pass

    return HealthResponse(
        status="ok",
        version=AppConfig.VERSION,
        uptime=time.time(),
        python_version=python_ver,
        dependencies=dependencies,
    )


@router.post("/api/auth", response_model=ApiResponse[AuthResponse])
async def auth_login(request: AuthRequest):
    """Token 认证（免认证）"""
    if request.token == AppConfig.TOKEN:
        return ApiResponse(
            success=True,
            data=AuthResponse(success=True, token=request.token),
            message="认证成功",
        )
    raise HTTPException(status_code=401, detail="Token 无效")


@router.get("/api/dashboard", response_model=ApiResponse)
async def get_dashboard(token: str = Depends(verify_token)):
    """仪表盘聚合数据"""
    data = await dashboard_service.aggregate_dashboard()
    return ApiResponse(success=True, data=data.model_dump())
