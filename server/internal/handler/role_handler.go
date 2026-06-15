package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct{ svc *service.RoleService }

func NewRoleHandler(s *service.RoleService) *RoleHandler { return &RoleHandler{svc: s} }

func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	roles, total, err := h.svc.List(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": roles, "total": total})
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.svc.Get(idParam(c))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "角色不存在")
		return
	}
	response.OK(c, role)
}

type roleReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	r, err := h.svc.Create(service.RoleInput(req))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, r)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Update(idParam(c), service.RoleInput(req)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(idParam(c)); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, nil)
}

type assignMenusReq struct {
	MenuIDs []uint `json:"menuIds"`
}

func (h *RoleHandler) AssignMenus(c *gin.Context) {
	var req assignMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.AssignMenus(idParam(c), req.MenuIDs); err != nil {
		response.Fail(c, http.StatusInternalServerError, "分配失败")
		return
	}
	response.OK(c, nil)
}
