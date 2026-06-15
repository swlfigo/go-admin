package middleware

import (
	"net/http"

	"go-admin/internal/auth"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequirePerm 在已认证（Auth 中间件之后）基础上做接口级授权。
func RequirePerm(e auth.Enforcer, perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint(CtxUserID)
		if uid == 0 {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		if !e.Can(uid, perm) {
			response.Fail(c, http.StatusForbidden, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
