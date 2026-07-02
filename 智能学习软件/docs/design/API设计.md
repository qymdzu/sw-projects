# API 设计 — 智能学习软件

> 版本：v1.0.0
> API 风格：RESTful JSON
> 基础路径：`/api/v1`
> 认证方式：JWT Bearer Token

---

## 1. API 规范

### 1.1 通用规范

| 维度 | 约定 |
|:-----|:------|
| 基础路径 | `/api/v1` |
| 命名规则 | 小写 + 连字符（kebab-case）|
| 数据格式 | 请求/响应均为 JSON，Content-Type: `application/json` |
| 字符编码 | UTF-8 |
| 认证方式 | `Authorization: Bearer <access_token>` |
| 分页参数 | `?page=1&page_size=20` |
| 时间格式 | ISO 8601（`2026-07-01T10:30:00+08:00`）|

### 1.2 统一响应格式

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**列表响应（带分页）：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

**错误响应：**
```json
{
  "code": 10001,
  "message": "参数校验失败",
  "detail": {
    "field": "email",
    "reason": "邮箱格式不正确"
  }
}
```

### 1.3 业务错误码

| 错误码 | HTTP 状态码 | 说明 |
|:-------|:------------|:------|
| 0 | 200 | 成功 |
| 10001 | 400 | 请求参数错误 |
| 10002 | 400 | 资源已存在 |
| 10003 | 400 | 资源状态不允许操作 |
| 20001 | 401 | 未认证 |
| 20002 | 401 | Token 已过期 |
| 20003 | 401 | Token 无效 |
| 20004 | 403 | 权限不足 |
| 30001 | 404 | 资源不存在 |
| 30002 | 409 | 资源冲突 |
| 40001 | 429 | 请求过于频繁 |
| 50001 | 500 | 服务器内部错误 |
| 50002 | 502 | 外部服务不可用 |

---

## 2. 用户模块 API

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| POST | `/api/v1/auth/register` | 用户注册 | 否 | ✅ |
| POST | `/api/v1/auth/login` | 用户登录 | 否 | ✅ |
| POST | `/api/v1/auth/refresh` | 刷新 Token | 否（需 refresh_token）| ✅ |
| GET | `/api/v1/users/me` | 获取当前用户信息 | 是 | ✅ |
| PUT | `/api/v1/users/me` | 更新个人信息 | 是 | ✅ |
| PUT | `/api/v1/users/me/password` | 修改密码 | 是 | ✅ |
| POST | `/api/v1/users/me/avatar` | 上传头像 | 是 | ✅ |
| POST | `/api/v1/users/me/bind-parent` | 绑定家长 | 是 | P1 |

### 2.1 注册

**请求：** `POST /api/v1/auth/register`
```json
{
  "name": "张三",
  "phone": "13800138000",
  "email": "zhangsan@example.com",
  "password": "Abc12345!",
  "role": "student"
}
```

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": "a1b2c3d4-...",
      "name": "张三",
      "phone": "138****8000",
      "email": "z****n@example.com",
      "role": "student",
      "avatar_url": null,
      "created_at": "2026-07-01T10:00:00+08:00"
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "dGhpcyBpcy..."
  }
}
```

### 2.2 登录

**请求：** `POST /api/v1/auth/login`
```json
{
  "account": "13800138000",
  "password": "Abc12345!"
}
```
> `account` 支持手机号或邮箱。

**响应：** `200`（同注册响应，含 access_token + refresh_token）

### 2.3 刷新 Token

**请求：** `POST /api/v1/auth/refresh`
```json
{
  "refresh_token": "dGhpcyBpcy..."
}
```

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "bmV3IHJlZnJl..."
  }
}
```

### 2.4 获取用户信息

**请求：** `GET /api/v1/users/me`
**Header：** `Authorization: Bearer eyJhbGciOi...`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "a1b2c3d4-...",
    "name": "张三",
    "phone": "138****8000",
    "email": "z****n@example.com",
    "role": "student",
    "avatar_url": "https://oss.example.com/avatars/a1b2.jpg",
    "created_at": "2026-07-01T10:00:00+08:00"
  }
}
```

### 2.5 更新个人信息

**请求：** `PUT /api/v1/users/me`
```json
{
  "name": "张三（更新）",
  "email": "newemail@example.com"
}
```

### 2.6 修改密码

**请求：** `PUT /api/v1/users/me/password`
```json
{
  "old_password": "Abc12345!",
  "new_password": "NewPass678!"
}
```

### 2.7 上传头像

**请求：** `POST /api/v1/users/me/avatar`
- Content-Type: `multipart/form-data`
- 字段：`file`（支持 jpg/png/gif，≤ 5MB）

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "avatar_url": "https://oss.example.com/avatars/a1b2.jpg"
  }
}
```

---

## 3. 学习规划模块 API

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| POST | `/api/v1/plans` | 创建学习计划 | 是 | ✅ |
| GET | `/api/v1/plans` | 获取计划列表 | 是 | ✅ |
| GET | `/api/v1/plans/:id` | 获取计划详情 | 是 | ✅ |
| PUT | `/api/v1/plans/:id` | 更新计划 | 是 | ✅ |
| POST | `/api/v1/plans/ai-generate` | AI 生成学习计划 | 是 | ✅ |
| POST | `/api/v1/plans/:id/checkin` | 打卡签到 | 是 | ✅ |

### 3.1 创建学习计划

**请求：** `POST /api/v1/plans`
```json
{
  "goal": "完成期末数学考试复习",
  "start_date": "2026-07-01",
  "end_date": "2026-07-30",
  "items": [
    {
      "date": "2026-07-01",
      "knowledge_point_ids": [1, 3],
      "duration_min": 60
    }
  ]
}
```

### 3.2 AI 生成学习计划

**请求：** `POST /api/v1/plans/ai-generate`
```json
{
  "goal": "完成期末数学考试复习",
  "start_date": "2026-07-01",
  "end_date": "2026-07-30",
  "daily_duration_min": 60
}
```

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "plan": { ... },
    "ai_used": true
  }
}
```
> LLM 不可用时，`ai_used: false`，使用规则引擎生成降级计划。

### 3.3 打卡签到

**请求：** `POST /api/v1/plans/:id/checkin`
```json
{
  "date": "2026-07-01",
  "duration_min": 60,
  "status": "done",
  "memo": "今天完成了一元二次方程练习"
}
```

**响应：** `201`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "record_id": 1001,
    "date": "2026-07-01",
    "status": "done"
  }
}
```

---

## 4. 练习模块 API

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| GET | `/api/v1/exercises` | 获取练习题目列表 | 是 | ✅ |
| GET | `/api/v1/exercises/random` | 获取随机题目 | 是 | ✅ |
| POST | `/api/v1/exercises/submit` | 提交练习答案 | 是 | ✅ |
| GET | `/api/v1/exercises/recommend` | 获取智能推荐题目 | 是 | ✅ |
| GET | `/api/v1/exercises/knowledge-points/:kp_id` | 按知识点获取题目 | 是 | ✅ |
| GET | `/api/v1/exercises/history` | 获取练习历史 | 是 | ✅ |

### 4.1 获取练习题目

**请求：** `GET /api/v1/exercises?subject_id=1&knowledge_point_id=3&difficulty=2&page=1&page_size=10`
**Header：** `Authorization: Bearer ...`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 101,
        "type": "choice",
        "difficulty": 2,
        "content": {
          "text": "一元二次方程 x² - 5x + 6 = 0 的解是？"
        },
        "options": [
          {"label": "A", "content": "x=2, x=3"},
          {"label": "B", "content": "x=-2, x=-3"},
          {"label": "C", "content": "x=1, x=6"},
          {"label": "D", "content": "x=-1, x=-6"}
        ],
        "knowledge_point": {
          "id": 3,
          "name": "一元二次方程"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```
> **注意：** 接口不返回答案和解析，防止学生直接查看。

### 4.2 提交练习答案

**请求：** `POST /api/v1/exercises/submit`
```json
{
  "question_id": 101,
  "answer": "A",
  "duration_sec": 45
}
```

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_correct": true,
    "correct_answer": "A",
    "analysis": "将方程因式分解为 (x-2)(x-3)=0，所以解为 x=2 或 x=3...",
    "mistake_recorded": false
  }
}
```
> `mistake_recorded`: 答错时自动收录错题。

### 4.3 智能推荐题目

**请求：** `GET /api/v1/exercises/recommend?count=10`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "recommend_strategy": "weak_point",
    "ai_used": true
  }
}
```

---

## 5. 错题模块 API

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| GET | `/api/v1/mistakes` | 获取错题列表 | 是 | ✅ |
| GET | `/api/v1/mistakes/groups` | 按知识点分组 | 是 | ✅ |
| PUT | `/api/v1/mistakes/:id/master` | 标记已掌握 | 是 | ✅ |
| POST | `/api/v1/mistakes/review` | 生成错题重练 | 是 | ✅ |
| DELETE | `/api/v1/mistakes/:id` | 删除错题记录 | 是 | ✅ |

### 5.1 获取错题列表

**请求：** `GET /api/v1/mistakes?knowledge_point_id=3&mastered=false&page=1&page_size=20`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 501,
        "question": {
          "id": 101,
          "type": "choice",
          "content": {...},
          "options": [...],
          "analysis": "..."
        },
        "wrong_answer": "B",
        "mistake_count": 3,
        "mastered": false,
        "knowledge_point": {
          "id": 3,
          "name": "一元二次方程"
        },
        "last_reviewed_at": "2026-07-01T15:00:00+08:00",
        "created_at": "2026-06-28T10:00:00+08:00"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20
  }
}
```

### 5.2 按知识点分组

**请求：** `GET /api/v1/mistakes/groups`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "groups": [
      {
        "knowledge_point": {"id": 3, "name": "一元二次方程"},
        "total": 5,
        "mastered": 2,
        "unmastered": 3
      }
    ]
  }
}
```

### 5.3 标记已掌握

**请求：** `PUT /api/v1/mistakes/:id/master`
```json
{
  "mastered": true
}
```

### 5.4 生成错题重练

**请求：** `POST /api/v1/mistakes/review`
```json
{
  "knowledge_point_ids": [3, 5],
  "count": 10
}
```

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "questions": [...],
    "total_unmastered": 8,
    "strategy": "unmastered_only"
  }
}
```
> 只返回未掌握的错题对应题目。

---

## 6. 数据看板模块 API

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| GET | `/api/v1/reports/summary` | 获取学习概览 | 是 | ✅ |
| GET | `/api/v1/reports/detail` | 获取学习报告详情 | 是 | ✅ |
| GET | `/api/v1/reports/mastery` | 获取知识点掌握度 | 是 | ✅ |
| GET | `/api/v1/reports/trend` | 获取学习趋势 | 是 | ✅ |

### 6.1 学习概览

**请求：** `GET /api/v1/reports/summary`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "today_duration_min": 45,
    "total_duration_min": 3600,
    "total_exercises": 250,
    "overall_correct_rate": 0.78,
    "streak_days": 5,
    "active_plan_count": 2,
    "unmastered_mistakes": 12
  }
}
```

### 6.2 学习报告详情

**请求：** `GET /api/v1/reports/detail?period_type=weekly&period_start=2026-06-24`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "type": "weekly",
      "start": "2026-06-24",
      "end": "2026-06-30"
    },
    "total_duration_min": 420,
    "total_exercises": 120,
    "correct_rate": 0.82,
    "mastery_by_kp": [
      {"kp_id": 3, "name": "一元二次方程", "rate": 0.85},
      {"kp_id": 5, "name": "二次函数", "rate": 0.72}
    ],
    "streak_days": 5,
    "daily_stats": [
      {"date": "2026-06-24", "duration_min": 60, "exercises": 20, "correct_rate": 0.75}
    ]
  }
}
```

### 6.3 知识点掌握度

**请求：** `GET /api/v1/reports/mastery?subject_id=1`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "subject": {"id": 1, "name": "数学"},
    "knowledge_points": [
      {"id": 1, "name": "实数", "level": 1, "rate": 0.95, "status": "mastered"},
      {"id": 3, "name": "一元二次方程", "level": 2, "rate": 0.78, "status": "learning"},
      {"id": 5, "name": "二次函数", "level": 2, "rate": 0.45, "status": "weak"}
    ]
  }
}
```

---

## 7. 内容管理模块 API（V1.5+）

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| GET | `/api/v1/subjects` | 获取科目列表 | 是 | ✅（基础） |
| GET | `/api/v1/knowledge-points` | 获取知识点树 | 是 | ✅（基础） |
| POST | `/api/v1/questions` | 创建题目 | 是（教师/管理员）| P1 |
| PUT | `/api/v1/questions/:id` | 更新题目 | 是（教师/管理员）| P1 |
| DELETE | `/api/v1/questions/:id` | 删除题目 | 是（教师/管理员）| P1 |
| POST | `/api/v1/questions/batch-import` | 批量导入题目 | 是（教师/管理员）| P1 |
| GET | `/api/v1/questions/export` | 导出题目 | 是（教师/管理员）| P1 |

### 7.1 获取科目列表

**请求：** `GET /api/v1/subjects`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {"id": 1, "name": "数学", "description": "初中数学"},
      {"id": 2, "name": "英语", "description": "初中英语"}
    ]
  }
}
```

### 7.2 获取知识点树

**请求：** `GET /api/v1/knowledge-points?subject_id=1`

**响应：** `200`
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "subject": {"id": 1, "name": "数学"},
    "tree": [
      {
        "id": 1, "name": "数与代数", "level": 1, "path": "1",
        "children": [
          {"id": 3, "name": "一元二次方程", "level": 2, "path": "1/3", "children": []}
        ]
      }
    ]
  }
}
```

---

## 8. 知识图谱模块 API（V1.5+）

| 方法 | 路径 | 描述 | 认证 | MVP |
|:-----|:-----|:-----|:-----|:----|
| GET | `/api/v1/knowledge-graph` | 获取知识图谱数据 | 是 | P1 |
| GET | `/api/v1/knowledge-graph/weak-points` | 获取薄弱点推荐 | 是 | P1 |

---

## 9. 认证鉴权流程

### 9.1 JWT Token 格式

```
Header: {"alg": "HS256", "typ": "JWT"}
Payload: {
  "user_id": "a1b2c3d4-...",
  "role": "student",
  "exp": 1719820800,
  "iat": 1719813600
}
Signature: HMAC-SHA256(base64(Header) + "." + base64(Payload), secret)
```

### 9.2 登录 → 认证 → 请求时序

```
客户端                    服务端                      数据库
  │                         │                          │
  │   POST /auth/login      │                          │
  │────────────────────────>│                          │
  │                         │  查询用户 + 验证密码      │
  │                         │─────────────────────────>│
  │                         │<─────────────────────────│
  │                         │  生成 access_token(2h)   │
  │                         │  生成 refresh_token(7d)  │
  │<────────────────────────│                          │
  │                         │                          │
  │  GET /api/v1/...        │                          │
  │  Authorization: Bearer  │                          │
  │────────────────────────>│                          │
  │                         │  JWT Middleware 校验      │
  │                         │  解析 Token → 提取 Claims│
  │                         │  校验角色权限             │
  │                         │  处理请求                 │
  │<────────────────────────│                          │
  │                         │                          │
  │  401 (Token 过期)       │                          │
  │<────────────────────────│                          │
  │                         │                          │
  │  POST /auth/refresh     │                          │
  │────────────────────────>│                          │
  │                         │  校验 refresh_token      │
  │                         │  生成新 token 对          │
  │<────────────────────────│                          │
```

### 9.3 角色权限校验

**中间件校验流程：**

```
请求到达 → JWT Middleware
  ├─ 无 Token → 返回 401
  ├─ Token 无效/过期 → 返回 401/20002
  └─ Token 有效 → 解析 Claims
      ├─ 提取 user_id, role
      └─ Role Middleware 校验
          ├─ 有权限 → 进入 Handler（user_id 注入 Context）
          └─ 无权限 → 返回 403/20004
```

### 9.4 角色-API 权限矩阵

| API | 学生 | 教师 | 管理员 | 家长 |
|:----|:-----|:-----|:-------|:-----|
| 用户注册/登录 | ✅ | ✅ | ✅ | ✅ |
| 个人信息管理 | ✅ | ✅ | ✅ | ✅ |
| 学习计划 CRUD | ✅ | ❌ | ❌ | ❌ |
| 练习与提交 | ✅ | ❌ | ❌ | ❌ |
| 错题本 | ✅ | ❌ | ❌ | 仅查看绑定学生 |
| 学习报告（自己） | ✅ | ❌ | ❌ | 仅查看绑定学生 |
| 班级学情 | ❌ | ✅ | ✅ | ❌ |
| 题库管理 | ❌ | ✅ | ✅ | ❌ |
| 知识点管理 | ❌ | ✅ | ✅ | ❌ |
| 系统配置 | ❌ | ❌ | ✅ | ❌ |
| 用户管理 | ❌ | ❌ | ✅ | ❌ |

---

## 10. 关键接口时序图

### 10.1 练习 → 批改 → 错题收录

```
学生                    练习模块                  错题模块                 数据库
  │                        │                        │                      │
  │  提交答案              │                        │                      │
  │───────────────────────>│                        │                      │
  │                        │  批改（匹配正确答案）   │                      │
  │                        │─────────────────────────────────────────────>│
  │                        │  ┌─ 答错 ─────────────┐                      │
  │                        │  │  写入练习记录       │                      │
  │                        │  │  通知错题模块       │                      │
  │                        │──┴────────────────────────>│                  │
  │                        │                        │  幂等检查            │
  │                        │                        │─────────────────────>│
  │                        │                        │  ┌─ 已存在：增加次数 │
  │                        │                        │  └─ 新记录：创建     │
  │                        │                        │<─────────────────────│
  │ 返回结果 + 解析       │                        │                      │
  │<───────────────────────│                        │                      │
```

### 10.2 学习报告生成

```
学生                    看板模块                   练习模块              数据库
  │                        │                        │                    │
  │  请求报告              │                        │                    │
  │───────────────────────>│                        │                    │
  │                        │  查 Redis 缓存         │                    │
  │                        │──────────────────────────────────────────>│
  │                        │  ┌─ 命中 → 返回缓存    │                    │
  │                        │  └─ 未命中 → 查询数据  │                    │
  │                        │                        │                    │
  │                        │─── 查询练习统计 ──────>│                    │
  │                        │<────────────────────────│                    │
  │                        │                        │                    │
  │                        │─── 查询计划完成 ──────>│                    │
  │                        │<────────────────────────│                    │
  │                        │                        │                    │
  │                        │  聚合数据 → 生成报告    │                    │
  │                        │  写入 Redis 缓存       │                    │
  │                        │  返回结果               │                    │
  │<───────────────────────│                        │                    │
```

---

## 11. 分页与过滤规范

### 11.1 分页参数

| 参数 | 类型 | 默认值 | 范围 | 说明 |
|:-----|:-----|:-------|:-----|:------|
| page | int | 1 | ≥ 1 | 页码 |
| page_size | int | 20 | 1~100 | 每页数量 |

### 11.2 通用过滤参数

| 参数 | 说明 | 示例 |
|:-----|:------|:------|
| `subject_id` | 按科目过滤 | `?subject_id=1` |
| `knowledge_point_id` | 按知识点过滤 | `?knowledge_point_id=3` |
| `difficulty` | 按难度过滤 | `?difficulty=2` |
| `type` | 按题目类型过滤 | `?type=choice` |
| `mastered` | 按掌握状态过滤 | `?mastered=false` |
| `start_date` / `end_date` | 按时间范围过滤 | `?start_date=2026-07-01&end_date=2026-07-30` |
| `sort_by` / `order` | 排序字段和方向 | `?sort_by=created_at&order=desc` |