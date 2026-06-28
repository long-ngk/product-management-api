package repository

import (
	"context"

	"github.com/product-management-server/internal/model"
)

// ProductRepository defines the contract between Service and Repository layers.
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByKeyword(ctx context.Context, keyword string) ([]model.Product, error)
	FindByID(ctx context.Context, id int) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) (*model.Product, error)
	Delete(ctx context.Context, id int) error
}
