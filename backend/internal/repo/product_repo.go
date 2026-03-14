package repo

import (
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{db: config.DB}
}

func (r *ProductRepo) Create(p *model.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepo) GetByID(id uint) (*model.Product, error) {
	var p model.Product
	if err := r.db.Preload("Category").Preload("Detail").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// ListParams 商品列表查询参数
type ListParams struct {
	CategoryID   *uint
	DeliveryType string
	Status       *int
	Keyword      string
	Page         int
	Size         int
	OrderBy      string // "sales" | "price_asc" | "price_desc" | "latest"
}

func (r *ProductRepo) List(params ListParams) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64

	q := r.db.Model(&model.Product{}).Preload("Category")

	if params.CategoryID != nil {
		q = q.Where("category_id = ?", *params.CategoryID)
	}
	if params.DeliveryType != "" {
		q = q.Where("delivery_type = ?", params.DeliveryType)
	}
	if params.Status != nil {
		q = q.Where("status = ?", *params.Status)
	}
	if params.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+params.Keyword+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch params.OrderBy {
	case "sales":
		q = q.Order("sales DESC")
	case "price_asc":
		q = q.Order("price ASC")
	case "price_desc":
		q = q.Order("price DESC")
	default:
		q = q.Order("sort DESC, id DESC")
	}

	// 分页
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 || params.Size > 100 {
		params.Size = 20
	}
	offset := (params.Page - 1) * params.Size
	if err := q.Offset(offset).Limit(params.Size).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *ProductRepo) Update(p *model.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepo) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// --- ProductDetail ---

func (r *ProductRepo) CreateDetail(d *model.ProductDetail) error {
	return r.db.Create(d).Error
}

func (r *ProductRepo) UpdateDetail(d *model.ProductDetail) error {
	return r.db.Save(d).Error
}

func (r *ProductRepo) GetDetailByProductID(productID uint) (*model.ProductDetail, error) {
	var d model.ProductDetail
	if err := r.db.Where("product_id = ?", productID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
