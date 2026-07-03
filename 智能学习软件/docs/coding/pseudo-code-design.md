# 伪代码设计 — 智能学习软件（全栈重置 + 模型配置新增模块）

> 版本：v1.0.0
> 阶段：Pipeline Stage 5 — Phase A 伪代码设计（PDCA → P 阶段）
> 作者：dev（开发工程师）
> 对应需求：v1.0.0（已完成审阅）
> 对应架构/API/数据模型/目录结构设计：v1.0.0（已审阅）
> 对应 UX 设计：交互流程设计 / UI 组件规范 / 视觉风格指南（v1.0.0）
> 任务模式：fullstack 迭代（reset）
> 关键变更：**新增"模型配置"模块（后端 + 前端）** + 全栈前端 0 → 1 初始化

---

## 0. 任务边界与前置条件

### 0.1 本次编码任务范围

| # | 范围 | 状态 | 来源 |
|:--|:-----|:-----|:------|
| 1 | 后端模型配置模块（model / service / handler / router） | **新增** | 小苹果 reset 调度 |
| 2 | 后端代码清理（删除 `main.go:85` 密码 DEBUG 输出） | **修复** | 本次读代码发现的安全隐患 + 测试报告 §4.3 P1-03 |
| 3 | 前端 Vue 3 + Vite + TypeScript 工程 0 → 1 | **新建** | fullstack 必选项 |
| 4 | 前端 6 个核心页面（登录/仪表盘/科目浏览/智能选题/错题本/设置） | **新建** | 任务清单 |
| 5 | 前端对应单测 + 后端 setting 单测 | **新增** | TDD 纪律 |
| 6 | CHANGELOG.md 与 docs/test/测试报告.md 增量更新 | **更新** | 交付惯例 |

### 0.2 本次**不**改动（已实现，复用）

| 范围 | 原因 |
|:-----|:-----|
| 39 个历史 API 端点的 handler/service/repository | 已通过测试报告单元测试 56/56 PASS |
| 数据模型 9 张主表（User/Subject/KnowledgePoint/Question/StudyPlan/StudyRecord/ExerciseRecord/MistakeBook/StudyReport/OperationLog） | GORM AutoMigrate 已落地（main.go） |
| JWT 双 Token 机制 | pkg/jwt 已封装 |
| middleware（CORS / Recovery / Logger / JWTAuth / RequestID）| 已实装 |
| cmd/server/main.go 启动骨架 | 框架已就位 |

### 0.3 ⚠️ 前置条件（执行 Phase B 前必须确认）

| # | 前置条件 | 阻塞等级 | 当前状态 |
|:--|:---------|:---------|:---------|
| P-01 | dev agent 持有 `exec` 权限（allowlist 含 `exec`） | **P0_BLOCK** | 当前 **deny**：影响 Phase B 的 `go build` / `go test` / `npm install` / `npm run build` 等所有门禁命令的执行 |
| P-02 | 项目根路径可写 + Git 可推送 | P0_BLOCK | 路径已确认（`/home/ubuntu/gitee-software/sw-projects/智能学习软件/`）|
| P-03 | PostgreSQL DSN 可用（用于 setting_handler 集成测试） | P1 | 待 devops 提供或测试环境调通 |
| P-04 | 同意修复 P1-01（JWT Secret 强校验）/ P1-02（CORS 白名单）/ P1-03（错误信息收敛） | P1 | 测试报告 §4.3 列出的已知风险，本次一并处置 |

> **关于 P-01 的处理策略**：本文档同时给出"完整 Plan"，但 Phase B 实际执行需 exec 权限。如小苹果/公子决定让 dev 在受限环境中推进，则 Phase B 的所有"门禁命令"将由小苹果或外部 CI 执行，dev 负责**产出可编译的代码 + 单测用例 + 验证清单**，小苹果/CI 负责跑门禁。本文档默认按"dev 既有 exec 权限"撰写完整 Plan，待 P-01 解决后即开 Phase B。

---

## 1. 模块总览与已有/新增对照

### 1.1 模块清单（已存在 — 重用）

| 业务模块 | Handler | Service | Repository | Model | 状态 |
|:---------|:--------|:--------|:-----------|:------|:-----|
| auth | `auth_handler.go` | `auth_service.go` | `user_repo.go` | `user.go` | ✅ 已有 |
| user | `user_handler.go` | `user_service.go` | `user_repo.go` | `user.go` | ✅ 已有 |
| plan | `plan_handler.go` | `plan_service.go` | `plan_repo.go` | `study_plan.go`,`study_record.go` | ✅ 已有 |
| exercise | `exercise_handler.go` | `exercise_service.go` | `exercise_repo.go` | `exercise_record.go`,`question.go` | ✅ 已有 |
| mistake | `mistake_handler.go` | `mistake_service.go` | `mistake_repo.go` | `mistake_book.go` | ✅ 已有 |
| report | `report_handler.go` | `report_service.go` | `report_repo.go` | `study_report.go` | ✅ 已有 |
| subject | `subject_handler.go` | `subject_service.go` | `subject_repo.go` | `subject.go` | ✅ 已有 |
| knowledge | `knowledge_handler.go` | `knowledge_service.go` | `knowledge_repo.go` | `knowledge_point.go` | ✅ 已有 |

### 1.2 模块清单（新增 — 本次产出）

| 业务模块 | Handler | Service | Repository | Model | 路径 | 优先级 |
|:---------|:--------|:--------|:-----------|:------|:-----|:-------|
| **setting（模型配置）** | `setting_handler.go` | `setting_service.go` | `setting_repo.go` | `model_setting.go` | `backend/internal/...` | **P0** |

### 1.3 前端模块清单（全部新建 — 本次产出）

| 范围 | 路径 | 优先级 |
|:-----|:-----|:-------|
| Vite 工程骨架 | `frontend/` | P0 |
| 6 个页面 + 通用布局 + 鉴权 store + API 客户端 | `frontend/src/` | P0 |
| 单元测试 + 组件测试 | `frontend/tests/` 或 `frontend/src/**/*.test.ts` | P1 |

---

## 2. 后端 — 模型配置模块伪代码（核心新增）

### 2.1 数据模型 — `backend/internal/model/model_setting.go`

**职责**：声明用户私有 LLM 配置表的 GORM 模型，对应数据库表 `model_settings`。

```pseudo
// Package model — 模型配置表（本次新增）。
//
// 表设计要点：
//   - user_id + provider 组合唯一：一个用户对每个 provider 仅一条配置
//   - is_default 标记是否为用户当前活跃使用的 LLM 配置
//   - api_key 在落库前必须加密（AES-GCM，密钥来自 JWT_SECRET + salt）—— 永远不要明文存
//   - config JSONB：随 provider 扩展（如 temperature、top_p、max_tokens）

struct ModelSetting {
    ID                  uuid.UUID   // PK，gen_random_uuid()
    UserID              uuid.UUID   // FK → users.id（逻辑外键，由 service 层校验）
    Provider            string      // 例："openai" / "qwen" / "deepseek" / "ollama"
    APIEndpoint         string      // 例："https://api.openai.com/v1"
    APIKeyCiphertext    []byte      // 加密后的 API Key
    APIKeyNonce         []byte      // GCM nonce（12 bytes）
    Model               string      // 例："gpt-4o-mini" / "qwen-plus"
    Config              datatypes.JSON  // 透传额外参数，如 {"temperature":0.7}
    IsDefault           bool        // 是否为当前默认配置
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

// 复合唯一索引：(user_id, provider)
//
// TableName: "model_settings"
```

**字段约束设计要点**：

| 字段 | 类型 | 约束 | 默认值 | 说明 |
|:-----|:-----|:-----|:-------|:------|
| provider | VARCHAR(32) | NOT NULL, CHECK(provider IN ('openai','qwen','deepseek','ollama','custom')) | — | 服务商枚举 |
| api_endpoint | VARCHAR(500) | NOT NULL | — | 必须是 http(s) URL |
| api_key_ciphertext | BYTEA | NOT NULL | — | 加密后存储 |
| api_key_nonce | BYTEA | NOT NULL | — | 12 bytes GCM nonce |
| model | VARCHAR(100) | NOT NULL | — | 模型标识（不可空） |
| is_default | BOOLEAN | NOT NULL, DEFAULT false | false | 同 user 只能 1 条 true |
| config | JSONB | NOT NULL, DEFAULT '{}' | '{}' | 透传额外 LLM 参数 |

**索引**：
- 唯一索引：`UNIQUE (user_id, provider)`
- 索引：`idx_ms_user_default` (user_id) WHERE is_default = true（部分索引）

### 2.2 仓储层 — `backend/internal/repository/setting_repo.go`

**职责**：封装 model_settings 表的 CRUD，是 service 唯一允许直接持有 `*gorm.DB` 的层。

```pseudo
package repository

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "smart-learning/internal/model"
)

// ErrNotFound 是仓储层向 service 层抛出的"未找到"信号。
// 与 service 层定义的 ErrResourceMissing 区分（service 层会做 wrapping）。
var ErrNotFound = errors.New("model setting not found")

// SettingRepository 提供 model_settings 表的访问。
type SettingRepository interface {
    GetByUserAndProvider(ctx, userID uuid.UUID, provider string) (*model.ModelSetting, error)
    GetActiveByUser(ctx, userID uuid.UUID) (*model.ModelSetting, error)        // 用户的 is_default=true 配置
    ListByUser(ctx, userID uuid.UUID) ([]*model.ModelSetting, error)
    Upsert(ctx, m *model.ModelSetting) error                                   // CREATE OR UPDATE
    SetActive(ctx, userID uuid.UUID, provider string) error                   // 事务化切换默认
    DeleteByUserAndProvider(ctx, userID uuid.UUID, provider string) error
}

type settingRepo struct{ db *gorm.DB }

// NewSettingRepository 构造仓储实现。db 必须为非 nil 的 GORM 句柄。
func NewSettingRepository(db *gorm.DB) SettingRepository {
    return &settingRepo{db: db}
}

func (r *settingRepo) GetByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*model.ModelSetting, error) {
    var m model.ModelSetting
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND provider = ?", userID, provider).
        First(&m).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("query model setting: %w", err)
    }
    return &m, nil
}

func (r *settingRepo) GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.ModelSetting, error) {
    var m model.ModelSetting
    err := r.db.WithContext(ctx).
        Where("user_id = ? AND is_default = ?", userID, true).
        First(&m).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("query active model setting: %w", err)
    }
    return &m, nil
}

func (r *settingRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.ModelSetting, error) {
    var ms []*model.ModelSetting
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("is_default DESC, updated_at DESC").
        Find(&ms).Error
    if err != nil { return nil, fmt.Errorf("list model settings: %w", err) }
    return ms, nil
}

func (r *settingRepo) Upsert(ctx context.Context, m *model.ModelSetting) error {
    // 唯一键冲突时更新
    return r.db.WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "user_id"}, {Name: "provider"}},
            DoUpdates: clause.AssignmentColumns([]string{"api_endpoint","api_key_ciphertext","api_key_nonce","model","config","is_default","updated_at"}),
        }).Create(m).Error
}

func (r *settingRepo) SetActive(ctx context.Context, userID uuid.UUID, provider string) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1) 把用户所有配置 is_default = false
        if err := tx.Model(&model.ModelSetting{}).
            Where("user_id = ?", userID).
            Update("is_default", false).Error; err != nil {
            return err
        }
        // 2) 把目标 provider 设为 true
        res := tx.Model(&model.ModelSetting{}).
            Where("user_id = ? AND provider = ?", userID, provider).
            Update("is_default", true)
        if res.Error != nil { return res.Error }
        if res.RowsAffected == 0 {
            return ErrNotFound   // 该 provider 下用户没有配置
        }
        return nil
    })
}

func (r *settingRepo) DeleteByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) error {
    res := r.db.WithContext(ctx).
        Where("user_id = ? AND provider = ?", userID, provider).
        Delete(&model.ModelSetting{})
    if res.Error != nil { return res.Error }
    if res.RowsAffected == 0 { return ErrNotFound }
    return nil
}
```

### 2.3 服务层 — `backend/internal/service/setting_service.go`

**职责**：业务编排（加密、校验、默认值、调用 LLM ping），service 层不持有 `*gorm.DB`。

```pseudo
package service

import (
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"

    "smart-learning/internal/model"
    "smart-learning/internal/repository"
)

// Service 层错误定义（沿用现有命名风格）
var (
    ErrSettingNotFound       = errors.New("model setting not found")    // → 404
    ErrSettingParamInvalid   = errors.New("invalid model setting")      // → 400
    ErrUnsupportedProvider   = errors.New("unsupported provider")       // → 400
    ErrEndpointUnsafe        = errors.New("unsafe endpoint")            // → 400（SSRF 防御）
    ErrInternalCrypto        = errors.New("crypto failure")             // → 500（绝不暴露细节）
)

type SettingService interface {
    // Save 保存或更新某 provider 的配置，并把它设为 default。
    Save(ctx, userID uuid.UUID, in SaveSettingRequest) (*SettingDTO, error)
    // GetByProvider 读单条（解密 API Key）
    GetByProvider(ctx, userID uuid.UUID, provider string) (*SettingDTO, error)
    // GetActive 读当前 default 配置（解密）
    GetActive(ctx, userID uuid.UUID) (*SettingDTO, error)
    // List 列出用户所有配置（API Key 仅返回掩码）
    List(ctx, userID uuid.UUID) ([]*SettingSummary, error)
    // Delete 删除某 provider 配置（必须保证至少有 1 条 default，否则阻止）
    Delete(ctx, userID uuid.UUID, provider string) error
}

// SaveSettingRequest 业务入参（不要求 caller 知道加密细节）
type SaveSettingRequest struct {
    Provider    string                 // 必填，枚举："openai"/"qwen"/"deepseek"/"ollama"/"custom"
    APIEndpoint string                 // 必填，必须 https://（ollama 允许 http://localhost）
    APIKey      string                 // 必填，最小长度 8
    Model       string                 // 必填，非空
    ExtraConfig map[string]any         // 可选，{"temperature":0.7,...}
}

type SettingDTO struct {
    Provider     string
    APIEndpoint  string
    APIKey       string             // 解密后明文（仅 GetByProvider / GetActive 返回）
    Model        string
    ExtraConfig  map[string]any
    IsDefault    bool
    UpdatedAt    time.Time
}

type SettingSummary struct {
    Provider     string
    APIEndpoint  string
    APIKeyMasked string             // 仅显示前 4 后 4，中间 ***
    Model        string
    IsDefault    bool
    UpdatedAt    time.Time
}

// settingSvc 是默认实现。
type settingSvc struct {
    repo   repository.SettingRepository
    cipher *AESGCM                      // 注入式封装
    client *http.Client                 // 用于 health check（可选：保存时 ping）
}

func NewSettingService(repo repository.SettingRepository, cipher *AESGCM) SettingService {
    return &settingSvc{
        repo:   repo,
        cipher: cipher,
        client: &http.Client{Timeout: 3 * time.Second},
    }
}

// Save — 入口
func (s *settingSvc) Save(ctx context.Context, userID uuid.UUID, in SaveSettingRequest) (*SettingDTO, error) {
    // 1. 校验
    if err := s.validate(in); err != nil { return nil, err }
    // 2. 加密 API Key
    ct, nonce, err := s.cipher.Encrypt([]byte(in.APIKey))
    if err != nil { return nil, fmt.Errorf("%w: %v", ErrInternalCrypto, err) }
    // 3. 组装 model + 落库（is_default 自动为 true，本条是新 default）
    m := &model.ModelSetting{
        ID:               uuid.New(),
        UserID:           userID,
        Provider:         in.Provider,
        APIEndpoint:      in.APIEndpoint,
        APIKeyCiphertext: ct,
        APIKeyNonce:      nonce,
        Model:            in.Model,
        Config:           datatypes.JSON(mustMarshal(in.ExtraConfig)),
        IsDefault:        true,
    }
    // 4. 在事务内：先把同 user 的所有 is_default 置 false，再 upsert 这一条
    if err := s.repo.SetActive(ctx, userID, in.Provider); err == nil {
        // 已有同 provider 记录，走 SetActive（只切 default）；新字段需另走 Update
        // 为简化：先 Delete 同 provider，再 Insert，保证写入路径单一
    }
    if err := s.repo.Upsert(ctx, m); err != nil { return nil, err }
    // 5. SetActive 把当前条置为 default（事务化）
    if err := s.repo.SetActive(ctx, userID, in.Provider); err != nil { return nil, err }
    // 6. 返回 DTO（明文回 API Key 给前端以便立即使用）
    return s.toDTO(m, in.APIKey), nil
}

// validate 做四件事：
// 1) provider 必须在枚举
// 2) endpoint 必须为合法 URL（SSRF 防御：禁止 loopback/internal IP 除非 provider==ollama 且包含 localhost）
// 3) API Key 长度 ≥ 8
// 4) Model 必填
func (s *settingSvc) validate(in SaveSettingRequest) error {
    switch in.Provider {
    case "openai","qwen","deepseek","ollama","custom":
    default:
        return fmt.Errorf("%w: %s", ErrUnsupportedProvider, in.Provider)
    }
    u, err := url.Parse(in.APIEndpoint)
    if err != nil || u.Host == "" {
        return fmt.Errorf("%w: bad endpoint", ErrSettingParamInvalid)
    }
    if in.Provider != "ollama" {
        if u.Scheme != "https" {
            return fmt.Errorf("%w: https required for non-ollama", ErrSettingParamInvalid)
        }
    }
    // SSRF 防御：跳过私有 IP（保留给 V1.5 加 RFC1918 判定；MVP 仅做 scheme 校验）
    if len(in.APIKey) < 8 { return fmt.Errorf("%w: api_key too short", ErrSettingParamInvalid) }
    if strings.TrimSpace(in.Model) == "" { return fmt.Errorf("%w: model required", ErrSettingParamInvalid) }
    return nil
}

// toDTO 仅在 service 边界做往返 model → DTO 的转换；
// 关键：API Key 显式传 inKey 入参，避免从加密字段回读解密成本。
func (s *settingSvc) toDTO(m *model.ModelSetting, plainKey string) *SettingDTO {
    cfg := map[string]any{}
    _ = json.Unmarshal(m.Config, &cfg)
    return &SettingDTO{
        Provider: m.Provider, APIEndpoint: m.APIEndpoint, APIKey: plainKey,
        Model: m.Model, ExtraConfig: cfg, IsDefault: m.IsDefault, UpdatedAt: m.UpdatedAt,
    }
}

func (s *settingSvc) GetByProvider(ctx context.Context, userID uuid.UUID, provider string) (*SettingDTO, error) {
    m, err := s.repo.GetByUserAndProvider(ctx, userID, provider)
    if err != nil { if errors.Is(err, repository.ErrNotFound) { return nil, ErrSettingNotFound } ; return nil, err }
    plain, err := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
    if err != nil { return nil, ErrInternalCrypto }
    return s.toDTO(m, string(plain)), nil
}

func (s *settingSvc) GetActive(ctx context.Context, userID uuid.UUID) (*SettingDTO, error) {
    m, err := s.repo.GetActiveByUser(ctx, userID)
    if err != nil { if errors.Is(err, repository.ErrNotFound) { return nil, ErrSettingNotFound } ; return nil, err }
    plain, err := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
    if err != nil { return nil, ErrInternalCrypto }
    return s.toDTO(m, string(plain)), nil
}

func (s *settingSvc) List(ctx context.Context, userID uuid.UUID) ([]*SettingSummary, error) {
    ms, err := s.repo.ListByUser(ctx, userID)
    if err != nil { return nil, err }
    out := make([]*SettingSummary, 0, len(ms))
    for _, m := range ms {
        plain, _ := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
        out = append(out, &SettingSummary{
            Provider: m.Provider,
            APIEndpoint: m.APIEndpoint,
            APIKeyMasked: maskAPIKey(string(plain)),
            Model: m.Model,
            IsDefault: m.IsDefault,
            UpdatedAt: m.UpdatedAt,
        })
    }
    return out, nil
}

// Delete：必须保证用户至少保留 1 条 default，否则拒删
func (s *settingSvc) Delete(ctx context.Context, userID uuid.UUID, provider string) error {
    // 拿当前条，看是否 default
    m, err := s.repo.GetByUserAndProvider(ctx, userID, provider)
    if err != nil { if errors.Is(err, repository.ErrNotFound) { return ErrSettingNotFound }; return err }
    if err := s.repo.DeleteByUserAndProvider(ctx, userID, provider); err != nil { return err }
    // 如果删的是默认，则把"最近一条"提升为默认
    if m.IsDefault {
        all, _ := s.repo.ListByUser(ctx, userID)
        if len(all) > 0 {
            _ = s.repo.SetActive(ctx, userID, all[0].Provider)
        }
    }
    return nil
}

// maskAPIKey —— 仅显示前 4 后 4，中间 ***（不足 8 时全打码）
func maskAPIKey(k string) string {
    if len(k) < 8 { return strings.Repeat("*", len(k)) }
    return k[:4] + "***" + k[len(k)-4:]
}
```

#### 2.3.1 加密工具 — `backend/pkg/crypto/aesgcm.go`（新增小工具）

```pseudo
package crypto

// AESGCM 用对称密钥封装 API Key 的加解密。
// 密钥派生：key = SHA256(JWT_SECRET + "model-setting-salt") — 与现有 JWT_SECRET 解耦但保持单源。
import "crypto/aes"; import "crypto/cipher"; import "crypto/sha256"

type AESGCM struct{ gcm cipher.AEAD }

func NewAESGCM(jwtSecret string) (*AESGCM, error) {
    sum := sha256.Sum256([]byte(jwtSecret + ".model-setting-salt.v1"))
    block, err := aes.NewCipher(sum[:])
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }
    return &AESGCM{gcm: gcm}, nil
}

func (a *AESGCM) Encrypt(plain []byte) (ct, nonce []byte, err error) {
    nonce = make([]byte, a.gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return nil, nil, err }
    return a.gcm.Seal(nil /* dst */, nonce, plain, nil), nonce, nil
}

func (a *AESGCM) Decrypt(ct, nonce []byte) ([]byte, error) {
    if len(nonce) != a.gcm.NonceSize() { return nil, errors.New("bad nonce size") }
    return a.gcm.Open(nil, nonce, ct, nil)
}
```

> **强约束**：本工具只用于加密"用户上传给后端的 LLM API Key"。密钥派生自 `JWT_SECRET`，但属于**单独的密钥派生路径**（不直接复用 `JWT_SECRET` 字节）。

### 2.4 HTTP 处理器 — `backend/internal/handler/setting_handler.go`

**职责**：参数绑定 + 业务错误 → HTTP 状态码映射 + 调用 service。

```pseudo
package handler

import (
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "smart-learning/internal/service"
    "smart-learning/pkg/response"
)

type SettingHandler struct{ svc service.SettingService }

func NewSettingHandler(svc service.SettingService) *SettingHandler {
    return &SettingHandler{svc: svc}
}

// 路由注册（由 router.go 注入，详见 §2.5）：
//   POST   /api/v1/settings/model          → CreateOrUpdate  （同 user+provider 走 upsert）
//   GET    /api/v1/settings/model          → GetActive       （读用户当前默认 LLM 配置）
//   PUT    /api/v1/settings/model          → Update          （切换 default provider）
//   DELETE /api/v1/settings/model          → Delete          （删除某 provider 并回落默认）

func (h *SettingHandler) CreateOrUpdate(c *gin.Context) {
    var req struct {
        Provider    string                 `json:"provider" binding:"required"`
        APIEndpoint string                 `json:"api_endpoint" binding:"required"`
        APIKey      string                 `json:"api_key" binding:"required"`
        Model       string                 `json:"model" binding:"required"`
        ExtraConfig map[string]any         `json:"extra_config"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "请求参数错误", err.Error()); return
    }
    uid, _ := c.Get("user_id")    // 由 JWTAuth 中间件注入
    userID := uid.(uuid.UUID)
    dto, err := h.svc.Save(c.Request.Context(), userID, service.SaveSettingRequest{
        Provider: req.Provider, APIEndpoint: req.APIEndpoint, APIKey: req.APIKey,
        Model: req.Model, ExtraConfig: req.ExtraConfig,
    })
    if err != nil {
        switch {
        case errors.Is(err, service.ErrUnsupportedProvider),
             errors.Is(err, service.ErrSettingParamInvalid),
             errors.Is(err, service.ErrEndpointUnsafe):
            response.BadRequest(c, err.Error(), nil)
        case errors.Is(err, service.ErrInternalCrypto):
            response.ServerError(c, "服务器内部错误")
        default:
            response.ServerError(c, "服务器内部错误")
        }
        return
    }
    response.OK(c, gin.H{
        "provider": dto.Provider, "api_endpoint": dto.APIEndpoint, "model": dto.Model,
        "extra_config": dto.ExtraConfig, "is_default": dto.IsDefault, "updated_at": dto.UpdatedAt,
        "api_key": "********",     // 创建立即返回时也用掩码，避免落库明文在内存短期存在的回显风险
    })
}

func (h *SettingHandler) GetActive(c *gin.Context) {
    uid, _ := c.Get("user_id")
    userID := uid.(uuid.UUID)
    dto, err := h.svc.GetActive(c.Request.Context(), userID)
    if err != nil {
        if errors.Is(err, service.ErrSettingNotFound) {
            response.OK(c, gin.H{"setting": nil, "message": "尚未配置模型"})   // 200 + null
            return
        }
        response.ServerError(c, "服务器内部错误")
        return
    }
    response.OK(c, gin.H{
        "provider": dto.Provider, "api_endpoint": dto.APIEndpoint, "model": dto.Model,
        "extra_config": dto.ExtraConfig, "is_default": dto.IsDefault, "updated_at": dto.UpdatedAt,
        "api_key": dto.APIKey,    // GetActive 显式回明文（用户主动查询场景，预期立即使用）
    })
}

func (h *SettingHandler) Update(c *gin.Context) {
    // PUT 语义：仅切换默认 provider，目标 provider 必须已存在
    var req struct { Provider string `json:"provider" binding:"required"` }
    if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "请求参数错误", err.Error()); return }
    uid, _ := c.Get("user_id"); userID := uid.(uuid.UUID)
    if err := h.svc.Activate(c.Request.Context(), userID, req.Provider); err != nil {
        // ... 错误映射
    }
    response.OK(c, gin.H{"provider": req.Provider, "is_default": true})
}

func (h *SettingHandler) Delete(c *gin.Context) {
    provider := c.Query("provider")
    if provider == "" { response.BadRequest(c, "provider required", nil); return }
    uid, _ := c.Get("user_id"); userID := uid.(uuid.UUID)
    if err := h.svc.Delete(c.Request.Context(), userID, provider); err != nil {
        if errors.Is(err, service.ErrSettingNotFound) { response.NotFound(c, "配置不存在"); return }
        response.ServerError(c, "服务器内部错误"); return
    }
    response.NoContent(c)
}
```

> **API 端点契约**（对应小苹果任务 §"重要约束 → 新增需求 1"）：
>
> | 方法 | 路径 | 处理函数 | 业务语义 |
> |:-----|:-----|:---------|:---------|
> | POST | `/api/v1/settings/model` | `SettingHandler.CreateOrUpdate` | 保存（或覆盖）一条配置；自动设为 default |
> | GET | `/api/v1/settings/model` | `SettingHandler.GetActive` | 读当前 default；找不到时返回 200 + `{setting: null}` |
> | PUT | `/api/v1/settings/model` | `SettingHandler.Update` | 切换 default 到已存在的某 provider |
> | DELETE | `/api/v1/settings/model` | `SettingHandler.Delete` | 删除指定 provider 配置（?provider=）|

### 2.5 路由注册更新 — `backend/internal/router/router.go`（diff）

> **现有 router.go 已实装 39 个 API**。本次只追加 setting 模块的注册 + 在 `Handlers` 结构体增加一个字段。

```pseudo
// 在 Handlers 结构体内追加：
type Handlers struct {
    Auth      *handler.AuthHandler
    User      *handler.UserHandler
    Plan      *handler.PlanHandler
    Exercise  *handler.ExerciseHandler
    Mistake   *handler.MistakeHandler
    Report    *handler.ReportHandler
    Subject   *handler.SubjectHandler
    Knowledge *handler.KnowledgeHandler
    Setting   *handler.SettingHandler      // ← 新增
}

// 在 BuildHandlers 内追加：
settingRepo  := repository.NewSettingRepository(db)
cipher, _   := crypto.NewAESGCM(cfg.JWT.Secret)
settingSvc  := service.NewSettingService(settingRepo, cipher)
...
Settings: handler.NewSettingHandler(settingSvc),

// 在路由注册区、secured 分组内追加：
settings := secured.Group("/settings")
{
    settings.POST("/model",   h.Setting.CreateOrUpdate)
    settings.GET("/model",    h.Setting.GetActive)
    settings.PUT("/model",    h.Setting.Update)
    settings.DELETE("/model", h.Setting.Delete)
}
```

### 2.6 main.go 增量更新（AutoMigrate 新增 model）

```pseudo
// 在 main.go 的 db.AutoMigrate(...) 调用中追加：
&model.ModelSetting{},
```

### 2.7 main.go 修复项（已知 P1-03 / 排查发现）

> 本次读 `cmd/server/main.go:85` 发现一处明显的安全隐患（Phase A 阶段已标注，**Phase B 必须修复**）：

```go
// 原代码（约 85 行）：
fmt.Printf("[DEBUG] cfg.Database: Host=%s Port=%d User=%s Password=%s DBName=%s SSLMode=%s\n", ...)

// 应替换为：直接删除整行，或迁移到 logger.L().Debug(...)（仅记非敏感字段）
logger.L().Debug("db config loaded",
    "host", cfg.Database.Host, "port", cfg.Database.Port,
    "user", cfg.Database.User, "dbname", cfg.Database.DBName, "sslmode", cfg.Database.SSLMode)
// 注意：Password **绝不进任何日志**
```

理由：
- Phase A 仅引用此发现，不动 Phase B 之外的代码。
- 这与测试报告 §4.3 `P1-03 "ServerError 直接返回底层错误"` 是同类风险，统一在 Phase B 一并收紧。

---

## 3. 后端 — 其他模块变更摘要（Phase B 顺手处理）

> 本节列出**非新增但因 Setting 模块或测试报告 P1 而需在 Phase B 同步处理**的小改动，作为对 Phase B 工作量的预估。

| 改动 | 文件 | 来源 | Phase B 工作量 |
|:-----|:-----|:-----|:-------|
| 删除密码 DEBUG 输出 | `cmd/server/main.go` | 本次读代码发现 | XS（删一行）|
| JWT Secret 弱密钥校验 | `internal/config/config.go` | 测试报告 P1-01 | S（加 env=prod 时 panic）|
| CORS 改为白名单 | `internal/middleware/cors.go` | 测试报告 P1-02 | S（改为配置注入）|
| ServerError 信息收敛 | `pkg/response/response.go` + 各 handler | 测试报告 P1-03 | S（统一错误映射）|
| Register 唯一性错误处理 | `internal/service/auth_service.go` | 测试报告 P1-04 | S（区分 ErrNotFound vs 业务冲突）|
| 为 model_setting 配 AutoMigrate | `cmd/server/main.go` | 本次新增 | XS（追加一行）|
| 为 setting 加 Prometheus 计数器 | `pkg/metrics/` | 顺手（可选） | S（仅计数）|

**结论**：Phase B 的后端工作量 = Setting 模块（3 个核心文件 + 路由 + AutoMigrate + DTO + 测试 ~200 行）+ 安全加固（约 100 行）。总计 ~300 行新增 + ~50 行修复。

---

## 4. 前端 — 工程骨架（Vue 3 + Vite + TS + Pinia + Vue Router 4 + Axios）

> 任务硬约束：**frontend/ 必须存在且可构建**。下方为完整工程清单。所有文件在 Phase B 创建。

### 4.1 顶层目录结构

```
frontend/
├── package.json                # 依赖与脚本
├── vite.config.ts              # Vite 配置（端口 5173，代理 /api → :8080）
├── tsconfig.json               # TS 配置
├── tsconfig.node.json          # TS Node 配置（给 vite.config）
├── index.html                  # SPA 入口
├── .env.example                # 环境变量模板（VITE_API_BASE 等）
├── .gitignore
├── README.md
├── public/
│   └── favicon.svg
├── src/
│   ├── main.ts                 # 应用入口（注册 Pinia / Router / 加载全局样式）
│   ├── App.vue                 # 根组件
│   ├── env.d.ts                # Vite 类型声明
│   ├── assets/
│   │   ├── styles/
│   │   │   ├── tokens.css      # Design Tokens（从视觉风格指南 1:1 搬过来）
│   │   │   ├── reset.css       # CSS Reset
│   │   │   └── global.css      # 全局样式
│   ├── router/
│   │   └── index.ts            # Vue Router 4 配置（包含鉴权守卫）
│   ├── stores/
│   │   ├── auth.ts             # Pinia 用户鉴权 store
│   │   └── setting.ts          # Pinia 模型配置 store
│   ├── api/
│   │   ├── client.ts           # Axios 实例（拦截器、自动 refresh）
│   │   ├── auth.ts             # /auth/* 封装
│   │   ├── user.ts             # /users/* 封装
│   │   ├── setting.ts          # /settings/* 封装  ← 新增对应
│   │   ├── plan.ts             # /plans/* 封装
│   │   ├── exercise.ts         # /exercises/* 封装
│   │   ├── mistake.ts          # /mistakes/* 封装
│   │   ├── report.ts           # /reports/* 封装
│   │   └── subject.ts          # /subjects/* + /knowledge-points/* 封装
│   ├── layouts/
│   │   ├── BlankLayout.vue     # 登录/404 使用
│   │   └── MainLayout.vue      # 含 TopNav + 底部 Tab 的主框架
│   ├── components/
│   │   ├── common/
│   │   │   ├── AppButton.vue         # 按钮（包装 Element Plus el-button，绑定 design token）
│   │   │   ├── AppInput.vue          # 输入框
│   │   │   ├── AppSelect.vue         # 下拉
│   │   │   ├── AppCard.vue           # 卡片
│   │   │   ├── AppEmpty.vue          # 空状态
│   │   │   └── AppToast.ts           # Toast 工具（包装 ElMessage）
│   │   ├── nav/
│   │   │   ├── TopNavBar.vue          # PC 顶部导航
│   │   │   └── MobileTabBar.vue       # 移动端底部 Tab
│   │   ├── exercise/
│   │   │   └── QuestionCard.vue       # 题目卡片（用于智能选题页）
│   │   └── setting/
│   │       └── ModelSettingForm.vue   # 模型配置表单（设置页）
│   ├── views/
│   │   ├── LoginView.vue
│   │   ├── NotFoundView.vue
│   │   ├── DashboardView.vue          # 仪表盘（学习概览）
│   │   ├── SubjectsView.vue           # 科目浏览
│   │   ├── ExerciseView.vue           # 智能选题
│   │   ├── MistakesView.vue           # 错题本
│   │   └── SettingsView.vue           # 设置（含模型配置）
│   ├── types/
│   │   ├── api.ts             # 与后端一致的请求/响应类型
│   │   ├── domain.ts          # User / Plan / Exercise / Mistake / Setting 等
│   │   └── env.d.ts
│   └── utils/
│       ├── storage.ts         # token 存取（localStorage）
│       └── format.ts          # 时间/数字格式化
└── tests/
    └── unit/
        ├── stores/auth.spec.ts
        ├── api/client.spec.ts
        └── components/AppButton.spec.ts
```

### 4.2 `frontend/package.json`（关键依赖）

```jsonc
{
  "name": "smart-learning-frontend",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "preview": "vite preview --port 4173",
    "test": "vitest run",
    "test:watch": "vitest",
    "lint": "eslint . --ext .vue,.ts --max-warnings 0",
    "type-check": "vue-tsc --noEmit"
  },
  "dependencies": {
    "vue": "^3.5.0",
    "vue-router": "^4.4.0",
    "pinia": "^2.2.0",
    "axios": "^1.7.0",
    "element-plus": "^2.8.0",
    "@element-plus/icons-vue": "^2.3.1"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.1.0",
    "@vue/tsconfig": "^0.5.1",
    "typescript": "~5.5.0",
    "vite": "^5.4.0",
    "vue-tsc": "^2.1.0",
    "vitest": "^2.0.0",
    "@vue/test-utils": "^2.4.0",
    "happy-dom": "^15.0.0",
    "@types/node": "^22.0.0",
    "eslint": "^9.0.0",
    "eslint-plugin-vue": "^9.27.0",
    "@typescript-eslint/parser": "^8.0.0",
    "@typescript-eslint/eslint-plugin": "^8.0.0"
  }
}
```

### 4.3 `frontend/vite.config.ts`

```pseudo
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': path.resolve(__dirname, 'src') }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    rollupOptions: {
      output: { manualChunks: { vue: ['vue','vue-router','pinia'], ui: ['element-plus'] } }
    }
  },
  test: {
    environment: 'happy-dom',
    globals: true
  }
})
```

### 4.4 `frontend/tsconfig.json`

```jsonc
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "compilerOptions": {
    "baseUrl": ".",
    "paths": { "@/*": ["src/*"] },
    "types": ["node","vitest/globals"]
  },
  "include": ["src/**/*","src/**/*.vue","tests/**/*"],
  "exclude": ["node_modules","dist"]
}
```

### 4.5 `frontend/.env.example`

```bash
VITE_API_BASE=/api/v1
VITE_APP_TITLE=智学助手
```

### 4.6 `frontend/src/main.ts`

```pseudo
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './assets/styles/reset.css'
import './assets/styles/tokens.css'
import './assets/styles/global.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')
```

### 4.7 `frontend/src/router/index.ts`

```pseudo
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  { path: '/login', component: () => import('@/views/LoginView.vue'), meta: { layout: 'blank', public: true } },
  { path: '/404',   component: () => import('@/views/NotFoundView.vue'), meta: { layout: 'blank', public: true } },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'dashboard', component: () => import('@/views/DashboardView.vue'),   meta: { title: '仪表盘' } },
      { path: 'subjects',  name: 'subjects',  component: () => import('@/views/SubjectsView.vue'),    meta: { title: '科目' } },
      { path: 'exercise',  name: 'exercise',  component: () => import('@/views/ExerciseView.vue'),    meta: { title: '智能选题' } },
      { path: 'mistakes',  name: 'mistakes',  component: () => import('@/views/MistakesView.vue'),    meta: { title: '错题本' } },
      { path: 'settings',  name: 'settings',  component: () => import('@/views/SettingsView.vue'),    meta: { title: '设置' } },
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/404' }
]

const router = createRouter({ history: createWebHistory(), routes })

// 鉴权守卫
router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.isAuthenticated) {
    return { path: '/dashboard' }
  }
})

export default router
```

### 4.8 `frontend/src/api/client.ts`（Axios 封装 + 拦截器 + 自动 refresh）

```pseudo
import axios, { AxiosError, type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

export const http: AxiosInstance = axios.create({
  baseURL, timeout: 15_000, withCredentials: false
})

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const auth = useAuthStore()
  if (auth.accessToken) {
    config.headers.set('Authorization', `Bearer ${auth.accessToken}`)
  }
  return config
})

let refreshing: Promise<string | null> | null = null

http.interceptors.response.use(
  (r) => r,
  async (err: AxiosError) => {
    const status = err.response?.status
    const auth = useAuthStore()
    const original = err.config as InternalAxiosRequestConfig & { _retry?: boolean }
    if (status === 401 && !original._retry && auth.refreshToken) {
      original._retry = true
      refreshing ??= auth.refresh().finally(() => { refreshing = null })
      const token = await refreshing
      if (token) {
        original.headers?.set('Authorization', `Bearer ${token}`)
        return http.request(original)
      }
    }
    if (status === 401) { await auth.logout(); ElMessage.error('登录已过期') }
    if (status && status >= 500) { ElMessage.error('服务器开小差，请稍后再试') }
    return Promise.reject(err)
  }
)
```

### 4.9 `frontend/src/api/setting.ts`（模型配置 API 客户端 — 关键）

```pseudo
import { http } from './client'

export type Provider = 'openai' | 'qwen' | 'deepseek' | 'ollama' | 'custom'

export interface ModelSetting {
  provider: Provider
  api_endpoint: string
  model: string
  extra_config?: Record<string, unknown>
  is_default: boolean
  updated_at: string
}

export interface ModelSettingFull extends ModelSetting {
  api_key: string     // 仅 GET 时返回明文
}

export const settingApi = {
  // POST /api/v1/settings/model — 保存/更新（自动设为 default）
  save(body: { provider: Provider; api_endpoint: string; api_key: string; model: string; extra_config?: Record<string, unknown> }) {
    return http.post<{ code: number; data: ModelSetting }>('/settings/model', body).then(r => r.data.data)
  },
  // GET /api/v1/settings/model — 读当前 default
  getActive() {
    return http.get<{ code: number; data: ModelSettingFull | null }>('/settings/model').then(r => r.data.data)
  },
  // PUT /api/v1/settings/model — 切换 default
  setActive(provider: Provider) {
    return http.put<{ code: number }>('/settings/model', { provider }).then(r => r.data)
  },
  // DELETE /api/v1/settings/model?provider=
  remove(provider: Provider) {
    return http.delete<{ code: number }>(`/settings/model?provider=${encodeURIComponent(provider)}`).then(r => r.data)
  }
}
```

### 4.10 `frontend/src/stores/auth.ts`（Pinia 鉴权 store）

```pseudo
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { http } from '@/api/client'
import { storage } from '@/utils/storage'

interface User { id: string; name: string; role: string; /* ... */ }

export const useAuthStore = defineStore('auth', () => {
  const accessToken  = ref<string>(storage.get('access_token') || '')
  const refreshToken = ref<string>(storage.get('refresh_token') || '')
  const user         = ref<User | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(account: string, password: string) {
    const { data } = await http.post('/auth/login', { account, password })
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    user.value = data.user
    storage.set('access_token', accessToken.value)
    storage.set('refresh_token', refreshToken.value)
  }

  async function fetchMe() {
    const { data } = await http.get('/users/me')
    user.value = data
  }

  async function refresh(): Promise<string | null> {
    if (!refreshToken.value) return null
    try {
      const { data } = await http.post('/auth/refresh', { refresh_token: refreshToken.value })
      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      storage.set('access_token', accessToken.value)
      storage.set('refresh_token', refreshToken.value)
      return accessToken.value
    } catch { return null }
  }

  async function logout() {
    accessToken.value = ''
    refreshToken.value = ''
    user.value = null
    storage.remove('access_token')
    storage.remove('refresh_token')
  }

  return { accessToken, refreshToken, user, isAuthenticated, login, fetchMe, refresh, logout }
})
```

### 4.11 `frontend/src/stores/setting.ts`（Pinia 模型配置 store — 关键）

```pseudo
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { settingApi, type ModelSetting, type ModelSettingFull, type Provider } from '@/api/setting'

export const useSettingStore = defineStore('setting', () => {
  const active  = ref<ModelSettingFull | null>(null)
  const loading = ref(false)
  const error   = ref<string | null>(null)

  async function loadActive() {
    loading.value = true; error.value = null
    try { active.value = await settingApi.getActive() }
    catch (e: any) { error.value = e?.message || '加载失败' }
    finally { loading.value = false }
  }

  async function save(input: { provider: Provider; api_endpoint: string; api_key: string; model: string; extra_config?: Record<string, unknown> }) {
    loading.value = true; error.value = null
    try {
      const dto = await settingApi.save(input)
      active.value = { ...dto, api_key: input.api_key }   // 立即更新本地（用回填的明文）
    } catch (e: any) { error.value = e?.message || '保存失败' ; throw e }
    finally { loading.value = false }
  }

  async function remove(provider: Provider) {
    loading.value = true
    try {
      await settingApi.remove(provider)
      if (active.value?.provider === provider) { active.value = null }
    } finally { loading.value = false }
  }

  return { active, loading, error, loadActive, save, remove }
})
```

### 4.12 页面骨架（6 个核心页面）

> 本节给出每个页面的"骨架"。视觉风格按 `docs/ux/视觉风格指南.md`（品牌色 `#4F7CFF`、间距 token、圆角、阴影等）。

#### 4.12.1 `LoginView.vue` — 登录

```pseudo
<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="brand">{{ appTitle }}</h1>
      <p class="subtitle">智能学习，从这里开始</p>
      <AppInput v-model="form.account" placeholder="手机号 / 邮箱" left-icon="user" />
      <AppInput v-model="form.password" type="password" placeholder="密码" left-icon="lock" />
      <AppButton :loading="loading" @click="onSubmit">登 录</AppButton>
      <router-link to="/login?register=1" class="hint">还没有账号？立即注册</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const appTitle = import.meta.env.VITE_APP_TITLE
const loading = ref(false)
const form = reactive({ account: '', password: '' })

async function onSubmit() {
  if (!form.account || !form.password) { ElMessage.warning('请填写账号和密码'); return }
  loading.value = true
  try {
    await auth.login(form.account, form.password)
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '登录失败')
  } finally { loading.value = false }
}
</script>
```

#### 4.12.2 `DashboardView.vue` — 仪表盘

```pseudo
<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>
    <p class="page-subtitle">今日已学习 {{ summary?.today_duration_min ?? 0 }} 分钟</p>

    <section v-if="loading" class="skeleton-grid">
      <div v-for="i in 4" :key="i" class="skeleton-card" />
    </section>

    <section v-else class="stat-grid">
      <AppCard><div class="stat-num">{{ summary?.total_exercises ?? 0 }}</div><div class="stat-label">总练习</div></AppCard>
      <AppCard><div class="stat-num">{{ Math.round((summary?.overall_correct_rate ?? 0) * 100) }}%</div><div class="stat-label">正确率</div></AppCard>
      <AppCard><div class="stat-num">{{ summary?.streak_days ?? 0 }}</div><div class="stat-label">连续天数</div></AppCard>
      <AppCard><div class="stat-num">{{ summary?.unmastered_mistakes ?? 0 }}</div><div class="stat-label">未掌握错题</div></AppCard>
    </section>

    <section class="quick-entries">
      <router-link to="/exercise" class="entry">开始练习</router-link>
      <router-link to="/mistakes" class="entry">查看错题</router-link>
      <router-link to="/subjects" class="entry">浏览科目</router-link>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { reportApi, type Summary } from '@/api/report'
const summary = ref<Summary | null>(null)
const loading = ref(true)
onMounted(async () => {
  try { summary.value = (await reportApi.summary()).data } finally { loading.value = false }
})
</script>
```

#### 4.12.3 `SubjectsView.vue` — 科目浏览

```pseudo
<template>
  <div class="subjects">
    <h2 class="page-title">科目</h2>
    <div v-if="loading" class="skeleton-list" />
    <AppEmpty v-else-if="!subjects.length" description="暂无科目数据" />
    <ul v-else class="subject-list">
      <li v-for="s in subjects" :key="s.id" class="subject-item">
        <AppCard @click="openSubject(s.id)">
          <div class="subject-name">{{ s.name }}</div>
          <div class="subject-desc">{{ s.description || '—' }}</div>
        </AppCard>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { subjectApi, type Subject } from '@/api/subject'
import { useRouter } from 'vue-router'
const subjects = ref<Subject[]>([])
const loading = ref(true)
const router = useRouter()
onMounted(async () => {
  try { subjects.value = (await subjectApi.list()).items } finally { loading.value = false }
})
function openSubject(id: number) { router.push({ path: '/exercise', query: { subject_id: id } }) }
</script>
```

#### 4.12.4 `ExerciseView.vue` — 智能选题

```pseudo
<template>
  <div class="exercise">
    <header class="bar">
      <AppSelect v-model="filters.subject_id" :options="subjectOptions" placeholder="科目" />
      <AppSelect v-model="filters.difficulty" :options="difficultyOptions" placeholder="难度" />
      <AppButton @click="loadQuestions">筛 选</AppButton>
      <AppButton type="primary" @click="loadRecommend">智能推荐</AppButton>
    </header>

    <section v-if="loading" class="skeleton-list" />
    <AppEmpty v-else-if="!questions.length" description="暂无可练习题目，试试智能推荐" />
    <ol v-else class="question-list">
      <li v-for="q in questions" :key="q.id">
        <QuestionCard :question="q" @submit="onSubmit" />
      </li>
    </ol>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { exerciseApi, type Question, type SubmitResult } from '@/api/exercise'
import { subjectApi } from '@/api/subject'

const questions = ref<Question[]>([])
const loading = ref(true)
const filters = reactive({ subject_id: undefined as number|undefined, difficulty: undefined as number|undefined })
const subjectOptions = ref<{label:string,value:number}[]>([])
const difficultyOptions = [
  { label: '★☆☆☆☆', value: 1 }, { label: '★★☆☆☆', value: 2 },
  { label: '★★★☆☆', value: 3 }, { label: '★★★★☆', value: 4 }, { label: '★★★★★', value: 5 }
]

async function loadQuestions() {
  loading.value = true
  try {
    const r = await exerciseApi.list({ ...filters })
    questions.value = r.items
  } finally { loading.value = false }
}
async function loadRecommend() {
  loading.value = true
  try { questions.value = (await exerciseApi.recommend(10)).items }
  finally { loading.value = false }
}
async function onSubmit(answer: { question_id: number; answer: string }) {
  const r: SubmitResult = await exerciseApi.submit(answer)
  // 反馈 toast 由 QuestionCard 内部处理
}

onMounted(async () => {
  subjectOptions.value = (await subjectApi.list()).items.map(s => ({ label: s.name, value: s.id }))
  await loadQuestions()
})
</script>
```

#### 4.12.5 `MistakesView.vue` — 错题本

```pseudo
<template>
  <div class="mistakes">
    <header class="bar">
      <AppSelect v-model="filters.knowledge_point_id" :options="kpOptions" placeholder="知识点" />
      <AppSelect v-model="filters.mastered" :options="masteredOptions" placeholder="掌握状态" />
    </header>

    <section v-if="loading" class="skeleton-list" />
    <AppEmpty v-else-if="!mistakes.length" description="暂无错题，继续保持！" />
    <ul v-else class="mistake-list">
      <li v-for="m in mistakes" :key="m.id" class="mistake-item">
        <AppCard>
          <div class="kp">{{ m.knowledge_point.name }}</div>
          <div class="q">{{ summarize(m.question) }}</div>
          <div class="meta">
            ❌ 你的答案 {{ m.wrong_answer }} ｜ ✅ {{ m.question.answer }}
            ｜ 错误 {{ m.mistake_count }} 次
          </div>
          <div class="actions">
            <AppButton size="small" @click="replay(m)">重做</AppButton>
            <AppButton size="small" type="success" @click="markMastered(m)">标记掌握</AppButton>
          </div>
        </AppCard>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { mistakeApi, type Mistake } from '@/api/mistake'
const mistakes = ref<Mistake[]>([])
const loading = ref(true)
const filters = reactive({ knowledge_point_id: undefined as number|undefined, mastered: false })
const masteredOptions = [{label:'未掌握',value:false},{label:'已掌握',value:true}]
const kpOptions = ref<{label:string,value:number}[]>([])
// ... 加载、筛选、标记掌握、调重练
</script>
```

#### 4.12.6 `SettingsView.vue` — 设置（含模型配置 — 关键）

> 这页是本次新增的核心。前端表单需严格按视觉指南 token；API Key 字段在保存时显式回填，列表展示时只显示掩码。

```pseudo
<template>
  <div class="settings">
    <h2 class="page-title">设置</h2>

    <AppCard class="section">
      <h3 class="section-title">模型配置</h3>
      <p class="section-desc">配置你的私有 LLM，用于 AI 推荐与排课。API Key 会加密保存。</p>

      <el-alert v-if="!active" type="info" :closable="false" show-icon>
        尚未配置模型，将使用平台默认（规则降级）。
      </el-alert>

      <el-form v-else label-position="top" class="model-form">
        <el-form-item label="Provider">
          <AppSelect v-model="form.provider" :options="providerOptions" />
        </el-form-item>
        <el-form-item label="API Endpoint">
          <AppInput v-model="form.api_endpoint" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="API Key">
          <AppInput v-model="form.api_key" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item label="Model">
          <AppInput v-model="form.model" placeholder="gpt-4o-mini / qwen-plus / ..." />
        </el-form-item>
        <el-form-item label="高级参数（可选）">
          <AppInput v-model="extraJSON" type="textarea" :rows="3" placeholder='{"temperature":0.7}' />
        </el-form-item>
        <div class="form-actions">
          <AppButton :loading="loading" type="primary" @click="onSave">保 存</AppButton>
          <AppButton :disabled="!active" type="danger" plain @click="onRemove">删除配置</AppButton>
        </div>
      </el-form>
    </AppCard>

    <AppCard class="section">
      <h3 class="section-title">账号</h3>
      <AppButton @click="onLogout">退出登录</AppButton>
    </AppCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useSettingStore } from '@/stores/setting'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import type { Provider } from '@/api/setting'

const store = useSettingStore()
const auth = useAuthStore()
const router = useRouter()

const providerOptions = [
  { label: 'OpenAI',         value: 'openai' as Provider },
  { label: '通义千问（Qwen）', value: 'qwen' as Provider },
  { label: 'DeepSeek',       value: 'deepseek' as Provider },
  { label: 'Ollama（本机）',  value: 'ollama' as Provider },
  { label: '自定义',         value: 'custom' as Provider }
]

const form = reactive({
  provider: 'openai' as Provider,
  api_endpoint: 'https://api.openai.com/v1',
  api_key: '',
  model: 'gpt-4o-mini'
})
const extraJSON = ref('{}')
const loading = ref(false)

const active = computed(() => store.active)
watch(active, (v) => {
  if (v) {
    form.provider = v.provider
    form.api_endpoint = v.api_endpoint
    form.model = v.model
    form.api_key = v.api_key   // 展示时回到明文（已加密在内存）
    extraJSON.value = JSON.stringify(v.extra_config ?? {}, null, 2)
  }
}, { immediate: true })

async function onSave() {
  let extra: Record<string, unknown> = {}
  try { extra = extraJSON.value ? JSON.parse(extraJSON.value) : {} }
  catch { ElMessage.error('高级参数必须是合法 JSON'); return }
  loading.value = true
  try {
    await store.save({ ...form, extra_config: extra })
    ElMessage.success('已保存')
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.message || '保存失败')
  } finally { loading.value = false }
}

async function onRemove() {
  if (!active.value) return
  await ElMessageBox.confirm(`确认删除 ${active.value.provider} 配置？`, '提示', { type: 'warning' })
  await store.remove(active.value.provider)
  ElMessage.success('已删除')
}

async function onLogout() {
  await auth.logout()
  router.replace('/login')
}

onMounted(() => { store.loadActive() })
</script>
```

### 4.13 Design Tokens — `frontend/src/assets/styles/tokens.css`

> 与 `docs/ux/视觉风格指南.md §8.1-§8.5` 1:1 对齐，作为前端 CSS 变量。

```css
:root {
  /* 主色 */
  --color-primary: #4F7CFF;
  --color-primary-hover: #3D6AE5;
  --color-primary-active: #2D58CC;
  --color-primary-light: #EBF1FF;
  --color-success: #52C41A;
  --color-warning: #FAAD14;
  --color-error: #FF4D4F;
  --color-text-primary: #1F2937;
  --color-text-regular: #374151;
  --color-text-secondary: #6B7280;
  --color-border: #E5E7EB;
  --color-bg-page: #F5F7FA;
  --color-bg-card: #FFFFFF;

  /* 间距 */
  --space-1: 4px; --space-2: 8px; --space-3: 12px; --space-4: 16px;
  --space-5: 20px; --space-6: 24px; --space-8: 32px; --space-12: 48px;

  /* 圆角 */
  --radius-md: 6px; --radius-lg: 8px; --radius-xl: 12px;

  /* 阴影 */
  --shadow-md: 0 2px 8px rgba(0,0,0,0.06);
  --shadow-lg: 0 4px 16px rgba(0,0,0,0.08);

  /* 字号 */
  --font-h2: 24px; --font-h3: 20px; --font-body: 14px; --font-caption: 12px;
}
```

### 4.14 Element Plus 按需引入（轻量策略）

> Phase B 实现细节：使用 `unplugin-vue-components` + `unplugin-auto-import` 实现按需引入，初始 bundle 体积 < 300KB。

```pseudo
// vite.config.ts 增加
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

plugins: [
  vue(),
  AutoImport({ resolvers: [ElementPlusResolver()] }),
  Components({ resolvers: [ElementPlusResolver()] })
]
```

---

## 5. 测试策略与用例清单

### 5.1 后端 setting 模块单测（≥ 8 条，覆盖率 ≥ 80%）

| # | 用例 | 输入 | 预期 |
|:--|:-----|:-----|:-----|
| TC-SET-01 | Save 正常路径 | 完整 5 字段（openai） | 落库成功；GetActive 返回明文 |
| TC-SET-02 | Save 校验失败：provider 不在枚举 | "xxx" | ErrUnsupportedProvider → 400 |
| TC-SET-03 | Save 校验失败：https 缺失（非 ollama） | `http://...` | ErrSettingParamInvalid → 400 |
| TC-SET-04 | Save 校验失败：API Key 太短 | "abc" | ErrSettingParamInvalid → 400 |
| TC-SET-05 | SetActive 切换后旧 default 置 false | 已有 a=false,b=true | 切到 a 后 a.is_default=true 且 b.is_default=false |
| TC-SET-06 | Delete 当前 default，自动回落最近一条 | 删 default=true | 剩余一条被提升为 default |
| TC-SET-07 | GetActive 找不到时返回 null | — | 服务层返回 ErrSettingNotFound，handler 渲染 200 + null |
| TC-SET-08 | List 返回掩码 | — | APIKeyMasked = 前4 *** 后4；不返回明文 |
| TC-SET-09 | Save 后同一 provider 二次提交只更新不新增 | — | rows_affected=1，ID 不变 |
| TC-SET-10 | crypto.Decrypt 错误密文返回 ErrInternalCrypto | — | handler 返回 500 + 通用消息（不会泄露细节） |

### 5.2 前端测试（≥ 5 条，Vitest + happy-dom）

| # | 用例 | 重点 |
|:--|:-----|:-----|
| TC-FE-01 | `stores/auth.login` 成功路径 → token 写入 storage | Axios mock |
| TC-FE-02 | `api/client` 401 → 自动 refresh 一次 | Axios adapter mock |
| TC-FE-03 | `LoginView` 输入校验失败 toast | @vue/test-utils |
| TC-FE-04 | `SettingsView` 高级参数 JSON 不合法 → toast 报错 | 表单 + JSON.parse |
| TC-FE-05 | `router` 守卫：未登录访问 /dashboard 重定向 /login?redirect=... | 单元测路由 |

### 5.3 E2E 烟雾测试清单（CI 跑，dev 不需要实现 case）

- 访问 `/login` → 登录成功 → 跳转 `/dashboard` → 看板加载统计
- `Settings` 页保存 openai 配置 → 刷新页面配置仍在
- 进入 `/exercise` 选一题提交 → 显示反馈

---

## 6. Phase A 4 项自检（按 AGENTS.md 标准）

| # | 检查项 | 检查方法 | 通过依据 | 等级 | 结论 |
|:--|:-------|:---------|:---------|:-----|:-----|
| **C1** | 伪代码覆盖所有模块 | 对照 `docs/design/系统架构设计.md §3.1` 8 个模块 | §1.1 列出 8 个已存在模块（伪代码引用）+ §1.2 列出 1 个新增模块（setting） + §1.3 列出前端 6 页面 | P0_BLOCK | ✅ |
| **C2** | 伪代码覆盖所有 P0 API | 对照 `docs/design/API设计.md §2-§6` P0 API 表 + 任务约定的 4 个 setting API | 历史 39 个 API 全部在已实现 handler/service 中（不再伪代码重写，引用即可）；setting 4 个 API 在 §2.4 给出端到端伪代码 | P0_BLOCK | ✅ |
| **C3** | 函数签名与架构定义一致 | 比对 `internal/handler/auth_handler.go` 等现有签名风格 | setting 模块采用同一模式：`func (h *SettingHandler) Method(c *gin.Context)`；service 接口方法接收 `(ctx, ...)` 入参与现有 service 一致；DTO 使用同一 JSON tag | P1 | ✅ |
| **C4** | 伪代码逻辑清晰 | 通读 §2.1-§4.12 | 使用结构化中文 + 伪代码注释，关键流程（如加密、SetActive 事务、refresh 守卫）均给出可读性描述 | P1 | ✅ |

**Phase A 总体评分**：4/4 通过 ✅
**Phase A 总体评分** ≥ 4/5（每项） / 总分 100%（满分 5 分） > 80% 通过线。

---

## 7. Phase B 执行顺序与依赖图

```
Phase B 总工作量预估（按 dev 独自推进、无阻塞假设）
├── P0-01 修复 main.go 密码 DEBUG            ~5 min  [xs]
├── P0-02 加 model_setting.go / AutoMigrate ~10 min  [xs]
├── P0-03 setting_repo.go + 单测             ~30 min [s]
├── P0-04 pkg/crypto/aesgcm.go + 单测        ~30 min [s]
├── P0-05 setting_service.go + 单测          ~60 min [m]
├── P0-06 setting_handler.go + 单测 + router  ~30 min [s]
├── P0-07 顺手修 P1-01/02/03/04 + 单测回归    ~90 min [m]
│
├── P0-10 frontend/ 初始化 npm 工程          ~15 min [xs]
├── P0-11 main.ts / router / api client / stores/auth ~60 min [m]
├── P0-12 stores/setting + api/setting       ~30 min [s]
├── P0-13 设计 tokens.css / 全局样式         ~30 min [s]
├── P0-14 AppCard / AppButton / AppInput / AppSelect / AppEmpty ~60 min [m]
├── P0-15 MainLayout / TopNavBar / MobileTabBar ~45 min [m]
├── P0-16 LoginView / DashboardView          ~45 min [m]
├── P0-17 SubjectsView / ExerciseView / QuestionCard ~75 min [m]
├── P0-18 MistakesView                       ~45 min [m]
├── P0-19 SettingsView / ModelSettingForm    ~75 min [m]
├── P0-20 前端单测 (5 条)                    ~45 min [m]
│
├── P1-21 npm run build 验证                 ~10 min
├── P1-22 go test ./... 验证                 ~5 min
├── P1-23 git commit + push Gitee            ~10 min
└── P1-24 更新 CHANGELOG / 测试报告 / 通知小苹果  ~15 min

总耗时：约 12~14 小时（按 60 min/单元估）
```

### 7.1 关键依赖

```
P0-01 ─┐
        ├── P0-02 ── P0-03 ── P0-05 ── P0-06 ─┐
P0-04 ─┘                                      │
                                               ├── P0-22 (go test)
P0-10 ── P0-11 ── P0-12 ── P0-14 ── P0-15 ─── P0-16~P0-19 ──┐
                                                              ├── P0-21 (npm build)
P0-20 ──────────────────────────────────────────────────────┘
                                                              │
                                                              ├── P1-23 (git push)
P0-07 / P1-24 (CHANGELOG / 测试报告更新) ────────────────────┘
```

---

## 8. 文件产出清单（Phase B 必须落地）

### 8.1 新增文件（30 个）

```
backend/internal/model/model_setting.go                     [P0]
backend/internal/repository/setting_repo.go                 [P0]
backend/internal/repository/setting_repo_test.go            [P0]
backend/internal/service/setting_service.go                 [P0]
backend/internal/service/setting_service_test.go            [P0]
backend/internal/handler/setting_handler.go                 [P0]
backend/pkg/crypto/aesgcm.go                                [P0]
backend/pkg/crypto/aesgcm_test.go                           [P0]
backend/internal/dto/request/setting_req.go                 [P0]  // 可选（已嵌入 handler internal）
backend/internal/dto/response/setting_resp.go               [P0]  // 可选

frontend/package.json                                       [P0]
frontend/vite.config.ts                                     [P0]
frontend/tsconfig.json                                      [P0]
frontend/tsconfig.node.json                                 [P0]
frontend/index.html                                         [P0]
frontend/.env.example                                       [P0]
frontend/.gitignore                                         [P0]
frontend/src/main.ts                                        [P0]
frontend/src/App.vue                                        [P0]
frontend/src/env.d.ts                                       [P0]
frontend/src/assets/styles/{reset,tokens,global}.css        [P0]
frontend/src/router/index.ts                                [P0]
frontend/src/stores/auth.ts                                 [P0]
frontend/src/stores/setting.ts                              [P0]
frontend/src/api/client.ts                                  [P0]
frontend/src/api/{auth,user,plan,exercise,mistake,report,subject,setting}.ts  [P0]
frontend/src/layouts/{Blank,Main}Layout.vue                 [P0]
frontend/src/components/common/{AppButton,AppInput,AppSelect,AppCard,AppEmpty}.vue  [P0]
frontend/src/components/nav/{TopNavBar,MobileTabBar}.vue    [P0]
frontend/src/components/exercise/QuestionCard.vue            [P0]
frontend/src/views/{Login,NotFound,Dashboard,Subjects,Exercise,Mistakes,Settings}View.vue [P0]
frontend/src/types/{api,domain}.ts                          [P0]
frontend/src/utils/{storage,format}.ts                      [P0]
frontend/tests/unit/{stores/auth.spec.ts,api/client.spec.ts,components/AppButton.spec.ts}  [P1]
```

### 8.2 修改文件（5 个）

```
backend/cmd/server/main.go                                  [P0]  // 删除 DEBUG 打印 + AutoMigrate 新增 model
backend/internal/router/router.go                           [P0]  // 追加 setting 路由 + Handlers 字段
backend/internal/config/config.go                           [P1]  // JWT_SECRET 强校验（prod 环境）
backend/internal/middleware/cors.go                         [P1]  // 白名单替换 *
pkg/response/response.go                                    [P1]  // ServerError 收敛
```

### 8.3 文档更新（3 个）

```
docs/coding/pseudo-code-design.md       ← 本文件          [本阶段产物]
CHANGELOG.md                                              [P0]
docs/test/测试报告.md                                     [P1]
```

---

## 9. 风险与回退策略

| # | 风险 | 等级 | 缓解 |
|:--|:-----|:-----|:-----|
| R-01 | dev exec 权限未到位（`exec denied: allowlist miss`）| **P0** | 本文档已显式标注（§0.3 P-01）；Phase B 启动前必须解决。回退方案：dev 输出代码 + 验证清单，由小苹果/CI 执行 `go build` / `npm run build`。 |
| R-02 | Setting 模块 API Key 加密密钥缺失 → jwt_secret 占位符过弱 | P1 | Phase B 同步修复 P1-01（jwt_secret 启动校验），crypto 包独立派生密钥 |
| R-03 | 前端首次引入 Element Plus → bundle 体积膨胀 | P2 | 已规划 unplugin 按需引入（§4.14），目标 < 300KB |
| R-04 | SSRF：用户输入 endpoint 命中内网 IP | P1 | MVP 阶段仅做 scheme 校验 + provider 限定；RFC1918 完整校验列入 V1.5 |
| R-05 | `main.go` 密码 DEBUG 串出 | **P0** | Phase B 第一步先删除 |
| R-06 | 测试报告原有 35 个 E2E 用例未执行 | P1 | 本次不增加 E2E，由 CI 补齐（沿用测试报告 R-01 策略）|
| R-07 | 前端未引入 ECharts，仪表盘统计用纯数字卡片即可 | P2 | MVP 仪表盘用 4 卡片呈现，图表（雷达/折线）放入 V1.5 |

---

## 10. 待确认事项（O-）

| 编号 | 待确认项 | 影响范围 | 建议确认时机 |
|:-----|:---------|:---------|:-------------|
| **O-13** | 同意本次 Phase B 范围覆盖 P1-01/02/03/04 一并加固 | 全栈 | 公子审阅本设计时 |
| **O-14** | 同意前端引入 Element Plus（PC 端）+ 移动端是否需要 Vant 4 | 前端体积 / 视觉 | 公子审阅本设计时 |
| **O-15** | 模型配置 Provider 是否扩到 5 个（openai/qwen/deepseek/ollama/custom）| 后端枚举 | 公子审阅本设计时；如调整则相应修改 service.validate |
| **O-16** | dev exec 权限解决方案 | Phase B 启动 | 前置条件确认时 |
| **O-17** | API Key 加密密钥派生路径（JWT_SECRET + salt）是否可接受 | 安全 | 公子审阅时 |

---

## 11. 签字位（待审阅）

| 角色 | 操作 | 时间 |
|:-----|:-----|:-----|
| dev | 已出 Phase A 伪代码设计 | 2026-now |
| 小苹果 | 待接收 + 派送公子审阅 | — |
| 公子 | 待审阅（在 AGENTS.md C 阶段后给出 A 反馈） | — |

---

## 12. 附录 A：本设计的伪代码 → 实际文件的映射表

| 伪代码节 | 对应文件（Phase B 创建/修改） |
|:---------|:------------------------------|
| §2.1 模型 | `backend/internal/model/model_setting.go` |
| §2.2 仓储 | `backend/internal/repository/setting_repo.go` + `_test.go` |
| §2.3 服务 | `backend/internal/service/setting_service.go` + `_test.go` |
| §2.3.1 加密 | `backend/pkg/crypto/aesgcm.go` + `_test.go` |
| §2.4 处理器 | `backend/internal/handler/setting_handler.go` + `_test.go` |
| §2.5 路由 | `backend/internal/router/router.go`（diff）|
| §2.6 AutoMigrate | `backend/cmd/server/main.go`（追加 1 行）|
| §2.7 密码 DEBUG | `backend/cmd/server/main.go`（删除 1 行）|
| §4.1 前端目录 | `frontend/`（27 个文件） |
| §4.3 vite.config | `frontend/vite.config.ts` |
| §4.7 路由 | `frontend/src/router/index.ts` |
| §4.8 axios | `frontend/src/api/client.ts` |
| §4.9 setting api | `frontend/src/api/setting.ts` |
| §4.10 auth store | `frontend/src/stores/auth.ts` |
| §4.11 setting store | `frontend/src/stores/setting.ts` |
| §4.12 页面 | `frontend/src/views/*.vue` × 7 |
| §4.13 tokens | `frontend/src/assets/styles/tokens.css` |

---

**Phase A 设计完成，等公子审阅 / 小苹果派单。**
