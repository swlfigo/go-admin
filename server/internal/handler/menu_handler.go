package handler

import (
	"net/http"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct{ svc *service.MenuAdminService }

func NewMenuHandler(s *service.MenuAdminService) *MenuHandler { return &MenuHandler{svc: s} }

func (h *MenuHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, tree)
}

type menuReq struct {
	ParentID  uint   `json:"parentId"`
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Perm      string `json:"perm"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	Visible   int    `json:"visible"`
	Status    int    `json:"status"`
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	m, err := h.svc.Create(service.MenuInput(req))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.OK(c, m)
}

func (h *MenuHandler) Update(c *gin.Context) {
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Update(idParam(c), service.MenuInput(req)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(idParam(c)); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, nil)
}
