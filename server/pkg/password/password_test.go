package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndCheck(t *testing.T) {
	h, err := Hash("admin123")
	require.NoError(t, err)
	require.NotEqual(t, "admin123", h)
	require.True(t, Check(h, "admin123"))
	require.False(t, Check(h, "wrong"))
}
