package database

import (
	"testing"

	"go-admin/internal/model"
	"go-admin/pkg/password"

	"github.com/stretchr/testify/require"
)

func TestMigrateAndSeed(t *testing.T) {
	db := connectTestDB(t)

	// 干净环境：删表重建
	require.NoError(t, db.Migrator().DropTable(
		&model.User{}, &model.Role{}, &model.Menu{},
		&model.DictType{}, &model.DictData{}, &model.OnlineSession{},
		"sys_user_role", "sys_role_menu",
	))

	require.NoError(t, MigrateAndSeed(db))

	// 超管账号存在且密码=admin123
	var admin model.User
	require.NoError(t, db.Where("username = ?", "admin").First(&admin).Error)
	require.True(t, password.Check(admin.Password, "admin123"))

	// 预置 3 角色
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	require.GreaterOrEqual(t, roleCount, int64(3))

	// 超管绑定了 super_admin 角色
	require.NoError(t, db.Preload("Roles").First(&admin, admin.ID).Error)
	require.Len(t, admin.Roles, 1)
	require.Equal(t, "super_admin", admin.Roles[0].Code)

	// 幂等：再次执行不报错、不重复插
	require.NoError(t, MigrateAndSeed(db))
	var roleCount2 int64
	db.Model(&model.Role{}).Count(&roleCount2)
	require.Equal(t, roleCount, roleCount2)
}
