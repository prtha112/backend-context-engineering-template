package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-context-engineering-template/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductUseCase struct {
	mock.Mock
}

func (m *MockProductUseCase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) UpdateProduct(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductUseCase) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api/v1")
	products := api.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.GET("/:id", handler.GetProduct)
		products.GET("", handler.GetProducts)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
	}

	return r
}

func TestProductHandler_CreateProduct(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		requestBody  interface{}
		mockFn       func(*MockProductUseCase)
		expectedCode int
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"store_id":    1,
				"name":        "Test Product",
				"description": "Test Description",
				"amount":      10,
				"price":       29.99,
			},
			mockFn: func(m *MockProductUseCase) {
				m.On("CreateProduct", mock.Anything, mock.Anything).Return(
					&domain.Product{
						ID:          1,
						StoreID:     1,
						Name:        "Test Product",
						Description: sql.NullString{String: "Test Description", Valid: true},
						Amount:      10,
						Price:       29.99,
					}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "validation error - missing required field",
			requestBody: map[string]interface{}{
				"name":        "Test Product",
				"description": "Test Description",
				"amount":      10,
				"price":       29.99,
			},
			mockFn:       func(m *MockProductUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid JSON",
			requestBody:  "invalid json",
			mockFn:       func(m *MockProductUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "domain error",
			requestBody: map[string]interface{}{
				"store_id":    1,
				"name":        "Test Product",
				"description": "Test Description",
				"amount":      10,
				"price":       29.99,
			},
			mockFn: func(m *MockProductUseCase) {
				m.On("CreateProduct", mock.Anything, mock.Anything).Return(
					(*domain.Product)(nil), domain.ErrInvalidProduct)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockProductUseCase{}
			tt.mockFn(mockUseCase)

			handler := NewProductHandler(mockUseCase, logger)
			router := setupTestRouter(handler)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		id           string
		mockFn       func(*MockProductUseCase)
		expectedCode int
	}{
		{
			name: "successful retrieval",
			id:   "1",
			mockFn: func(m *MockProductUseCase) {
				m.On("GetProduct", mock.Anything, int64(1)).Return(
					&domain.Product{
						ID:      1,
						StoreID: 1,
						Name:    "Test Product",
						Amount:  10,
						Price:   29.99,
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid ID",
			id:           "invalid",
			mockFn:       func(m *MockProductUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "product not found",
			id:   "999",
			mockFn: func(m *MockProductUseCase) {
				m.On("GetProduct", mock.Anything, int64(999)).Return(
					(*domain.Product)(nil), domain.ErrProductNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockProductUseCase{}
			tt.mockFn(mockUseCase)

			handler := NewProductHandler(mockUseCase, logger)
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetProducts(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		query        string
		mockFn       func(*MockProductUseCase)
		expectedCode int
	}{
		{
			name:  "successful retrieval",
			query: "",
			mockFn: func(m *MockProductUseCase) {
				m.On("GetProducts", mock.Anything, 10, 0).Return(
					[]*domain.Product{
						{ID: 1, Name: "Product 1", StoreID: 1, Amount: 5, Price: 19.99},
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "with pagination",
			query: "?limit=5&offset=10",
			mockFn: func(m *MockProductUseCase) {
				m.On("GetProducts", mock.Anything, 5, 10).Return(
					[]*domain.Product{}, nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockProductUseCase{}
			tt.mockFn(mockUseCase)

			handler := NewProductHandler(mockUseCase, logger)
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products"+tt.query, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		id           string
		requestBody  interface{}
		mockFn       func(*MockProductUseCase)
		expectedCode int
	}{
		{
			name: "successful update",
			id:   "1",
			requestBody: map[string]interface{}{
				"store_id":    1,
				"name":        "Updated Product",
				"description": "Updated Description",
				"amount":      15,
				"price":       39.99,
			},
			mockFn: func(m *MockProductUseCase) {
				m.On("UpdateProduct", mock.Anything, int64(1), mock.Anything).Return(
					&domain.Product{
						ID:          1,
						StoreID:     1,
						Name:        "Updated Product",
						Description: sql.NullString{String: "Updated Description", Valid: true},
						Amount:      15,
						Price:       39.99,
					}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid ID",
			id:           "invalid",
			requestBody:  map[string]interface{}{},
			mockFn:       func(m *MockProductUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "product not found",
			id:   "999",
			requestBody: map[string]interface{}{
				"store_id":    1,
				"name":        "Updated Product",
				"description": "Updated Description",
				"amount":      15,
				"price":       39.99,
			},
			mockFn: func(m *MockProductUseCase) {
				m.On("UpdateProduct", mock.Anything, int64(999), mock.Anything).Return(
					(*domain.Product)(nil), domain.ErrProductNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockProductUseCase{}
			tt.mockFn(mockUseCase)

			handler := NewProductHandler(mockUseCase, logger)
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	logger := logrus.New()

	tests := []struct {
		name         string
		id           string
		mockFn       func(*MockProductUseCase)
		expectedCode int
	}{
		{
			name: "successful deletion",
			id:   "1",
			mockFn: func(m *MockProductUseCase) {
				m.On("DeleteProduct", mock.Anything, int64(1)).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "invalid ID",
			id:           "invalid",
			mockFn:       func(m *MockProductUseCase) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "product not found",
			id:   "999",
			mockFn: func(m *MockProductUseCase) {
				m.On("DeleteProduct", mock.Anything, int64(999)).Return(domain.ErrProductNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockProductUseCase{}
			tt.mockFn(mockUseCase)

			handler := NewProductHandler(mockUseCase, logger)
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}
