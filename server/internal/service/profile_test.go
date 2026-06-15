package service

import (
	"testing"

	"go-admin/internal/repository"
	"go-admin/pkg/password"

	"github.com/stretchr/testify/require"
)

func TestChangeOwnPassword(t *testing.T) {
	db := userSvcDB(t) // 复用 user_service_test 的 helper（含 seed）
	svc := NewUserService(repository.NewUserRepo(db))

	created, err := svc.Create(CreateUserInput{Username: "pw_" + t.Name(), Password: "old123"})
	require.NoError(t, err)

	// 旧密码错 → 失败
	require.ErrorIs(t, svc.ChangeOwnPassword(created.ID, "wrong", "new123"), ErrBadCredential)

	// 旧密码对 → 成功，且新密码可校验
	require.NoError(t, svc.ChangeOwnPassword(created.ID, "old123", "new123"))
	u, err := svc.users.FindByID(created.ID)
	require.NoError(t, err)
	require.True(t, password.Check(u.Password, "new123"))
}

func TestChangeOwnPasswordUserNotFound(t *testing.T) {
	db := userSvcDB(t)
	svc := NewUserService(repository.NewUserRepo(db))

	// Use a non-existent ID (very large number unlikely to exist)
	err := svc.ChangeOwnPassword(99999999, "anypassword", "newpassword")
	require.Error(t, err, "ChangeOwnPassword on a non-existent user id must return an error")
}
