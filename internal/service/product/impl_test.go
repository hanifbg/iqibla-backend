package product

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/service/product/mocks"
	"github.com/stretchr/testify/assert"
)

func createTestProductService(productRepo *mocks.MockProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	t.Run("Success - Get all products", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		expectedProducts := []entity.Product{
			{
				ID:          "product-1",
				Name:        "Test Product 1",
				Description: "Test Description 1",
				Category:    "Electronics",
				Brand:       "TestBrand",
				Features:    entity.JSONArray{"Feature 1", "Feature 2"},
				InBoxItems:  entity.JSONArray{"Item 1", "Item 2"},
				ImageURLs:   entity.JSONArray{"https://example.com/image1.jpg"},
				IsActive:    true,
			},
			{
				ID:          "product-2",
				Name:        "Test Product 2",
				Description: "Test Description 2",
				Category:    "Clothing",
				Brand:       "TestBrand2",
				Features:    entity.JSONArray{"Feature A", "Feature B"},
				InBoxItems:  entity.JSONArray{"Item A", "Item B"},
				ImageURLs:   entity.JSONArray{"https://example.com/image2.jpg"},
				IsActive:    true,
			},
		}

		mockProductRepo.EXPECT().GetAllProducts().Return(expectedProducts, nil)

		// Act
		result, err := service.GetAllProducts()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedProducts, result)
	})

	t.Run("Success - Empty products list", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		expectedProducts := []entity.Product{}

		mockProductRepo.EXPECT().GetAllProducts().Return(expectedProducts, nil)

		// Act
		result, err := service.GetAllProducts()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("Error - Repository failure", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		mockProductRepo.EXPECT().GetAllProducts().Return(nil, errors.New("database error"))

		// Act
		result, err := service.GetAllProducts()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestProductService_GetProductByID(t *testing.T) {
	t.Run("Success - Get product by ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		productID := "product-123"
		expectedProduct := &entity.Product{
			ID:          productID,
			Name:        "Test Product",
			Description: "Test Description",
			Category:    "Books",
			Brand:       "TestBrand",
			Features:    entity.JSONArray{"Feature 1"},
			InBoxItems:  entity.JSONArray{"Item 1"},
			ImageURLs:   entity.JSONArray{"https://example.com/image.jpg"},
			IsActive:    true,
		}

		mockProductRepo.EXPECT().GetProductByID(productID).Return(expectedProduct, nil)

		// Act
		result, err := service.GetProductByID(productID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct, result)
		assert.Equal(t, productID, result.ID)
	})

	t.Run("Error - Product not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		productID := "non-existent-product"

		mockProductRepo.EXPECT().GetProductByID(productID).Return(nil, errors.New("product not found"))

		// Act
		result, err := service.GetProductByID(productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found")
	})

	t.Run("Error - Repository failure", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		productID := "product-123"

		mockProductRepo.EXPECT().GetProductByID(productID).Return(nil, errors.New("database connection failed"))

		// Act
		result, err := service.GetProductByID(productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("Error - Empty product ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProductRepo := mocks.NewMockProductRepository(ctrl)
		service := createTestProductService(mockProductRepo)

		productID := ""

		mockProductRepo.EXPECT().GetProductByID(productID).Return(nil, errors.New("invalid product ID"))

		// Act
		result, err := service.GetProductByID(productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid product ID")
	})
}