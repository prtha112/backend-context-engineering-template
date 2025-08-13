package dto

import (
	"database/sql"
	"time"

	"backend-context-engineering-template/internal/domain"
)

type CreateProductRequest struct {
	StoreID     int64   `json:"store_id" binding:"required,min=1"`
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=1000"`
	Amount      int64   `json:"amount" binding:"required,min=0"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type UpdateProductRequest struct {
	StoreID     int64   `json:"store_id" binding:"required,min=1"`
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=1000"`
	Amount      int64   `json:"amount" binding:"required,min=0"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type ProductResponse struct {
	ID          int64   `json:"id"`
	StoreID     int64   `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      int64   `json:"amount"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (r *CreateProductRequest) ToDomain() *domain.Product {
	description := sql.NullString{}
	if r.Description != "" {
		description = sql.NullString{String: r.Description, Valid: true}
	}

	return &domain.Product{
		StoreID:     r.StoreID,
		Name:        r.Name,
		Description: description,
		Amount:      r.Amount,
		Price:       r.Price,
	}
}

func (r *UpdateProductRequest) ToDomain() *domain.Product {
	description := sql.NullString{}
	if r.Description != "" {
		description = sql.NullString{String: r.Description, Valid: true}
	}

	return &domain.Product{
		StoreID:     r.StoreID,
		Name:        r.Name,
		Description: description,
		Amount:      r.Amount,
		Price:       r.Price,
	}
}

func ToProductResponse(product *domain.Product) ProductResponse {
	description := ""
	if product.Description.Valid {
		description = product.Description.String
	}

	return ProductResponse{
		ID:          product.ID,
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: description,
		Amount:      product.Amount,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

func ToProductListResponse(products []*domain.Product, limit, offset int) ProductListResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ToProductResponse(product)
	}

	return ProductListResponse{
		Products: productResponses,
		Total:    len(products),
		Limit:    limit,
		Offset:   offset,
	}
}
