package service

import (
	"errors"

	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
)

type ProductService struct {
	repo *repo.ProductRepo
}

func NewProductService() *ProductService {
	return &ProductService{repo: repo.NewProductRepo()}
}

type CreateProductReq struct {
	CategoryID    uint    `json:"category_id" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Cover         string  `json:"cover"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	OriginalPrice float64 `json:"original_price"`
	DeliveryType  string  `json:"delivery_type"`
	Tags          string  `json:"tags"`
	Sort          int     `json:"sort"`
	IsNew         bool    `json:"is_new"`
	IsHot         bool    `json:"is_hot"`
	// 详情
	Description string `json:"description"`
	Notice      string `json:"notice"`
}

type UpdateProductReq struct {
	CategoryID    *uint    `json:"category_id"`
	Title         *string  `json:"title"`
	Cover         *string  `json:"cover"`
	Price         *float64 `json:"price"`
	OriginalPrice *float64 `json:"original_price"`
	DeliveryType  *string  `json:"delivery_type"`
	Tags          *string  `json:"tags"`
	Sort          *int     `json:"sort"`
	Status        *int     `json:"status"`
	IsNew         *bool    `json:"is_new"`
	IsHot         *bool    `json:"is_hot"`
	// 详情
	Description *string `json:"description"`
	Notice      *string `json:"notice"`
}

func (s *ProductService) Create(req CreateProductReq) (*model.Product, error) {
	if req.DeliveryType == "" {
		req.DeliveryType = "auto"
	}
	p := &model.Product{
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Cover:         req.Cover,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		DeliveryType:  req.DeliveryType,
		Tags:          req.Tags,
		Sort:          req.Sort,
		Status:        1,
		IsNew:         req.IsNew,
		IsHot:         req.IsHot,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	// 创建商品详情
	detail := &model.ProductDetail{
		ProductID:   p.ID,
		Description: req.Description,
		Notice:      req.Notice,
	}
	if err := s.repo.CreateDetail(detail); err != nil {
		return nil, err
	}
	p.Detail = detail

	return p, nil
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) List(params repo.ListParams) ([]model.Product, int64, error) {
	return s.repo.List(params)
}

func (s *ProductService) Update(id uint, req UpdateProductReq) (*model.Product, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}

	if req.CategoryID != nil {
		p.CategoryID = *req.CategoryID
	}
	if req.Title != nil {
		p.Title = *req.Title
	}
	if req.Cover != nil {
		p.Cover = *req.Cover
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.OriginalPrice != nil {
		p.OriginalPrice = *req.OriginalPrice
	}
	if req.DeliveryType != nil {
		p.DeliveryType = *req.DeliveryType
	}
	if req.Tags != nil {
		p.Tags = *req.Tags
	}
	if req.Sort != nil {
		p.Sort = *req.Sort
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.IsNew != nil {
		p.IsNew = *req.IsNew
	}
	if req.IsHot != nil {
		p.IsHot = *req.IsHot
	}
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	// 更新详情
	if req.Description != nil || req.Notice != nil {
		detail, _ := s.repo.GetDetailByProductID(id)
		if detail == nil {
			detail = &model.ProductDetail{ProductID: id}
		}
		if req.Description != nil {
			detail.Description = *req.Description
		}
		if req.Notice != nil {
			detail.Notice = *req.Notice
		}
		if detail.ID == 0 {
			s.repo.CreateDetail(detail)
		} else {
			s.repo.UpdateDetail(detail)
		}
		p.Detail = detail
	}

	return p, nil
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
