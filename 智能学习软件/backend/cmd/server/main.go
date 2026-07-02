// Package main 是智能学习软件后端服务的入口。
//
// 启动流程：
//  1. 加载配置
//  2. 初始化 logger
//  3. 初始化数据库连接 + AutoMigrate
//  4. 构建路由
//  5. 启动 HTTP 服务
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"smart-learning/internal/config"
	"smart-learning/internal/model"
	"smart-learning/internal/router"
	"smart-learning/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	logger.Init(cfg.Logger.Level, cfg.Logger.Env)
	logger.L().Info("config loaded", "port", cfg.Server.Port, "mode", cfg.Server.Mode)

	db, err := initDB(cfg)
	if err != nil {
		logger.L().Error("init db failed", "err", err)
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// AutoMigrate（MVP 阶段；生产建议使用 golang-migrate 版本化迁移）
	if err := db.AutoMigrate(
		&model.User{},
		&model.Subject{},
		&model.KnowledgePoint{},
		&model.Question{},
		&model.StudyPlan{},
		&model.StudyRecord{},
		&model.ExerciseRecord{},
		&model.MistakeBook{},
		&model.StudyReport{},
	); err != nil {
		logger.L().Error("automigrate failed", "err", err)
		log.Fatalf("自动迁移失败: %v", err)
	}
	logger.L().Info("automigrate completed")

	handlers := router.BuildHandlers(db, cfg)
	engine := router.New(cfg, handlers)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 优雅关闭
	go func() {
		logger.L().Info("http server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Error("http server error", "err", err)
			log.Fatalf("HTTP 服务异常: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.L().Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error("server forced to shutdown", "err", err)
	}
}

// initDB 初始化 PostgreSQL 数据库连接。
func initDB(cfg *config.Config) (*gorm.DB, error) {
	fmt.Printf("[DEBUG] cfg.Database: Host=%s Port=%d User=%s Password=%s DBName=%s SSLMode=%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 PostgreSQL 失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)
	return db, nil
}