package middleware

import (
	"net/http"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// RequireRole 角色校验。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == nil || !allowed[role.(string)] {
			util.Fail(c, http.StatusForbidden, constants.CodeForbidden, constants.MsgForbidden)
			return
		}
		c.Next()
	}
}
