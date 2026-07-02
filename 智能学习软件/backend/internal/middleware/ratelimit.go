package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"smart-learning/pkg/response"
)

// ipBucket 是基于 IP 的简单令牌桶。
type ipBucket struct {
	tokens     float64
	lastRefill time.Time
}

// RateLimit 内存版基于 IP 的令牌桶限流。
// rps: 每秒补充令牌数；burst: 桶容量上限。
func RateLimit(rps, burst int) gin.HandlerFunc {
	mu := sync.Mutex{}
	buckets := make(map[string]*ipBucket)
	capacity := float64(burst)
	rate := float64(rps)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &ipBucket{tokens: capacity, lastRefill: now}
			buckets[ip] = b
		}
		elapsed := now.Sub(b.lastRefill).Seconds()
		b.tokens += elapsed * rate
		if b.tokens > capacity {
			b.tokens = capacity
		}
		b.lastRefill = now
		if b.tokens < 1 {
			mu.Unlock()
			response.Fail(c, 429, response.CodeRateLimited, "请求过于频繁", nil)
			return
		}
		b.tokens--
		mu.Unlock()
		c.Next()
	}
}