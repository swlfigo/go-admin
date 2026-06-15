package service

import (
	"errors"

	"go-admin/internal/model"
	"go-admin/internal/repository"
)

var ErrMenuHasChildren = errors.New("存在子菜单，不能删除")

type MenuInput struct {
	ParentID  uint
	Name      string
	Type      string
	Path      string
	Component string
	Perm      string
	Icon      string
	Sort      int
	Visible   int
	Status    int
}

type MenuAdminService struct{ menus *repository.MenuRepo }

func NewMenuAdminService(m *repository.MenuRepo) *MenuAdminService {
	return &MenuAdminService{menus: m}
}

// Tree 返回包含按钮(F)的完整管理树。
func (s *MenuAdminService) Tree() ([]MenuNode, error) {
	ms, err := s.menus.All()
	if err != nil {
		return nil, err
	}
	nodes := make([]MenuNode, 0, len(ms))
	for _, m := range ms {
		nodes = append(nodes, MenuNode{
			ID: m.ID, ParentID: m.ParentID, Name: m.Name, Type: m.Type,
			Path: m.Path, Component: m.Component, Perm: m.Perm, Icon: m.Icon, Sort: m.Sort,
			Visible: m.Visible, Status: m.Status,
		})
	}
	return buildTree(nodes, 0), nil
}

func (s *MenuAdminService) Create(in MenuInput) (*model.Menu, error) {
	m := &model.Menu{
		ParentID: in.ParentID, Name: in.Name, Type: in.Type, Path: in.Path,
		Component: in.Component, Perm: in.Perm, Icon: in.Icon, Sort: in.Sort,
		Visible: in.Visible, Status: 1,
	}
	if err := s.menus.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MenuAdminService) Update(id uint, in MenuInput) error {
	m := &model.Menu{
		ParentID: in.ParentID, Name: in.Name, Type: in.Type, Path: in.Path,
		Component: in.Component, Perm: in.Perm, Icon: in.Icon, Sort: in.Sort,
		Visible: in.Visible, Status: in.Status,
	}
	m.ID = id
	return s.menus.Update(m)
}

func (s *MenuAdminService) Delete(id uint) error {
	has, err := s.menus.HasChildren(id)
	if err != nil {
		return err
	}
	if has {
		return ErrMenuHasChildren
	}
	return s.menus.Delete(id)
}
