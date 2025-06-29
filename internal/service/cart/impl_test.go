package cart

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/service/cart/mocks"
	"github.com/stretchr/testify/assert"
)

func createTestCartService(cartRepo *mocks.MockCartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

func createTestProductVariant() *entity.ProductVariant {
	return &entity.ProductVariant{
		ID:            "variant-123",
		ProductID:     "product-123",
		SKU:           "TEST-SKU-001",
		Name:          "Test Variant",
		Price:         100.00,
		StockQuantity: 10,
		ImageURL:      "https://example.com/image.jpg",
		Weight:        0.5,
		AttributeValues: entity.JSONMap{
			"color": "red",
			"size":  "M",
		},
		IsActive: true,
	}
}

func createTestCart() *entity.Cart {
	return &entity.Cart{
		ID:        "cart-123",
		IsActive:  true,
		CartItems: []entity.CartItem{},
	}
}

func createTestCartWithItems() *entity.Cart {
	variant := createTestProductVariant()
	return &entity.Cart{
		ID:       "cart-123",
		IsActive: true,
		CartItems: []entity.CartItem{
			{
				ID:               "item-1",
				CartID:           "cart-123",
				ProductVariantID: "variant-123",
				ProductVariant:   variant,
				Quantity:         2,
			},
		},
	}
}

func createTestDiscount() *entity.Discount {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	return &entity.Discount{
		ID:                 "discount-123",
		Code:               "TEST10",
		Type:               "percentage",
		Value:              10.0,
		MinimumOrderAmount: 50.0,
		StartsAt:           now.Add(-1 * time.Hour),
		ExpiresAt:          &future,
		UsageLimit:         100,
		UsesCount:          5,
		IsActive:           true,
	}
}

func TestCartService_AddItem(t *testing.T) {
	t.Run("Success - Add item to new cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "", // Empty cart ID to create new cart
			VariantID: "variant-123",
			Quantity:  2,
		}

		variant := createTestProductVariant()
		cart := createTestCartWithItems()

		mockCartRepo.EXPECT().CreateCart(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)
		mockCartRepo.EXPECT().FindCartItem(gomock.Any(), "variant-123").Return(nil, errors.New("not found"))
		mockCartRepo.EXPECT().CreateCartItem(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems(gomock.Any()).Return(cart, nil)

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "cart-123", result.CartID)
		assert.Equal(t, 2, result.TotalItems)
		assert.Equal(t, 200.0, result.SubtotalAmount)
	})

	t.Run("Success - Add item to existing cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  1,
		}

		cart := createTestCart()
		variant := createTestProductVariant()
		updatedCart := createTestCartWithItems()

		mockCartRepo.EXPECT().FindCartByID("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)
		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(nil, errors.New("not found"))
		mockCartRepo.EXPECT().CreateCartItem(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(updatedCart, nil)

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "cart-123", result.CartID)
	})

	t.Run("Success - Update existing item quantity", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  1,
		}

		cart := createTestCart()
		variant := createTestProductVariant()
		existingItem := &entity.CartItem{
			CartID:           "cart-123",
			ProductVariantID: "variant-123",
			Quantity:         2,
		}
		updatedCart := createTestCartWithItems()

		mockCartRepo.EXPECT().FindCartByID("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)
		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(existingItem, nil)
		mockCartRepo.EXPECT().UpdateCartItem(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(updatedCart, nil)

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Insufficient stock", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  15, // More than available stock
		}

		cart := createTestCart()
		variant := createTestProductVariant() // Stock is 10

		mockCartRepo.EXPECT().FindCartByID("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient stock")
	})

	t.Run("Error - Cart not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "non-existent-cart",
			VariantID: "variant-123",
			Quantity:  1,
		}

		mockCartRepo.EXPECT().FindCartByID("non-existent-cart").Return(nil, errors.New("cart not found"))

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find cart")
	})

	t.Run("Error - Product variant not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.AddItemRequest{
			CartID:    "cart-123",
			VariantID: "non-existent-variant",
			Quantity:  1,
		}

		cart := createTestCart()

		mockCartRepo.EXPECT().FindCartByID("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("non-existent-variant").Return(nil, errors.New("variant not found"))

		// Act
		result, err := service.AddItem(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get product variant")
	})
}

func TestCartService_UpdateItemQuantity(t *testing.T) {
	t.Run("Success - Update item quantity", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.UpdateItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  3,
		}

		cartItem := &entity.CartItem{
			CartID:           "cart-123",
			ProductVariantID: "variant-123",
			Quantity:         2,
		}
		variant := createTestProductVariant()
		updatedCart := createTestCartWithItems()

		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(cartItem, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)
		mockCartRepo.EXPECT().UpdateCartItem(gomock.Any()).Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(updatedCart, nil)

		// Act
		result, err := service.UpdateItemQuantity(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Success - Remove item when quantity is 0", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.UpdateItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  0,
		}

		cartItem := &entity.CartItem{
			CartID:           "cart-123",
			ProductVariantID: "variant-123",
			Quantity:         2,
		}
		variant := createTestProductVariant()
		updatedCart := createTestCart()

		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(cartItem, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)
		mockCartRepo.EXPECT().DeleteCartItem("cart-123", "variant-123").Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(updatedCart, nil)

		// Act
		result, err := service.UpdateItemQuantity(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Cart item not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.UpdateItemRequest{
			CartID:    "cart-123",
			VariantID: "non-existent-variant",
			Quantity:  1,
		}

		mockCartRepo.EXPECT().FindCartItem("cart-123", "non-existent-variant").Return(nil, errors.New("item not found"))

		// Act
		result, err := service.UpdateItemQuantity(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart item not found")
	})

	t.Run("Error - Insufficient stock for update", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.UpdateItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
			Quantity:  15, // More than available stock
		}

		cartItem := &entity.CartItem{
			CartID:           "cart-123",
			ProductVariantID: "variant-123",
			Quantity:         2,
		}
		variant := createTestProductVariant() // Stock is 10

		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(cartItem, nil)
		mockCartRepo.EXPECT().GetProductVariantByID("variant-123").Return(variant, nil)

		// Act
		result, err := service.UpdateItemQuantity(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "insufficient stock")
	})
}

func TestCartService_RemoveItem(t *testing.T) {
	t.Run("Success - Remove item from cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.RemoveItemRequest{
			CartID:    "cart-123",
			VariantID: "variant-123",
		}

		cartItem := &entity.CartItem{
			CartID:           "cart-123",
			ProductVariantID: "variant-123",
			Quantity:         2,
		}
		updatedCart := createTestCart()

		mockCartRepo.EXPECT().FindCartItem("cart-123", "variant-123").Return(cartItem, nil)
		mockCartRepo.EXPECT().DeleteCartItem("cart-123", "variant-123").Return(nil)
		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(updatedCart, nil)

		// Act
		result, err := service.RemoveItem(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Error - Cart item not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.RemoveItemRequest{
			CartID:    "cart-123",
			VariantID: "non-existent-variant",
		}

		mockCartRepo.EXPECT().FindCartItem("cart-123", "non-existent-variant").Return(nil, errors.New("item not found"))

		// Act
		result, err := service.RemoveItem(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart item not found")
	})
}

func TestCartService_GetCart(t *testing.T) {
	t.Run("Success - Get cart with items", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		cart := createTestCartWithItems()

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)

		// Act
		result, err := service.GetCart("cart-123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "cart-123", result.CartID)
		assert.Equal(t, 2, result.TotalItems)
		assert.Equal(t, 200.0, result.SubtotalAmount)
	})

	t.Run("Error - Cart not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		mockCartRepo.EXPECT().GetCartWithItems("non-existent-cart").Return(nil, errors.New("cart not found"))

		// Act
		result, err := service.GetCart("non-existent-cart")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get cart")
	})
}

func TestCartService_ApplyDiscount(t *testing.T) {
	t.Run("Success - Apply percentage discount", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "TEST10",
		}

		cart := createTestCartWithItems()
		discount := createTestDiscount()

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("TEST10").Return(discount, nil)

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.DiscountAmount)
		assert.Equal(t, 20.0, *result.DiscountAmount) // 10% of 200
		assert.NotNil(t, result.DiscountCodeApplied)
		assert.Equal(t, "TEST10", *result.DiscountCodeApplied)
	})

	t.Run("Success - Apply fixed amount discount", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "FIXED20",
		}

		cart := createTestCartWithItems()
		discount := createTestDiscount()
		discount.Type = "fixed_amount"
		discount.Value = 20.0

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("FIXED20").Return(discount, nil)

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.DiscountAmount)
		assert.Equal(t, 20.0, *result.DiscountAmount)
	})

	t.Run("Error - Discount not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "INVALID",
		}

		cart := createTestCartWithItems()

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("INVALID").Return(nil, errors.New("discount not found"))

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "discount not found")
	})

	t.Run("Error - Discount not active", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "INACTIVE",
		}

		cart := createTestCartWithItems()
		discount := createTestDiscount()
		discount.IsActive = false

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("INACTIVE").Return(discount, nil)

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "discount is not active")
	})

	t.Run("Error - Discount has expired", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "EXPIRED",
		}

		cart := createTestCartWithItems()
		discount := createTestDiscount()
		past := time.Now().Add(-1 * time.Hour)
		discount.ExpiresAt = &past

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("EXPIRED").Return(discount, nil)

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "discount has expired")
	})

	t.Run("Error - Minimum order amount not met", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		req := request.ApplyDiscountRequest{
			CartID:       "cart-123",
			DiscountCode: "HIGHMIN",
		}

		cart := createTestCartWithItems() // Subtotal is 200
		discount := createTestDiscount()
		discount.MinimumOrderAmount = 300.0 // Higher than cart subtotal

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockCartRepo.EXPECT().GetDiscountByCode("HIGHMIN").Return(discount, nil)

		// Act
		result, err := service.ApplyDiscount(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart subtotal does not meet minimum order amount")
	})
}

func TestCartService_calculateCartTotals(t *testing.T) {
	t.Run("Error - Cart is nil", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		// Act
		result, err := service.calculateCartTotals(nil, nil)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart is nil")
	})

	t.Run("Error - Product variant not loaded", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		service := createTestCartService(mockCartRepo)

		cart := &entity.Cart{
			ID:       "cart-123",
			IsActive: true,
			CartItems: []entity.CartItem{
				{
					ID:               "item-1",
					CartID:           "cart-123",
					ProductVariantID: "variant-123",
					ProductVariant:   nil, // Not loaded
					Quantity:         2,
				},
			},
		}

		// Act
		result, err := service.calculateCartTotals(cart, nil)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product variant not loaded")
	})
}