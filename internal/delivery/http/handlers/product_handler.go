package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend-context-engineering-template/internal/delivery/http/dto"
	"backend-context-engineering-template/internal/domain"
	"backend-context-engineering-template/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCaseInterface
	logger         *logrus.Logger
}

func NewProductHandler(productUseCase usecase.ProductUseCaseInterface, logger *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind create product request")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product := req.ToDomain()
	createdProduct, err := h.productUseCase.CreateProduct(ctx, product)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dto.ToProductResponse(createdProduct)
	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Product ID must be a valid number",
		})
		return
	}

	product, err := h.productUseCase.GetProduct(ctx, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dto.ToProductResponse(product)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	products, err := h.productUseCase.GetProducts(ctx, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dto.ToProductListResponse(products, limit, offset)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Product ID must be a valid number",
		})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind update product request")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product := req.ToDomain()
	updatedProduct, err := h.productUseCase.UpdateProduct(ctx, id, product)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dto.ToProductResponse(updatedProduct)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Product ID must be a valid number",
		})
		return
	}

	if err := h.productUseCase.DeleteProduct(ctx, id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProductHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "product_not_found",
			Message: "Product not found",
		})
	case errors.Is(err, domain.ErrInvalidProduct):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_product",
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrDuplicateProduct):
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "duplicate_product",
			Message: "Product with this name already exists",
		})
	default:
		h.logger.WithError(err).Error("Internal server error")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_server_error",
			Message: "An internal error occurred",
		})
	}
}
