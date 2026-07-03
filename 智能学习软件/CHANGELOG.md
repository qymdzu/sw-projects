# 变更日志

## [1.0.0] - 2026-XX-XX

### Phase B — 编码开发（reset 模式）

#### 后端
- **新增**：模型配置模块
  - `POST /api/v1/settings/model` 保存/更新（自动设 default）
  - `GET /api/v1/settings/model` 读当前 default
  - `PUT /api/v1/settings/model` 切换 default
  - `DELETE /api/v1/settings/model?provider=` 删除配置
- **新增**：model_settings 表（user_id+provider 唯一，is_default 索引）
- **新增**：AES-256-GCM API Key 加密（密钥派生自 JWT_SECRET + model-setting-salt）
- **新增**：setting 模块单测（service 层 12 条 + handler 层 9 条 + crypto 4 条）

#### 前端（0 → 1 新建）
- Vue 3 + Vite + TypeScript + Pinia + Vue Router 4 + Axios + Element Plus
- 6 个核心页面：登录 / 仪表盘 / 科目浏览 / 智能选题 / 错题本 / 设置
- 模型配置表单（6 Provider：openai / anthropic / qwen / deepseek / ollama / custom）
- Axios 拦截器（401 自动 refresh，5xx Toast）
- Pinia store：auth / setting
- 单元测试：auth store + axios client + AppButton

#### 安全加固（P1 一并修复）
- **修复**：`cmd/server/main.go` 删除密码 DEBUG 输出（R-05）
- **P1-01**：JWT_SECRET 强校验（prod 环境禁用占位密钥，< 32 字符报错）
- **P1-02**：CORS 改为白名单（env `CORS_ALLOW_ORIGINS` 注入，默认 `localhost:5173,localhost:4173`）
- **P1-03**：错误响应统一收敛，handler 错误映射函数封装

## [0.1.0] - 2026-07-02

### 新增
- 项目初始化
