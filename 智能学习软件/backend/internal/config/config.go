// Package config 提供应用配置加载。
//
// 加载优先级：环境变量 > config.{env}.yaml > 默认值。
// MVP 阶段使用纯环境变量方式以减少外部依赖（详见 docs/design/目录结构.md 第 5.1 节）。
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 是应用配置总入口。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

// ServerConfig 是 HTTP 服务配置。
type ServerConfig struct {
	Port         int
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig 是 PostgreSQL 数据库配置。
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

// JWTConfig 是 JWT 配置。
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LoggerConfig 是日志配置。
type LoggerConfig struct {
	Level string
	Env   string
}

// Load 从环境变量加载配置，应用默认值。
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			Mode:         getEnvStr("SERVER_MODE", "debug"),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnvStr("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			User:         getEnvStr("DB_USER", "smart_learning"),
			Password:     getEnvStr("DB_PASSWORD", ""),
			DBName:       getEnvStr("DB_NAME", "smart_learning"),
			SSLMode:      getEnvStr("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 30*time.Minute),
		},
		JWT: JWTConfig{
			Secret:          getEnvStr("JWT_SECRET", "dev-secret-please-replace-in-production"),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TTL", 2*time.Hour),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		Logger: LoggerConfig{
			Level: getEnvStr("LOG_LEVEL", "info"),
			Env:   getEnvStr("APP_ENV", "dev"),
		},
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET 不能为空")
	}
	return cfg, nil
}

func getEnvStr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}