package usecase

import (
	"context"
	"fmt"

	"backend-context-engineering-template/internal/domain"
	"github.com/sirupsen/logrus"
)

type ProductUseCase struct {
	productRepo ProductRepository
	logger      *logrus.Logger
}

func NewProductUseCase(productRepo ProductRepository, logger *logrus.Logger) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		logger:      logger,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	uc.logger.WithFields(logrus.Fields{
		"action":   "create_product",
		"store_id": product.StoreID,
		"name":     product.Name,
	}).Info("Creating new product")

	if err := product.Validate(); err != nil {
		uc.logger.WithError(err).Error("Product validation failed")
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidProduct, err.Error())
	}

	createdProduct, err := uc.productRepo.Create(ctx, product)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create product in repository")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"action":     "create_product",
		"product_id": createdProduct.ID,
	}).Info("Product created successfully")

	return createdProduct, nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	uc.logger.WithFields(logrus.Fields{
		"action":     "get_product",
		"product_id": id,
	}).Info("Retrieving product")

	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid product ID", domain.ErrInvalidProduct)
	}

	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get product from repository")
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	uc.logger.WithFields(logrus.Fields{
		"action": "get_products",
		"limit":  limit,
		"offset": offset,
	}).Info("Retrieving products")

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	products, err := uc.productRepo.GetAll(ctx, limit, offset)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get products from repository")
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error) {
	uc.logger.WithFields(logrus.Fields{
		"action":     "update_product",
		"product_id": id,
	}).Info("Updating product")

	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid product ID", domain.ErrInvalidProduct)
	}

	if err := product.Validate(); err != nil {
		uc.logger.WithError(err).Error("Product validation failed")
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidProduct, err.Error())
	}

	updatedProduct, err := uc.productRepo.Update(ctx, id, product)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to update product in repository")
		return nil, err
	}

	uc.logger.WithFields(logrus.Fields{
		"action":     "update_product",
		"product_id": updatedProduct.ID,
	}).Info("Product updated successfully")

	return updatedProduct, nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id int64) error {
	uc.logger.WithFields(logrus.Fields{
		"action":     "delete_product",
		"product_id": id,
	}).Info("Deleting product")

	if id <= 0 {
		return fmt.Errorf("%w: invalid product ID", domain.ErrInvalidProduct)
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		uc.logger.WithError(err).Error("Failed to delete product from repository")
		return err
	}

	uc.logger.WithFields(logrus.Fields{
		"action":     "delete_product",
		"product_id": id,
	}).Info("Product deleted successfully")

	return nil
}
