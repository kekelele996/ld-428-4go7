package middleware

import (
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// RequestLogger 结构化请求日志。
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = util.NewRequestID()
		}
		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		start := time.Now()
		c.Next()
		logger.Info(constants.LogRequest,
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
