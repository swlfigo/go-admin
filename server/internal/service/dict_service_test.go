package service

import (
	"testing"

	"go-admin/internal/model"
	"go-admin/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestDictServiceCRUD(t *testing.T) {
	db := menuTestDB(t) // 复用 helper（含 seed）
	svc := NewDictService(repository.NewDictRepo(db))

	// 新增类型
	dt, err := svc.CreateType(DictTypeInput{Name: "测试类型_" + t.Name(), Type: "test_" + t.Name()})
	require.NoError(t, err)
	require.NotZero(t, dt.ID)

	// 重复 type 报错
	_, err = svc.CreateType(DictTypeInput{Name: "x", Type: "test_" + t.Name()})
	require.Error(t, err)

	// 新增数据
	dd, err := svc.CreateData(DictDataInput{DictType: dt.Type, Label: "启用", Value: "1", Sort: 1})
	require.NoError(t, err)
	require.NotZero(t, dd.ID)

	// 按类型查数据
	list, err := svc.DataByType(dt.Type)
	require.NoError(t, err)
	require.Len(t, list, 1)

	// 删除数据、删除类型
	require.NoError(t, svc.DeleteData(dd.ID))
	require.NoError(t, svc.DeleteType(dt.ID))
}

// TestUpdateDictDataPersists creates a type + a data row, UpdateData its
// label+value, reloads via DataByType and asserts the new label/value.
func TestUpdateDictDataPersists(t *testing.T) {
	db := menuTestDB(t)
	svc := NewDictService(repository.NewDictRepo(db))

	// Create a dict type
	dt, err := svc.CreateType(DictTypeInput{
		Name: "persist_type_" + t.Name(),
		Type: "persist_" + t.Name(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Unscoped().Where("dict_type = ?", dt.Type).Delete(&model.DictData{})
		db.Unscoped().Delete(dt)
	})

	// Create a data row
	dd, err := svc.CreateData(DictDataInput{
		DictType: dt.Type,
		Label:    "original_label",
		Value:    "original_value",
		Sort:     1,
	})
	require.NoError(t, err)
	require.NotZero(t, dd.ID)

	t.Cleanup(func() {
		svc.DeleteData(dd.ID) //nolint:errcheck
	})

	// Update label and value
	err = svc.UpdateData(dd.ID, DictDataInput{
		DictType: dt.Type,
		Label:    "updated_label",
		Value:    "updated_value",
		Sort:     2,
	})
	require.NoError(t, err)

	// Reload via DataByType and assert new label/value
	list, err := svc.DataByType(dt.Type)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "updated_label", list[0].Label)
	require.Equal(t, "updated_value", list[0].Value)
}
