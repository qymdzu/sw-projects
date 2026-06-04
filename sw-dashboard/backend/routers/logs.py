"""
日志查看路由 — 日志文件列表、tail、搜索、归档
"""

from fastapi import APIRouter, Depends, HTTPException, Query

from config import AppConfig
from dependencies import verify_token
from models.schemas import ApiResponse
from services import log_service
from utils.path_utils import is_within_any_safe_path

router = APIRouter(prefix="/api/logs", tags=["logs"])


@router.get("/files", response_model=ApiResponse)
async def get_log_files(token: str = Depends(verify_token)):
    """获取日志文件列表"""
    log_dir = AppConfig.LOG_DIR
    result = await log_service.list_log_files(log_dir)
    return ApiResponse(success=True, data=result.model_dump())


@router.get("/tail", response_model=ApiResponse)
async def get_log_tail(
    file: str = Query("hermes.log", description="日志文件名"),
    lines: int = Query(100, ge=1, le=2000, description="返回行数"),
    token: str = Depends(verify_token),
):
    """获取日志尾部内容"""
    log_dir = AppConfig.LOG_DIR
    full_path = str(__import__("pathlib").Path(log_dir).resolve() / file)

    # 安全校验
    safe_base = is_within_any_safe_path(full_path)
    if safe_base is None:
        raise HTTPException(status_code=403, detail="路径不在安全白名单内")

    try:
        result = await log_service.tail_log(full_path, num_lines=lines)
        return ApiResponse(success=True, data=result.model_dump())
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.get("/search", response_model=ApiResponse)
async def get_log_search(
    keyword: str = Query(..., min_length=1, description="搜索关键字"),
    file: str = Query("hermes.log", description="日志文件名"),
    max_lines: int = Query(50, ge=1, le=500, description="最大匹配行数"),
    case_sensitive: bool = Query(False, description="是否大小写敏感"),
    context_lines: int = Query(0, ge=0, le=10, description="上下文行数"),
    token: str = Depends(verify_token),
):
    """搜索日志"""
    if not keyword:
        raise HTTPException(status_code=400, detail="keyword 不能为空")

    log_dir = AppConfig.LOG_DIR
    full_path = str(__import__("pathlib").Path(log_dir).resolve() / file)

    safe_base = is_within_any_safe_path(full_path)
    if safe_base is None:
        raise HTTPException(status_code=403, detail="路径不在安全白名单内")

    try:
        result = await log_service.search_logs(
            keyword=keyword,
            file_path=full_path,
            max_lines=max_lines,
            case_sensitive=case_sensitive,
            context_lines=context_lines,
        )
        return ApiResponse(success=True, data=result.model_dump())
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.get("/archive", response_model=ApiResponse)
async def get_log_archive(token: str = Depends(verify_token)):
    """获取历史归档列表"""
    log_dir = AppConfig.LOG_DIR
    archives = await log_service.list_archives(log_dir)
    return ApiResponse(
        success=True,
        data=[a.model_dump() for a in archives],
    )
