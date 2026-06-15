package repository

import (
	"time"

	"go-admin/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var u model.User
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.db.Preload("Roles").First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) IncFailCount(id uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("login_fail_count", gorm.Expr("login_fail_count + 1")).Error
}

func (r *UserRepo) Lock(id uint, until time.Time) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("lock_until", until).Error
}

func (r *UserRepo) ResetFailCount(id uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]any{"login_fail_count": 0, "lock_until": nil}).Error
}

func (r *UserRepo) UpdateLastLogin(id uint) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("last_login_at", now).Error
}

// ListUsers 分页 + 可选用户名模糊。
func (r *UserRepo) ListUsers(keyword string, page, size int) ([]model.User, int64, error) {
	q := r.db.Model(&model.User{}).Preload("Roles")
	if keyword != "" {
		q = q.Where("username ILIKE ?", "%"+keyword+"%")
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
	var users []model.User
	err := q.Order("id asc").Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}

func (r *UserRepo) Create(u *model.User) error { return r.db.Create(u).Error }

// UpdateProfile 更新非敏感字段。
func (r *UserRepo) UpdateProfile(id uint, fields map[string]any) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepo) Delete(id uint) error { return r.db.Delete(&model.User{}, id).Error }

func (r *UserRepo) SetPassword(id uint, hashed string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashed).Error
}

func (r *UserRepo) AssignRoles(id uint, roleIDs []uint) error {
	var roles []model.Role
	if len(roleIDs) > 0 {
		if err := r.db.Find(&roles, roleIDs).Error; err != nil {
			return err
		}
	}
	u := model.User{}
	u.ID = id
	return r.db.Model(&u).Association("Roles").Replace(&roles)
}
