package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/middleware"
	"go-admin/internal/repository"
	"go-admin/internal/service"
	"go-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestKickEndToEnd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := handlerTestDB(t) // reuses helper from auth_handler_http_test.go

	jm := jwt.New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	userRepo := repository.NewUserRepo(db)
	sessRepo := repository.NewSessionRepo(db)

	authSvc := service.NewAuthService(userRepo, sessRepo, jm, config.LoginConfig{MaxFailCount: 5, LockMinutes: 10})
	onlineH := NewOnlineHandler(sessRepo)

	checker := testSessionAdapter{repo: sessRepo} // reuses adapter from auth_handler_http_test.go

	// Log in as admin to create a real session + token
	loginOut, err := authSvc.Login("admin", "admin123", service.LoginMeta{IP: "127.0.0.1"})
	require.NoError(t, err)
	token := loginOut.AccessToken

	// Parse the access token to extract its JTI (session ID)
	claims, err := jm.ParseAccess(token)
	require.NoError(t, err)
	jti := claims.ID

	// Cleanup: remove the session we created
	t.Cleanup(func() {
		sessRepo.Delete(jti)
	})

	// Build router: Auth-protected ping + kick endpoint
	r := gin.New()
	r.GET("/ping", middleware.Auth(jm, checker), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.DELETE("/online/:id", onlineH.Kick)

	// Assert that the token works (200) before kick
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusOK, w1.Code, "token should be valid before kick")

	// Kick the session
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/online/%s", jti), nil)
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusOK, w2.Code, "kick should succeed")

	// Assert the SAME token now returns 401
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w3, req3)
	require.Equal(t, http.StatusUnauthorized, w3.Code, "kicked token should be rejected")
}
