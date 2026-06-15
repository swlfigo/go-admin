package database

import (
	"testing"

	"go-admin/internal/model"

	"github.com/stretchr/testify/require"
)

func TestSeedMenusAndLinks(t *testing.T) {
	db := connectTestDB(t)
	require.NoError(t, MigrateAndSeed(db))

	// 至少种了：dashboard、system 目录、user/role/menu/dict 菜单、online
	var menuCount int64
	db.Model(&model.Menu{}).Count(&menuCount)
	require.GreaterOrEqual(t, menuCount, int64(7))

	// 存在按钮权限节点
	var btnCount int64
	db.Model(&model.Menu{}).Where("type = ?", "F").Count(&btnCount)
	require.Greater(t, btnCount, int64(0))

	// super_admin 关联了菜单（用于 /menus/me；DBEnforcer 另会短路）
	var superRole model.Role
	require.NoError(t, db.Preload("Menus").Where("code = ?", "super_admin").First(&superRole).Error)
	require.Greater(t, len(superRole.Menus), 0)

	// 幂等
	require.NoError(t, MigrateAndSeed(db))
	var menuCount2 int64
	db.Model(&model.Menu{}).Count(&menuCount2)
	require.Equal(t, menuCount, menuCount2)
}
