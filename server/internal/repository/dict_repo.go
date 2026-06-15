package repository

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

type DictRepo struct{ db *gorm.DB }

func NewDictRepo(db *gorm.DB) *DictRepo { return &DictRepo{db: db} }

func (r *DictRepo) ListTypes(keyword string, page, size int) ([]model.DictType, int64, error) {
	q := r.db.Model(&model.DictType{})
	if keyword != "" {
		q = q.Where("name ILIKE ? OR type ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
	var types []model.DictType
	err := q.Order("id asc").Offset((page - 1) * size).Limit(size).Find(&types).Error
	return types, total, err
}

func (r *DictRepo) FindType(typ string) (*model.DictType, error) {
	var dt model.DictType
	if err := r.db.Where("type = ?", typ).First(&dt).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}

func (r *DictRepo) CreateType(dt *model.DictType) error { return r.db.Create(dt).Error }
func (r *DictRepo) UpdateType(id uint, fields map[string]any) error {
	return r.db.Model(&model.DictType{}).Where("id = ?", id).Updates(fields).Error
}
func (r *DictRepo) DeleteType(id uint) error { return r.db.Delete(&model.DictType{}, id).Error }

func (r *DictRepo) DataByType(typ string) ([]model.DictData, error) {
	var data []model.DictData
	err := r.db.Where("dict_type = ?", typ).Order("sort asc, id asc").Find(&data).Error
	return data, err
}

func (r *DictRepo) CreateData(d *model.DictData) error { return r.db.Create(d).Error }
func (r *DictRepo) UpdateData(id uint, fields map[string]any) error {
	return r.db.Model(&model.DictData{}).Where("id = ?", id).Updates(fields).Error
}
func (r *DictRepo) DeleteData(id uint) error { return r.db.Delete(&model.DictData{}, id).Error }
