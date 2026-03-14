package repo

import (
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"gorm.io/gorm"
)

type OAuthProviderRepo struct {
	db *gorm.DB
}

func NewOAuthProviderRepo() *OAuthProviderRepo {
	return &OAuthProviderRepo{db: config.DB}
}

func (r *OAuthProviderRepo) List() ([]model.OAuthProvider, error) {
	var providers []model.OAuthProvider
	if err := r.db.Order("sort DESC, id ASC").Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *OAuthProviderRepo) GetByProvider(provider string) (*model.OAuthProvider, error) {
	var p model.OAuthProvider
	if err := r.db.Where("provider = ?", provider).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *OAuthProviderRepo) GetByID(id uint) (*model.OAuthProvider, error) {
	var p model.OAuthProvider
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *OAuthProviderRepo) Update(provider *model.OAuthProvider) error {
	return r.db.Save(provider).Error
}

func (r *OAuthProviderRepo) Create(provider *model.OAuthProvider) error {
	return r.db.Create(provider).Error
}

func (r *OAuthProviderRepo) ListEnabled() ([]model.OAuthProvider, error) {
	var providers []model.OAuthProvider
	if err := r.db.Where("enabled = ?", true).Order("sort DESC, id ASC").Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}
