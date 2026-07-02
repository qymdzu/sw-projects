#!/bin/bash
# ============================================================
# 智能学习软件 - 调试环境启动脚本
# 项目路径: /home/ubuntu/gitee-software/sw-projects/智能学习软件
# ============================================================

set -e
PROJECT_DIR="/home/ubuntu/gitee-software/sw-projects/智能学习软件"
BACKEND_DIR="$PROJECT_DIR/backend"
SERVER_BIN="/tmp/smart-server"
LOG_FILE="/tmp/smart-server.log"
API_PORT=8081
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=smart_learning
DB_USER=smart_admin
DB_PASS=smart_dev_2026

echo "============================================"
echo "  智能学习软件 - 调试环境启动"
echo "============================================"

# Step 1: Check Go
echo "[1/5] 检查 Go 环境..."
if ! /usr/local/go/bin/go version >/dev/null 2>&1; then
    echo "❌ Go 未安装，请先安装 Go 1.23+"
    exit 1
fi
echo "✅ Go $(/usr/local/go/bin/go version)"

# Step 2: Check PostgreSQL
echo "[2/5] 检查 PostgreSQL..."
if ! pg_isready -h $DB_HOST -p $DB_PORT >/dev/null 2>&1; then
    echo "⚠️  PostgreSQL 未运行，尝试启动..."
    sudo systemctl start postgresql
    sleep 2
    if ! pg_isready -h $DB_HOST -p $DB_PORT; then
        echo "❌ PostgreSQL 启动失败"
        exit 1
    fi
fi
echo "✅ PostgreSQL 运行中 ($DB_HOST:$DB_PORT/$DB_NAME)"

# Step 3: Check database exists
echo "[3/5] 检查数据库..."
if ! PGPASSWORD=$DB_PASS psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;" >/dev/null 2>&1; then
    echo "⚠️  数据库不存在，创建..."
    sudo -u postgres psql -c "CREATE USER smart_admin WITH PASSWORD '$DB_PASS';" 2>/dev/null || true
    sudo -u postgres psql -c "CREATE DATABASE smart_learning OWNER smart_admin;" 2>/dev/null || true
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE smart_learning TO smart_admin;" 2>/dev/null || true
    echo "✅ 数据库已创建"
fi

# Step 4: Build backend
echo "[4/5] 编译后端..."
cd $BACKEND_DIR
export GOPROXY=https://goproxy.cn,direct
go build -o $SERVER_BIN ./cmd/server 2>&1
echo "✅ 编译完成 ($SERVER_BIN)"

# Step 5: Start server
echo "[5/5] 启动服务..."
# Kill old instance if any
pkill -f "smart-server" 2>/dev/null || true
sleep 1

nohup $SERVER_BIN > $LOG_FILE 2>&1 &
SERVER_PID=$!
echo "✅ 服务已启动 (PID: $SERVER_PID)"

# Wait for server to be ready
echo "等待服务就绪..."
for i in $(seq 1 10); do
    if curl -s http://127.0.0.1:$API_PORT/health >/dev/null 2>&1; then
        echo "✅ 服务就绪!"
        break
    fi
    sleep 2
done

echo ""
echo "============================================"
echo "  调试信息"
echo "============================================"
echo "  API 地址:  http://127.0.0.1:$API_PORT"
echo "  服务器公网: http://100.117.73.92:$API_PORT"
echo "  iOS App:   配置 API 地址为 http://100.117.73.92:$API_PORT"
echo "  日志文件:  $LOG_FILE"
echo "  停止服务:  pkill -f smart-server"
echo "============================================"
