package service

import (
	"testing"

	"go-admin/internal/database"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/internal/testsupport"
	"go-admin/pkg/password"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func menuTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := testsupport.Connect()
	require.NoError(t, err)
	require.NoError(t, database.MigrateAndSeed(db))
	return db
}

func TestSuperAdminMenuTreeAndPerms(t *testing.T) {
	db := menuTestDB(t)
	svc := NewMenuService(repository.NewMenuRepo(db), repository.NewUserRepo(db))
	var admin model.User
	require.NoError(t, db.Where("username = ?", "admin").First(&admin).Error)

	tree, err := svc.UserMenuTree(admin.ID)
	require.NoError(t, err)
	require.Greater(t, len(tree), 0)
	// 树里只含目录/菜单（M/C），不含按钮（F）
	for _, n := range tree {
		require.NotEqual(t, "F", n.Type)
	}
	// system 目录有子节点
	var sysNode *MenuNode
	for i := range tree {
		if tree[i].Name == "系统管理" {
			sysNode = &tree[i]
		}
	}
	require.NotNil(t, sysNode)
	require.Greater(t, len(sysNode.Children), 0)

	// 权限码包含按钮权限
	perms, err := svc.UserPermCodes(admin.ID)
	require.NoError(t, err)
	require.Contains(t, perms, "system:user:delete")
}

// TestNormalUserMenuTreeFiltered verifies that a normal user with a role linked
// only to "用户管理" (system:user:list) and its parent "系统管理" directory sees
// only those menus and only that perm code.
func TestNormalUserMenuTreeFiltered(t *testing.T) {
	db := menuTestDB(t)
	svc := NewMenuService(repository.NewMenuRepo(db), repository.NewUserRepo(db))

	// Create a normal user
	h, err := password.Hash("testpass")
	require.NoError(t, err)
	u := model.User{Username: "filtered_user_" + t.Name(), Password: h, Status: 1}
	require.NoError(t, db.Create(&u).Error)
	t.Cleanup(func() { db.Unscoped().Delete(&u) })

	// Find the "系统管理" directory and "用户管理" menu from seeds
	var sysMenu model.Menu
	require.NoError(t, db.Where("name = ? AND type = ?", "系统管理", "M").First(&sysMenu).Error)
	var userMenu model.Menu
	require.NoError(t, db.Where("name = ? AND type = ?", "用户管理", "C").First(&userMenu).Error)

	// Create a role linked only to these two menus
	role := model.Role{
		Name:   "role_" + t.Name(),
		Code:   "code_" + t.Name(),
		Status: 1,
	}
	require.NoError(t, db.Create(&role).Error)
	t.Cleanup(func() { db.Unscoped().Delete(&role) })
	require.NoError(t, db.Model(&role).Association("Menus").Append(&sysMenu, &userMenu))
	t.Cleanup(func() { db.Model(&role).Association("Menus").Clear() })

	// Assign role to user
	require.NoError(t, db.Model(&u).Association("Roles").Append(&role))
	t.Cleanup(func() { db.Model(&u).Association("Roles").Clear() })

	// UserMenuTree should return system dir containing exactly user-management child
	tree, err := svc.UserMenuTree(u.ID)
	require.NoError(t, err)
	require.Len(t, tree, 1, "only system dir at root level")
	require.Equal(t, "系统管理", tree[0].Name)
	require.Len(t, tree[0].Children, 1, "exactly one child")
	require.Equal(t, "用户管理", tree[0].Children[0].Name)

	// Role/menu/dict should NOT appear in children
	for _, child := range tree[0].Children {
		require.NotEqual(t, "角色管理", child.Name)
		require.NotEqual(t, "菜单管理", child.Name)
		require.NotEqual(t, "字典管理", child.Name)
	}

	// UserPermCodes: only system:user:list (not system:role:list etc.)
	perms, err := svc.UserPermCodes(u.ID)
	require.NoError(t, err)
	require.Contains(t, perms, "system:user:list")
	require.NotContains(t, perms, "system:role:list")
}

// TestBuildTreeOrItemHandlesEmpty verifies that a freshly-created user with NO
// roles returns an empty (non-nil) slice for UserPermCodes and empty tree for
// UserMenuTree.
func TestBuildTreeOrItemHandlesEmpty(t *testing.T) {
	db := menuTestDB(t)
	svc := NewMenuService(repository.NewMenuRepo(db), repository.NewUserRepo(db))

	h, err := password.Hash("emptypass")
	require.NoError(t, err)
	u := model.User{Username: "noroles_user_" + t.Name(), Password: h, Status: 1}
	require.NoError(t, db.Create(&u).Error)
	t.Cleanup(func() { db.Unscoped().Delete(&u) })

	// UserPermCodes should return non-nil empty slice
	perms, err := svc.UserPermCodes(u.ID)
	require.NoError(t, err)
	require.NotNil(t, perms)
	require.Len(t, perms, 0)

	// UserMenuTree should return empty (non-nil) slice
	tree, err := svc.UserMenuTree(u.ID)
	require.NoError(t, err)
	require.NotNil(t, tree)
	require.Len(t, tree, 0)
}
