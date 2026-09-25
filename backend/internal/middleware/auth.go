package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/artvault/artvault/internal/config"
	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth JWT 认证中间件，将用户信息注入 Context。
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		claims, err := util.ParseToken(cfg.JWTSecret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			code := constants.CodeUnauthorized
			if errors.Is(err, jwt.ErrTokenExpired) {
				code = constants.CodeTokenExpired
			}
			util.Fail(c, http.StatusUnauthorized, code, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
