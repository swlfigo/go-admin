package service

import (
	"errors"

	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/password"

	"gorm.io/gorm"
)

var ErrUsernameTaken = errors.New("用户名已存在")

type CreateUserInput struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	RoleIDs  []uint
}

type UpdateUserInput struct {
	Nickname string
	Email    string
	Phone    string
	Status   int
}

type UserService struct{ users *repository.UserRepo }

func NewUserService(u *repository.UserRepo) *UserService { return &UserService{users: u} }

func (s *UserService) List(keyword string, page, size int) ([]model.User, int64, error) {
	return s.users.ListUsers(keyword, page, size)
}

func (s *UserService) Create(in CreateUserInput) (*model.User, error) {
	if _, err := s.users.FindByUsername(in.Username); err == nil {
		return nil, ErrUsernameTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashed, err := password.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Username: in.Username, Password: hashed, Nickname: in.Nickname,
		Email: in.Email, Phone: in.Phone, Status: 1,
	}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	if len(in.RoleIDs) > 0 {
		if err := s.users.AssignRoles(u.ID, in.RoleIDs); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (s *UserService) Update(id uint, in UpdateUserInput) error {
	return s.users.UpdateProfile(id, map[string]any{
		"nickname": in.Nickname, "email": in.Email, "phone": in.Phone, "status": in.Status,
	})
}

func (s *UserService) Delete(id uint) error { return s.users.Delete(id) }

func (s *UserService) ResetPassword(id uint, newPlain string) error {
	hashed, err := password.Hash(newPlain)
	if err != nil {
		return err
	}
	return s.users.SetPassword(id, hashed)
}

func (s *UserService) AssignRoles(id uint, roleIDs []uint) error {
	return s.users.AssignRoles(id, roleIDs)
}

// Unlock 清空失败计数与锁定（超管手动解锁）。
func (s *UserService) Unlock(id uint) error { return s.users.ResetFailCount(id) }

// ChangeOwnPassword 校验旧密码后改为新密码（个人中心用）。
func (s *UserService) ChangeOwnPassword(id uint, oldPlain, newPlain string) error {
	u, err := s.users.FindByID(id)
	if err != nil {
		return err
	}
	if !password.Check(u.Password, oldPlain) {
		return ErrBadCredential
	}
	hashed, err := password.Hash(newPlain)
	if err != nil {
		return err
	}
	return s.users.SetPassword(id, hashed)
}
