package captcha

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndVerify(t *testing.T) {
	s := NewStore()
	id, b64, answer, err := s.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Contains(t, b64, "data:image")
	// 正确答案校验通过（且一次性，校验后即删除）
	require.True(t, s.Verify(id, answer))
	require.False(t, s.Verify(id, answer)) // 已被消费
}

func TestVerifyWrong(t *testing.T) {
	s := NewStore()
	id, _, _, err := s.Generate()
	require.NoError(t, err)
	require.False(t, s.Verify(id, "definitely-wrong"))
}
