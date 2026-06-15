package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOKAndFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	OK(c, gin.H{"hello": "world"})
	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, float64(0), body["code"])
	require.Equal(t, "ok", body["msg"])

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	Fail(c2, http.StatusUnauthorized, "bad token")
	require.Equal(t, http.StatusUnauthorized, w2.Code)
	var body2 map[string]any
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &body2))
	require.Equal(t, "bad token", body2["msg"])
}
