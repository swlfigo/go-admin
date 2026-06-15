package service

import (
	"errors"

	"go-admin/internal/model"
	"go-admin/internal/repository"

	"gorm.io/gorm"
)

var ErrDictTypeTaken = errors.New("字典类型已存在")

type DictTypeInput struct {
	Name   string
	Type   string
	Status int
	Remark string
}

type DictDataInput struct {
	DictType string
	Label    string
	Value    string
	TagType  string
	Sort     int
	Status   int
	Remark   string
}

type DictService struct{ dict *repository.DictRepo }

func NewDictService(d *repository.DictRepo) *DictService { return &DictService{dict: d} }

func (s *DictService) ListTypes(keyword string, page, size int) ([]model.DictType, int64, error) {
	return s.dict.ListTypes(keyword, page, size)
}

func (s *DictService) CreateType(in DictTypeInput) (*model.DictType, error) {
	if _, err := s.dict.FindType(in.Type); err == nil {
		return nil, ErrDictTypeTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	dt := &model.DictType{Name: in.Name, Type: in.Type, Status: 1, Remark: in.Remark}
	if err := s.dict.CreateType(dt); err != nil {
		return nil, err
	}
	return dt, nil
}

func (s *DictService) UpdateType(id uint, in DictTypeInput) error {
	return s.dict.UpdateType(id, map[string]any{"name": in.Name, "status": in.Status, "remark": in.Remark})
}

func (s *DictService) DeleteType(id uint) error { return s.dict.DeleteType(id) }

func (s *DictService) DataByType(typ string) ([]model.DictData, error) {
	return s.dict.DataByType(typ)
}

func (s *DictService) CreateData(in DictDataInput) (*model.DictData, error) {
	d := &model.DictData{
		DictType: in.DictType, Label: in.Label, Value: in.Value,
		TagType: in.TagType, Sort: in.Sort, Status: 1, Remark: in.Remark,
	}
	if err := s.dict.CreateData(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DictService) UpdateData(id uint, in DictDataInput) error {
	return s.dict.UpdateData(id, map[string]any{
		"label": in.Label, "value": in.Value, "tag_type": in.TagType,
		"sort": in.Sort, "status": in.Status, "remark": in.Remark,
	})
}

func (s *DictService) DeleteData(id uint) error { return s.dict.DeleteData(id) }
