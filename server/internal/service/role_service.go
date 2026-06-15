package service

import (
	"errors"

	"go-admin/internal/auth"
	"go-admin/internal/model"
	"go-admin/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrRoleCodeTaken = errors.New("角色编码已存在")
	ErrProtectedRole = errors.New("内置角色不可删除")
)

type RoleInput struct {
	Name   string
	Code   string
	Sort   int
	Status int
	Remark string
}

type RoleService struct{ roles *repository.RoleRepo }

func NewRoleService(r *repository.RoleRepo) *RoleService { return &RoleService{roles: r} }

func (s *RoleService) List(keyword string, page, size int) ([]model.Role, int64, error) {
	return s.roles.List(keyword, page, size)
}

func (s *RoleService) Get(id uint) (*model.Role, error) { return s.roles.GetByID(id) }

func (s *RoleService) Create(in RoleInput) (*model.Role, error) {
	if _, err := s.roles.FindByCode(in.Code); err == nil {
		return nil, ErrRoleCodeTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	role := &model.Role{Name: in.Name, Code: in.Code, Sort: in.Sort, Status: 1, Remark: in.Remark}
	if err := s.roles.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Update(id uint, in RoleInput) error {
	return s.roles.Update(id, map[string]any{
		"name": in.Name, "sort": in.Sort, "status": in.Status, "remark": in.Remark,
	})
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roles.GetByID(id)
	if err != nil {
		return err
	}
	if role.Code == auth.SuperAdminCode {
		return ErrProtectedRole
	}
	if err := s.roles.ClearAssociations(id); err != nil {
		return err
	}
	return s.roles.Delete(id)
}

func (s *RoleService) AssignMenus(id uint, menuIDs []uint) error {
	return s.roles.AssignMenus(id, menuIDs)
}
