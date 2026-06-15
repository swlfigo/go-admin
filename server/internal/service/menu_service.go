package service

import (
	"go-admin/internal/auth"
	"go-admin/internal/model"
	"go-admin/internal/repository"
)

// MenuNode 是返回给前端建路由的树节点。
type MenuNode struct {
	ID        uint       `json:"id"`
	ParentID  uint       `json:"parentId"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Perm      string     `json:"perm"`
	Icon      string     `json:"icon"`
	Sort      int        `json:"sort"`
	Visible   int        `json:"visible"`
	Status    int        `json:"status"`
	Children  []MenuNode `json:"children,omitempty"`
}

type MenuService struct {
	menus *repository.MenuRepo
	users *repository.UserRepo
}

func NewMenuService(m *repository.MenuRepo, u *repository.UserRepo) *MenuService {
	return &MenuService{menus: m, users: u}
}

// isSuper 判断用户是否超管。
func (s *MenuService) isSuper(userID uint) (bool, error) {
	u, err := s.users.FindByID(userID)
	if err != nil {
		return false, err
	}
	for _, r := range u.Roles {
		if r.Code == auth.SuperAdminCode {
			return true, nil
		}
	}
	return false, nil
}

func (s *MenuService) userMenus(userID uint) ([]model.Menu, error) {
	super, err := s.isSuper(userID)
	if err != nil {
		return nil, err
	}
	if super {
		return s.menus.All()
	}
	return s.menus.ByUser(userID)
}

// UserMenuTree 返回去掉按钮(F)、按 parent 组装的目录/菜单树。
func (s *MenuService) UserMenuTree(userID uint) ([]MenuNode, error) {
	ms, err := s.userMenus(userID)
	if err != nil {
		return nil, err
	}
	nodes := make([]MenuNode, 0, len(ms))
	for _, m := range ms {
		if m.Type == "F" {
			continue // 按钮不进路由树
		}
		nodes = append(nodes, MenuNode{
			ID: m.ID, ParentID: m.ParentID, Name: m.Name, Type: m.Type,
			Path: m.Path, Component: m.Component, Perm: m.Perm, Icon: m.Icon, Sort: m.Sort,
			Visible: m.Visible, Status: m.Status,
		})
	}
	result := buildTree(nodes, 0)
	if result == nil {
		result = []MenuNode{}
	}
	return result, nil
}

func buildTree(nodes []MenuNode, parentID uint) []MenuNode {
	var out []MenuNode
	for _, n := range nodes {
		if n.ParentID == parentID {
			n.Children = buildTree(nodes, n.ID)
			out = append(out, n)
		}
	}
	return out
}

// UserPermCodes 返回用户全部按钮/菜单权限码（去重，非空）。
func (s *MenuService) UserPermCodes(userID uint) ([]string, error) {
	ms, err := s.userMenus(userID)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	perms := make([]string, 0)
	for _, m := range ms {
		if m.Perm == "" {
			continue
		}
		if _, ok := seen[m.Perm]; ok {
			continue
		}
		seen[m.Perm] = struct{}{}
		perms = append(perms, m.Perm)
	}
	return perms, nil
}
