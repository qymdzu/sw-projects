#!/usr/bin/env bash
# =============================================================================
# deploy.sh — sw-dashboard 一键部署脚本 (systemd 直接启动)
# 适用环境: 旧云服务器 (512MB 内存), Ubuntu 22.04+
# 部署方式: systemd + uvicorn 直接启动 (非 Docker)
# =============================================================================

set -euo pipefail

# ── 颜色输出 ─────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info()  { echo -e "${CYAN}[INFO]${NC} $*"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}   $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ── 配置 ─────────────────────────────────────────────────────────────────────
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SERVICE_NAME="sw-dashboard"
SERVICE_PORT=8900
UVICORN_WORKERS=1

PYTHON_BIN="${PYTHON_BIN:-python3}"
NPM_BIN="${NPM_BIN:-npm}"
NODE_BIN="${NODE_BIN:-node}"

VENV_DIR="${PROJECT_DIR}/venv"
FRONTEND_DIR="${PROJECT_DIR}/frontend"
BACKEND_DIR="${PROJECT_DIR}/backend"
SYSTEMD_SERVICE="/etc/systemd/system/${SERVICE_NAME}.service"

# ── 环境检查 ─────────────────────────────────────────────────────────────────
check_env() {
    log_info "检查运行环境..."

    # Token 检查
    if [ -z "${HERMES_DASHBOARD_TOKEN:-}" ]; then
        log_warn "HERMES_DASHBOARD_TOKEN 未设置，将使用默认开发 Token"
        log_warn "生产环境请务必设置: export HERMES_DASHBOARD_TOKEN='your-strong-token-here'"
    fi

    # Python 3.12+
    if ! command -v "$PYTHON_BIN" &>/dev/null; then
        log_error "Python3 未安装，请先安装: sudo apt install python3 python3-venv python3-pip"
        exit 1
    fi
    PY_VER=$("$PYTHON_BIN" --version 2>&1 | grep -oP '\d+\.\d+')
    log_ok "Python: $("$PYTHON_BIN" --version 2>&1)"

    # Node 18+
    if ! command -v "$NODE_BIN" &>/dev/null; then
        log_error "Node.js 未安装，请先安装: sudo apt install nodejs npm"
        exit 1
    fi
    NODE_VER=$("$NODE_BIN" --version 2>&1 | grep -oP '\d+' | head -1)
    if [ "$NODE_VER" -lt 18 ]; then
        log_error "Node.js 版本过低 (当前 $("$NODE_BIN" --version))，需要 18+"
        exit 1
    fi
    log_ok "Node.js: $("$NODE_BIN" --version)"

    # 检查虚拟环境
    if [ ! -d "$VENV_DIR" ]; then
        log_warn "Python 虚拟环境未创建，将自动创建"
    fi

    # 检查 systemd
    if ! command -v systemctl &>/dev/null; then
        log_error "systemd 不可用，当前环境不支持此部署方式"
        exit 1
    fi
    log_ok "systemd 可用"

    return 0
}

# ── 安装后端依赖 ─────────────────────────────────────────────────────────────
install_backend_deps() {
    log_info "安装后端 Python 依赖..."

    if [ ! -d "$VENV_DIR" ]; then
        "$PYTHON_BIN" -m venv "$VENV_DIR"
        log_ok "Python 虚拟环境已创建: $VENV_DIR"
    fi

    source "${VENV_DIR}/bin/activate"
    pip install --upgrade pip setuptools wheel -q
    pip install -r "${BACKEND_DIR}/requirements.txt" -q
    deactivate

    log_ok "后端依赖安装完成"
}

# ── 构建前端 ─────────────────────────────────────────────────────────────────
build_frontend() {
    log_info "构建前端..."

    if [ ! -d "${FRONTEND_DIR}/node_modules" ]; then
        log_info "安装前端 Node 依赖..."
        cd "$FRONTEND_DIR"
        $NPM_BIN ci
        cd "$PROJECT_DIR"
    fi

    cd "$FRONTEND_DIR"
    $NPM_BIN run build
    cd "$PROJECT_DIR"

    if [ ! -d "${PROJECT_DIR}/frontend/dist" ]; then
        log_error "前端构建失败：dist 目录未生成"
        exit 1
    fi

    log_ok "前端构建完成"
}

# ── 配置 systemd 服务 ────────────────────────────────────────────────────────
setup_systemd() {
    log_info "配置 systemd 服务..."

    # 构建 ExecStart 命令
    local VENV_PYTHON="${VENV_DIR}/bin/python"
    local EXEC_START="${VENV_PYTHON} -m uvicorn main:app --host 0.0.0.0 --port ${SERVICE_PORT} --workers ${UVICORN_WORKERS}"

    # 服务配置
    sudo tee "$SYSTEMD_SERVICE" > /dev/null <<EOF
[Unit]
Description=sw-dashboard — 翠花集群管理面板
After=network.target

[Service]
Type=simple
User=${USER}
WorkingDirectory=${BACKEND_DIR}
ExecStart=${EXEC_START}
Restart=on-failure
RestartSec=5
StartLimitIntervalSec=60
StartLimitBurst=3

# 环境变量
Environment=HERMES_DASHBOARD_TOKEN=${HERMES_DASHBOARD_TOKEN:-}

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${SERVICE_NAME}

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload

    log_ok "systemd 服务配置完成: ${SYSTEMD_SERVICE}"
}

# ── 启动服务 ─────────────────────────────────────────────────────────────────
start_service() {
    log_info "启动服务..."

    sudo systemctl enable "${SERVICE_NAME}" 2>/dev/null || true
    sudo systemctl restart "${SERVICE_NAME}"

    # 等待服务启动
    sleep 3
    if sudo systemctl is-active --quiet "${SERVICE_NAME}"; then
        log_ok "服务启动成功！"
        sudo systemctl status "${SERVICE_NAME}" --no-pager -l | head -15
    else
        log_error "服务启动失败，查看日志: sudo journalctl -u ${SERVICE_NAME} -n 50 --no-pager"
        sudo journalctl -u "${SERVICE_NAME}" -n 20 --no-pager
        exit 1
    fi
}

# ── 检查服务健康 ──────────────────────────────────────────────────────────────
health_check() {
    log_info "执行健康检查..."
    sleep 2
    local STATUS_URL="http://localhost:${SERVICE_PORT}/api/health"

    if command -v curl &>/dev/null; then
        if curl -sf "$STATUS_URL" > /dev/null 2>&1; then
            log_ok "健康检查通过: ${STATUS_URL}"
        else
            log_warn "健康检查未通过，请检查服务日志"
        fi
    else
        log_warn "curl 未安装，跳过健康检查"
    fi
}

# ── 主流程 ────────────────────────────────────────────────────────────────────
main() {
    echo ""
    echo -e "${GREEN}=================================================${NC}"
    echo -e "${GREEN}  sw-dashboard 一键部署脚本${NC}"
    echo -e "${GREEN}  项目: ${PROJECT_DIR}${NC}"
    echo -e "${GREEN}  端口: ${SERVICE_PORT}${NC}"
    echo -e "${GREEN}=================================================${NC}"
    echo ""

    check_env
    install_backend_deps
    build_frontend
    setup_systemd
    start_service
    health_check

    echo ""
    echo -e "${GREEN}=================================================${NC}"
    echo -e "${GREEN}  部署完成！${NC}"
    echo -e "${GREEN}  访问地址: http://localhost:${SERVICE_PORT}${NC}"
    echo -e "${GREEN}  服务管理: sudo systemctl {start|stop|restart|status} ${SERVICE_NAME}${NC}"
    echo -e "${GREEN}  查看日志: sudo journalctl -u ${SERVICE_NAME} -f${NC}"
    echo -e "${GREEN}=================================================${NC}"
    echo ""
}

# ── 执行 ──────────────────────────────────────────────────────────────────────
main "$@"
