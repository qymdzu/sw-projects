"""
记忆浏览路由 — 翠花 MEMORY.md 读取、解析、搜索
"""

from fastapi import APIRouter, Depends, HTTPException, Query

from dependencies import verify_token
from models.schemas import ApiResponse
from services import memory_service

router = APIRouter(prefix="/api/memory", tags=["memory"])


@router.get("/file", response_model=ApiResponse)
async def get_memory_file(
    token: str = Depends(verify_token),
):
    """读取翠花 MEMORY.md 全文"""
    try:
        content = memory_service.read_memory()
        return ApiResponse(success=True, data=content.model_dump())
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except PermissionError as e:
        raise HTTPException(status_code=403, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"读取失败: {e}")


@router.get("/sections", response_model=ApiResponse)
async def get_memory_sections(
    token: str = Depends(verify_token),
):
    """解析 MEMORY.md 为段落结构（供面板渲染）"""
    try:
        sections = memory_service.parse_memory_sections()
        return ApiResponse(success=True, data=sections)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"解析失败: {e}")


@router.get("/search", response_model=ApiResponse)
async def search_memory(
    keyword: str = Query(..., min_length=1, description="搜索关键词"),
    token: str = Depends(verify_token),
):
    """在 MEMORY.md 中搜索关键词"""
    if not keyword:
        raise HTTPException(status_code=400, detail="keyword 不能为空")
    try:
        results = memory_service.search_memory(keyword)
        return ApiResponse(success=True, data={
            "keyword": keyword,
            "count": len(results),
            "results": results,
        })
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"搜索失败: {e}")
