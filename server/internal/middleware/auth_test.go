package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-admin/internal/model"
	"go-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 假 session 校验器：内存集合
type fakeSessions struct{ live map[string]bool }

func (f *fakeSessions) Exists(jti string) bool { return f.live[jti] }
func (f *fakeSessions) Touch(jti string) error { return nil }

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := jwt.New("test-secret-at-least-16", 30*time.Minute, time.Hour)
	tok, claims, _ := m.IssueAccess(7, "admin")

	sessions := &fakeSessions{live: map[string]bool{claims.ID: true}}
	r := gin.New()
	r.GET("/p", Auth(m, sessions), func(c *gin.Context) {
		uid := c.GetUint(CtxUserID)
		c.JSON(200, gin.H{"uid": uid})
	})

	// 合法 token + 活跃会话 → 200
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	// 会话被踢（删除）→ 401
	delete(sessions.live, claims.ID)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/p", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r.ServeHTTP(w2, req2)
	require.Equal(t, 401, w2.Code)

	// 无 token → 401
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/p", nil)
	r.ServeHTTP(w3, req3)
	require.Equal(t, 401, w3.Code)

	_ = model.User{} // 保持 import
}
