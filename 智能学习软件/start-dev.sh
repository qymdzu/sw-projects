#!/bin/bash
# === 智能学习软件 · 调试环境一键启动脚本 ===
#
# 用法：./start-dev.sh
# 功能：创建 .env → 检查环境 → 前端构建检查 → 启动 Docker Compose → 健康检查
#
# 环境变量说明（可通过 .env 或 export 覆盖）：
#   API_PORT     后端 API 映射端口   默认 8081
#   WEB_PORT     前端 Web 映射端口    默认 8080
#   DB_PORT      数据库映射端口       默认 5432
#   DEBUG_PORT   Delve 调试端口       默认 2345
#   DB_USER      数据库用户           默认 smart_learning
#   DB_PASSWORD  数据库密码           默认 smart_learning
#   DB_NAME      数据库名称           默认 smart_learning
#   JWT_SECRET   JWT 签名密钥         默认 dev-secret-please-replace-in-production
#   LOG_LEVEL    日志级别             默认 debug

set -euo pipefail

cd "$(dirname "$0")"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR" && pwd)"

echo "╔════════════════════════════════════════════╗"
echo "║   智能学习软件 · 调试环境启动             ║"
echo "╚════════════════════════════════════════════╝"

# ─── 1. 环境检查 ────────────────────────────────────
echo ""
echo "📋 步骤 1/5：环境检查"

if ! command -v docker &>/dev/null; then
  echo "❌ 错误：未找到 Docker 命令"
  echo "   请先安装 Docker：https://docs.docker.com/engine/install/"
  exit 1
fi
echo "  ✅ Docker: $(docker --version)"

if ! docker compose version &>/dev/null; then
  echo "⚠️  建议升级到新版 Docker Compose（V2）"
  echo "   当前版本: $(docker compose version 2>&1 || echo '未知')"
else
  echo "  ✅ Docker Compose: $(docker compose version)"
fi

# ─── 2. .env 检查 ────────────────────────────────────
echo ""
echo "📋 步骤 2/5：环境变量配置"

ENV_FILE="deploy/config/development/.env"
if [ ! -f "$ENV_FILE" ]; then
  if [ -f "deploy/config/development/.env.example" ]; then
    cp "deploy/config/development/.env.example" "$ENV_FILE"
    echo "  📝 已从 .env.example 创建 .env 文件"
    echo "  ⚠️  请根据需要修改 .env 中的配置（特别是 JWT_SECRET）"
  else
    echo "  ⚠️  未找到 .env.example，将使用默认配置"
  fi
else
  echo "  ✅ .env 文件已存在"
fi

# 加载 .env（如果存在）
if [ -f "$ENV_FILE" ]; then
  set -a
  source "$ENV_FILE"
  set +a
  echo "  ✅ .env 已加载"
fi

# ─── 3. 前端构建检查 ────────────────────────────────
echo ""
echo "📋 步骤 3/5：前端构建检查"

if [ -d "frontend/dist" ]; then
  DIST_SIZE=$(du -sh "frontend/dist" 2>/dev/null | cut -f1)
  echo "  ✅ 前端构建产物存在 (${DIST_SIZE:-未知})"
  if [ -f "frontend/dist/index.html" ]; then
    echo "  ✅ index.html 确认存在"
  else
    echo "  ⚠️  dist 目录存在但缺少 index.html，可能构建不完整"
  fi
else
  echo "  ⚠️  前端未构建，尝试构建..."
  if command -v npm &>/dev/null && [ -f "frontend/package.json" ]; then
    cd frontend
    npm install --silent
    npm run build
    cd "$PROJECT_ROOT"
    echo "  ✅ 前端构建完成"
  else
    echo "  ⚠️  无法自动构建前端，请手动构建：cd frontend && npm install && npm run build"
  fi
fi

# ─── 4. 构建并启动服务 ──────────────────────────────
echo ""
echo "📋 步骤 4/5：构建并启动服务"

echo "  🚀 启动 Docker Compose..."
docker compose -f deploy/docker/docker-compose.yml up -d --build 2>&1 | sed 's/^/  /'

echo "  ✅ 服务已启动"

# ─── 5. 健康检查 ────────────────────────────────────
echo ""
echo "📋 步骤 5/5：健康检查"

API_PORT="${API_PORT:-8081}"
WEB_PORT="${WEB_PORT:-8080}"

echo "  ⏳ 等待后端就绪..."
for i in $(seq 1 20); do
  if curl -sf "http://localhost:${API_PORT}/health" >/dev/null 2>&1; then
    echo ""
    echo "╔════════════════════════════════════════════╗"
    echo "║  ✅ 调试环境就绪！                        ║"
    echo "╠════════════════════════════════════════════╣"
    echo "║  后端 API:  http://localhost:${API_PORT}      ║"
    echo "║  前端 Web:  http://localhost:${WEB_PORT}      ║"
    echo "║  数据库:    localhost:${DB_PORT:-5432}         ║"
    echo "║  调试端口:  localhost:${DEBUG_PORT:-2345}      ║"
    echo "╚════════════════════════════════════════════╝"
    echo ""
    echo "📖 调试手册见：docs/deploy/调试手册.md"
    exit 0
  fi
  echo -n "."
  sleep 2
done

echo ""
echo "⚠️  健康检查超时（40 秒），请查看日志："
echo "   docker compose -f deploy/docker/docker-compose.yml logs"
echo ""
echo "📋 常用排查命令："
echo "   查看所有日志：  docker compose -f deploy/docker/docker-compose.yml logs -f"
echo "   查看后端日志：  docker compose -f deploy/docker/docker-compose.yml logs backend"
echo "   查看数据库日志：docker compose -f deploy/docker/docker-compose.yml logs db"
echo "   重启服务：      docker compose -f deploy/docker/docker-compose.yml restart"
exit 1