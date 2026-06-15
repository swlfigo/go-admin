package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type fakeEnforcer struct{ allow bool }

func (f fakeEnforcer) Can(userID uint, perm string) bool { return f.allow }

func TestRequirePerm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	run := func(allow bool, setUser bool) int {
		r := gin.New()
		r.GET("/x", func(c *gin.Context) {
			if setUser {
				c.Set(CtxUserID, uint(1))
			}
			c.Next()
		}, RequirePerm(fakeEnforcer{allow: allow}, "p:x"), func(c *gin.Context) { c.Status(200) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
		return w.Code
	}

	require.Equal(t, 200, run(true, true))  // 有权限
	require.Equal(t, 403, run(false, true)) // 无权限
	require.Equal(t, 401, run(true, false)) // 未认证（无 userID）
}
