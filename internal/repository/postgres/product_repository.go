package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"backend-context-engineering-template/internal/domain"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type ProductRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewProductRepository(db *sql.DB, logger *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (store_id, name, description, amount, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, store_id, name, description, amount, price, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query,
		product.StoreID,
		product.Name,
		nullStringFromString(product.Description.String),
		product.Amount,
		product.Price,
	)

	result := &domain.Product{}
	err := row.Scan(
		&result.ID,
		&result.StoreID,
		&result.Name,
		&result.Description,
		&result.Amount,
		&result.Price,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return nil, domain.ErrDuplicateProduct
			}
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return result, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, store_id, name, description, amount, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	product := &domain.Product{}
	err := row.Scan(
		&product.ID,
		&product.StoreID,
		&product.Name,
		&product.Description,
		&product.Amount,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, store_id, name, description, amount, price, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.StoreID,
			&product.Name,
			&product.Description,
			&product.Amount,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over products: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error) {
	query := `
		UPDATE products
		SET store_id = $1, name = $2, description = $3, amount = $4, price = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, store_id, name, description, amount, price, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query,
		product.StoreID,
		product.Name,
		nullStringFromString(product.Description.String),
		product.Amount,
		product.Price,
		id,
	)

	result := &domain.Product{}
	err := row.Scan(
		&result.ID,
		&result.StoreID,
		&result.Name,
		&result.Description,
		&result.Amount,
		&result.Price,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return nil, domain.ErrDuplicateProduct
			}
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return result, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func nullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
