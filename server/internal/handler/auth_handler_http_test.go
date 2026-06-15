package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-admin/internal/auth"
	"go-admin/internal/config"
	"go-admin/internal/database"
	"go-admin/internal/middleware"
	"go-admin/internal/repository"
	"go-admin/internal/service"
	"go-admin/internal/testsupport"
	"go-admin/pkg/captcha"
	"go-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func handlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, database.MigrateAndSeed(db))
	return db
}

type testSessionAdapter struct{ repo *repository.SessionRepo }

func (a testSessionAdapter) Exists(jti string) bool {
	_, err := a.repo.GetByID(jti)
	return err == nil
}
func (a testSessionAdapter) Touch(jti string) error { return a.repo.Touch(jti) }

func TestMeReturnsUserAndPerms(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := handlerTestDB(t)

	jm := jwt.New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	userRepo := repository.NewUserRepo(db)
	sessRepo := repository.NewSessionRepo(db)
	menuRepo := repository.NewMenuRepo(db)

	authSvc := service.NewAuthService(userRepo, sessRepo, jm, config.LoginConfig{MaxFailCount: 5, LockMinutes: 10})
	menuSvc := service.NewMenuService(menuRepo, userRepo)
	enforcer := auth.NewDBEnforcer(db)
	_ = enforcer // used implicitly via RequirePerm; Auth middleware suffices for this test

	authH := NewAuthHandler(authSvc, menuSvc, captcha.NewStore(), false)
	checker := testSessionAdapter{repo: sessRepo}

	// login admin through service
	loginOut, err := authSvc.Login("admin", "admin123", service.LoginMeta{IP: "127.0.0.1"})
	require.NoError(t, err)
	token := loginOut.AccessToken

	// build minimal router
	r := gin.New()
	r.GET("/api/auth/me", middleware.Auth(jm, checker), authH.Me)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data struct {
			User struct {
				Username string `json:"username"`
			} `json:"user"`
			Perms []string `json:"perms"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "admin", resp.Data.User.Username)
	require.Contains(t, resp.Data.Perms, "system:user:delete")
}
