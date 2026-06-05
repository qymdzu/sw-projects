"""
存储浏览路由 — 目录树、文件读取、文件写入
"""

from fastapi import APIRouter, Depends, HTTPException, Query

from dependencies import verify_token
from models.schemas import ApiResponse, FileWriteRequest
from services import file_service
from utils.path_utils import get_repo_path, is_within_any_safe_path

router = APIRouter(prefix="/api/storage", tags=["storage"])


@router.get("/tree", response_model=ApiResponse)
async def get_storage_tree(
    root: str = Query("all", description="仓库名称: workspace|skills-library|projects|hermes|all"),
    path: str = Query("/", description="子路径"),
    max_depth: int = Query(3, ge=1, le=10, description="最大深度"),
    token: str = Depends(verify_token),
):
    """获取目录树"""
    repo_paths = get_repo_path(root)

    if repo_paths is None:
        raise HTTPException(status_code=400, detail=f"无效的仓库名称: {root}")

    if isinstance(repo_paths, list):
        # "all" — 并行获取三个仓库的树
        import asyncio

        async def get_single_tree(base_path: str, repo_name: str) -> dict:
            nodes = await file_service.build_tree(
                base_path=base_path,
                sub_path=path,
                max_depth=max_depth,
            )
            # 修复 2026-06-05：顶层虚拟节点的 path 改成 repo 名（短 key），
            # 避免与内层 children 的绝对路径冲突导致 el-tree node-key 错乱
            # （"workspace" 嵌套 8 层就是 node-key 重复引发的渲染错位）
            return {
                "name": repo_name,  # 用于 el-tree 顶层 label
                "repo": repo_name,
                "path": f"/{repo_name}",  # 用相对路径，el-tree node-key 不冲突
                "type": "directory",
                "isLeaf": False,
                "children": [n.model_dump() for n in nodes],
            }

        tasks = []
        repo_names = ["workspace", "skills-library", "projects"]
        for rp, rn in zip(repo_paths, repo_names):
            if rp:
                tasks.append(get_single_tree(rp, rn))

        trees = await asyncio.gather(*tasks, return_exceptions=True)
        results = []
        for t in trees:
            if isinstance(t, dict):
                results.append(t)

        return ApiResponse(success=True, data=results)

    # Single repo
    nodes = await file_service.build_tree(
        base_path=repo_paths,
        sub_path=path,
        max_depth=max_depth,
    )
    return ApiResponse(
        success=True,
        data=[n.model_dump() for n in nodes],
    )


@router.get("/file", response_model=ApiResponse)
async def get_storage_file(
    path: str = Query(..., description="文件绝对路径"),
    token: str = Depends(verify_token),
):
    """读取文件内容"""
    if not path:
        raise HTTPException(status_code=400, detail="path 参数不能为空")

    try:
        content = await file_service.read_file(path)
        return ApiResponse(success=True, data=content.model_dump())
    except file_service.PathTraversalError as e:
        raise HTTPException(status_code=403, detail=str(e))
    except FileNotFoundError as e:
        raise HTTPException(status_code=404, detail=str(e))
    except file_service.FileTooLargeError as e:
        raise HTTPException(status_code=413, detail=str(e))


@router.put("/file", response_model=ApiResponse)
async def put_storage_file(
    request: FileWriteRequest,
    token: str = Depends(verify_token),
):
    """写入文件（自动备份 .bak）"""
    try:
        result = await file_service.write_with_backup(
            file_path=request.path,
            content=request.content,
            create_backup=request.create_backup,
        )
        return ApiResponse(success=True, data=result.model_dump(), message="文件保存成功")
    except file_service.PathTraversalError as e:
        raise HTTPException(status_code=403, detail=str(e))
    except file_service.WriteError as e:
        raise HTTPException(status_code=500, detail=str(e))
