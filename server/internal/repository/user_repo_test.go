package repository

import (
	"testing"
	"time"

	"go-admin/internal/model"
	"go-admin/internal/testsupport"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func userTestGorm(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Role{}))
	db.Where("username = ?", "tuser").Delete(&model.User{})
	return db
}

func TestUserRepoFindAndFailCount(t *testing.T) {
	db := userTestGorm(t)
	repo := NewUserRepo(db)
	require.NoError(t, db.Create(&model.User{Username: "tuser", Password: "h", Status: 1}).Error)

	u, err := repo.FindByUsername("tuser")
	require.NoError(t, err)
	require.Equal(t, "tuser", u.Username)

	require.NoError(t, repo.IncFailCount(u.ID))
	u2, _ := repo.FindByUsername("tuser")
	require.Equal(t, 1, u2.LoginFailCount)

	until := time.Now().Add(10 * time.Minute)
	require.NoError(t, repo.Lock(u.ID, until))
	u3, _ := repo.FindByUsername("tuser")
	require.NotNil(t, u3.LockUntil)

	require.NoError(t, repo.ResetFailCount(u.ID))
	u4, _ := repo.FindByUsername("tuser")
	require.Equal(t, 0, u4.LoginFailCount)
	require.Nil(t, u4.LockUntil)
}
