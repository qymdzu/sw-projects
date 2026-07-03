// Package router 提供 HTTP 路由注册。
package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smart-learning/internal/config"
	"smart-learning/internal/handler"
	"smart-learning/internal/middleware"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
	"smart-learning/pkg/crypto"
	"smart-learning/pkg/jwt"
	"smart-learning/pkg/logger"
)

// Handlers 集中持有所有 HTTP 处理器。
type Handlers struct {
	Auth      *handler.AuthHandler
	User      *handler.UserHandler
	Plan      *handler.PlanHandler
	Exercise  *handler.ExerciseHandler
	Mistake   *handler.MistakeHandler
	Report    *handler.ReportHandler
	Subject   *handler.SubjectHandler
	Knowledge *handler.KnowledgeHandler
	Setting   *handler.SettingHandler // Phase B 新增：模型配置
}

// BuildHandlers 根据 DB 与 cfg 构造所有 handler（共享 JWT 密钥）。
func BuildHandlers(db *gorm.DB, cfg *config.Config) (*Handlers, error) {
	jwtMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	userRepo := repository.NewUserRepository(db)
	planRepo := repository.NewPlanRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	mistakeRepo := repository.NewMistakeRepository(db)
	reportRepo := repository.NewReportRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	settingRepo := repository.NewSettingRepository(db) // Phase B 新增

	authSvc := service.NewAuthService(userRepo, jwtMgr)
	userSvc := service.NewUserService(userRepo)
	planSvc := service.NewPlanService(planRepo)
	exerciseSvc := service.NewExerciseService(exerciseRepo, mistakeRepo)
	mistakeSvc := service.NewMistakeService(mistakeRepo, exerciseRepo)
	reportSvc := service.NewReportService(planRepo, exerciseRepo, mistakeRepo, subjectRepo, knowledgeRepo, reportRepo)
	subjectSvc := service.NewSubjectService(subjectRepo)
	knowledgeSvc := service.NewKnowledgeService(knowledgeRepo)

	// Phase B：模型配置 AES-GCM 密钥（派生自 JWT_SECRET，与签名密钥分离但同源）
	cipher, err := crypto.NewAESGCM(cfg.JWT.Secret)
	if err != nil {
		return nil, nil // 上层 main 应当 panic
	}
	settingSvc := service.NewSettingService(settingRepo, cipher) // Phase B 新增

	return &Handlers{
		Auth:      handler.NewAuthHandler(authSvc),
		User:      handler.NewUserHandler(userSvc),
		Plan:      handler.NewPlanHandler(planSvc),
		Exercise:  handler.NewExerciseHandler(exerciseSvc),
		Mistake:   handler.NewMistakeHandler(mistakeSvc),
		Report:    handler.NewReportHandler(reportSvc),
		Subject:   handler.NewSubjectHandler(subjectSvc),
		Knowledge: handler.NewKnowledgeHandler(knowledgeSvc),
		Setting:   handler.NewSettingHandler(settingSvc), // Phase B 新增
	}, nil
}

// New 构建 *gin.Engine 并注册所有路由。
func New(cfg *config.Config, h *Handlers) *gin.Engine {
	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORSWithConfig(&cfg.CORS)) // Phase B P1-02：白名单 CORS

	jwtMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "ts": time.Now().Unix()})
	})

	api := r.Group("/api/v1")
	{
		// 公开路由：认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.Refresh)
		}

		// 鉴权路由
		secured := api.Group("")
		secured.Use(middleware.JWTAuth(jwtMgr))
		{
			// 用户
			users := secured.Group("/users")
			{
				users.GET("/me", h.User.GetMe)
				users.PUT("/me", h.User.UpdateMe)
				users.PUT("/me/password", h.User.ChangePassword)
				users.POST("/me/avatar", h.User.UpdateAvatar)
			}

			// 学习计划
			plans := secured.Group("/plans")
			{
				plans.GET("", h.Plan.List)
				plans.POST("", h.Plan.Create)
				plans.POST("/ai-generate", h.Plan.AIGenerate)
				plans.GET("/:id", h.Plan.GetByID)
				plans.PUT("/:id", h.Plan.Update)
				plans.DELETE("/:id", h.Plan.Delete)
				plans.POST("/:id/checkin", h.Plan.Checkin)
			}

			// 练习
			exercises := secured.Group("/exercises")
			{
				exercises.GET("", h.Exercise.List)
				exercises.GET("/random", h.Exercise.Random)
				exercises.POST("/submit", h.Exercise.Submit)
				exercises.GET("/recommend", h.Exercise.Recommend)
				exercises.GET("/knowledge-points/:kp_id", h.Exercise.ByKnowledgePoint)
				exercises.GET("/history", h.Exercise.History)
			}

			// 错题
			mistakes := secured.Group("/mistakes")
			{
				mistakes.GET("", h.Mistake.List)
				mistakes.GET("/groups", h.Mistake.Groups)
				mistakes.POST("/review", h.Mistake.Review)
				mistakes.PUT("/:id/master", h.Mistake.MarkMastered)
				mistakes.DELETE("/:id", h.Mistake.Delete)
			}

			// 看板
			reports := secured.Group("/reports")
			{
				reports.GET("/summary", h.Report.Summary)
				reports.GET("/detail", h.Report.Detail)
				reports.GET("/mastery", h.Report.Mastery)
				reports.GET("/trend", h.Report.Trend)
			}

			// 科目 & 知识点
			api.GET("/subjects", h.Subject.List)
			api.GET("/knowledge-points", h.Knowledge.Tree)

			// Phase B 新增：模型配置（4 个端点）
			settings := secured.Group("/settings")
			{
				settings.POST("/model", h.Setting.CreateOrUpdate)
				settings.GET("/model", h.Setting.GetActive)
				settings.PUT("/model", h.Setting.Update)
				settings.DELETE("/model", h.Setting.Delete)
			}
		}
	}

	logger.L().Info("router initialized", "mode", cfg.Server.Mode, "port", cfg.Server.Port)
	return r
}