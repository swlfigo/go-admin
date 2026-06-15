package auth

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

const SuperAdminCode = "super_admin"

// Enforcer 抽象接口级鉴权，便于未来替换为 Casbin/redis 实现。
type Enforcer interface {
	Can(userID uint, perm string) bool
}

type DBEnforcer struct{ db *gorm.DB }

func NewDBEnforcer(db *gorm.DB) *DBEnforcer { return &DBEnforcer{db: db} }

func (e *DBEnforcer) Can(userID uint, perm string) bool {
	var user model.User
	if err := e.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return false
	}
	roleIDs := make([]uint, 0, len(user.Roles))
	for _, r := range user.Roles {
		if r.Code == SuperAdminCode {
			return true // 超管短路
		}
		roleIDs = append(roleIDs, r.ID)
	}
	if perm == "" || len(roleIDs) == 0 {
		return false
	}
	var count int64
	e.db.Table("sys_role_menu rm").
		Joins("JOIN sys_menu m ON m.id = rm.menu_id AND m.deleted_at IS NULL").
		Joins("JOIN sys_role r ON r.id = rm.role_id AND r.deleted_at IS NULL").
		Where("rm.role_id IN ? AND m.perm = ?", roleIDs, perm).
		Count(&count)
	return count > 0
}
