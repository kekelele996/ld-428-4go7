package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

type visitor struct {
	count    int
	windowAt time.Time
}

// RateLimiter 每 IP 每分钟最多 limit 次请求。
func RateLimiter(limit int) gin.HandlerFunc {
	return rateLimiter(limit, time.Minute)
}

// RateLimiterWindow 使用自定义窗口的敏感接口限流器。
func RateLimiterWindow(limit int, window time.Duration) gin.HandlerFunc {
	if window <= 0 {
		window = time.Minute
	}
	return rateLimiter(limit, window)
}

func rateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 120
	}
	var mu sync.Mutex
	visitors := map[string]*visitor{}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, ok := visitors[ip]
		now := time.Now()
		if !ok || now.Sub(v.windowAt) > window {
			v = &visitor{windowAt: now}
			visitors[ip] = v
		}
		v.count++
		tooMany := v.count > limit
		mu.Unlock()
		if tooMany {
			util.Fail(c, http.StatusTooManyRequests, constants.CodeRateLimited, constants.MsgRateLimited)
			c.Abort()
			return
		}
		c.Next()
	}
}
