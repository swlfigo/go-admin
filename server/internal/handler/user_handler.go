package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ svc *service.UserService }

func NewUserHandler(s *service.UserService) *UserHandler { return &UserHandler{svc: s} }

func idParam(c *gin.Context) uint {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	users, total, err := h.svc.List(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": users, "total": total})
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleIDs  []uint `json:"roleIds"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	u, err := h.svc.Create(service.CreateUserInput{
		Username: req.Username, Password: req.Password, Nickname: req.Nickname,
		Email: req.Email, Phone: req.Phone, RoleIDs: req.RoleIDs,
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, u)
}

type updateUserReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

func (h *UserHandler) Update(c *gin.Context) {
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.Update(idParam(c), service.UpdateUserInput(req)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(idParam(c)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	response.OK(c, nil)
}

type assignRolesReq struct {
	RoleIDs []uint `json:"roleIds"`
}

func (h *UserHandler) AssignRoles(c *gin.Context) {
	var req assignRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.AssignRoles(idParam(c), req.RoleIDs); err != nil {
		response.Fail(c, http.StatusInternalServerError, "分配失败")
		return
	}
	response.OK(c, nil)
}

type resetPwdReq struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req resetPwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.ResetPassword(idParam(c), req.Password); err != nil {
		response.Fail(c, http.StatusInternalServerError, "重置失败")
		return
	}
	response.OK(c, nil)
}

func (h *UserHandler) Unlock(c *gin.Context) {
	if err := h.svc.Unlock(idParam(c)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "解锁失败")
		return
	}
	response.OK(c, nil)
}
