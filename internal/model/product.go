package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Product represents the domain model for a product.
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MarshalJSON implements custom JSON marshaling for Product.
// - Price is serialized with exactly 2 decimal places.
// - Timestamps are serialized in RFC 3339 format with "Z" suffix (UTC).
// - Description appears as JSON null when nil (not omitted).
func (p Product) MarshalJSON() ([]byte, error) {
	type Alias struct {
		ID          int             `json:"id"`
		Name        string          `json:"name"`
		Description *string         `json:"description"`
		Price       json.RawMessage `json:"price"`
		Quantity    int             `json:"quantity"`
		CreatedAt   string          `json:"created_at"`
		UpdatedAt   string          `json:"updated_at"`
	}

	priceStr := fmt.Sprintf("%.2f", p.Price)

	a := Alias{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       json.RawMessage(priceStr),
		Quantity:    p.Quantity,
		CreatedAt:   p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return json.Marshal(a)
}

// CreateProductRequest represents the request body for creating a product.
// Price and Quantity are pointers to detect missing fields vs zero values.
type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Quantity    *int     `json:"quantity"`
}

// UpdateProductRequest represents the request body for updating a product.
// Price and Quantity are pointers to detect missing fields vs zero values.
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Quantity    *int     `json:"quantity"`
}
