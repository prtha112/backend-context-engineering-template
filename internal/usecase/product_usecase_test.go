package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"backend-context-engineering-template/internal/domain"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, id int64, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, id, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductUseCase_CreateProduct(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()

	tests := []struct {
		name    string
		product *domain.Product
		mockFn  func(*MockProductRepository)
		want    *domain.Product
		wantErr bool
		errType error
	}{
		{
			name: "successful creation",
			product: &domain.Product{
				StoreID:     1,
				Name:        "Test Product",
				Description: sql.NullString{String: "Test Description", Valid: true},
				Amount:      10,
				Price:       29.99,
			},
			mockFn: func(m *MockProductRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(
					&domain.Product{
						ID:          1,
						StoreID:     1,
						Name:        "Test Product",
						Description: sql.NullString{String: "Test Description", Valid: true},
						Amount:      10,
						Price:       29.99,
					}, nil)
			},
			want: &domain.Product{
				ID:          1,
				StoreID:     1,
				Name:        "Test Product",
				Description: sql.NullString{String: "Test Description", Valid: true},
				Amount:      10,
				Price:       29.99,
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			product: &domain.Product{
				StoreID: 1,
				Name:    "",
				Amount:  10,
				Price:   29.99,
			},
			mockFn:  func(m *MockProductRepository) {},
			want:    nil,
			wantErr: true,
			errType: domain.ErrInvalidProduct,
		},
		{
			name: "validation error - negative price",
			product: &domain.Product{
				StoreID: 1,
				Name:    "Test Product",
				Amount:  10,
				Price:   -5.0,
			},
			mockFn:  func(m *MockProductRepository) {},
			want:    nil,
			wantErr: true,
			errType: domain.ErrInvalidProduct,
		},
		{
			name: "repository error",
			product: &domain.Product{
				StoreID: 1,
				Name:    "Test Product",
				Amount:  10,
				Price:   29.99,
			},
			mockFn: func(m *MockProductRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(
					(*domain.Product)(nil), errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockProductRepository{}
			tt.mockFn(repo)

			uc := NewProductUseCase(repo, logger)
			got, err := uc.CreateProduct(ctx, tt.product)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestProductUseCase_GetProduct(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		mockFn  func(*MockProductRepository)
		want    *domain.Product
		wantErr bool
		errType error
	}{
		{
			name: "successful retrieval",
			id:   1,
			mockFn: func(m *MockProductRepository) {
				m.On("GetByID", mock.Anything, int64(1)).Return(
					&domain.Product{
						ID:      1,
						StoreID: 1,
						Name:    "Test Product",
						Amount:  10,
						Price:   29.99,
					}, nil)
			},
			want: &domain.Product{
				ID:      1,
				StoreID: 1,
				Name:    "Test Product",
				Amount:  10,
				Price:   29.99,
			},
			wantErr: false,
		},
		{
			name:    "invalid ID",
			id:      0,
			mockFn:  func(m *MockProductRepository) {},
			want:    nil,
			wantErr: true,
			errType: domain.ErrInvalidProduct,
		},
		{
			name: "product not found",
			id:   999,
			mockFn: func(m *MockProductRepository) {
				m.On("GetByID", mock.Anything, int64(999)).Return(
					(*domain.Product)(nil), domain.ErrProductNotFound)
			},
			want:    nil,
			wantErr: true,
			errType: domain.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockProductRepository{}
			tt.mockFn(repo)

			uc := NewProductUseCase(repo, logger)
			got, err := uc.GetProduct(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestProductUseCase_GetProducts(t *testing.T) {
	logger := logrus.New()
	ctx := context.Background()

	tests := []struct {
		name    string
		limit   int
		offset  int
		mockFn  func(*MockProductRepository)
		want    []*domain.Product
		wantErr bool
	}{
		{
			name:   "successful retrieval",
			limit:  10,
			offset: 0,
			mockFn: func(m *MockProductRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return(
					[]*domain.Product{
						{ID: 1, Name: "Product 1", StoreID: 1, Amount: 5, Price: 19.99},
						{ID: 2, Name: "Product 2", StoreID: 1, Amount: 10, Price: 29.99},
					}, nil)
			},
			want: []*domain.Product{
				{ID: 1, Name: "Product 1", StoreID: 1, Amount: 5, Price: 19.99},
				{ID: 2, Name: "Product 2", StoreID: 1, Amount: 10, Price: 29.99},
			},
			wantErr: false,
		},
		{
			name:   "invalid limit - should default to 10",
			limit:  0,
			offset: 0,
			mockFn: func(m *MockProductRepository) {
				m.On("GetAll", mock.Anything, 10, 0).Return([]*domain.Product{}, nil)
			},
			want:    []*domain.Product{},
			wantErr: false,
		},
		{
			name:   "limit too large - should cap at 100",
			limit:  150,
			offset: 0,
			mockFn: func(m *MockProductRepository) {
				m.On("GetAll", mock.Anything, 100, 0).Return([]*domain.Product{}, nil)
			},
			want:    []*domain.Product{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockProductRepository{}
			tt.mockFn(repo)

			uc := NewProductUseCase(repo, logger)
			got, err := uc.GetProducts(ctx, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			repo.AssertExpectations(t)
		})
	}
}
