package postgres

import (
	"context"
	"database/sql"
	"testing"

	"backend-context-engineering-template/internal/domain"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Use environment variables or test database configuration
	dsn := "host=localhost port=5432 user=test_user password=test_password dbname=test_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Cannot connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Cannot ping test database: %v", err)
	}

	// Create test table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			store_id INTEGER NOT NULL,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			amount INTEGER NOT NULL DEFAULT 0,
			price NUMERIC(12,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		TRUNCATE TABLE products RESTART IDENTITY;
	`

	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	return db
}

func TestProductRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := logrus.New()
	repo := NewProductRepository(db, logger)
	ctx := context.Background()

	t.Run("Create and Get Product", func(t *testing.T) {
		product := &domain.Product{
			StoreID:     1,
			Name:        "Integration Test Product",
			Description: sql.NullString{String: "Test Description", Valid: true},
			Amount:      5,
			Price:       19.99,
		}

		// Test Create
		created, err := repo.Create(ctx, product)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, product.StoreID, created.StoreID)
		assert.Equal(t, product.Name, created.Name)
		assert.Equal(t, product.Description, created.Description)
		assert.Equal(t, product.Amount, created.Amount)
		assert.Equal(t, product.Price, created.Price)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.StoreID, retrieved.StoreID)
		assert.Equal(t, created.Name, retrieved.Name)
		assert.Equal(t, created.Description, retrieved.Description)
		assert.Equal(t, created.Amount, retrieved.Amount)
		assert.Equal(t, created.Price, retrieved.Price)
	})

	t.Run("Get Nonexistent Product", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.ErrorIs(t, err, domain.ErrProductNotFound)
	})

	t.Run("Update Product", func(t *testing.T) {
		// Create a product first
		product := &domain.Product{
			StoreID:     1,
			Name:        "Original Product",
			Description: sql.NullString{String: "Original Description", Valid: true},
			Amount:      10,
			Price:       29.99,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)

		// Update the product
		updateData := &domain.Product{
			StoreID:     1,
			Name:        "Updated Product",
			Description: sql.NullString{String: "Updated Description", Valid: true},
			Amount:      15,
			Price:       39.99,
		}

		updated, err := repo.Update(ctx, created.ID, updateData)
		require.NoError(t, err)
		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, updateData.Name, updated.Name)
		assert.Equal(t, updateData.Description, updated.Description)
		assert.Equal(t, updateData.Amount, updated.Amount)
		assert.Equal(t, updateData.Price, updated.Price)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt) || updated.UpdatedAt.Equal(updated.CreatedAt))
	})

	t.Run("Update Nonexistent Product", func(t *testing.T) {
		updateData := &domain.Product{
			StoreID: 1,
			Name:    "Updated Product",
			Amount:  15,
			Price:   39.99,
		}

		_, err := repo.Update(ctx, 99999, updateData)
		assert.ErrorIs(t, err, domain.ErrProductNotFound)
	})

	t.Run("Delete Product", func(t *testing.T) {
		// Create a product first
		product := &domain.Product{
			StoreID: 1,
			Name:    "Product to Delete",
			Amount:  5,
			Price:   19.99,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)

		// Delete the product
		err = repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Verify it's deleted
		_, err = repo.GetByID(ctx, created.ID)
		assert.ErrorIs(t, err, domain.ErrProductNotFound)
	})

	t.Run("Delete Nonexistent Product", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		assert.ErrorIs(t, err, domain.ErrProductNotFound)
	})

	t.Run("Get All Products", func(t *testing.T) {
		// Clean up first
		db.Exec("TRUNCATE TABLE products RESTART IDENTITY")

		// Create multiple products
		products := []*domain.Product{
			{StoreID: 1, Name: "Product 1", Amount: 5, Price: 19.99},
			{StoreID: 1, Name: "Product 2", Amount: 10, Price: 29.99},
			{StoreID: 2, Name: "Product 3", Amount: 15, Price: 39.99},
		}

		for _, p := range products {
			_, err := repo.Create(ctx, p)
			require.NoError(t, err)
		}

		// Test GetAll with no limit
		all, err := repo.GetAll(ctx, 10, 0)
		require.NoError(t, err)
		assert.Len(t, all, 3)

		// Test GetAll with limit
		limited, err := repo.GetAll(ctx, 2, 0)
		require.NoError(t, err)
		assert.Len(t, limited, 2)

		// Test GetAll with offset
		offset, err := repo.GetAll(ctx, 10, 1)
		require.NoError(t, err)
		assert.Len(t, offset, 2)

		// Verify ordering (should be by created_at DESC)
		assert.True(t, all[0].CreatedAt.After(all[1].CreatedAt) || all[0].CreatedAt.Equal(all[1].CreatedAt))
	})

	t.Run("Product with Null Description", func(t *testing.T) {
		product := &domain.Product{
			StoreID:     1,
			Name:        "Product with No Description",
			Description: sql.NullString{Valid: false},
			Amount:      5,
			Price:       19.99,
		}

		created, err := repo.Create(ctx, product)
		require.NoError(t, err)
		assert.False(t, created.Description.Valid)

		retrieved, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.False(t, retrieved.Description.Valid)
	})
}
