# sw-projects 仓库工具脚本

翠花软件开发集群的**仓库级共享工具**，位于 `sw-projects/scripts/`。

> ⚠️ 这些脚本不属于任何单个项目，是仓库级别的工具。

## 脚本清单

| 脚本 | 用途 | 何时用 |
|:-----|:-----|:-------|
| **hotfix-sync.sh** | 本地修改 → add → commit → push | 改完代码，公子说"改吧"后 |
| **status-check.sh** | 三库一库状态总览 | 排查问题、每日开工 |
| **pull-sync.sh** | 拉取三库最新到本地 | 开机、重启服务前、怀疑本地过时 |

## 快速上手

```bash
cd ~/gitee-software/sw-projects

# 看当前状态
./scripts/status-check.sh

# 改完代码，推到 Gitee
./scripts/hotfix-sync.sh sw-dashboard "修复日志路径不存在"

# 拉取最新
./scripts/pull-sync.sh
```

## hotfix-sync.sh 详解

**核心设计**：每步都暂停等你确认，**绝不自动 commit/push**。

### 工作流

```
1. 检查 git 状态
   ├── 工作区干净 → 提示无需同步
   ├── 仅有目标项目改动 → 继续
   └── 有非目标项目改动 → 拒绝（防误推）
       ↓
2. 显示 diff（stat + 前 50 行）
   ↓
3. 暂停等你按 Enter
   ↓
4. 输入 commit message
   ↓
5. 第二次确认（y/N）
   ↓
6. 执行 add → commit → push
```

### 高级用法

```bash
# Dry-run（只显示，不执行）
DRY_RUN=1 ./scripts/hotfix-sync.sh sw-dashboard "测试"

# 同步仓库根文件（scripts/、.gitignore 等）
ALLOW_NONPROJECT=1 ./scripts/hotfix-sync.sh sw-dashboard "新增 hotfix 脚本"

# 跳过所有确认（仅翠花自动调用，慎用）
SKIP_CONFIRM=1 ./scripts/hotfix-sync.sh sw-dashboard "..."
```

### 安全机制

- ✅ **拒绝混推**：工作区有非目标项目改动 → 中止
- ✅ **拒绝大改**：diff 超过 300 行 → 警告
- ✅ **必须 confirm**：y/N 二次确认才能 commit
- ✅ **可回滚**：每次都打印 commit hash，方便 `git reset`

## 必坑（2026-06-04 教训）

- **绝不自动 push**：脚本需要你显式调用
- **绝不自动重启服务**：commit 完成后告诉你怎么重启，但不替你做
- **紧急修复**：可先改 + 重启，事后用 hotfix-sync.sh 同步
- **工作区不干净**：pull-sync 会跳过，hotfix-sync 会拒绝——这是设计
