package auth

import (
	"testing"

	"go-admin/internal/database"
	"go-admin/internal/model"
	"go-admin/internal/testsupport"
	"go-admin/pkg/password"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func enforcerDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, database.MigrateAndSeed(db))
	return db
}

func TestSuperAdminCanEverything(t *testing.T) {
	db := enforcerDB(t)
	e := NewDBEnforcer(db)
	var admin model.User
	require.NoError(t, db.Where("username = ?", "admin").First(&admin).Error)
	require.True(t, e.Can(admin.ID, "system:user:delete"))
	require.True(t, e.Can(admin.ID, "anything:at:all"))
}

func TestNormalUserLimitedByRoleMenu(t *testing.T) {
	db := enforcerDB(t)
	// 清理上次残留（幂等），测试前后都清理确保不污染其它测试
	cleanup := func() {
		db.Where("username = ?", "limited_"+t.Name()).Delete(&model.User{})
		db.Where("name = ?", "t-list-"+t.Name()).Delete(&model.Menu{})
		db.Where("code = ?", "r_"+t.Name()).Delete(&model.Role{})
	}
	cleanup()
	t.Cleanup(cleanup)
	// 造一个普通用户，只给一个仅含 system:user:list 的角色
	h, _ := password.Hash("x")
	u := model.User{Username: "limited_" + t.Name(), Password: h, Status: 1}
	require.NoError(t, db.Create(&u).Error)
	menu := model.Menu{Name: "t-list-" + t.Name(), Type: "F", Perm: "system:user:list", Status: 1}
	require.NoError(t, db.Create(&menu).Error)
	role := model.Role{Name: "r-" + t.Name(), Code: "r_" + t.Name(), Status: 1}
	require.NoError(t, db.Create(&role).Error)
	require.NoError(t, db.Model(&role).Association("Menus").Append(&menu))
	require.NoError(t, db.Model(&u).Association("Roles").Append(&role))

	e := NewDBEnforcer(db)
	require.True(t, e.Can(u.ID, "system:user:list"))
	require.False(t, e.Can(u.ID, "system:user:delete"))
}

func TestSoftDeletedRoleGrantsNothing(t *testing.T) {
	db := enforcerDB(t)
	cleanup := func() {
		// Must clear join tables before deleting parent rows (FK constraints)
		var userIDs []uint
		db.Unscoped().Model(&model.User{}).Where("username = ?", "sd_user_"+t.Name()).Pluck("id", &userIDs)
		if len(userIDs) > 0 {
			db.Exec("DELETE FROM sys_user_role WHERE user_id IN ?", userIDs)
		}
		var roleIDs []uint
		db.Unscoped().Model(&model.Role{}).Where("code = ?", "sd_role_"+t.Name()).Pluck("id", &roleIDs)
		if len(roleIDs) > 0 {
			db.Exec("DELETE FROM sys_role_menu WHERE role_id IN ?", roleIDs)
			db.Exec("DELETE FROM sys_user_role WHERE role_id IN ?", roleIDs)
		}
		var menuIDs []uint
		db.Unscoped().Model(&model.Menu{}).Where("name = ?", "sd_menu_"+t.Name()).Pluck("id", &menuIDs)
		if len(menuIDs) > 0 {
			db.Exec("DELETE FROM sys_role_menu WHERE menu_id IN ?", menuIDs)
		}
		db.Unscoped().Where("username = ?", "sd_user_"+t.Name()).Delete(&model.User{})
		db.Unscoped().Where("name = ?", "sd_menu_"+t.Name()).Delete(&model.Menu{})
		db.Unscoped().Where("code = ?", "sd_role_"+t.Name()).Delete(&model.Role{})
	}
	cleanup()
	t.Cleanup(cleanup)

	h, _ := password.Hash("x")
	u := model.User{Username: "sd_user_" + t.Name(), Password: h, Status: 1}
	require.NoError(t, db.Create(&u).Error)
	menu := model.Menu{Name: "sd_menu_" + t.Name(), Type: "F", Perm: "t:soft:perm", Status: 1}
	require.NoError(t, db.Create(&menu).Error)
	role := model.Role{Name: "sd_role_" + t.Name(), Code: "sd_role_" + t.Name(), Status: 1}
	require.NoError(t, db.Create(&role).Error)
	require.NoError(t, db.Model(&role).Association("Menus").Append(&menu))
	require.NoError(t, db.Model(&u).Association("Roles").Append(&role))

	e := NewDBEnforcer(db)
	// Before soft-delete: permission must be granted
	require.True(t, e.Can(u.ID, "t:soft:perm"))

	// Soft-delete the role
	require.NoError(t, db.Delete(&role).Error)

	// After soft-delete: permission must be revoked
	require.False(t, e.Can(u.ID, "t:soft:perm"))
}
