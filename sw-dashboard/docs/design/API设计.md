# API 设计 — 翠花集群管理面板 (sw-dashboard)

> 版本：v1.0.0  
> 日期：2026-06-03  
> 状态：设计定稿  
> 对应阶段：Stage 3 — 系统设计

---

## 1. 设计规范

### 1.1 通用约定

| 约定 | 规则 |
|:-----|:------|
| **基础路径** | 所有 API 以 `/api` 为前缀 |
| **协议** | HTTP (Tailscale 内网，无需 HTTPS) |
| **请求体格式** | `application/json` |
| **响应格式** | JSON，统一包装 `{"success": true, "data": ..., "error": ...}` |
| **字符编码** | UTF-8 |
| **时间格式** | ISO 8601 (`2026-06-03T14:30:00+08:00`) |
| **分页** | 暂不需要（数据量小），预留 `offset`/`limit` 参数 |

### 1.2 认证方式

**Token 认证**：所有 API（除 `/api/health` 外）需要在 Header 中携带 `X-Token`。

```http
GET /api/dashboard
X-Token: your-token-here
```

**Token 验证逻辑**：后端读取 `config.py` 中配置的 Token 值（从环境变量或 `.env` 读取），与请求 Header 比对。匹配则放行，不匹配返回 401。

### 1.3 HTTP 状态码

| 状态码 | 含义 | 适用场景 |
|:-------|:-----|:---------|
| 200 | 请求成功 | 数据返回、写入成功 |
| 400 | 请求参数错误 | YAML 格式错误、参数缺失 |
| 401 | 未认证 | Token 缺失或无效 |
| 403 | 无权限 | 路径穿越尝试、操作被禁止 |
| 404 | 资源不存在 | 文件不存在、日志文件不存在 |
| 413 | 请求实体过大 | 文件超过大小限制 |
| 500 | 服务器内部错误 | 文件写入失败、系统命令执行异常 |

---

## 2. 接口列表总览

| 方法 | 路径 | 功能 | 模块 | 优先级 |
|:-----|:-----|:-----|:-----|:-------|
| `GET` | `/api/health` | 健康检查 | 系统 | P0 |
| `POST` | `/api/auth` | Token 认证 | 认证 | P0 |
| `GET` | `/api/dashboard` | 仪表盘聚合数据 | 仪表盘 | P0 |
| `GET` | `/api/storage/tree` | 三库目录树 | 存储 | P0 |
| `GET` | `/api/storage/file` | 获取文件内容 | 存储 | P0 |
| `PUT` | `/api/storage/file` | 保存文件内容 | 存储 | P1 |
| `GET` | `/api/configs` | 配置列表 | 配置 | P0 |
| `GET` | `/api/configs/{name}` | 读取配置文件 | 配置 | P0 |
| `PUT` | `/api/configs/{name}` | 保存配置文件 | 配置 | P0 |
| `GET` | `/api/configs/cron` | Cron 任务列表 | 配置 | P1 |
| `POST` | `/api/configs/cron/toggle` | 启停 Cron 任务 | 配置 | P1 |
| `GET` | `/api/logs/files` | 日志文件列表 | 日志 | P1 |
| `GET` | `/api/logs/tail` | 日志尾部读取 | 日志 | P1 |
| `GET` | `/api/logs/search` | 日志搜索 | 日志 | P1 |
| `GET` | `/api/logs/archive` | 历史归档列表 | 日志 | P1 |
| `GET` | `/api/skills/tree` | 技能目录树 | 技能 | P2 |
| `GET` | `/api/skills/file` | 技能文件预览 | 技能 | P2 |

---

## 3. 接口详细设计

### 3.1 健康检查

```
GET /api/health
```

**认证**：❌ 免认证

**请求参数**：无

**响应示例**：
```json
{
  "status": "ok",
  "version": "0.1.0",
  "python_version": "3.12.0",
  "uptime": "2 hours, 15 minutes",
  "dependencies": {
    "fastapi": "0.136.3",
    "uvicorn": "0.34.0",
    "pyyaml": "6.0.3",
    "aiofiles": "25.1.0"
  }
}
```

---

### 3.2 用户认证

```
POST /api/auth
```

**认证**：❌ 免认证（登录接口）

**请求体**：
```json
{
  "token": "your-token-here"
}
```

**成功响应 (200)**：
```json
{
  "success": true,
  "data": {
    "token": "your-token-here",
    "expires_in": null
  }
}
```

**失败响应 (401)**：
```json
{
  "detail": "Invalid token"
}
```

---

### 3.3 仪表盘数据聚合

```
GET /api/dashboard
```

**认证**：✅ 需要 Token

**请求参数**：无

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "gateway": {
      "status": "active",
      "uptime": "2 days, 3 hours, 15 minutes",
      "pid": 684033,
      "memory_mb": 156.7,
      "cpu_percent": 2.3
    },
    "sessions": {
      "total_sessions": 12,
      "total_messages": 87,
      "active_sessions": 1,
      "avg_duration_sec": 184.5
    },
    "pipeline": {
      "current_stage": 3,
      "total_stages": 8,
      "overall_progress": 0.375,
      "stages": [
        {"stage_id": 1, "name": "需求分析", "status": "completed", "progress": 1.0},
        {"stage_id": 2, "name": "技术调研", "status": "completed", "progress": 1.0},
        {"stage_id": 3, "name": "系统设计", "status": "running", "progress": 0.5},
        {"stage_id": 4, "name": "UI/UX 设计", "status": "pending", "progress": 0.0},
        {"stage_id": 5, "name": "编码", "status": "pending", "progress": 0.0},
        {"stage_id": 6, "name": "审查", "status": "pending", "progress": 0.0},
        {"stage_id": 7, "name": "测试", "status": "pending", "progress": 0.0},
        {"stage_id": 8, "name": "部署", "status": "pending", "progress": 0.0}
      ]
    },
    "cron_jobs": [
      {
        "name": "午间快报",
        "schedule": "0 12 * * *",
        "last_run": "2026-06-03T12:00:00+08:00",
        "last_status": "success",
        "enabled": true
      }
    ],
    "server_time": "2026-06-03T14:30:00+08:00",
    "version": "0.1.0"
  }
}
```

---

### 3.4 存储目录树

```
GET /api/storage/tree
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:-----|
| `repo` | string | 否 | `all` | 仓库名：`workspace` / `skills-library` / `projects` / `all` |
| `path` | string | 否 | `/` | 子路径，用于懒加载子目录 |
| `max_depth` | int | 否 | `5` | 最大递归深度 |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": [
    {
      "name": "sw-agents-workspace",
      "path": "/home/ubuntu/gitee-software/sw-agents-workspace",
      "type": "directory",
      "isLeaf": false,
      "children": [
        {
          "name": "stage-1-analyst",
          "path": "/home/ubuntu/gitee-software/sw-agents-workspace/stage-1-analyst",
          "type": "directory",
          "isLeaf": false,
          "children": [
            {
              "name": "SKILL.md",
              "path": "/home/ubuntu/gitee-software/sw-agents-workspace/stage-1-analyst/SKILL.md",
              "type": "file",
              "isLeaf": true,
              "size": 1234,
              "mtime": "2026-06-01T10:00:00+08:00"
            }
          ]
        }
      ]
    }
  ]
}
```

---

### 3.5 获取文件内容

```
GET /api/storage/file
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:-----|
| `path` | string | 是 | — | 文件绝对路径 |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "path": "/home/ubuntu/gitee-software/sw-projects/sw-dashboard/README.md",
    "name": "README.md",
    "content": "# sw-dashboard\n\n翠花集群管理面板...",
    "size": 2048,
    "mtime": "2026-06-03T10:00:00+08:00",
    "language": "markdown",
    "editable": true
  }
}
```

**错误响应 (404)**：
```json
{
  "detail": "File not found: /path/to/nonexistent"
}
```

---

### 3.6 保存文件内容

```
PUT /api/storage/file
```

**认证**：✅ 需要 Token

**请求体**：
```json
{
  "path": "/home/ubuntu/gitee-software/sw-projects/sw-dashboard/README.md",
  "content": "# sw-dashboard\n\n新内容...",
  "create_backup": true
}
```

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "path": "/home/ubuntu/gitee-software/sw-projects/sw-dashboard/README.md",
    "size": 1536,
    "backup_path": "/home/ubuntu/gitee-software/sw-projects/sw-dashboard/README.md.bak",
    "timestamp": "2026-06-03T14:35:00+08:00"
  }
}
```

**错误响应 (403)**：
```json
{
  "detail": "Path traversal detected"
}
```

**错误响应 (413)**：
```json
{
  "detail": "File too large (max 10MB)"
}
```

---

### 3.7 获取配置列表

```
GET /api/configs
```

**认证**：✅ 需要 Token

**请求参数**：无

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "configs": [
      {"name": "config.yaml", "path": "...", "language": "yaml", "editable": true, "last_modified": "..."},
      {"name": ".env", "path": "...", "language": "env", "editable": true, "last_modified": "..."},
      {"name": "CLAUDE.md", "path": "...", "language": "markdown", "editable": true, "last_modified": "..."},
      {"name": "SOUL.md", "path": "...", "language": "markdown", "editable": true, "last_modified": "..."}
    ]
  }
}
```

---

### 3.8 读取配置文件

```
GET /api/configs/{name}
```

**认证**：✅ 需要 Token

**路径参数**：

| 参数 | 类型 | 说明 |
|:-----|:-----|:------|
| `name` | string | 配置文件名：`config.yaml` / `.env` / `CLAUDE.md` / `SOUL.md` |

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| `masked` | bool | 否 | `true` | 是否掩码敏感信息（仅 .env 有效） |

**响应示例 (200)** — `config.yaml`：
```json
{
  "success": true,
  "data": {
    "name": "config.yaml",
    "path": "/home/ubuntu/.hermes/config.yaml",
    "content": "api_key: sk-...\nmodel: deepseek-v4-flash",
    "language": "yaml",
    "editable": true,
    "last_modified": "2026-06-02T08:00:00+08:00"
  }
}
```

**响应示例 (200)** — `.env`（掩码模式）：
```json
{
  "success": true,
  "data": {
    "name": ".env",
    "path": "/home/ubuntu/.hermes/.env",
    "content": "DEEPSEEK_API_KEY=deep*********ash\nMINIMAX_API_KEY=mini*******max",
    "original_content": null,
    "language": "env",
    "editable": true,
    "last_modified": "2026-06-01T12:00:00+08:00"
  }
}
```

**说明**：`.env` 默认返回掩码版本。仅当 `masked=false` 参数传递且用户确认后，才返回原始内容。

---

### 3.9 保存配置文件

```
PUT /api/configs/{name}
```

**认证**：✅ 需要 Token

**路径参数**：

| 参数 | 类型 | 说明 |
|:-----|:-----|:------|
| `name` | string | 配置文件名 |

**请求体**：
```json
{
  "content": "api_key: sk-new-key\nmodel: deepseek-v4-flash",
  "create_backup": true,
  "validate": true
}
```

**成功响应 (200)**：
```json
{
  "success": true,
  "data": {
    "path": "/home/ubuntu/.hermes/config.yaml",
    "size": 64,
    "backup_path": "/home/ubuntu/.hermes/config.yaml.bak",
    "timestamp": "2026-06-03T14:40:00+08:00"
  }
}
```

**错误响应 (400)** — YAML 格式错误：
```json
{
  "detail": "YAML parse error: mapping values are not allowed here"
}
```

---

### 3.10 Cron 任务列表

```
GET /api/configs/cron
```

**认证**：✅ 需要 Token

**请求参数**：无

**响应示例 (200)**：
```json
{
  "success": true,
  "data": [
    {
      "raw": "0 12 * * * /home/ubuntu/scripts/daily-log.py --period=midday",
      "schedule": "0 12 * * *",
      "command": "/home/ubuntu/scripts/daily-log.py --period=midday",
      "comment": null,
      "enabled": true,
      "last_run": "2026-06-03T12:00:00+08:00",
      "last_output": "✅ 午间快报完成"
    },
    {
      "raw": "#0 17 * * * /home/ubuntu/scripts/daily-log.py --period=afternoon",
      "schedule": "0 17 * * *",
      "command": "/home/ubuntu/scripts/daily-log.py --period=afternoon",
      "comment": null,
      "enabled": false,
      "last_run": null,
      "last_output": null
    }
  ]
}
```

---

### 3.11 启停 Cron 任务

```
POST /api/configs/cron/toggle
```

**认证**：✅ 需要 Token

**请求体**：
```json
{
  "command": "/home/ubuntu/scripts/daily-log.py --period=afternoon",
  "enable": true
}
```

**成功响应 (200)**：
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Cron job enabled: /home/ubuntu/scripts/daily-log.py --period=afternoon"
  }
}
```

---

### 3.12 日志文件列表

```
GET /api/logs/files
```

**认证**：✅ 需要 Token

**请求参数**：无

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "files": [
      {"name": "hermes.log", "path": "...", "size": 1048576, "mtime": "2026-06-03T14:00:00+08:00"},
      {"name": "hermes.log.1", "path": "...", "size": 2097152, "mtime": "2026-06-02T23:59:00+08:00"}
    ],
    "log_path": "/home/ubuntu/.hermes/hermes-agent/venv/var/log",
    "default_file": "hermes.log"
  }
}
```

---

### 3.13 日志尾部读取

```
GET /api/logs/tail
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| `file` | string | 否 | `hermes.log` | 日志文件名 |
| `lines` | int | 否 | `100` | 返回行数 (1-2000) |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "file": "hermes.log",
    "content": "[2026-06-03 14:29:00] INFO: Session started\n[2026-06-03 14:30:00] INFO: Message processed\n",
    "lines": 100,
    "total_size": 1048576,
    "truncated": true
  }
}
```

---

### 3.14 日志搜索

```
GET /api/logs/search
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| `keyword` | string | 是 | — | 搜索关键字 |
| `file` | string | 否 | `hermes.log` | 目标日志文件 |
| `max_lines` | int | 否 | `50` | 最大返回行数 |
| `case_sensitive` | bool | 否 | `false` | 是否大小写敏感 |
| `context_lines` | int | 否 | `0` | 上下文行数 |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "keyword": "ERROR",
    "file": "hermes.log",
    "total_matches": 3,
    "matches": [
      "[2026-06-03 10:15:00] ERROR: Connection timeout",
      "[2026-06-03 11:30:00] ERROR: API rate limit exceeded",
      "[2026-06-03 13:45:00] ERROR: File not found"
    ],
    "truncated": false,
    "time_cost_ms": 12.5
  }
}
```

---

### 3.15 日志归档列表

```
GET /api/logs/archive
```

**认证**：✅ 需要 Token

**请求参数**：无

**响应示例 (200)**：
```json
{
  "success": true,
  "data": [
    {
      "name": "hermes.log.1.gz",
      "path": "...",
      "size": 512000,
      "mtime": "2026-06-02T23:59:00+08:00",
      "lines": null
    }
  ]
}
```

---

### 3.16 技能目录树

```
GET /api/skills/tree
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| `path` | string | 否 | `/` | 子路径（懒加载） |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": [
    {
      "name": "analyst",
      "path": "/home/ubuntu/.hermes/skills/analyst",
      "type": "directory",
      "isLeaf": false,
      "has_skill_md": true,
      "children": [
        {
          "name": "SKILL.md",
          "path": "/home/ubuntu/.hermes/skills/analyst/SKILL.md",
          "type": "file",
          "isLeaf": true,
          "has_skill_md": false
        }
      ]
    }
  ]
}
```

---

### 3.17 技能文件预览

```
GET /api/skills/file
```

**认证**：✅ 需要 Token

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| `path` | string | 是 | — | 技能文件绝对路径 |

**响应示例 (200)**：
```json
{
  "success": true,
  "data": {
    "path": "/home/ubuntu/.hermes/skills/analyst/SKILL.md",
    "name": "SKILL.md",
    "content": "# Analyst Skill\n\n## 职责\n...",
    "language": "markdown",
    "size": 3072,
    "mtime": "2026-06-01T10:00:00+08:00"
  }
}
```

---

## 4. 错误处理

### 4.1 统一错误响应格式

所有错误响应遵循 FastAPI 默认格式：

```json
{
  "detail": "人类可读的错误描述"
}
```

### 4.2 错误码映射

| 场景 | HTTP 状态码 | detail 示例 |
|:-----|:-----------|:-----------|
| Token 缺失 | 401 | `"Missing X-Token header"` |
| Token 无效 | 401 | `"Invalid token"` |
| 路径穿越 | 403 | `"Path traversal detected: path is outside allowed base"` |
| 文件不存在 | 404 | `"File not found: /path/to/file"` |
| 文件过大 | 413 | `"File too large (max 10485760 bytes)"` |
| 格式错误 | 400 | `"YAML parse error at line 3: mapping values are not allowed here"` |
| 写入失败 | 500 | `"Failed to write file: Permission denied"` |
| 命令超时 | 500 | `"System command timed out after 5s"` |

---

## 5. 前端 API 调用示例

### 5.1 axios 实例配置

```typescript
// src/api/index.ts
import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' }
})

// 请求拦截器：自动注入 Token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('dashboard_token')
  if (token) {
    config.headers['X-Token'] = token
  }
  return config
})

// 响应拦截器：统一错误处理
api.interceptors.response.use(
  response => response,
  error => {
    const msg = error.response?.data?.detail || error.message
    ElMessage.error(msg)
    if (error.response?.status === 401) {
      // 跳转到登录页
      window.location.hash = '#/login'
    }
    return Promise.reject(error)
  }
)

export default api
```

### 5.2 各模块 API 调用

```typescript
// src/api/dashboard.ts
import api from './index'

export const getDashboard = () => api.get('/dashboard')
export const getHealth = () => api.get('/health')
export const verifyToken = (token: string) => api.post('/auth', { token })
```

```typescript
// src/api/storage.ts
export const getTree = (repo?: string, path?: string) =>
  api.get('/storage/tree', { params: { repo, path } })
export const getFile = (path: string) =>
  api.get('/storage/file', { params: { path } })
export const saveFile = (path: string, content: string) =>
  api.put('/storage/file', { path, content })
```

```typescript
// src/api/configs.ts
export const getConfigs = () => api.get('/configs')
export const getConfig = (name: string, masked?: boolean) =>
  api.get(`/configs/${name}`, { params: { masked } })
export const saveConfig = (name: string, content: string) =>
  api.put(`/configs/${name}`, { content })
export const getCronJobs = () => api.get('/configs/cron')
export const toggleCronJob = (command: string, enable: boolean) =>
  api.post('/configs/cron/toggle', { command, enable })
```

```typescript
// src/api/logs.ts
export const getLogFiles = () => api.get('/logs/files')
export const tailLog = (file?: string, lines?: number) =>
  api.get('/logs/tail', { params: { file, lines } })
export const searchLogs = (keyword: string, file?: string) =>
  api.get('/logs/search', { params: { keyword, file } })
export const getArchives = () => api.get('/logs/archive')
```

---

## 6. 接口安全要求

| 要求 | 实现方式 |
|:-----|:---------|
| 认证 | X-Token Header，所有非健康检查接口都需要 |
| 路径安全 | 后端 `validate_path()` 校验，防止目录遍历 |
| 输入校验 | Pydantic v2 Schema 严格校验请求体 |
| 大小限制 | 文件读取 ≤10MB，文件编辑 ≤1MB |
| 敏感信息 | .env 默认掩码返回 |
| 操作确认 | 前端二次确认弹窗（非 API 层面） |
| 速率限制 | 暂不需要（仅 1 用户，Tailscale 内网） |
