# =============================================================================
# Dockerfile — sw-dashboard 翠花集群管理面板
# 多阶段构建：前端 (Node 18) → 后端 (Python 3.12)
# 备选构建方式，生产环境以 systemd 直接部署为主
# =============================================================================

# ---- Stage 1: 前端构建 ----
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# ---- Stage 2: 后端运行 ----
FROM python:3.12-slim

WORKDIR /app

# 系统依赖（如有需要可在此添加）
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 复制后端代码
COPY backend/ ./backend/

# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# 安装 Python 依赖
RUN pip install --no-cache-dir -r backend/requirements.txt

# 工作目录设为 backend，以便相对路径 frontend/dist 正确
WORKDIR /app/backend

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD curl -f http://localhost:8900/api/health || exit 1

# 暴露端口
EXPOSE 8900

# 启动命令
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8900", "--workers", "1"]
