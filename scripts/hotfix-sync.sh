#!/usr/bin/env bash
# =============================================================================
# hotfix-sync.sh — 翠花热修复同步工具
# 用途：单文件/小改小补后，把本地修改 → add → commit → push 到 Gitee
# 设计：每步都暂停等你确认（绝不自动 commit/push）
# 范围：仅处理 sw-projects 仓库（项目代码），技能/配置走其他脚本
# =============================================================================

set -euo pipefail

# ── 颜色 ─────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
DIM='\033[2m'
NC='\033[0m'

log_info()  { echo -e "${CYAN}[INFO]${NC} $*"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}   $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ── 帮助 ─────────────────────────────────────────────────────────────────────
usage() {
    cat <<EOF
用法: hotfix-sync.sh <project-name> [commit-message]

例子:
    hotfix-sync.sh sw-dashboard "修复日志路径不存在"
    hotfix-sync.sh sw-dashboard

参数:
    project-name        项目子目录名（必须存在于 sw-projects/ 下）
    commit-message      commit 信息（可选，不传则进入交互式输入）

流程:
    1. 检查 git 状态（必须干净或仅有目标项目改动）
    2. 显示当前未提交改动（diff + 状态）
    3. 暂停等你看，按 Enter 继续
    4. 再次确认（y/n）→ add → commit → push
    5. 显示推送结果

环境变量:
    SKIP_CONFIRM=1      跳过所有确认（仅翠花自动调用时使用，慎用）
    DRY_RUN=1           只显示操作，不实际执行
EOF
}

# ── 参数检查 ─────────────────────────────────────────────────────────────────
if [ $# -lt 1 ]; then
    usage
    exit 1
fi

PROJECT="$1"
COMMIT_MSG="${2:-}"

REPO_ROOT="/home/ubuntu/gitee-software/sw-projects"
PROJECT_DIR="${REPO_ROOT}/${PROJECT}"

# ── 前置检查 ─────────────────────────────────────────────────────────────────
log_info "═══ 热修复同步工具 ═══"
log_info "项目: ${PROJECT}"
log_info "仓库: ${REPO_ROOT}"
echo ""

# 检查仓库目录
if [ ! -d "${REPO_ROOT}/.git" ]; then
    log_error "sw-projects 仓库不存在: ${REPO_ROOT}"
    exit 1
fi

# 检查项目目录（如果存在）；不存在但有 ALLOW_NONPROJECT=1 则放行
if [ ! -d "${PROJECT_DIR}" ]; then
    if [ "${ALLOW_NONPROJECT:-0}" = "1" ]; then
        log_warn "项目目录不存在: ${PROJECT_DIR}（ALLOW_NONPROJECT=1，继续）"
    else
        log_error "项目目录不存在: ${PROJECT_DIR}"
        log_info "现有项目:"
        ls -1 "${REPO_ROOT}" | grep -v "^\." | sed 's/^/  - /'
        log_info ""
        log_info "如果要同步仓库根的文件（如 scripts/、.gitignore 等），请设置:"
        log_info "  ALLOW_NONPROJECT=1 hotfix-sync.sh sw-dashboard '...'"
        exit 1
    fi
fi

cd "${REPO_ROOT}"

# 检查 git 状态
log_info "检查 git 状态..."
GIT_STATUS=$(git status -s 2>&1)

if [ -z "${GIT_STATUS}" ]; then
    log_warn "工作区无任何改动，无需同步"
    log_info "（这通常意味着你还没改代码，或已经提交过了）"
    exit 0
fi

# 检查是否仅有目标项目改动（防止误推别人的工作）
# 当 ALLOW_NONPROJECT=1 时跳过此检查
if [ "${ALLOW_NONPROJECT:-0}" != "1" ]; then
    OTHER_CHANGES=$(echo "${GIT_STATUS}" | grep -v "^.. ${PROJECT}/" | grep -v "^?? ${PROJECT}/" || true)
    if [ -n "${OTHER_CHANGES}" ]; then
        log_error "工作区有非 '${PROJECT}/' 的改动，为防止误推已中止:"
        echo "${OTHER_CHANGES}" | sed 's/^/  /'
        echo ""
        log_info "如确认是同一项目的多个子目录，请手动 git add"
        log_info "如要同步仓库根文件（如 scripts/、.gitignore），请设置 ALLOW_NONPROJECT=1"
        exit 1
    fi
fi

# ── 显示当前改动 ─────────────────────────────────────────────────────────────
log_ok "检测到 '${PROJECT}/' 下的改动:"
echo ""
echo -e "${DIM}───── git status ─────${NC}"
git status -s "${PROJECT}/" | sed 's/^/  /'
echo ""
echo -e "${DIM}───── diff stat ─────${NC}"
git diff --stat "${PROJECT}/" 2>/dev/null | sed 's/^/  /' || echo "  (无 tracked 文件改动，仅新增)"
git diff --cached --stat "${PROJECT}/" 2>/dev/null | sed 's/^/  /' || true
echo ""

# 显示完整 diff 预览（前 50 行）
log_info "diff 预览（前 50 行）:"
echo -e "${DIM}────────────────────────────────────────${NC}"
git diff "${PROJECT}/" 2>/dev/null | head -50 || true
git diff --cached "${PROJECT}/" 2>/dev/null | head -50 || true
echo -e "${DIM}────────────────────────────────────────${NC}"
echo ""

# 统计行数
DIFF_LINES=$(git diff "${PROJECT}/" 2>/dev/null | wc -l || echo 0)
DIFF_CACHED_LINES=$(git diff --cached "${PROJECT}/" 2>/dev/null | wc -l || echo 0)
TOTAL_LINES=$((DIFF_LINES + DIFF_CACHED_LINES))
log_info "差异总行数: ${TOTAL_LINES}"

if [ "${TOTAL_LINES}" -gt 300 ]; then
    log_warn "差异较大（${TOTAL_LINES} 行），请确认是否合理"
fi

# ── 第一次确认：看 diff ──────────────────────────────────────────────────────
if [ "${SKIP_CONFIRM:-0}" != "1" ]; then
    read -rp "$(echo -e "${CYAN}→ 按 Enter 继续，Ctrl+C 中止...${NC}")" || {
        log_warn "已中止"
        exit 1
    }
fi

# ── 决定 commit message ──────────────────────────────────────────────────────
if [ -z "${COMMIT_MSG}" ]; then
    echo ""
    log_info "请输入 commit 信息（feat/fix/chore/docs 开头）:"
    read -rp "→ " COMMIT_MSG
    if [ -z "${COMMIT_MSG}" ]; then
        log_error "commit 信息不能为空"
        exit 1
    fi
fi

# ── 第二次确认：是否 commit + push ───────────────────────────────────────────
echo ""
log_info "准备执行:"
echo -e "  ${CYAN}git add ${PROJECT}/${NC}"
echo -e "  ${CYAN}git commit -m \"${COMMIT_MSG}\"${NC}"
echo -e "  ${CYAN}git push origin master${NC}"
echo ""

if [ "${SKIP_CONFIRM:-0}" != "1" ] && [ "${DRY_RUN:-0}" != "1" ]; then
    read -rp "$(echo -e "${YELLOW}→ 确认执行？(y/N)${NC} ")" CONFIRM
    if [ "${CONFIRM}" != "y" ] && [ "${CONFIRM}" != "Y" ]; then
        log_warn "已中止（未执行任何操作）"
        exit 0
    fi
fi

# ── Dry-run 模式 ─────────────────────────────────────────────────────────────
if [ "${DRY_RUN:-0}" = "1" ]; then
    log_info "[DRY-RUN] 模式，不实际执行"
    log_info "将执行:"
    echo "  git add ${PROJECT}/"
    echo "  git commit -m \"${COMMIT_MSG}\""
    echo "  git push origin master"
    exit 0
fi

# ── 实际执行 ─────────────────────────────────────────────────────────────────
log_info "执行 add..."

# 决定 add 哪个路径
# 默认 add 项目目录；ALLOW_NONPROJECT=1 时 add 整个工作区
if [ "${ALLOW_NONPROJECT:-0}" = "1" ]; then
    git add -A
else
    git add "${PROJECT}/"
fi

log_info "执行 commit..."
git commit -m "${COMMIT_MSG}"

log_info "执行 push..."
git push origin master

# ── 完成 ─────────────────────────────────────────────────────────────────────
echo ""
log_ok "═══ 同步完成 ═══"
COMMIT_HASH=$(git log -1 --format='%h')
log_ok "Commit: ${COMMIT_HASH} - ${COMMIT_MSG}"
log_ok "远端:   https://gitee.com/qymdzu/sw-projects/commit/${COMMIT_HASH}"
echo ""
log_info "后续:"
echo "  - 如果服务需要重启:  sudo systemctl restart ${PROJECT}"
echo "  - 查看推送结果:      cd ${REPO_ROOT} && git log --oneline -3"
