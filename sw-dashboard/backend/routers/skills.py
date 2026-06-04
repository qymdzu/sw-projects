"""
技能浏览路由 — 技能目录树、文件预览（只读）
"""

from fastapi import APIRouter, Depends, HTTPException, Query

from config import AppConfig
from dependencies import verify_token
from models.schemas import ApiResponse
from services import file_service
from utils.path_utils import is_within_base

router = APIRouter(prefix="/api/skills", tags=["skills"])


@router.get("/tree", response_model=ApiResponse)
async def get_skills_tree(
    path: str = Query("/", description="子路径（懒加载）"),
    token: str = Depends(verify_token),
):
    """获取技能目录树"""
    skills_path = AppConfig.skills_path

    nodes = await file_service.build_tree(
        base_path=skills_path,
        sub_path=path,
        max_depth=5,
        check_skill_md=True,
    )
    return ApiResponse(
        success=True,
        data=[n.model_dump() for n in nodes],
    )


@router.get("/file", response_model=ApiResponse)
async def get_skills_file(
    path: str = Query(..., description="技能文件绝对路径"),
    token: str = Depends(verify_token),
):
    """预览技能文件（只读）"""
    if not path:
        raise HTTPException(status_code=400, detail="path 参数不能为空")

    # 校验路径在 skills_path 白名单内
    skills_path = AppConfig.skills_path
    if not is_within_base(path, skills_path):
        raise HTTPException(status_code=403, detail="路径不在技能目录内")

    try:
        content = await file_service.read_file(path)
        return ApiResponse(success=True, data=content.model_dump())
    except file_service.PathTraversalError as e:
        raise HTTPException(status_code=403, detail=str(e))
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except file_service.FileTooLargeError as e:
        raise HTTPException(status_code=413, detail=str(e))
