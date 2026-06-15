package service

import (
	"errors"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/jwt"
	"go-admin/pkg/password"
)

var (
	ErrBadCredential = errors.New("用户名或密码错误")
	ErrAccountLocked = errors.New("账号已锁定，请稍后再试")
	ErrUserDisabled  = errors.New("账号已停用")
	ErrInvalidToken  = errors.New("无效的令牌")
)

type LoginMeta struct {
	IP      string
	Browser string
	OS      string
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthService struct {
	users    *repository.UserRepo
	sessions *repository.SessionRepo
	jwt      *jwt.Manager
	loginCfg config.LoginConfig
}

func NewAuthService(u *repository.UserRepo, s *repository.SessionRepo, j *jwt.Manager, lc config.LoginConfig) *AuthService {
	return &AuthService{users: u, sessions: s, jwt: j, loginCfg: lc}
}

func (a *AuthService) Login(username, plain string, meta LoginMeta) (*TokenPair, error) {
	u, err := a.users.FindByUsername(username)
	if err != nil {
		return nil, ErrBadCredential
	}
	if u.LockUntil != nil && u.LockUntil.After(time.Now()) {
		return nil, ErrAccountLocked
	}
	if u.Status == 0 {
		return nil, ErrUserDisabled
	}
	if !password.Check(u.Password, plain) {
		_ = a.users.IncFailCount(u.ID)
		if u.LoginFailCount+1 >= a.loginCfg.MaxFailCount {
			_ = a.users.Lock(u.ID, time.Now().Add(time.Duration(a.loginCfg.LockMinutes)*time.Minute))
		}
		return nil, ErrBadCredential
	}
	_ = a.users.ResetFailCount(u.ID)
	_ = a.users.UpdateLastLogin(u.ID)
	return a.issueSession(u, meta, "")
}

// issueSession 签发 access+refresh 并落 session；oldID!="" 时走轮换替换。
func (a *AuthService) issueSession(u *model.User, meta LoginMeta, oldID string) (*TokenPair, error) {
	access, ac, err := a.jwt.IssueAccess(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	refresh, rc, err := a.jwt.IssueRefresh(u.ID, u.Username)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	sess := &model.OnlineSession{
		ID: ac.ID, UserID: u.ID, Username: u.Username,
		LoginIP: meta.IP, Browser: meta.Browser, OS: meta.OS,
		RefreshID: rc.ID, RefreshExpireAt: rc.ExpiresAt.Time,
		LoginAt: now, LastActiveAt: now,
	}
	if oldID != "" {
		if err := a.sessions.Rotate(oldID, sess); err != nil {
			return nil, err
		}
	} else if err := a.sessions.Create(sess); err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (a *AuthService) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := a.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	// refresh 必须对应一条仍有效的会话（防重放：轮换后旧 refresh_id 已不存在）
	sess, err := a.sessions.GetByRefreshID(claims.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	u, err := a.users.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return a.issueSession(u, LoginMeta{IP: sess.LoginIP, Browser: sess.Browser, OS: sess.OS}, sess.ID)
}

func (a *AuthService) Logout(jti string) error {
	return a.sessions.Delete(jti)
}

func (a *AuthService) Me(userID uint) (*model.User, error) {
	return a.users.FindByID(userID)
}
