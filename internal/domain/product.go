package domain

import (
	"database/sql"
	"errors"
	"time"
)

type Product struct {
	ID          int64          `json:"id" db:"id"`
	StoreID     int64          `json:"store_id" db:"store_id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	Amount      int64          `json:"amount" db:"amount"`
	Price       float64        `json:"price" db:"price"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

func (p *Product) Validate() error {
	if p.StoreID <= 0 {
		return errors.New("store_id must be positive")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if len(p.Name) > 100 {
		return errors.New("name must not exceed 100 characters")
	}

	if p.Description.Valid && len(p.Description.String) > 1000 {
		return errors.New("description must not exceed 1000 characters")
	}

	if p.Amount < 0 {
		return errors.New("amount must be non-negative")
	}

	if !p.IsValidPrice() {
		return errors.New("price must be positive")
	}

	return nil
}

func (p *Product) IsValidPrice() bool {
	return p.Price > 0
}
