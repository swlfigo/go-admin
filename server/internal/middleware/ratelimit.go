package middleware

import (
	"net/http"
	"sync"

	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitByIP 每个 IP 每分钟最多 perMinute 次（突发=perMinute）。
func RateLimitByIP(perMinute int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)
	// 每分钟补充 perMinute 个令牌 => 速率 perMinute/60 每秒
	every := rate.Limit(float64(perMinute) / 60.0)

	get := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		l, ok := limiters[ip]
		if !ok {
			l = rate.NewLimiter(every, perMinute)
			limiters[ip] = l
		}
		return l
	}

	return func(c *gin.Context) {
		if !get(c.ClientIP()).Allow() {
			response.Fail(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
