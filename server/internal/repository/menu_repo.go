package repository

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

type MenuRepo struct{ db *gorm.DB }

func NewMenuRepo(db *gorm.DB) *MenuRepo { return &MenuRepo{db: db} }

// AllVisible 返回全部菜单（管理 + 超管菜单树用）。
func (r *MenuRepo) All() ([]model.Menu, error) {
	var ms []model.Menu
	err := r.db.Order("sort asc, id asc").Find(&ms).Error
	return ms, err
}

// ByUser 返回某用户经角色可见的全部菜单（去重）。
func (r *MenuRepo) ByUser(userID uint) ([]model.Menu, error) {
	var ms []model.Menu
	err := r.db.
		Distinct("m.*").
		Table("sys_menu m").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.deleted_at IS NULL", userID).
		Order("m.sort asc, m.id asc").
		Scan(&ms).Error
	return ms, err
}

func (r *MenuRepo) Create(m *model.Menu) error { return r.db.Create(m).Error }
func (r *MenuRepo) Update(m *model.Menu) error { return r.db.Save(m).Error }
func (r *MenuRepo) Delete(id uint) error       { return r.db.Delete(&model.Menu{}, id).Error }

// HasChildren 判断是否有子菜单（删除前校验）。
func (r *MenuRepo) HasChildren(id uint) (bool, error) {
	var n int64
	err := r.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&n).Error
	return n > 0, err
}
