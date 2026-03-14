package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/internal/service"
	"github.com/x-store/backend/pkg/response"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{svc: service.NewProductService()}
}

// Create 创建商品 [POST /api/admin/products]
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	p, err := h.svc.Create(req)
	if err != nil {
		response.ServerError(c, "创建商品失败: "+err.Error())
		return
	}
	response.OK(c, p)
}

// Get 获取商品详情 [GET /api/products/:id]
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	p, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "商品不存在")
		return
	}
	response.OK(c, p)
}

// List 获取商品列表 [GET /api/products]
func (h *ProductHandler) List(c *gin.Context) {
	var params repo.ListParams

	if cid := c.Query("category_id"); cid != "" {
		if v, err := strconv.ParseUint(cid, 10, 64); err == nil {
			id := uint(v)
			params.CategoryID = &id
		}
	}
	params.DeliveryType = c.Query("delivery_type")
	params.Keyword = c.Query("keyword")
	params.OrderBy = c.Query("order_by")

	if c.Query("all") != "true" {
		v := 1
		params.Status = &v
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	params.Page = page
	params.Size = size

	list, total, err := h.svc.List(params)
	if err != nil {
		response.ServerError(c, "获取商品列表失败")
		return
	}
	response.OKPage(c, list, total, page, size)
}

// Update 更新商品 [PUT /api/admin/products/:id]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req service.UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	p, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, p)
}

// Delete 删除商品 [DELETE /api/admin/products/:id]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, "删除商品失败")
		return
	}
	response.OKMsg(c, "删除成功")
}
