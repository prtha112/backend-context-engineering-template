package domain

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInvalidProduct   = errors.New("invalid product data")
	ErrDuplicateProduct = errors.New("product with this name already exists")
)
