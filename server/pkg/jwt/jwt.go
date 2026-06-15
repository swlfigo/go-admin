package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"uname"`
	Kind     string `json:"kind"` // access | refresh
	jwt.RegisteredClaims
}

type Manager struct {
	secret  []byte
	access  time.Duration
	refresh time.Duration
}

func New(secret string, access, refresh time.Duration) *Manager {
	return &Manager{secret: []byte(secret), access: access, refresh: refresh}
}

func (m *Manager) issue(uid uint, uname, kind string, ttl time.Duration) (string, *Claims, error) {
	now := time.Now()
	claims := &Claims{
		UserID: uid, Username: uname, Kind: kind,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "go-admin",
		},
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	return tok, claims, err
}

func (m *Manager) IssueAccess(uid uint, uname string) (string, *Claims, error) {
	return m.issue(uid, uname, "access", m.access)
}

func (m *Manager) IssueRefresh(uid uint, uname string) (string, *Claims, error) {
	return m.issue(uid, uname, "refresh", m.refresh)
}

func (m *Manager) parse(tok, wantKind string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tok, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer("go-admin"))
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid || claims.Kind != wantKind {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *Manager) ParseAccess(tok string) (*Claims, error)  { return m.parse(tok, "access") }
func (m *Manager) ParseRefresh(tok string) (*Claims, error) { return m.parse(tok, "refresh") }
