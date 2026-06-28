package service

import (
	"context"

	"github.com/product-management-server/internal/model"
)

// ProductService defines the contract between Handler and Service layers.
type ProductService interface {
	CreateProduct(ctx context.Context, req model.CreateProductRequest) (*model.Product, error)
	GetProducts(ctx context.Context, keyword string) ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (*model.Product, error)
	UpdateProduct(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
