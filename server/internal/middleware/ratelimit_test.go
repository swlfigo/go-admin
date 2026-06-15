package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRateLimitBlocksAfterBurst(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// 每分钟 2 次：第 3 次应 429
	r.POST("/login", RateLimitByIP(2), func(c *gin.Context) { c.Status(200) })

	do := func() int {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "1.2.3.4:5555"
		r.ServeHTTP(w, req)
		return w.Code
	}
	require.Equal(t, 200, do())
	require.Equal(t, 200, do())
	require.Equal(t, 429, do())
}
