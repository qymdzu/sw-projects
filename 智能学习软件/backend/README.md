# 智能学习软件 — 后端服务

> 版本：v0.1.0  
> 语言：Go 1.22+  
> 框架：Gin + GORM  
> 数据库：PostgreSQL 15

---

## 快速开始

### 1. 本地开发

```bash
# 复制环境变量模板
cp .env.example .env

# 下载依赖
go mod tidy

# 运行单元测试
make test

# 启动服务（需要先有 PostgreSQL）
make run
```

### 2. Docker 构建与运行

```bash
docker build -t smart-learning-server:latest .
docker run --rm -p 8080:8080 --env-file .env smart-learning-server:latest
```

### 3. 健康检查

```bash
curl http://localhost:8080/health
```

---

## 目录结构

```
backend/
├── cmd/server/              程序入口
├── internal/
│   ├── config/              配置加载
│   ├── model/               GORM 数据模型（9 张表）
│   ├── repository/          数据访问层（GORM 实现）
│   ├── service/             业务逻辑层
│   ├── handler/             HTTP 处理器
│   ├── middleware/          JWT / CORS / Logger / Recovery / RateLimit
│   ├── router/              路由注册
│   └── dto/                 （预留）请求/响应 DTO
├── pkg/
│   ├── jwt/                 JWT 生成与校验
│   ├── hash/                bcrypt 密码哈希
│   ├── response/            统一响应封装
│   ├── pagination/          分页参数解析
│   ├── validator/           入参校验
│   └── logger/              slog 封装
├── go.mod
├── Makefile
└── Dockerfile
```

详见 `docs/design/目录结构.md`。

---

## 模块清单（MVP）

| 模块 | Handler | Service | Repository |
|:-----|:--------|:--------|:-----------|
| 用户管理 | auth_handler, user_handler | auth_service, user_service | user_repo |
| 学习规划 | plan_handler | plan_service | plan_repo |
| 智能练习 | exercise_handler | exercise_service | exercise_repo |
| 错题本 | mistake_handler | mistake_service | mistake_repo |
| 数据看板 | report_handler | report_service | report_repo, exercise_repo, plan_repo, mistake_repo |
| 科目 | subject_handler | subject_service | subject_repo |
| 知识点 | knowledge_handler | knowledge_service | knowledge_repo |

---

## API 端点

详见 `docs/design/API设计.md`（30 个 MVP 端点全部实现）。

---

## 测试

- `make test` — 跑全量单元测试
- `make cover` — 生成覆盖率报告（service 层目标 ≥70%）

测试覆盖：
- pkg/jwt：Token 生成 / 解析 / 刷新 / 过期 / 类型错误
- pkg/hash：bcrypt 加解密 / 密码校验
- pkg/validator：手机号 / 邮箱 / 密码强度 / 角色
- pkg/pagination：默认值 / 边界 / 上限
- service/auth：注册 / 登录 / 刷新 / 重复 / 弱密码
- service/user：个人信息 / 修改密码 / 头像
- service/plan：CRUD / AI 生成 / 打卡
- service/exercise：批改 / 自动收录错题 / 推荐
- service/mistake：错题 CRUD / 按知识点分组 / 重练
- service/report：概览 / 详情 / 掌握度 / 趋势
- handler：HTTP 端到端测试（鉴权 / 错误码 / 边界）

---

## 环境变量

| 变量 | 默认值 | 说明 |
|:-----|:-------|:------|
| `SERVER_PORT` | 8080 | HTTP 端口 |
| `SERVER_MODE` | debug | Gin 模式（debug/release/test） |
| `DB_HOST` | localhost | PostgreSQL 主机 |
| `DB_PORT` | 5432 | PostgreSQL 端口 |
| `DB_USER` | smart_learning | 数据库用户 |
| `DB_PASSWORD` | — | 数据库密码（必填） |
| `DB_NAME` | smart_learning | 数据库名 |
| `JWT_SECRET` | — | JWT 密钥（必填，生产环境 ≥64 字符） |
| `JWT_ACCESS_TTL` | 2h | access_token 有效期 |
| `JWT_REFRESH_TTL` | 168h | refresh_token 有效期 |
| `LOG_LEVEL` | info | 日志级别 |

---

## 后续演进（V1.5）

- [ ] 对接 LLM 真实推荐（替换规则降级）
- [ ] 知识点树可视化（V1.5 模块）
- [ ] 操作日志表（operation_logs）
- [ ] Redis 缓存层（NF-02/NF-04）
- [ ] 读写分离
- [ ] 异步报告生成