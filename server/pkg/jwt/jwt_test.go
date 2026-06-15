package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIssueAndParseAccess(t *testing.T) {
	m := New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	tok, claims, err := m.IssueAccess(7, "admin")
	require.NoError(t, err)
	require.NotEmpty(t, tok)
	require.NotEmpty(t, claims.ID) // jti

	parsed, err := m.ParseAccess(tok)
	require.NoError(t, err)
	require.Equal(t, uint(7), parsed.UserID)
	require.Equal(t, "admin", parsed.Username)
	require.Equal(t, claims.ID, parsed.ID)
	require.Equal(t, "access", parsed.Kind)
}

func TestParseRejectsTampered(t *testing.T) {
	m := New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	_, err := m.ParseAccess("not.a.jwt")
	require.Error(t, err)
}

func TestRefreshKind(t *testing.T) {
	m := New("test-secret-at-least-16", 30*time.Minute, 24*time.Hour)
	tok, claims, err := m.IssueRefresh(7, "admin")
	require.NoError(t, err)
	parsed, err := m.ParseRefresh(tok)
	require.NoError(t, err)
	require.Equal(t, "refresh", parsed.Kind)
	require.Equal(t, claims.ID, parsed.ID)
}
