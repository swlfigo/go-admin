package repository

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) List(keyword string, page, size int) ([]model.Role, int64, error) {
	q := r.db.Model(&model.Role{})
	if keyword != "" {
		q = q.Where("name ILIKE ? OR code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var roles []model.Role
	err := q.Order("sort asc, id asc").Offset((page - 1) * size).Limit(size).Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) Create(role *model.Role) error { return r.db.Create(role).Error }

func (r *RoleRepo) Update(id uint, fields map[string]any) error {
	return r.db.Model(&model.Role{}).Where("id = ?", id).Updates(fields).Error
}

func (r *RoleRepo) Delete(id uint) error { return r.db.Delete(&model.Role{}, id).Error }

func (r *RoleRepo) AssignMenus(id uint, menuIDs []uint) error {
	var menus []model.Menu
	if len(menuIDs) > 0 {
		if err := r.db.Find(&menus, menuIDs).Error; err != nil {
			return err
		}
	}
	role := model.Role{}
	role.ID = id
	return r.db.Model(&role).Association("Menus").Replace(&menus)
}

// ClearAssociations removes all menu and user join rows for the given role (defense in depth before delete).
func (r *RoleRepo) ClearAssociations(id uint) error {
	role := model.Role{}
	role.ID = id
	if err := r.db.Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}
	return r.db.Exec("DELETE FROM sys_user_role WHERE role_id = ?", id).Error
}
