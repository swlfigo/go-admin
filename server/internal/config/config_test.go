package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadReadsYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
server:
  port: 9090
database:
  host: db.local
  port: 5432
  user: u
  password: p
  name: n
  sslmode: disable
jwt:
  secret: s
  accessMinutes: 15
  refreshDays: 3
captcha:
  enabled: false
login:
  maxFailCount: 7
  lockMinutes: 20
  rateLimitPerMinute: 30
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	cfg, err := Load(path)
	require.NoError(t, err)
	require.Equal(t, 9090, cfg.Server.Port)
	require.Equal(t, "db.local", cfg.Database.Host)
	require.Equal(t, 15, cfg.JWT.AccessMinutes)
	require.False(t, cfg.Captcha.Enabled)
	require.Equal(t, 7, cfg.Login.MaxFailCount)
}
