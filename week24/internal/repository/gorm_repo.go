package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
	"gorm.io/gorm"
)

type GORMProductRepo struct{ db *gorm.DB }

func NewGORMProductRepo(db *gorm.DB) *GORMProductRepo {
	return &GORMProductRepo{db: db}
}

func (r *GORMProductRepo) Create(ctx context.Context, p *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("%w: sku %q already exists", domain.ErrInvalidInput, p.SKU)
		}
		return err
	}
	return nil
}

func (r *GORMProductRepo) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var p domain.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *GORMProductRepo) List(ctx context.Context) ([]*domain.Product, error) {
	var ps []*domain.Product
	if err := r.db.WithContext(ctx).Order("id").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (r *GORMProductRepo) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *GORMProductRepo) Delete(ctx context.Context, id int) error {
	res := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}