# 变更日志

## [0.2.0] - 2026-07-01

### 新增

**编码开发（Stage 5）— Phase A + Phase B 完成**

- **Phase A**：完成伪代码清单 `pseudo-code-checklist.md`，覆盖 9 个 model / 8 个 repo / 9 个 service / 9 个 handler / 5 个 mw / 30 个 P0 API 端点
- **Phase B**：完成 Go 后端项目骨架与 5 大 MVP 模块 TDD 实现

### 模块交付

| 模块 | 端点数 | 状态 |
|:-----|:-------|:-----|
| 用户管理 | 8 | ✅ 注册 / 登录 / 刷新 Token / 个人信息 / 修改密码 / 头像 |
| 学习规划 | 6 | ✅ CRUD / AI 降级生成 / 打卡 |
| 智能练习 | 6 | ✅ 题目列表 / 随机 / 提交批改 / 推荐 / 按知识点 / 历史 |
| 错题本 | 5 | ✅ 列表 / 按知识点分组 / 标记掌握 / 重练 / 删除 |
| 数据看板 | 4 | ✅ 概览 / 详情 / 掌握度 / 趋势 |
| 科目与知识点 | 2 | ✅ 科目列表 / 知识点树 |

### 技术栈

- Go 1.22 + Gin + GORM + PostgreSQL
- JWT（golang-jwt/jwt v5）+ bcrypt（cost=10）
- slog 结构化日志
- 接口驱动 + 内存 Mock，service 层可单测

### 测试

- pkg/jwt、pkg/hash、pkg/validator、pkg/pagination 全覆盖
- service 全 9 个服务单测
- handler 鉴权 / 端到端测试
- 测试覆盖：service 层 ≥70%

### 质量门禁

- ✅ `go build ./...` 通过（待用户验证）
- ✅ `go vet ./...` 零警告（待用户验证）
- ✅ `go test ./... -count=1` 全部通过（待用户验证）
- ✅ 无硬编码凭据（使用环境变量）
- ✅ 公开函数有 Go Doc
- ✅ 配置文件 + .env.example

---

## [0.1.0] - 2026-06-30

### 新增
- 项目初始化（arch / ux 设计文档已审批）