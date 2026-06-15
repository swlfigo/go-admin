package repository

import (
	"time"

	"go-admin/internal/model"

	"gorm.io/gorm"
)

type SessionRepo struct{ db *gorm.DB }

func NewSessionRepo(db *gorm.DB) *SessionRepo { return &SessionRepo{db: db} }

func (r *SessionRepo) Create(s *model.OnlineSession) error {
	return r.db.Create(s).Error
}

func (r *SessionRepo) GetByID(id string) (*model.OnlineSession, error) {
	var s model.OnlineSession
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepo) GetByRefreshID(rid string) (*model.OnlineSession, error) {
	var s model.OnlineSession
	if err := r.db.First(&s, "refresh_id = ?", rid).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepo) List() ([]model.OnlineSession, error) {
	var list []model.OnlineSession
	err := r.db.Order("last_active_at desc").Find(&list).Error
	return list, err
}

func (r *SessionRepo) Touch(id string) error {
	return r.db.Model(&model.OnlineSession{}).
		Where("id = ?", id).Update("last_active_at", time.Now()).Error
}

func (r *SessionRepo) Delete(id string) error {
	return r.db.Delete(&model.OnlineSession{}, "id = ?", id).Error
}

// Rotate 用新 jti/refresh 替换旧会话（删旧建新，原子事务）
func (r *SessionRepo) Rotate(oldID string, ns *model.OnlineSession) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.OnlineSession{}, "id = ?", oldID).Error; err != nil {
			return err
		}
		return tx.Create(ns).Error
	})
}
