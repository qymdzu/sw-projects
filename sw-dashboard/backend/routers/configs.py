"""
配置管理路由 — 配置列表、读取、保存、Cron 管理
"""

from fastapi import APIRouter, Depends, HTTPException, Query

from dependencies import verify_token
from models.schemas import ApiResponse, ConfigWriteRequest, CronToggleRequest
from services import config_service

router = APIRouter(prefix="/api/configs", tags=["configs"])

# ⚠ 路由顺序：FastAPI 按注册顺序匹配，/cron 系列必须注册在 /{name} 之前
# 否则 GET /api/configs/cron 会被 /{name} 拦截（name="cron"），触发
# "不支持的配置文件: cron" 错误。Bug 修复 2026-06-05。

@router.get("", response_model=ApiResponse)
async def get_configs_list(token: str = Depends(verify_token)):
    """获取配置文件列表"""
    configs = await config_service.list_configs()
    return ApiResponse(
        success=True,
        data={"configs": [c.model_dump() for c in configs]},
    )


# --- Cron 系列：必须在 /{name} 之前 ---
@router.get("/cron", response_model=ApiResponse)
async def get_cron_list(token: str = Depends(verify_token)):
    """获取 Cron 任务列表"""
    jobs = await config_service.list_crontab()
    return ApiResponse(
        success=True,
        data=[j.model_dump() for j in jobs],
    )


@router.post("/cron/toggle", response_model=ApiResponse)
async def post_cron_toggle(
    request: CronToggleRequest,
    token: str = Depends(verify_token),
):
    """启停 Cron 任务"""
    result = await config_service.toggle_cron(
        command=request.command,
        enable=request.enable,
    )
    if not result.success:
        raise HTTPException(status_code=404, detail=result.message)
    return ApiResponse(
        success=True,
        data=result.model_dump(),
        message=result.message,
    )


# --- 通用 /{name}：必须在最后 ---
@router.get("/{name}", response_model=ApiResponse)
async def get_config(
    name: str,
    masked: bool = Query(True, description="是否掩码敏感信息（仅 .env 有效）"),
    token: str = Depends(verify_token),
):
    """读取配置文件"""
    try:
        if name == ".env":
            config_data = await config_service.read_env(masked=masked)
        else:
            config_data = await config_service.read_config(name)
        return ApiResponse(success=True, data=config_data.model_dump())
    except config_service.ConfigError as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.put("/{name}", response_model=ApiResponse)
async def put_config(
    name: str,
    request: ConfigWriteRequest,
    token: str = Depends(verify_token),
):
    """保存配置文件"""
    try:
        result = await config_service.write_config(
            config_name=name,
            content=request.content,
            validate=request.do_validate,
            create_backup=request.create_backup,
        )
        return ApiResponse(
            success=True,
            data=result.model_dump(),
            message="配置保存成功",
        )
    except config_service.ConfigError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
