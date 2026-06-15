package handler

import (
	"net/http"

	"go-admin/internal/repository"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type OnlineHandler struct{ sessions *repository.SessionRepo }

func NewOnlineHandler(s *repository.SessionRepo) *OnlineHandler { return &OnlineHandler{sessions: s} }

func (h *OnlineHandler) List(c *gin.Context) {
	list, err := h.sessions.List()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, list)
}

// Kick 删除会话 → 对应 token 下次请求即 401（即时踢人）。
func (h *OnlineHandler) Kick(c *gin.Context) {
	if err := h.sessions.Delete(c.Param("id")); err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作失败")
		return
	}
	response.OK(c, nil)
}
