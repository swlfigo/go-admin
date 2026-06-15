package database

import (
	"testing"

	"go-admin/internal/testsupport"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// connectTestDB 确保测试库存在后连接（独立于开发库；单独跑本包也能自建库）。
func connectTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	require.NoError(t, testsupport.EnsureDB())
	db, err := Connect(testsupport.DBConfig())
	require.NoError(t, err)
	return db
}

func TestConnectPings(t *testing.T) {
	db := connectTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())
}
