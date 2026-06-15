package repository

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

type LogRepo struct{ db *gorm.DB }

func NewLogRepo(db *gorm.DB) *LogRepo { return &LogRepo{db: db} }

func (r *LogRepo) Create(l *model.OperationLog) error { return r.db.Create(l).Error }

func (r *LogRepo) List(keyword string, page, size int) ([]model.OperationLog, int64, error) {
	q := r.db.Model(&model.OperationLog{})
	if keyword != "" {
		q = q.Where("username ILIKE ? OR path ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
	var logs []model.OperationLog
	err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&logs).Error
	return logs, total, err
}

func (r *LogRepo) Clear() error {
	return r.db.Where("1=1").Delete(&model.OperationLog{}).Error
}
