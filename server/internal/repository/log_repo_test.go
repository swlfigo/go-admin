package repository

import (
	"testing"
	"time"

	"go-admin/internal/model"
	"go-admin/internal/testsupport"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func logTestGorm(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.OperationLog{}))
	require.NoError(t, db.Exec("DELETE FROM sys_operation_log").Error)
	return db
}

func TestLogRepoListAndClear(t *testing.T) {
	repo := NewLogRepo(logTestGorm(t))
	now := time.Now()
	logs := []model.OperationLog{
		{Username: "admin", Method: "POST", Path: "/api/users", IP: "127.0.0.1", Status: 200, LatencyMs: 5, CreatedAt: now},
		{Username: "admin", Method: "DELETE", Path: "/api/roles/1", IP: "127.0.0.1", Status: 200, LatencyMs: 3, CreatedAt: now},
		{Username: "bob", Method: "PUT", Path: "/api/profile", IP: "127.0.0.1", Status: 200, LatencyMs: 7, CreatedAt: now},
	}
	for i := range logs {
		require.NoError(t, repo.Create(&logs[i]))
	}

	// no keyword -> all three
	all, total, err := repo.List("", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, all, 3)
	// order by id desc -> newest first
	require.Equal(t, "bob", all[0].Username)

	// keyword filters username
	byUser, total, err := repo.List("bob", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, byUser, 1)
	require.Equal(t, "bob", byUser[0].Username)

	// keyword filters path
	byPath, total, err := repo.List("roles", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, byPath, 1)
	require.Equal(t, "/api/roles/1", byPath[0].Path)

	// clear empties
	require.NoError(t, repo.Clear())
	_, total, err = repo.List("", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(0), total)
}
