# sw-projects — 软件开发项目源代码

> Pipeline 8 阶段产出的项目源代码仓库。由 `cuihua_sw_pipeline__scaffold` 创建骨架。

## 架构目录

```
sw-projects/
├── <project-name>/           # 每项目独立子目录
│   ├── .project-meta.yaml    # 项目元数据（类型 / Stage / 历史）
│   │
│   ├── docs/
│   │   ├── requirements/     # 需求规格说明书.md + 待确认问题.md
│   │   ├── research/         # 技术方案说明书.md + 技术调研清单.md
│   │   ├── design/           # 系统架构设计.md / 数据模型.md / API.md / 目录结构.md
│   │   ├── ux/               # 交互流程设计.md / UI组件规范.md
│   │   ├── review/           # 代码审查报告.md
│   │   ├── test/             # 测试计划.md / 测试用例.md / 测试报告.md
│   │   └── deploy/           # 部署手册.md
│   │
│   ├── design/               # 架构师设计文档
│   ├── backend/              # Go 源码 + _test.go
│   ├── scripts/              # 部署脚本
│   ├── deploy/docker/        # Dockerfile + docker-compose.yml
│   └── CHANGELOG.md          # 版本变更记录
│
├── ai-smart-learner/         # Flutter iOS app（智能错题熔断）
├── sw-dashboard/             # 系统监控面板
└── 项目存储结构规范.md       # 骨架创建规范
```

## 项目形态

| 形态 | 保留目录 | 示例 |
|:-----|:---------|:-----|
| `fullstack` | 全部 | sw-dashboard |
| `backend-only` | 去 frontend/ + ai/ | — |
| `frontend-only` | 去 backend/ + ai/ | ai-smart-learner |
| `tool-script` | 仅 scripts/ + docs/ | — |
| `ai-service` | 去 frontend/ | — |

## 生命周期

```
公子提需求
  → Scaffold 创建骨架
  → Pipeline 8 阶段执行（契约校验 + 门禁 + 审阅）
  → 每阶段结束推 Gitee
  → 全流程完成 → 交付公子
```

## Pipeline 工具（sw-agents-workspace）

| 命令 | 用途 |
|:-----|:------|
| `sw-pipeline-runner.py validate <项目>` | 门禁校验 |
| `sw-pipeline-runner.py advance <项目>` | Stage 推进 |
| `sw-pipeline-runner.py status <项目>` | 状态查看 |
| `sw-dispatch.py next <项目>` | 下一动作指令 |
| `sw-pipeline-dashboard.py` | 全局仪表板 |

## 相关仓库

- [sw-agents-workspace](https://gitee.com/qymdzu/sw-agents-workspace) — Pipeline 脚本 + 契约
- [sw-mcp-servers](https://gitee.com/qymdzu/sw-mcp-servers) — MCP 工具
- [sw-dev-knowledge](https://gitee.com/qymdzu/sw-dev-knowledge) — 架构文档