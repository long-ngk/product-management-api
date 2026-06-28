package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/product-management-server/internal/model"
)

// productRepository implements ProductRepository using PostgreSQL.
type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository backed by PostgreSQL.
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database and returns the created product
// with its generated id, created_at, and updated_at fields.
func (r *productRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := `
		INSERT INTO products (name, description, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, quantity, created_at, updated_at`

	created := &model.Product{}
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.Price,
		&created.Quantity,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// FindAll retrieves all products from the database ordered by id ascending.
func (r *productRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// FindByKeyword retrieves products whose name or description contains the keyword
// (case-insensitive), ordered by id ascending.
func (r *productRepository) FindByKeyword(ctx context.Context, keyword string) ([]model.Product, error) {
	query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY id ASC`

	pattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// FindByID retrieves a single product by its id.
// Returns an error if no product is found with the given id.
func (r *productRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		WHERE id = $1`

	var p model.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}

	return &p, nil
}

// Update modifies an existing product in the database and returns the updated product
// with its new updated_at timestamp.
func (r *productRepository) Update(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, quantity = $4
		WHERE id = $5
		RETURNING id, name, description, price, quantity, created_at, updated_at`

	updated := &model.Product{}
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.Price,
		&updated.Quantity,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", product.ID)
		}
		return nil, err
	}

	return updated, nil
}

// Delete removes a product from the database by its id.
// Returns an error if no product exists with the given id.
func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}
