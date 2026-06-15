package service

import (
	"testing"

	"go-admin/internal/database"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/internal/testsupport"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func roleSvcDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, database.MigrateAndSeed(db))
	return db
}

func TestRoleServiceCRUDAndAssign(t *testing.T) {
	db := roleSvcDB(t)
	svc := NewRoleService(repository.NewRoleRepo(db))

	r, err := svc.Create(RoleInput{Name: "运营", Code: "ops_" + t.Name(), Sort: 5})
	require.NoError(t, err)
	require.NotZero(t, r.ID)

	// 重复 code 报错
	_, err = svc.Create(RoleInput{Name: "运营2", Code: "ops_" + t.Name()})
	require.Error(t, err)

	// 分配两个菜单权限
	var menus []model.Menu
	require.NoError(t, db.Where("perm IN ?", []string{"system:user:list", "system:user:create"}).Find(&menus).Error)
	require.Len(t, menus, 2)
	ids := []uint{menus[0].ID, menus[1].ID}
	require.NoError(t, svc.AssignMenus(r.ID, ids))

	got, err := svc.Get(r.ID)
	require.NoError(t, err)
	require.Len(t, got.Menus, 2)

	// 不能删除内置 super_admin
	var sa model.Role
	require.NoError(t, db.Where("code = ?", "super_admin").First(&sa).Error)
	require.ErrorIs(t, svc.Delete(sa.ID), ErrProtectedRole)

	// 可删自建角色
	require.NoError(t, svc.Delete(r.ID))
}

func TestUpdateRolePersists(t *testing.T) {
	db := roleSvcDB(t)
	svc := NewRoleService(repository.NewRoleRepo(db))

	// create a role with unique code
	code := "upd_role_" + t.Name()
	r, err := svc.Create(RoleInput{Name: "OrigName_" + t.Name(), Code: code, Sort: 1, Remark: "orig remark"})
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM sys_role_menu WHERE role_id = ?", r.ID)
		db.Exec("DELETE FROM sys_role WHERE id = ?", r.ID)
	})

	// update name and remark
	newName := "UpdatedName_" + t.Name()
	newRemark := "updated remark"
	require.NoError(t, svc.Update(r.ID, RoleInput{Name: newName, Code: code, Sort: 2, Remark: newRemark}))

	// reload via Get and assert the new name/remark
	updated, err := svc.Get(r.ID)
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)
	require.Equal(t, newRemark, updated.Remark)
}
