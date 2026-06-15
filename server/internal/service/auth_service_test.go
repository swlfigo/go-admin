package service

import (
	"testing"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/internal/testsupport"
	"go-admin/pkg/jwt"
	"go-admin/pkg/password"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func authTestGorm(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Role{}, &model.OnlineSession{}))
	db.Where("username = ?", "alice").Delete(&model.User{})
	db.Where("1 = 1").Delete(&model.OnlineSession{})
	h, _ := password.Hash("pass123")
	require.NoError(t, db.Create(&model.User{Username: "alice", Password: h, Status: 1}).Error)
	return db
}

func newAuthService(db *gorm.DB) *AuthService {
	return NewAuthService(
		repository.NewUserRepo(db),
		repository.NewSessionRepo(db),
		jwt.New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour),
		config.LoginConfig{MaxFailCount: 3, LockMinutes: 10},
	)
}

func TestLoginSuccess(t *testing.T) {
	svc := newAuthService(authTestGorm(t))
	out, err := svc.Login("alice", "pass123", LoginMeta{IP: "127.0.0.1"})
	require.NoError(t, err)
	require.NotEmpty(t, out.AccessToken)
	require.NotEmpty(t, out.RefreshToken)
}

func TestLoginWrongPasswordCountsAndLocks(t *testing.T) {
	svc := newAuthService(authTestGorm(t))
	for i := 0; i < 3; i++ {
		_, err := svc.Login("alice", "wrong", LoginMeta{IP: "127.0.0.1"})
		require.Error(t, err)
	}
	// 第 4 次即使密码正确也被锁
	_, err := svc.Login("alice", "pass123", LoginMeta{IP: "127.0.0.1"})
	require.ErrorIs(t, err, ErrAccountLocked)
}

func TestRefreshRotates(t *testing.T) {
	svc := newAuthService(authTestGorm(t))
	out, err := svc.Login("alice", "pass123", LoginMeta{IP: "127.0.0.1"})
	require.NoError(t, err)

	out2, err := svc.Refresh(out.RefreshToken)
	require.NoError(t, err)
	require.NotEqual(t, out.AccessToken, out2.AccessToken)
	require.NotEqual(t, out.RefreshToken, out2.RefreshToken)

	// 旧 refresh 已作废（轮换）→ 再用报错
	_, err = svc.Refresh(out.RefreshToken)
	require.Error(t, err)
}

func TestLogoutRemovesSession(t *testing.T) {
	db := authTestGorm(t)
	svc := newAuthService(db)
	out, _ := svc.Login("alice", "pass123", LoginMeta{IP: "127.0.0.1"})
	m := jwt.New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	claims, _ := m.ParseAccess(out.AccessToken)
	require.NoError(t, svc.Logout(claims.ID))

	var cnt int64
	db.Model(&model.OnlineSession{}).Where("id = ?", claims.ID).Count(&cnt)
	require.Equal(t, int64(0), cnt)
}
