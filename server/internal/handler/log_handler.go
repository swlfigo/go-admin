package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type LogHandler struct{ svc *service.LogService }

func NewLogHandler(s *service.LogService) *LogHandler { return &LogHandler{svc: s} }

func (h *LogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := h.svc.List(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list, "total": total})
}

func (h *LogHandler) Clear(c *gin.Context) {
	if err := h.svc.Clear(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作失败")
		return
	}
	response.OK(c, nil)
}
