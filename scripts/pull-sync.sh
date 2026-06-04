#!/usr/bin/env bash
# =============================================================================
# pull-sync.sh — 拉取三库最新到本地
# 用途：开机 / 重启服务前 / 怀疑本地过时
# 设计：每个仓库单独拉，失败不阻塞其他
# =============================================================================

set -uo pipefail

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

ok()   { echo -e "  ${GREEN}✓${NC} $1"; }
warn() { echo -e "  ${YELLOW}⚠${NC} $1"; }
err()  { echo -e "  ${RED}✗${NC} $1"; }
info() { echo -e "  ${CYAN}ℹ${NC} $1"; }

BASE="/home/ubuntu/gitee-software"
REPOS=(sw-agents-workspace sw-skills-library sw-projects)

echo -e "${CYAN}═══ 拉取三库最新 ═══${NC}"
echo ""

PULLED_TOTAL=0
FAILED_TOTAL=0

for repo in "${REPOS[@]}"; do
    DIR="${BASE}/${repo}"
    echo -e "📦 ${repo}"

    if [ ! -d "${DIR}/.git" ]; then
        err "不是 git 仓库，跳过"
        FAILED_TOTAL=$((FAILED_TOTAL + 1))
        echo ""
        continue
    fi

    cd "${DIR}"

    # 检查工作区是否干净（拉取 rebase 要求干净）
    if [ -n "$(git status -s 2>&1)" ]; then
        warn "工作区有未提交改动，跳过（先 commit 或 stash）"
        git status -s | head -3 | sed 's/^/      /'
        FAILED_TOTAL=$((FAILED_TOTAL + 1))
        echo ""
        continue
    fi

    # 拉取
    BEFORE=$(git rev-parse HEAD 2>/dev/null || echo "none")
    if git pull --rebase origin master 2>&1 | tail -5; then
        AFTER=$(git rev-parse HEAD 2>/dev/null || echo "none")
        if [ "${BEFORE}" = "${AFTER}" ]; then
            info "已是最新（${BEFORE:0:7}）"
        else
            ok "已更新: ${BEFORE:0:7} → ${AFTER:0:7}"
            PULLED_TOTAL=$((PULLED_TOTAL + 1))
        fi
    else
        err "拉取失败（可能冲突）"
        FAILED_TOTAL=$((FAILED_TOTAL + 1))
    fi
    echo ""
done

echo -e "${CYAN}═══ 总结 ═══${NC}"
echo "  拉取更新: ${PULLED_TOTAL} 个仓库"
echo "  跳过/失败: ${FAILED_TOTAL} 个仓库"

if [ "${FAILED_TOTAL}" -gt 0 ]; then
    echo ""
    warn "部分仓库未拉取，请检查"
    exit 1
fi
