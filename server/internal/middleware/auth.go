package middleware

import (
	"net/http"
	"strings"

	"go-admin/pkg/jwt"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID   = "ctx_user_id"
	CtxUsername = "ctx_username"
	CtxJTI      = "ctx_jti"
)

// SessionChecker 抽象会话存在性校验，便于测试与未来换 redis。
type SessionChecker interface {
	Exists(jti string) bool
	Touch(jti string) error
}

func Auth(m *jwt.Manager, sessions SessionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		claims, err := m.ParseAccess(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if !sessions.Exists(claims.ID) {
			response.Fail(c, http.StatusUnauthorized, "登录已失效或被强制下线")
			c.Abort()
			return
		}
		_ = sessions.Touch(claims.ID)
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxJTI, claims.ID)
		c.Next()
	}
}
