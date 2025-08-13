package usecase

import (
	"context"

	"backend-context-engineering-template/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductUseCaseInterface interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}
