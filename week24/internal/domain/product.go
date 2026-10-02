package domain

import (
	"errors"
	"time"
)

// Domain errors
var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidInput     = errors.New("invalid input")
)

type Product struct {
	ID          int       `json:"id"          gorm:"primaryKey"`
	Name        string    `json:"name"        gorm:"not null"`
	SKU         string    `json:"sku"         gorm:"uniqueIndex;not null"`
	CategoryID  int       `json:"category_id"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) Validate() error {
	if p.Name == ""    { return fmt.Errorf("%w: name is required", ErrInvalidInput) }
	if p.SKU == ""     { return fmt.Errorf("%w: sku is required",  ErrInvalidInput) }
	if p.Price < 0    { return fmt.Errorf("%w: price cannot be negative", ErrInvalidInput) }
	return nil
}

func (p *Product) AdjustStock(delta int) error {
	newStock := p.Stock + delta
	if newStock < 0 { return fmt.Errorf("%w: have %d, need %d", ErrInsufficientStock, p.Stock, -delta) }
	p.Stock = newStock
	return nil
}