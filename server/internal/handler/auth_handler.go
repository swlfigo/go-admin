package handler

import (
	"net/http"
	"strings"

	"go-admin/internal/middleware"
	"go-admin/internal/service"
	"go-admin/pkg/captcha"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

type AuthHandler struct {
	svc       *service.AuthService
	menus     *service.MenuService
	captcha   *captcha.Store
	captchaOn bool
}

func NewAuthHandler(svc *service.AuthService, menus *service.MenuService, cap *captcha.Store, captchaOn bool) *AuthHandler {
	return &AuthHandler{svc: svc, menus: menus, captcha: cap, captchaOn: captchaOn}
}

func (h *AuthHandler) Captcha(c *gin.Context) {
	if !h.captchaOn {
		response.OK(c, gin.H{"enabled": false})
		return
	}
	id, b64, _, err := h.captcha.Generate()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "验证码生成失败")
		return
	}
	response.OK(c, gin.H{"enabled": true, "captchaId": id, "img": b64})
}

type loginReq struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captchaId"`
	CaptchaCode string `json:"captchaCode"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if h.captchaOn && !h.captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		response.Fail(c, http.StatusBadRequest, "验证码错误")
		return
	}
	ua := useragent.New(c.Request.UserAgent())
	browser, _ := ua.Browser()
	out, err := h.svc.Login(req.Username, req.Password, service.LoginMeta{
		IP: c.ClientIP(), Browser: browser, OS: normalizeOS(ua.OS()),
	})
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, out)
}

// normalizeOS 把 UA 里冗长/误导的系统串收敛成简洁名称。
// 注意：macOS 上的浏览器（含 Apple Silicon）UA 一律写死 "Intel Mac OS X 10_15_7"，
// 这是浏览器防指纹的固定值，无法从 UA 区分架构，这里统一显示为 "macOS"。
func normalizeOS(raw string) string {
	s := strings.ToLower(raw)
	switch {
	case strings.Contains(s, "mac os") || strings.Contains(s, "macos"):
		return "macOS"
	case strings.Contains(s, "windows"):
		return "Windows"
	case strings.Contains(s, "android"):
		return "Android"
	case strings.Contains(s, "iphone"), strings.Contains(s, "ipad"), strings.Contains(s, "ios"):
		return "iOS"
	case strings.Contains(s, "linux"):
		return "Linux"
	case raw == "":
		return "未知"
	default:
		return raw
	}
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	out, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, out)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	jti := c.GetString(middleware.CtxJTI)
	_ = h.svc.Logout(jti)
	response.OK(c, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	u, err := h.svc.Me(uid)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "用户不存在")
		return
	}
	perms, err := h.menus.UserPermCodes(uid)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载权限失败")
		return
	}
	response.OK(c, gin.H{"user": u, "perms": perms})
}

func (h *AuthHandler) MenusMe(c *gin.Context) {
	uid := c.GetUint(middleware.CtxUserID)
	tree, err := h.menus.UserMenuTree(uid)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载菜单失败")
		return
	}
	response.OK(c, tree)
}
