package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/config"
)

// CORS 跨域中间件（Phase B P1-02 重写：白名单替换 *）。
//
// 行为：
//   - AllowOrigins 通过 cfg.CORS.AllowOrigins 注入，逗号分隔
//   - 当请求 Origin 在白名单内时回写对应的 Access-Control-Allow-Origin
//   - 不在白名单时**不**回写该头，浏览器会按 CORS 规范阻断
//   - 允许携带 credentials（Cookie / Authorization）
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfgAny, exists := c.Get("cors_allow_origins")
		var allowed []string
		if exists {
			if list, ok := cfgAny.([]string); ok {
				allowed = list
			}
		}
		origin := c.GetHeader("Origin")
		if origin != "" && isAllowedOrigin(origin, allowed) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// isAllowedOrigin 检查 origin 是否在白名单内（含通配符 *）。
func isAllowedOrigin(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}
		if strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}

// CORSWithConfig 返回带显式配置的 CORS 中间件（推荐用法，Phase B）。
func CORSWithConfig(cfg *config.CORSConfig) gin.HandlerFunc {
	allow := cfg.AllowOrigins
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isAllowedOrigin(origin, allow) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Vary", "Origin")
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}