package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{svc: service.NewCategoryService()}
}

// Create 创建分类 [POST /api/admin/categories]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req service.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	cat, err := h.svc.Create(req)
	if err != nil {
		response.ServerError(c, "创建分类失败: "+err.Error())
		return
	}
	response.OK(c, cat)
}

// Get 获取分类详情 [GET /api/categories/:id]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	cat, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}
	response.OK(c, cat)
}

// List 获取分类列表 [GET /api/categories]
func (h *CategoryHandler) List(c *gin.Context) {
	// 管理端传 all=true 可查看全部，否则只返回启用的
	onlyEnabled := c.Query("all") != "true"
	list, err := h.svc.List(onlyEnabled)
	if err != nil {
		response.ServerError(c, "获取分类列表失败")
		return
	}
	response.OK(c, list)
}

// Update 更新分类 [PUT /api/admin/categories/:id]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req service.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	cat, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, cat)
}

// Delete 删除分类 [DELETE /api/admin/categories/:id]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, "删除分类失败")
		return
	}
	response.OKMsg(c, "删除成功")
}
