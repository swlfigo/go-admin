package repository

import (
	"testing"
	"time"

	"go-admin/internal/model"
	"go-admin/internal/testsupport"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func testGorm(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.OnlineSession{}))
	db.Where("1 = 1").Delete(&model.OnlineSession{})
	return db
}

func TestSessionCRUD(t *testing.T) {
	repo := NewSessionRepo(testGorm(t))
	now := time.Now()
	s := &model.OnlineSession{
		ID: "jti-1", UserID: 1, Username: "admin", LoginIP: "127.0.0.1",
		RefreshID: "rid-1", RefreshExpireAt: now.Add(time.Hour),
		LoginAt: now, LastActiveAt: now,
	}
	require.NoError(t, repo.Create(s))

	got, err := repo.GetByID("jti-1")
	require.NoError(t, err)
	require.Equal(t, "admin", got.Username)

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 1)

	require.NoError(t, repo.Delete("jti-1"))
	_, err = repo.GetByID("jti-1")
	require.Error(t, err) // 已删除
}
