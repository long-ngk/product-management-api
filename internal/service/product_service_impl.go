package service

import (
	"context"
	"strings"

	"github.com/product-management-server/internal/model"
	"github.com/product-management-server/internal/repository"
)

// productService implements ProductService interface.
type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new ProductService instance.
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// validateProductInput validates the common fields for create and update requests.
// Checks are ordered: name required → name >= 3 chars → price required → price > 0 → quantity required → quantity >= 0.
// Returns the first validation error encountered (first error wins).
func validateProductInput(name string, price *float64, quantity *int) error {
	if strings.TrimSpace(name) == "" {
		return &model.ValidationError{Message: "name is required"}
	}
	if len(name) < 3 {
		return &model.ValidationError{Message: "name must be at least 3 characters"}
	}
	if price == nil {
		return &model.ValidationError{Message: "price is required"}
	}
	if *price <= 0 {
		return &model.ValidationError{Message: "price must be greater than 0"}
	}
	if quantity == nil {
		return &model.ValidationError{Message: "quantity is required"}
	}
	if *quantity < 0 {
		return &model.ValidationError{Message: "quantity must be greater than or equal to 0"}
	}
	return nil
}

// CreateProduct validates input and creates a new product.
func (s *productService) CreateProduct(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
	if err := validateProductInput(req.Name, req.Price, req.Quantity); err != nil {
		return nil, err
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       *req.Price,
		Quantity:    *req.Quantity,
	}

	return s.repo.Create(ctx, product)
}

// GetProducts returns all products or filters by keyword.
func (s *productService) GetProducts(ctx context.Context, keyword string) ([]model.Product, error) {
	if keyword == "" {
		return s.repo.FindAll(ctx)
	}
	return s.repo.FindByKeyword(ctx, keyword)
}

// GetProductByID retrieves a product by ID, mapping not-found errors to NotFoundError.
func (s *productService) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, &model.NotFoundError{Message: "product not found"}
		}
		return nil, err
	}
	return product, nil
}

// UpdateProduct validates input, checks existence, and updates the product.
func (s *productService) UpdateProduct(ctx context.Context, id int, req model.UpdateProductRequest) (*model.Product, error) {
	if err := validateProductInput(req.Name, req.Price, req.Quantity); err != nil {
		return nil, err
	}

	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, &model.NotFoundError{Message: "product not found"}
		}
		return nil, err
	}

	product := &model.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       *req.Price,
		Quantity:    *req.Quantity,
	}

	return s.repo.Update(ctx, product)
}

// DeleteProduct deletes a product by ID, mapping not-found errors to NotFoundError.
func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &model.NotFoundError{Message: "product not found"}
		}
		return err
	}
	return nil
}
