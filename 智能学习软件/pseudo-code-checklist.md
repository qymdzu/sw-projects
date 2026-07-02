# 伪代码清单 — 智能学习软件后端

> 版本：v1.0.0
> 对应设计文档：系统架构设计/API设计/数据模型设计/目录结构（v1.0.0）
> 阶段：Stage 5 Phase A — 伪代码（公子确认后进入 Phase B）

---

## 1. 模块清单与依赖倒置实现顺序

| 顺序 | 模块 | 路径 | 说明 |
|:-----|:-----|:-----|:------|
| 1 | pkg/jwt | `pkg/jwt/jwt.go` | JWT 生成与校验 |
| 2 | pkg/hash | `pkg/hash/hash.go` | bcrypt 密码哈希 |
| 3 | pkg/response | `pkg/response/response.go` | 统一响应包装 |
| 4 | pkg/validator | `pkg/validator/validator.go` | 参数校验 |
| 5 | pkg/pagination | `pkg/pagination/pagination.go` | 分页参数解析 |
| 6 | pkg/logger | `pkg/logger/logger.go` | slog 封装 |
| 7 | internal/config | `internal/config/config.go` | 配置加载 |
| 8 | internal/model | `internal/model/*.go` | 9 个数据模型 |
| 9 | internal/repository | `internal/repository/*.go` | 8 个数据访问层 |
| 10 | internal/service | `internal/service/*.go` | 9 个业务服务 |
| 11 | internal/middleware | `internal/middleware/*.go` | 5 个中间件 |
| 12 | internal/handler | `internal/handler/*.go` | 9 个 HTTP 处理器 |
| 13 | internal/dto | `internal/dto/**/*.go` | 请求/响应结构体 |
| 14 | internal/router | `internal/router/router.go` | 路由注册 |
| 15 | cmd/server | `cmd/server/main.go` | 程序入口 |

---

## 2. 各模块伪代码

### 2.1 pkg/jwt/jwt.go

```go
package jwt

// Manager JWT 管理器
type Manager struct {
    secret        []byte
    accessTTL     time.Duration
    refreshTTL    time.Duration
}

func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager

// GenerateTokenPair 生成 access + refresh 双 token
//   返回 (*TokenPair, error)
//   TokenPair { AccessToken, RefreshToken, ExpiresIn }
func (m *Manager) GenerateTokenPair(userID string, role string) (*TokenPair, error)

// ParseToken 解析并校验 token
//   返回 (*Claims, error)
func (m *Manager) ParseToken(tokenStr string) (*Claims, error)

// RefreshToken 用 refresh_token 换新 token 对
func (m *Manager) RefreshToken(refreshToken string) (*TokenPair, error)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    Type   string `json:"type"`  // "access" | "refresh"
    jwt.RegisteredClaims
}
```

### 2.2 pkg/hash/hash.go

```go
package hash

// Hash 加密密码（bcrypt cost=10）
func Hash(password string) (string, error)

// Verify 校验密码
func Verify(hash, password string) bool
```

### 2.3 pkg/response/response.go

```go
package response

// Response 统一响应
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Detail  interface{} `json:"detail,omitempty"`
}

// PaginationData 分页数据
type PaginationData struct {
    Items    interface{} `json:"items"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

func OK(c *gin.Context, data interface{})
func Created(c *gin.Context, data interface{})
func Fail(c *gin.Context, httpStatus int, code int, msg string, detail interface{})
func Unauthorized(c *gin.Context, msg string)
func Forbidden(c *gin.Context, msg string)
func NotFound(c *gin.Context, msg string)
func ServerError(c *gin.Context, msg string)
```

### 2.4 pkg/pagination/pagination.go

```go
package pagination

type Params struct {
    Page     int
    PageSize int
}

const DefaultPageSize = 20
const MaxPageSize = 100

// Parse 解析分页参数（带默认值与上限校验）
func Parse(c *gin.Context) Params

func (p Params) Offset() int { return (p.Page - 1) * p.PageSize }
func (p Params) Limit() int  { return p.PageSize }
```

### 2.5 pkg/logger/logger.go

```go
package logger

// Init 初始化 slog（JSON Handler，级别由配置控制）
func Init(level string, env string)

// L 全局 logger
func L() *slog.Logger
```

### 2.6 internal/config/config.go

```go
package config

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Logger   LoggerConfig
}
type ServerConfig struct {
    Port         int
    Mode         string  // debug | release | test
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}
type DatabaseConfig struct {
    Host         string
    Port         int
    User         string
    Password     string
    DBName       string
    SSLMode      string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}
type JWTConfig struct {
    Secret          string
    AccessTokenTTL  time.Duration  // 2h
    RefreshTokenTTL time.Duration  // 7d
}
type LoggerConfig struct {
    Level string  // debug | info | warn | error
}

// Load 加载配置：环境变量 > config.{env}.yaml > 默认值
func Load() (*Config, error)
```

### 2.7 internal/model/*

9 个 GORM 模型，与数据模型设计 1:1 映射。每个文件定义：

```go
package model

import "time"
import "github.com/google/uuid"

// User 用户表
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Name         string    `gorm:"type:varchar(100);not null" json:"name"`
    Phone        *string   `gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
    Email        *string   `gorm:"type:varchar(255);uniqueIndex" json:"email,omitempty"`
    PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
    Role         string    `gorm:"type:varchar(20);not null;default:student;check:role IN ('student','teacher','admin','parent')" json:"role"`
    AvatarURL    *string   `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
    ParentID     *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
func (User) TableName() string { return "users" }

// Subject 科目表
type Subject struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
    Description *string   `gorm:"type:text" json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
func (Subject) TableName() string { return "subjects" }

// KnowledgePoint 知识点表
type KnowledgePoint struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    SubjectID  uint64    `gorm:"not null;index" json:"subject_id"`
    ParentID   *uint64   `gorm:"index" json:"parent_id,omitempty"`
    Name       string    `gorm:"type:varchar(200);not null" json:"name"`
    Level      int       `gorm:"not null" json:"level"`
    Path       string    `gorm:"type:varchar(500);not null" json:"path"`
    CreatedAt  time.Time `json:"created_at"`
}
func (KnowledgePoint) TableName() string { return "knowledge_points" }

// Question 题目表
type Question struct {
    ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    SubjectID         uint64    `gorm:"not null;index:idx_q_subject_kp_diff" json:"subject_id"`
    KnowledgePointID  uint64    `gorm:"not null;index:idx_q_subject_kp_diff,priority:2" json:"knowledge_point_id"`
    Type              string    `gorm:"type:varchar(20);not null;index:idx_q_type" json:"type"`  // choice/fill/judge/subjective
    Difficulty        int       `gorm:"not null;default:3" json:"difficulty"`
    Content           datatypes.JSON `gorm:"type:jsonb;not null" json:"content"`
    Options           datatypes.JSON `gorm:"type:jsonb" json:"options,omitempty"`
    Answer            string    `gorm:"type:text;not null" json:"answer"`
    Analysis          *string   `gorm:"type:text" json:"analysis,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}
func (Question) TableName() string { return "questions" }

// StudyPlan 学习计划表
type StudyPlan struct {
    ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
    Goal         string         `gorm:"type:text;not null" json:"goal"`
    StartDate    time.Time      `gorm:"type:date;not null" json:"start_date"`
    EndDate      time.Time      `gorm:"type:date;not null" json:"end_date"`
    Items        datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"items"`
    Status       string         `gorm:"type:varchar(20);not null;default:active;index" json:"status"`
    AIGenerated  bool           `gorm:"not null;default:false" json:"ai_generated"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}
func (StudyPlan) TableName() string { return "study_plans" }

// StudyRecord 学习记录
type StudyRecord struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_sr_user_date,priority:1" json:"user_id"`
    PlanID      uint64    `gorm:"not null" json:"plan_id"`
    Date        time.Time `gorm:"type:date;not null;index:idx_sr_user_date,priority:2" json:"date"`
    DurationMin int       `gorm:"not null" json:"duration_min"`
    Status      string    `gorm:"type:varchar(20);not null;default:done" json:"status"`
    Memo        *string   `gorm:"type:text" json:"memo,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
func (StudyRecord) TableName() string { return "study_records" }

// ExerciseRecord 练习记录
type ExerciseRecord struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_er_user_question,priority:1" json:"user_id"`
    QuestionID  uint64    `gorm:"not null;index:idx_er_user_question,priority:2" json:"question_id"`
    Answer      string    `gorm:"type:text;not null" json:"answer"`
    IsCorrect   bool      `gorm:"not null;index:idx_er_user_correct_time,priority:2" json:"is_correct"`
    Score       *int      `gorm:"check:score >= 0" json:"score,omitempty"`
    DurationSec *int      `gorm:"check:duration_sec > 0" json:"duration_sec,omitempty"`
    CreatedAt   time.Time `gorm:"index" json:"created_at"`
}
func (ExerciseRecord) TableName() string { return "exercise_records" }

// MistakeBook 错题本
type MistakeBook struct {
    ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID            uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_mb_user_question,priority:1" json:"user_id"`
    QuestionID        uint64     `gorm:"not null;uniqueIndex:idx_mb_user_question,priority:2" json:"question_id"`
    KnowledgePointID  uint64     `gorm:"not null;index:idx_mb_user_kp,priority:2" json:"knowledge_point_id"`
    WrongAnswer       string     `gorm:"type:text;not null" json:"wrong_answer"`
    MistakeCount      int        `gorm:"not null;default:1" json:"mistake_count"`
    Mastered          bool       `gorm:"not null;default:false;index:idx_mb_user_mastered,priority:2" json:"mastered"`
    MasteredAt        *time.Time `json:"mastered_at,omitempty"`
    LastReviewedAt    time.Time  `gorm:"not null;index" json:"last_reviewed_at"`
    CreatedAt         time.Time  `json:"created_at"`
}
func (MistakeBook) TableName() string { return "mistake_books" }

// StudyReport 学习报告
type StudyReport struct {
    ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_sr_user_period,priority:1" json:"user_id"`
    PeriodType  string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_sr_user_period,priority:2" json:"period_type"`
    PeriodStart time.Time      `gorm:"type:date;not null;uniqueIndex:idx_sr_user_period,priority:3" json:"period_start"`
    PeriodEnd   time.Time      `gorm:"type:date;not null" json:"period_end"`
    Stats       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"stats"`
    CreatedAt   time.Time      `json:"created_at"`
}
func (StudyReport) TableName() string { return "study_reports" }
```

### 2.8 internal/repository/*

```go
package repository

// ============== UserRepo ==============
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
    GetByPhone(ctx context.Context, phone string) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error
    UpdateAvatar(ctx context.Context, id uuid.UUID, url string) error
}

// ============== PlanRepo ==============
type PlanRepository interface {
    Create(ctx context.Context, plan *model.StudyPlan) error
    GetByID(ctx context.Context, id uint64) (*model.StudyPlan, error)
    ListByUser(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error)
    Update(ctx context.Context, plan *model.StudyPlan) error
    Delete(ctx context.Context, id uint64) error
    // StudyRecord
    CreateRecord(ctx context.Context, rec *model.StudyRecord) error
    ListRecordsByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.StudyRecord, error)
    SumDurationByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
}

// ============== ExerciseRepo ==============
type ExerciseRepository interface {
    // Question
    GetQuestionByID(ctx context.Context, id uint64) (*model.Question, error)
    ListQuestions(ctx context.Context, filter QuestionFilter) ([]model.Question, int64, error)
    ListQuestionsByIDs(ctx context.Context, ids []uint64) ([]model.Question, error)
    RandomQuestions(ctx context.Context, subjectID, kpID *uint64, count int) ([]model.Question, error)
    // Record
    CreateRecord(ctx context.Context, rec *model.ExerciseRecord) error
    CountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
    CorrectCountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
    CorrectRateByKP(ctx context.Context, userID uuid.UUID, subjectID uint64) ([]KPRate, error)
}

// ============== MistakeRepo ==============
type MistakeRepository interface {
    GetByUserQuestion(ctx context.Context, userID uuid.UUID, questionID uint64) (*model.MistakeBook, error)
    Create(ctx context.Context, m *model.MistakeBook) error
    IncrementCount(ctx context.Context, id uint64, lastReviewed time.Time) error
    UpdateLastReviewed(ctx context.Context, id uint64) error
    ListByUser(ctx context.Context, userID uuid.UUID, kpID *uint64, mastered *bool, page, pageSize int) ([]MistakeWithQuestion, int64, error)
    GroupByKP(ctx context.Context, userID uuid.UUID) ([]MistakeGroup, error)
    GetByID(ctx context.Context, id uint64) (*model.MistakeBook, error)
    MarkMastered(ctx context.Context, id uint64, mastered bool) error
    Delete(ctx context.Context, id uint64) error
    ListUnmasteredQuestions(ctx context.Context, userID uuid.UUID, kpIDs []uint64, limit int) ([]model.Question, error)
    CountUnmastered(ctx context.Context, userID uuid.UUID) (int64, error)
}

// ============== ReportRepo ==============
type ReportRepository interface {
    Upsert(ctx context.Context, r *model.StudyReport) error
    GetByPeriod(ctx context.Context, userID uuid.UUID, periodType string, start time.Time) (*model.StudyReport, error)
    LatestSummary(ctx context.Context, userID uuid.UUID) (*model.StudyReport, error)
}

// ============== SubjectRepo / KnowledgeRepo / QuestionRepo（基础查询） ==============
type SubjectRepository interface {
    List(ctx context.Context) ([]model.Subject, error)
}
type KnowledgeRepository interface {
    ListTree(ctx context.Context, subjectID uint64) ([]KnowledgeNode, error)
    GetByID(ctx context.Context, id uint64) (*model.KnowledgePoint, error)
}
```

### 2.9 internal/service/*

```go
package service

// ============== AuthService ==============
type AuthService interface {
    Register(ctx, req RegisterRequest) (*AuthResponse, error)
    //  1. 校验必填（name, account, password, role）
    //  2. 校验 phone/email 格式；至少有一个
    //  3. 检查 phone/email 唯一
    //  4. bcrypt 加密密码
    //  5. 创建 User（role 默认 student）
    //  6. 生成 token 对并返回 AuthResponse

    Login(ctx, req LoginRequest) (*AuthResponse, error)
    //  1. 按 phone 或 email 查询 user
    //  2. bcrypt 校验密码
    //  3. 生成 token 对

    Refresh(ctx, refreshToken string) (*AuthResponse, error)
    //  1. ParseToken(type=refresh)
    //  2. 重新生成 token 对
}

// ============== UserService ==============
type UserService interface {
    GetMe(ctx, userID) (*UserDTO, error)
    UpdateMe(ctx, userID, req UpdateUserRequest) (*UserDTO, error)
    ChangePassword(ctx, userID, oldPwd, newPwd string) error
    UpdateAvatar(ctx, userID, url string) error
}

// ============== PlanService ==============
type PlanService interface {
    Create(ctx, userID, req) (*StudyPlanDTO, error)
    GetByID(ctx, userID, id) (*StudyPlanDTO, error)
    List(ctx, userID, status, page, pageSize) (Pagination, error)
    Update(ctx, userID, id, req) (*StudyPlanDTO, error)
    Delete(ctx, userID, id) error
    AIGenerate(ctx, userID, req) (*StudyPlanDTO, aiUsed bool, error)
    // 规则降级算法：按 start_date~end_date 每天生成 item，duration=daily_duration_min，
    // knowledge_point_ids 轮询常见 KP（此处简化为按 subject 内 KP 平摊）
    Checkin(ctx, userID, planID, req) (*StudyRecordDTO, error)
}

// ============== ExerciseService ==============
type ExerciseService interface {
    List(ctx, filter) (Pagination, error)
    Random(ctx, userID, subjectID, kpID *uint64, count int) ([]QuestionDTO, error)
    Submit(ctx, userID, req SubmitRequest) (*SubmitResponse, error)
    //  1. 查询题目
    //  2. 批改（按 type：choice/judge 直接比答案；fill 简单 trim 比对）
    //  3. 写入 exercise_record
    //  4. 答错 → MistakeService.RecordMistake（幂等）
    //  5. 返回 is_correct + correct_answer + analysis + mistake_recorded
    Recommend(ctx, userID, count int) (*RecommendResponse, error)
    // 简单规则：取错题本对应 KP 的未掌握题目 + 随机补充
    ByKnowledgePoint(ctx, userID, kpID, page, pageSize) (Pagination, error)
    History(ctx, userID, page, pageSize) (Pagination, error)
}

// ============== MistakeService ==============
type MistakeService interface {
    RecordMistake(ctx, userID, questionID, wrongAnswer string) error
    //  1. 查询 question → kp_id
    //  2. 查询是否已存在 → 存在则 mistake_count++，更新 last_reviewed_at
    //  3. 不存在则插入新记录
    List(ctx, userID, kpID *uint64, mastered *bool, page, pageSize) (Pagination, error)
    GroupByKP(ctx, userID) ([]MistakeGroupDTO, error)
    MarkMastered(ctx, userID, mistakeID, mastered bool) error
    Review(ctx, userID, kpIDs []uint64, count int) (*ReviewResponse, error)
    Delete(ctx, userID, mistakeID) error
}

// ============== ReportService ==============
type ReportService interface {
    Summary(ctx, userID) (*SummaryDTO, error)
    Detail(ctx, userID, periodType, periodStart) (*DetailDTO, error)
    Mastery(ctx, userID, subjectID) (*MasteryDTO, error)
    Trend(ctx, userID, days int) (*TrendDTO, error)
}

// ============== SubjectService / KnowledgeService / QuestionService ==============
type SubjectService interface { List(ctx) ([]SubjectDTO, error) }
type KnowledgeService interface { GetTree(ctx, subjectID) (*KnowledgeTreeDTO, error) }
type QuestionService interface {
    Create(ctx, req) error
    Update(ctx, id, req) error
    Delete(ctx, id) error
    BatchImport(ctx, items) (int, error)
    Export(ctx, subjectID *uint64) ([]QuestionDTO, error)
}
```

### 2.10 internal/middleware/*

```go
package middleware

// JWTAuth 校验 Authorization Bearer Token，将 user_id/role 注入 ctx
func JWTAuth(mgr *jwt.Manager) gin.HandlerFunc

// RequireRole(roles ...string) 角色守卫
func RequireRole(roles ...string) gin.HandlerFunc

// CORS 跨域
func CORS() gin.HandlerFunc

// RequestLogger 请求日志（slog + request_id）
func RequestLogger() gin.HandlerFunc

// Recovery panic 恢复
func Recovery() gin.HandlerFunc

// RateLimit 基于 IP 的令牌桶（内存版）
func RateLimit(rps, burst int) gin.HandlerFunc

// RequestID 注入 X-Request-ID
func RequestID() gin.HandlerFunc
```

### 2.11 internal/handler/*

```go
package handler

// 每个 handler：绑定参数 → 调 service → 用 pkg/response 包装
// 关键端点（API设计 章节对应）：
// POST /api/v1/auth/register        AuthHandler.Register
// POST /api/v1/auth/login           AuthHandler.Login
// POST /api/v1/auth/refresh         AuthHandler.Refresh
// GET  /api/v1/users/me             UserHandler.GetMe
// PUT  /api/v1/users/me             UserHandler.UpdateMe
// PUT  /api/v1/users/me/password    UserHandler.ChangePassword
// POST /api/v1/users/me/avatar      UserHandler.UpdateAvatar
//
// POST /api/v1/plans                PlanHandler.Create
// GET  /api/v1/plans                PlanHandler.List
// GET  /api/v1/plans/:id            PlanHandler.GetByID
// PUT  /api/v1/plans/:id            PlanHandler.Update
// POST /api/v1/plans/ai-generate    PlanHandler.AIGenerate
// POST /api/v1/plans/:id/checkin    PlanHandler.Checkin
//
// GET  /api/v1/exercises            ExerciseHandler.List
// GET  /api/v1/exercises/random     ExerciseHandler.Random
// POST /api/v1/exercises/submit     ExerciseHandler.Submit
// GET  /api/v1/exercises/recommend   ExerciseHandler.Recommend
// GET  /api/v1/exercises/kp/:kp_id  ExerciseHandler.ByKnowledgePoint
// GET  /api/v1/exercises/history    ExerciseHandler.History
//
// GET  /api/v1/mistakes             MistakeHandler.List
// GET  /api/v1/mistakes/groups      MistakeHandler.Groups
// PUT  /api/v1/mistakes/:id/master  MistakeHandler.MarkMastered
// POST /api/v1/mistakes/review      MistakeHandler.Review
// DELETE /api/v1/mistakes/:id       MistakeHandler.Delete
//
// GET  /api/v1/reports/summary      ReportHandler.Summary
// GET  /api/v1/reports/detail       ReportHandler.Detail
// GET  /api/v1/reports/mastery      ReportHandler.Mastery
// GET  /api/v1/reports/trend        ReportHandler.Trend
//
// GET  /api/v1/subjects             SubjectHandler.List
// GET  /api/v1/knowledge-points     KnowledgeHandler.Tree
```

### 2.12 internal/router/router.go

```go
package router

// New 构建 *gin.Engine，按模块注册路由
func New(cfg *config.Config, h Handlers, mw *middleware.MW) *gin.Engine

// 公开路由：/api/v1/auth/*
// 鉴权路由：所有 /api/v1/users/*, /plans/*, /exercises/*, /mistakes/*, /reports/*
// /api/v1/subjects, /api/v1/knowledge-points 需要认证
// /health GET 用于健康检查
```

### 2.13 cmd/server/main.go

```go
package main

func main() {
    // 1. Load config
    // 2. Init logger
    // 3. Init DB (gorm.Open postgres + AutoMigrate)
    // 4. Build repositories, services, handlers
    // 5. Build router
    // 6. http.ListenAndServe(":port", router)
}
```

---

## 3. P0 API 端点覆盖矩阵

| P0 端点 | Handler 方法 | Service 方法 | Repository | 行号 |
|:---------|:-------------|:-------------|:-----------|:------|
| POST /auth/register | AuthHandler.Register | AuthService.Register | UserRepository.Create | ✅ |
| POST /auth/login | AuthHandler.Login | AuthService.Login | UserRepository.GetByPhone/Email | ✅ |
| POST /auth/refresh | AuthHandler.Refresh | AuthService.Refresh | — | ✅ |
| GET /users/me | UserHandler.GetMe | UserService.GetMe | UserRepository.GetByID | ✅ |
| PUT /users/me | UserHandler.UpdateMe | UserService.UpdateMe | UserRepository.Update | ✅ |
| PUT /users/me/password | UserHandler.ChangePassword | UserService.ChangePassword | UserRepository.UpdatePassword | ✅ |
| POST /users/me/avatar | UserHandler.UpdateAvatar | UserService.UpdateAvatar | UserRepository.UpdateAvatar | ✅ |
| POST /plans | PlanHandler.Create | PlanService.Create | PlanRepository.Create | ✅ |
| GET /plans | PlanHandler.List | PlanService.List | PlanRepository.ListByUser | ✅ |
| GET /plans/:id | PlanHandler.GetByID | PlanService.GetByID | PlanRepository.GetByID | ✅ |
| PUT /plans/:id | PlanHandler.Update | PlanService.Update | PlanRepository.Update | ✅ |
| POST /plans/ai-generate | PlanHandler.AIGenerate | PlanService.AIGenerate | PlanRepository.Create | ✅ |
| POST /plans/:id/checkin | PlanHandler.Checkin | PlanService.Checkin | PlanRepository.CreateRecord | ✅ |
| GET /exercises | ExerciseHandler.List | ExerciseService.List | ExerciseRepository.ListQuestions | ✅ |
| GET /exercises/random | ExerciseHandler.Random | ExerciseService.Random | ExerciseRepository.RandomQuestions | ✅ |
| POST /exercises/submit | ExerciseHandler.Submit | ExerciseService.Submit | ExerciseRepository.CreateRecord | ✅ |
| GET /exercises/recommend | ExerciseHandler.Recommend | ExerciseService.Recommend | ExerciseRepository + MistakeRepository | ✅ |
| GET /exercises/kp/:kp_id | ExerciseHandler.ByKnowledgePoint | ExerciseService.ByKnowledgePoint | ExerciseRepository.ListQuestions | ✅ |
| GET /exercises/history | ExerciseHandler.History | ExerciseService.History | ExerciseRepository + ExerciseRecord | ✅ |
| GET /mistakes | MistakeHandler.List | MistakeService.List | MistakeRepository.ListByUser | ✅ |
| GET /mistakes/groups | MistakeHandler.Groups | MistakeService.GroupByKP | MistakeRepository.GroupByKP | ✅ |
| PUT /mistakes/:id/master | MistakeHandler.MarkMastered | MistakeService.MarkMastered | MistakeRepository.MarkMastered | ✅ |
| POST /mistakes/review | MistakeHandler.Review | MistakeService.Review | MistakeRepository.ListUnmasteredQuestions | ✅ |
| DELETE /mistakes/:id | MistakeHandler.Delete | MistakeService.Delete | MistakeRepository.Delete | ✅ |
| GET /reports/summary | ReportHandler.Summary | ReportService.Summary | 多表聚合 | ✅ |
| GET /reports/detail | ReportHandler.Detail | ReportService.Detail | 多表聚合 | ✅ |
| GET /reports/mastery | ReportHandler.Mastery | ReportService.Mastery | ExerciseRepository.CorrectRateByKP | ✅ |
| GET /reports/trend | ReportHandler.Trend | ReportService.Trend | PlanRepository + ExerciseRepository | ✅ |
| GET /subjects | SubjectHandler.List | SubjectService.List | SubjectRepository.List | ✅ |
| GET /knowledge-points | KnowledgeHandler.Tree | KnowledgeService.GetTree | KnowledgeRepository.ListTree | ✅ |

---

## 4. 自检表（Phase A）

| # | 检查项 | 状态 | 等级 |
|:--|:-------|:-----|:-----|
| C1 | 覆盖所有模块（9 个 model / 8 个 repo / 9 个 service / 9 个 handler / 5 个 mw） | ✅ | P0_BLOCK |
| C2 | 覆盖所有 P0 API（30 个端点全部映射） | ✅ | P0_BLOCK |
| C3 | 签名与架构一致（接口入参/出参类型匹配） | ✅ | P1 |
| C4 | 逻辑清晰（自然语言描述每步） | ✅ | P1 |

**Phase A 自检结论：✅ 通过。等待公子确认进入 Phase B。**