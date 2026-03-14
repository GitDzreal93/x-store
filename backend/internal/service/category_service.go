package service

import (
	"errors"

	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
)

type CategoryService struct {
	repo *repo.CategoryRepo
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repo: repo.NewCategoryRepo()}
}

type CreateCategoryReq struct {
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateCategoryReq struct {
	Name   *string `json:"name"`
	Icon   *string `json:"icon"`
	Sort   *int    `json:"sort"`
	Status *int    `json:"status"`
}

func (s *CategoryService) Create(req CreateCategoryReq) (*model.Category, error) {
	cat := &model.Category{
		Name:     req.Name,
		Icon:     req.Icon,
		Sort:     req.Sort,
		ParentID: req.ParentID,
		Status:   1,
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) GetByID(id uint) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) List(onlyEnabled bool) ([]model.Category, error) {
	var status *int
	if onlyEnabled {
		v := 1
		status = &v
	}
	return s.repo.List(status)
}

func (s *CategoryService) Update(id uint, req UpdateCategoryReq) (*model.Category, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("分类不存在")
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Icon != nil {
		cat.Icon = *req.Icon
	}
	if req.Sort != nil {
		cat.Sort = *req.Sort
	}
	if req.Status != nil {
		cat.Status = *req.Status
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
