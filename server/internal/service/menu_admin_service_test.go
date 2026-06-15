package service

import (
	"testing"

	"go-admin/internal/model"
	"go-admin/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestMenuAdminCRUD(t *testing.T) {
	db := menuTestDB(t) // 复用 Task 4 的 helper
	svc := NewMenuAdminService(repository.NewMenuRepo(db))

	// 全量树
	tree, err := svc.Tree()
	require.NoError(t, err)
	require.Greater(t, len(tree), 0)

	// 新增一个目录
	m, err := svc.Create(MenuInput{ParentID: 0, Name: "报表中心_" + t.Name(), Type: "M", Path: "/report", Sort: 9})
	require.NoError(t, err)
	require.NotZero(t, m.ID)

	// 改名
	require.NoError(t, svc.Update(m.ID, MenuInput{Name: "报表_" + t.Name(), Type: "M", Path: "/report", Sort: 9}))

	// 有子节点不可删
	child := model.Menu{ParentID: m.ID, Name: "子_" + t.Name(), Type: "C", Status: 1}
	require.NoError(t, db.Create(&child).Error)
	require.ErrorIs(t, svc.Delete(m.ID), ErrMenuHasChildren)

	// 删子后可删父
	require.NoError(t, svc.Delete(child.ID))
	require.NoError(t, svc.Delete(m.ID))
}

// TestMenuUpdatePreservesStatus creates a menu with Status=1/Visible=1, calls
// Update setting Status=0 and Visible=0, reloads via repo All() and asserts
// that Status==0 and Visible==0 (i.e. Update no longer forces 1).
func TestMenuUpdatePreservesStatus(t *testing.T) {
	db := menuTestDB(t)
	svc := NewMenuAdminService(repository.NewMenuRepo(db))
	repo := repository.NewMenuRepo(db)

	created, err := svc.Create(MenuInput{
		ParentID: 0,
		Name:     "status_test_" + t.Name(),
		Type:     "M",
		Path:     "/status-test",
		Sort:     1,
		Visible:  1,
		Status:   1,
	})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		db.Unscoped().Delete(&model.Menu{}, created.ID)
	})

	// Update with Status=0 and Visible=0
	err = svc.Update(created.ID, MenuInput{
		Name:    "status_test_" + t.Name(),
		Type:    "M",
		Path:    "/status-test",
		Sort:    1,
		Visible: 0,
		Status:  0,
	})
	require.NoError(t, err)

	// Reload via repo All() and find by id
	all, err := repo.All()
	require.NoError(t, err)

	var updated *model.Menu
	for i := range all {
		if all[i].ID == created.ID {
			updated = &all[i]
			break
		}
	}
	require.NotNil(t, updated, "menu should still exist after update")
	require.Equal(t, 0, updated.Status, "Status should be 0 after update, not forced to 1")
	require.Equal(t, 0, updated.Visible, "Visible should be 0 after update, not forced to 1")
}

// TestMenuUpdateChangesFields creates a menu, updates its Path+Sort, reloads via
// repo All() (finds by id) and asserts the new Path/Sort are persisted.
func TestMenuUpdateChangesFields(t *testing.T) {
	db := menuTestDB(t)
	svc := NewMenuAdminService(repository.NewMenuRepo(db))
	repo := repository.NewMenuRepo(db)

	// Create a new menu
	created, err := svc.Create(MenuInput{
		ParentID: 0,
		Name:     "update_test_" + t.Name(),
		Type:     "M",
		Path:     "/original-path",
		Sort:     50,
	})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		db.Unscoped().Delete(&model.Menu{}, created.ID)
	})

	// Update Path and Sort
	err = svc.Update(created.ID, MenuInput{
		Name: "update_test_" + t.Name(),
		Type: "M",
		Path: "/updated-path",
		Sort: 99,
	})
	require.NoError(t, err)

	// Reload via repo All() and find by id
	all, err := repo.All()
	require.NoError(t, err)

	var updated *model.Menu
	for i := range all {
		if all[i].ID == created.ID {
			updated = &all[i]
			break
		}
	}
	require.NotNil(t, updated, "menu should still exist after update")
	require.Equal(t, "/updated-path", updated.Path)
	require.Equal(t, 99, updated.Sort)
}
