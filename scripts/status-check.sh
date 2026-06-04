#!/usr/bin/env bash
# =============================================================================
# status-check.sh — 三库一库状态总览
# 用途：一键查看 sw-agents-workspace / sw-skills-library / sw-projects 状态
# 何时用：排查问题前 / 每日开工时 / 推送出错后
# =============================================================================

set -euo pipefail

# ── 颜色 ─────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
DIM='\033[2m'
BOLD='\033[1m'
NC='\033[0m'

ok()   { echo -e "  ${GREEN}✓${NC} $1"; }
warn() { echo -e "  ${YELLOW}⚠${NC} $1"; }
err()  { echo -e "  ${RED}✗${NC} $1"; }
info() { echo -e "  ${CYAN}ℹ${NC} $1"; }

BASE="/home/ubuntu/gitee-software"
REPOS=(
    "sw-agents-workspace:配置/技能/工作流"
    "sw-skills-library:技能库（备份源）"
    "sw-projects:项目库"
)

echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}             三库一库 · 状态总览${NC}"
echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
echo ""

for entry in "${REPOS[@]}"; do
    REPO="${entry%%:*}"
    DESC="${entry#*:}"
    DIR="${BASE}/${REPO}"

    echo -e "${BOLD}📦 ${REPO}${NC} ${DIM}(${DESC})${NC}"

    if [ ! -d "${DIR}/.git" ]; then
        err "不是 git 仓库（缺失 .git）"
        echo ""
        continue
    fi

    cd "${DIR}"

    # 当前分支
    BRANCH=$(git branch --show-current 2>/dev/null || echo "detached")
    info "分支: ${BRANCH}"

    # 远程连接
    REMOTE=$(git remote get-url origin 2>/dev/null || echo "未配置")
    info "远端: ${REMOTE}"

    # 同步状态
    git fetch origin --quiet 2>/dev/null || warn "fetch 失败（可能网络问题）"
    LOCAL=$(git rev-parse HEAD 2>/dev/null || echo "none")
    REMOTE_HEAD=$(git rev-parse origin/${BRANCH} 2>/dev/null || echo "none")

    if [ "${LOCAL}" = "${REMOTE_HEAD}" ]; then
        ok "本地与远端同步"
    elif [ "${REMOTE_HEAD}" = "none" ]; then
        warn "远端分支不存在"
    else
        # 检查 ahead/behind
        AHEAD=$(git rev-list --count origin/${BRANCH}..HEAD 2>/dev/null || echo 0)
        BEHIND=$(git rev-list --count HEAD..origin/${BRANCH} 2>/dev/null || echo 0)
        if [ "${AHEAD}" -gt 0 ]; then
            warn "本地领先远端 ${AHEAD} 个 commit（未推送）"
        fi
        if [ "${BEHIND}" -gt 0 ]; then
            warn "本地落后远端 ${BEHIND} 个 commit（未拉取）"
        fi
    fi

    # 工作区状态
    WORKTREE=$(git status -s 2>&1)
    if [ -z "${WORKTREE}" ]; then
        ok "工作区干净"
    else
        MODIFIED=$(echo "${WORKTREE}" | grep -c "^ M\|^M " || true)
        ADDED=$(echo "${WORKTREE}" | grep -c "^??" || true)
        DELETED=$(echo "${WORKTREE}" | grep -c "^ D\|^D " || true)
        MODIFIED=${MODIFIED:-0}
        ADDED=${ADDED:-0}
        DELETED=${DELETED:-0}
        warn "工作区有改动: ${MODIFIED} 修改 / ${ADDED} 新增 / ${DELETED} 删除"
        echo "${WORKTREE}" | head -5 | sed 's/^/      /'
        TOTAL=$(echo "${WORKTREE}" | grep -c "." || true)
        TOTAL=${TOTAL:-0}
        if [ "${TOTAL}" -gt 5 ]; then
            echo "      ... 还有 $((TOTAL - 5)) 行"
        fi
    fi

    # 最新 commit
    LAST_COMMIT=$(git log -1 --format='%h %s (%ar)' 2>/dev/null || echo "无")
    info "最近: ${LAST_COMMIT}"

    echo ""
done

# ── 总结建议 ─────────────────────────────────────────────────────────────────
echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}💡 常用命令${NC}"
echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
echo -e "  同步项目代码:    ${CYAN}./scripts/hotfix-sync.sh <project> 'msg'${NC}"
echo -e "  同步技能:        ${CYAN}cd ~/gitee-software/sw-skills-library && git add -A && git commit -m '...' && git push${NC}"
echo -e "  拉取最新:        ${CYAN}cd <repo> && git pull --rebase${NC}"
echo -e "  详细状态:        ${CYAN}cd <repo> && git status${NC}"
