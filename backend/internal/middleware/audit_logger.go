package middleware

import (
	"context"
	"log/slog"

	"github.com/artvault/artvault/internal/service"
	"github.com/gin-gonic/gin"
)

// AuditLogger 写操作审计日志。
func AuditLogger(auditSvc *service.AuditLogService, logger *slog.Logger) gin.HandlerFunc {
	writeMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true, "DELETE": true}
	return func(c *gin.Context) {
		c.Next()
		if !writeMethods[c.Request.Method] || c.Writer.Status() >= 400 {
			return
		}
		userID, _ := c.Get("userID")
		uid, _ := userID.(string)
		auditSvc.Write(context.Background(), uid, c.Request.Method+" "+c.Request.URL.Path, "request", "", "")
	}
}
