package service

import (
	"context"
	"fmt"

	"github.com/OluwaTobba/golang-siwes-training/week24/internal/domain"
)

type ProductRepo interface {
	Create(ctx context.Context, p *domain.Product) error
	GetByID(ctx context.Context, id int) (*domain.Product, error)
	List(ctx context.Context) ([]*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id int) error
}

type InventoryService struct{ repo ProductRepo }

func NewInventoryService(repo ProductRepo) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) CreateProduct(ctx context.Context, p *domain.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, p)
}

func (s *InventoryService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *InventoryService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.List(ctx)
}

func (s *InventoryService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *InventoryService) AdjustStock(ctx context.Context, id, delta int) (*domain.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.AdjustStock(delta); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("updating stock: %w", err)
	}
	return p, nil
}