package handler

import (
	"net/http"

	"go-admin/internal/middleware"
	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{ svc *service.UserService }

func NewProfileHandler(s *service.UserService) *ProfileHandler { return &ProfileHandler{svc: s} }

type profileReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// UpdateProfile 改自己的昵称/邮箱/电话（不改状态/角色）。
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req profileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	uid := c.GetUint(middleware.CtxUserID)
	if err := h.svc.Update(uid, service.UpdateUserInput{
		Nickname: req.Nickname, Email: req.Email, Phone: req.Phone, Status: 1,
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

type changePwdReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	var req changePwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	uid := c.GetUint(middleware.CtxUserID)
	if err := h.svc.ChangeOwnPassword(uid, req.OldPassword, req.NewPassword); err != nil {
		response.Fail(c, http.StatusBadRequest, "原密码错误或修改失败")
		return
	}
	response.OK(c, nil)
}
