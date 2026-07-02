package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"smart-learning/pkg/logger"
	"smart-learning/pkg/response"
)

// Recovery panic 恢复中间件。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L().Error("panic recovered",
					"error", fmt.Sprintf("%v", r),
					"stack", string(debug.Stack()),
				)
				response.ServerError(c, "服务器内部错误")
			}
		}()
		c.Next()
	}
}