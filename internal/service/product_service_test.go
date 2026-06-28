package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/product-management-server/internal/model"
	"github.com/product-management-server/internal/repository"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

// --- Shared mocks ---

// mockProductRepository is an in-memory mock implementation of repository.ProductRepository.
type mockProductRepository struct {
	products map[int]*model.Product
	nextID   int
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[int]*model.Product),
		nextID:   1,
	}
}

func (m *mockProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.ID = m.nextID
	m.nextID++
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
	stored := *product
	m.products[product.ID] = &stored
	return product, nil
}

func (m *mockProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	var result []model.Product
	for _, p := range m.products {
		result = append(result, *p)
	}
	return result, nil
}

func (m *mockProductRepository) FindByKeyword(ctx context.Context, keyword string) ([]model.Product, error) {
	return nil, nil
}

func (m *mockProductRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, fmt.Errorf("product not found")
	}
	result := *p
	return &result, nil
}

func (m *mockProductRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	existing, ok := m.products[product.ID]
	if !ok {
		return nil, fmt.Errorf("product not found")
	}
	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now().UTC()
	stored := *product
	m.products[product.ID] = &stored
	return product, nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id int) error {
	if _, ok := m.products[id]; !ok {
		return fmt.Errorf("product not found")
	}
	delete(m.products, id)
	return nil
}

// panicRepository is a mock repository that panics if any method is called.
// This proves validation rejects input BEFORE any repository method is reached.
type panicRepository struct{}

var _ repository.ProductRepository = (*panicRepository)(nil)

func (r *panicRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	panic("Create should not be called when validation fails")
}

func (r *panicRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	panic("FindAll should not be called when validation fails")
}

func (r *panicRepository) FindByKeyword(ctx context.Context, keyword string) ([]model.Product, error) {
	panic("FindByKeyword should not be called when validation fails")
}

func (r *panicRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	panic("FindByID should not be called when validation fails")
}

func (r *panicRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	panic("Update should not be called when validation fails")
}

func (r *panicRepository) Delete(ctx context.Context, id int) error {
	panic("Delete should not be called when validation fails")
}


// --- Property 1 Tests ---

// TestProperty1_CreateProduct_ValidationRejectsAllInvalidInputs tests that CreateProduct
// rejects all invalid inputs with appropriate ValidationError.
func TestProperty1_CreateProduct_ValidationRejectsAllInvalidInputs(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		repo := &panicRepository{}
		svc := NewProductService(repo)

		// Choose which validation rule to violate
		violation := rapid.IntRange(0, 5).Draw(t, "violation")

		var req model.CreateProductRequest
		validPrice := rapid.Float64Range(0.01, 999999.99).Draw(t, "validPrice")
		validQty := rapid.IntRange(0, 100000).Draw(t, "validQty")

		switch violation {
		case 0: // empty name
			req = model.CreateProductRequest{Name: "", Price: &validPrice, Quantity: &validQty}
		case 1: // name too short (1-2 chars)
			shortName := rapid.StringMatching(`[a-zA-Z]{1,2}`).Draw(t, "shortName")
			req = model.CreateProductRequest{Name: shortName, Price: &validPrice, Quantity: &validQty}
		case 2: // nil price
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			req = model.CreateProductRequest{Name: validName, Price: nil, Quantity: &validQty}
		case 3: // price <= 0
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			badPrice := rapid.Float64Range(-9999.0, 0.0).Draw(t, "badPrice")
			req = model.CreateProductRequest{Name: validName, Price: &badPrice, Quantity: &validQty}
		case 4: // nil quantity
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			req = model.CreateProductRequest{Name: validName, Price: &validPrice, Quantity: nil}
		case 5: // negative quantity
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			badQty := rapid.IntRange(-10000, -1).Draw(t, "badQty")
			req = model.CreateProductRequest{Name: validName, Price: &validPrice, Quantity: &badQty}
		}

		_, err := svc.CreateProduct(context.Background(), req)
		assert.Error(t, err, "CreateProduct should return error for invalid input")
		var valErr *model.ValidationError
		assert.ErrorAs(t, err, &valErr, "error should be ValidationError")
	})
}

// TestProperty1_UpdateProduct_ValidationRejectsAllInvalidInputs tests that UpdateProduct
// rejects all invalid inputs with appropriate ValidationError.
func TestProperty1_UpdateProduct_ValidationRejectsAllInvalidInputs(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		repo := &panicRepository{}
		svc := NewProductService(repo)

		violation := rapid.IntRange(0, 5).Draw(t, "violation")

		var req model.UpdateProductRequest
		validPrice := rapid.Float64Range(0.01, 999999.99).Draw(t, "validPrice")
		validQty := rapid.IntRange(0, 100000).Draw(t, "validQty")

		switch violation {
		case 0:
			req = model.UpdateProductRequest{Name: "", Price: &validPrice, Quantity: &validQty}
		case 1:
			shortName := rapid.StringMatching(`[a-zA-Z]{1,2}`).Draw(t, "shortName")
			req = model.UpdateProductRequest{Name: shortName, Price: &validPrice, Quantity: &validQty}
		case 2:
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			req = model.UpdateProductRequest{Name: validName, Price: nil, Quantity: &validQty}
		case 3:
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			badPrice := rapid.Float64Range(-9999.0, 0.0).Draw(t, "badPrice")
			req = model.UpdateProductRequest{Name: validName, Price: &badPrice, Quantity: &validQty}
		case 4:
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			req = model.UpdateProductRequest{Name: validName, Price: &validPrice, Quantity: nil}
		case 5:
			validName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "validName")
			badQty := rapid.IntRange(-10000, -1).Draw(t, "badQty")
			req = model.UpdateProductRequest{Name: validName, Price: &validPrice, Quantity: &badQty}
		}

		_, err := svc.UpdateProduct(context.Background(), 1, req)
		assert.Error(t, err, "UpdateProduct should return error for invalid input")
		var valErr *model.ValidationError
		assert.ErrorAs(t, err, &valErr, "error should be ValidationError")
	})
}

// --- Property 5 Test ---

// TestProperty5_UpdateReflectsNewValues tests Property 5: Update reflects new values.
// For any existing product and any valid UpdateProductRequest, calling UpdateProduct
// SHALL return a Product whose name, description, price, and quantity match the request
// values, and whose id remains unchanged.
func TestProperty5_UpdateReflectsNewValues(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Set up mock repository
		repo := newMockProductRepository()
		svc := NewProductService(repo)
		ctx := context.Background()

		// Generate an existing product and pre-populate the repo
		existingName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "existingName")
		existingDesc := rapid.Ptr(rapid.StringMatching(`[a-zA-Z0-9 ]{0,50}`), true).Draw(t, "existingDesc")
		existingPrice := rapid.Float64Range(0.01, 99999.99).Draw(t, "existingPrice")
		existingQty := rapid.IntRange(0, 10000).Draw(t, "existingQty")

		existingProduct := &model.Product{
			ID:          1,
			Name:        existingName,
			Description: existingDesc,
			Price:       existingPrice,
			Quantity:    existingQty,
			CreatedAt:   time.Now().UTC().Add(-time.Hour),
			UpdatedAt:   time.Now().UTC().Add(-time.Hour),
		}
		repo.products[1] = existingProduct
		repo.nextID = 2

		// Generate a valid UpdateProductRequest with new values
		newName := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "newName")
		newDesc := rapid.Ptr(rapid.StringMatching(`[a-zA-Z0-9 ]{0,50}`), true).Draw(t, "newDesc")
		newPrice := rapid.Float64Range(0.01, 99999.99).Draw(t, "newPrice")
		newQty := rapid.IntRange(0, 10000).Draw(t, "newQty")

		updateReq := model.UpdateProductRequest{
			Name:        newName,
			Description: newDesc,
			Price:       &newPrice,
			Quantity:    &newQty,
		}

		// Call UpdateProduct
		result, err := svc.UpdateProduct(ctx, 1, updateReq)

		// Assertions
		if err != nil {
			t.Fatalf("UpdateProduct returned unexpected error: %v", err)
		}

		if result.Name != newName {
			t.Fatalf("expected name %q, got %q", newName, result.Name)
		}
		if (result.Description == nil) != (newDesc == nil) {
			t.Fatalf("description nil mismatch: expected %v, got %v", newDesc, result.Description)
		}
		if result.Description != nil && newDesc != nil && *result.Description != *newDesc {
			t.Fatalf("expected description %q, got %q", *newDesc, *result.Description)
		}
		if result.Price != newPrice {
			t.Fatalf("expected price %v, got %v", newPrice, result.Price)
		}
		if result.Quantity != newQty {
			t.Fatalf("expected quantity %d, got %d", newQty, result.Quantity)
		}

		// Returned product ID is unchanged
		if result.ID != 1 {
			t.Fatalf("expected id %d to remain unchanged, got %d", 1, result.ID)
		}
	})
}

// --- Property 2 Test ---

// TestProperty2_ValidInputProducesCompleteProduct verifies that for any valid
// CreateProductRequest (name >= 3 chars, price > 0, quantity >= 0), calling CreateProduct
// returns a Product with same name, description, price, quantity; positive id; non-zero timestamps.
func TestProperty2_ValidInputProducesCompleteProduct(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		repo := newMockProductRepository()
		svc := NewProductService(repo)

		// Generate valid name (3+ characters)
		name := rapid.StringMatching(`[a-zA-Z0-9]{3,50}`).Draw(t, "name")

		// Generate valid price (> 0)
		price := rapid.Float64Range(0.01, 999999.99).Draw(t, "price")

		// Generate valid quantity (>= 0)
		quantity := rapid.IntRange(0, 100000).Draw(t, "quantity")

		// Generate optional description (nil or non-nil)
		hasDescription := rapid.Bool().Draw(t, "hasDescription")
		var description *string
		if hasDescription {
			desc := rapid.StringMatching(`[a-zA-Z0-9 ]{0,200}`).Draw(t, "description")
			description = &desc
		}

		req := model.CreateProductRequest{
			Name:        name,
			Description: description,
			Price:       &price,
			Quantity:    &quantity,
		}

		// Execute
		product, err := svc.CreateProduct(context.Background(), req)

		// Assert: no error
		assert.NoError(t, err, "CreateProduct should not return error for valid input")
		assert.NotNil(t, product, "CreateProduct should return a non-nil product")

		if product == nil {
			return
		}

		// Assert: returned product has same name, description, price, quantity as request
		assert.Equal(t, name, product.Name, "product name should match request")
		assert.Equal(t, description, product.Description, "product description should match request")
		assert.Equal(t, price, product.Price, "product price should match request")
		assert.Equal(t, quantity, product.Quantity, "product quantity should match request")

		// Assert: returned product has positive id (> 0)
		assert.Greater(t, product.ID, 0, "product ID should be positive")

		// Assert: returned product has non-zero created_at and updated_at
		assert.False(t, product.CreatedAt.IsZero(), "created_at should be non-zero")
		assert.False(t, product.UpdatedAt.IsZero(), "updated_at should be non-zero")
	})
}

// --- Property 3 Test ---

// TestProperty3_CreateThenGetRoundTrip tests Property 3: Create then get by ID round-trip.
// For any valid product that is successfully created, calling GetProductByID with the returned id
// SHALL return a Product with identical field values (name, description, price, quantity).
func TestProperty3_CreateThenGetRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Set up a fresh in-memory repo and service for each iteration
		repo := newMockProductRepository()
		svc := NewProductService(repo)
		ctx := context.Background()

		// Generate valid CreateProductRequest
		name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9 ]{2,49}`).Draw(t, "name")
		price := rapid.Float64Range(0.01, 99999.99).Draw(t, "price")
		quantity := rapid.IntRange(0, 10000).Draw(t, "quantity")

		var description *string
		if rapid.Bool().Draw(t, "hasDescription") {
			desc := rapid.StringMatching(`[a-zA-Z0-9 ]{0,200}`).Draw(t, "description")
			description = &desc
		}

		req := model.CreateProductRequest{
			Name:        name,
			Description: description,
			Price:       &price,
			Quantity:    &quantity,
		}

		// Step 1: CreateProduct
		created, err := svc.CreateProduct(ctx, req)
		assert.NoError(t, err, "CreateProduct should not return error for valid input")
		assert.NotNil(t, created, "Created product should not be nil")
		if created == nil {
			return
		}

		// Step 2: GetProductByID with the returned id
		retrieved, err := svc.GetProductByID(ctx, created.ID)
		assert.NoError(t, err, "GetProductByID should not return error for existing product")
		assert.NotNil(t, retrieved, "Retrieved product should not be nil")
		if retrieved == nil {
			return
		}

		// Assert: identical field values after round-trip
		assert.Equal(t, created.Name, retrieved.Name, "Name should be identical after round-trip")
		assert.Equal(t, created.Description, retrieved.Description, "Description should be identical after round-trip")
		assert.Equal(t, created.Price, retrieved.Price, "Price should be identical after round-trip")
		assert.Equal(t, created.Quantity, retrieved.Quantity, "Quantity should be identical after round-trip")
		assert.Equal(t, created.ID, retrieved.ID, "ID should be identical after round-trip")
	})
}

// --- Property 6 Test ---

// TestProperty6_DeleteMakesProductUnretrievable tests Property 6: Delete makes product unretrievable.
// For any product that exists in the repository, after calling DeleteProduct with that product's id,
// calling GetProductByID with the same id SHALL return a NotFoundError.
func TestProperty6_DeleteMakesProductUnretrievable(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Set up a fresh in-memory repo and service
		repo := newMockProductRepository()
		svc := NewProductService(repo)
		ctx := context.Background()

		// Generate a random product id and pre-populate the repo
		id := rapid.IntRange(1, 100000).Draw(t, "productID")
		name := rapid.StringMatching(`[a-zA-Z]{3,20}`).Draw(t, "name")
		price := rapid.Float64Range(0.01, 99999.99).Draw(t, "price")
		quantity := rapid.IntRange(0, 10000).Draw(t, "quantity")

		var description *string
		if rapid.Bool().Draw(t, "hasDescription") {
			desc := rapid.StringMatching(`[a-zA-Z0-9 ]{0,100}`).Draw(t, "description")
			description = &desc
		}

		// Pre-populate the mock repo with the product
		repo.products[id] = &model.Product{
			ID:          id,
			Name:        name,
			Description: description,
			Price:       price,
			Quantity:    quantity,
			CreatedAt:   time.Now().UTC().Add(-time.Hour),
			UpdatedAt:   time.Now().UTC().Add(-time.Hour),
		}

		// Step 1: DeleteProduct should succeed
		err := svc.DeleteProduct(ctx, id)
		assert.NoError(t, err, "DeleteProduct should not return error for existing product")

		// Step 2: GetProductByID should return NotFoundError
		_, err = svc.GetProductByID(ctx, id)
		assert.Error(t, err, "GetProductByID should return error after deletion")

		var notFoundErr *model.NotFoundError
		assert.ErrorAs(t, err, &notFoundErr, "error should be NotFoundError after deletion")
	})
}
