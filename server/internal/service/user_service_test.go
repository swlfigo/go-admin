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

func userSvcDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, database.MigrateAndSeed(db))
	return db
}

func TestUserServiceCRUD(t *testing.T) {
	db := userSvcDB(t)
	svc := NewUserService(repository.NewUserRepo(db))

	created, err := svc.Create(CreateUserInput{Username: "u_crud_" + t.Name(), Password: "pass123", Nickname: "U"})
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	require.NotEqual(t, "pass123", created.Password) // 已哈希

	// 重复用户名报错
	_, err = svc.Create(CreateUserInput{Username: created.Username, Password: "pass123"})
	require.Error(t, err)

	// 列表能查到
	list, total, err := svc.List("u_crud_"+t.Name(), 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	// 改昵称
	require.NoError(t, svc.Update(created.ID, UpdateUserInput{Nickname: "U2", Status: 1}))

	// 重置密码
	require.NoError(t, svc.ResetPassword(created.ID, "newpass"))

	// 删除
	require.NoError(t, svc.Delete(created.ID))
	_, _, err2 := svc.List("u_crud_"+t.Name(), 1, 10)
	require.NoError(t, err2)
}

func TestAssignRolesAndStatus(t *testing.T) {
	db := userSvcDB(t)
	svc := NewUserService(repository.NewUserRepo(db))

	// create a fresh user
	username := "u_roles_" + t.Name()
	created, err := svc.Create(CreateUserInput{Username: username, Password: "pass123", Nickname: "RolesUser"})
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM sys_user_role WHERE user_id = ?", created.ID)
		db.Exec("DELETE FROM sys_user WHERE id = ?", created.ID)
	})

	// find the seeded "admin" role
	var adminRole struct{ ID uint }
	require.NoError(t, db.Raw("SELECT id FROM sys_role WHERE code = 'admin' LIMIT 1").Scan(&adminRole).Error)
	require.NotZero(t, adminRole.ID)

	// assign the admin role
	require.NoError(t, svc.AssignRoles(created.ID, []uint{adminRole.ID}))

	// reload user and assert role is attached
	repo := repository.NewUserRepo(db)
	reloaded, err := repo.FindByID(created.ID)
	require.NoError(t, err)
	require.Len(t, reloaded.Roles, 1)
	require.Equal(t, adminRole.ID, reloaded.Roles[0].ID)
}

func TestRecreateAfterSoftDelete(t *testing.T) {
	db := userSvcDB(t)
	svc := NewUserService(repository.NewUserRepo(db))

	username := "recreate_" + t.Name()

	// Cleanup any residual records from previous runs (unscoped to remove soft-deleted rows too)
	t.Cleanup(func() {
		db.Unscoped().Where("username = ?", username).Delete(&model.User{})
	})
	db.Unscoped().Where("username = ?", username).Delete(&model.User{})

	// Create user
	u1, err := svc.Create(CreateUserInput{Username: username, Password: "pass123"})
	require.NoError(t, err)
	require.NotZero(t, u1.ID)

	// Soft-delete it
	require.NoError(t, svc.Delete(u1.ID))

	// Create another user with the SAME username — must succeed (partial unique index allows this)
	u2, err := svc.Create(CreateUserInput{Username: username, Password: "pass456"})
	require.NoError(t, err, "recreating user with same username after soft delete must succeed")
	require.NotZero(t, u2.ID)
	require.NotEqual(t, u1.ID, u2.ID)
}
