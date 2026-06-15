package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type DictHandler struct{ svc *service.DictService }

func NewDictHandler(s *service.DictService) *DictHandler { return &DictHandler{svc: s} }

func (h *DictHandler) ListTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	types, total, err := h.svc.ListTypes(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": types, "total": total})
}

type dictTypeReq struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

func (h *DictHandler) CreateType(c *gin.Context) {
	var req dictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	dt, err := h.svc.CreateType(service.DictTypeInput(req))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, dt)
}

func (h *DictHandler) UpdateType(c *gin.Context) {
	var req dictTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.UpdateType(idParam(c), service.DictTypeInput(req)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

func (h *DictHandler) DeleteType(c *gin.Context) {
	if err := h.svc.DeleteType(idParam(c)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	response.OK(c, nil)
}

func (h *DictHandler) DataByType(c *gin.Context) {
	data, err := h.svc.DataByType(c.Query("type"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, data)
}

type dictDataReq struct {
	DictType string `json:"dictType" binding:"required"`
	Label    string `json:"label" binding:"required"`
	Value    string `json:"value" binding:"required"`
	TagType  string `json:"tagType"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

func (h *DictHandler) CreateData(c *gin.Context) {
	var req dictDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	d, err := h.svc.CreateData(service.DictDataInput(req))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.OK(c, d)
}

func (h *DictHandler) UpdateData(c *gin.Context) {
	var req dictDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.UpdateData(idParam(c), service.DictDataInput(req)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}
	response.OK(c, nil)
}

func (h *DictHandler) DeleteData(c *gin.Context) {
	if err := h.svc.DeleteData(idParam(c)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	response.OK(c, nil)
}
