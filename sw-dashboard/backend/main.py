"""
应用入口 — sw-dashboard 翠花集群管理面板

启动命令:
    python main.py
    # 或
    uvicorn main:app --host 0.0.0.0 --port 8900
"""

import sys
from pathlib import Path

# 确保项目根目录在 sys.path 中
backend_dir = Path(__file__).parent.resolve()
if str(backend_dir) not in sys.path:
    sys.path.insert(0, str(backend_dir))

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles

from config import AppConfig

# ---------- 路由 ----------
from routers import dashboard as dashboard_router
from routers import storage as storage_router
from routers import configs as configs_router
from routers import logs as logs_router
from routers import skills as skills_router
from routers import memory as memory_router


def create_app() -> FastAPI:
    """创建并配置 FastAPI 应用实例"""
    app = FastAPI(
        title="sw-dashboard — 翠花集群管理面板",
        version=AppConfig.VERSION,
        description="翠花集群管理面板后端 API",
    )

    # CORS 中间件（Tailscale 内网，允许所有来源）
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # 注册 API 路由
    app.include_router(dashboard_router.router)
    app.include_router(storage_router.router)
    app.include_router(configs_router.router)
    app.include_router(logs_router.router)
    app.include_router(skills_router.router)
    app.include_router(memory_router.router)

    # 挂载前端静态文件（如果存在）
    frontend_dist = backend_dir.parent / "frontend" / "dist"
    if frontend_dist.exists():
        app.mount("/", StaticFiles(directory=str(frontend_dist), html=True), name="frontend")

    @app.on_event("startup")
    async def startup():
        """应用启动时创建必要目录"""
        # 确保 hermes home 目录存在（供配置编辑用）
        hermes_home = Path(AppConfig.base_paths.get("hermes", ""))
        hermes_home.mkdir(parents=True, exist_ok=True)

    return app


# 创建应用实例
app = create_app()


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=AppConfig.HOST,
        port=AppConfig.PORT,
        workers=1,
        reload=False,
    )
