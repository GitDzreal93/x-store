package repo

import (
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{db: config.DB}
}

func (r *CategoryRepo) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepo) GetByID(id uint) (*model.Category, error) {
	var cat model.Category
	if err := r.db.First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepo) List(status *int) ([]model.Category, error) {
	var list []model.Category
	q := r.db.Order("sort DESC, id ASC")
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CategoryRepo) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
