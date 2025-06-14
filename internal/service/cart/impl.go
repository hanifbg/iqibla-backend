package cart

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

func (s *CartService) calculateCartTotals(cart *entity.Cart, discount *entity.Discount) (*response.CartResponse, error) {
	if cart == nil {
		return nil, fmt.Errorf("cart is nil")
	}

	var totalItems int
	subtotalAmount := 0.0
	itemResponses := make([]response.CartItemResponse, 0)

	for _, item := range cart.CartItems {
		if item.ProductVariant == nil {
			return nil, fmt.Errorf("product variant not loaded for cart item")
		}

		totalItems += item.Quantity
		subtotalAmount += float64(item.Quantity) * item.ProductVariant.Price

		itemResponses = append(itemResponses, response.CartItemResponse{
			VariantID:         item.ProductVariant.ID,
			VariantName:       item.ProductVariant.Name,
			VariantPrice:      item.ProductVariant.Price,
			Quantity:          item.Quantity,
			ImageURL:          item.ProductVariant.ImageURL,
			ProductAttributes: item.ProductVariant.AttributeValues,
		})
	}

	response := &response.CartResponse{
		CartID:         cart.ID,
		TotalItems:     totalItems,
		SubtotalAmount: subtotalAmount,
		Items:          itemResponses,
	}

	if discount != nil {
		var discountAmount float64
		if discount.Type == "percentage" {
			discountAmount = subtotalAmount * (discount.Value / 100)
		} else { // fixed_amount
			discountAmount = discount.Value
		}
		response.DiscountAmount = &discountAmount
		response.DiscountCodeApplied = &discount.Code
	}

	return response, nil
}

func (s *CartService) AddItem(req request.AddItemRequest) (*response.CartResponse, error) {
	var cart *entity.Cart
	var err error

	if req.CartID == "" {
		cart = &entity.Cart{ID: uuid.New().String()}
		if err := s.cartRepo.CreateCart(cart); err != nil {
			return nil, fmt.Errorf("failed to create cart: %v", err)
		}
		req.CartID = cart.ID
	} else {
		cart, err = s.cartRepo.FindCartByID(req.CartID)
		if err != nil {
			return nil, fmt.Errorf("failed to find cart: %v", err)
		}
	}

	variant, err := s.cartRepo.GetProductVariantByID(req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variant: %v", err)
	}

	if variant.StockQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", variant.StockQuantity, req.Quantity)
	}

	existingItem, err := s.cartRepo.FindCartItem(req.CartID, req.VariantID)
	if err == nil { // Item exists
		newQuantity := existingItem.Quantity + req.Quantity
		if variant.StockQuantity < newQuantity {
			return nil, fmt.Errorf("insufficient stock for total quantity: available %d, total requested %d", variant.StockQuantity, newQuantity)
		}
		existingItem.Quantity = newQuantity
		if err := s.cartRepo.UpdateCartItem(existingItem); err != nil {
			return nil, fmt.Errorf("failed to update cart item: %v", err)
		}
	} else { // New item
		newItem := &entity.CartItem{
			CartID:           req.CartID,
			ProductVariantID: req.VariantID,
			Quantity:         req.Quantity,
		}
		if err := s.cartRepo.CreateCartItem(newItem); err != nil {
			return nil, fmt.Errorf("failed to create cart item: %v", err)
		}
	}

	// Reload cart with items
	updatedCart, err := s.cartRepo.GetCartWithItems(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated cart: %v", err)
	}

	return s.calculateCartTotals(updatedCart, nil)
}

func (s *CartService) UpdateItemQuantity(req request.UpdateItemRequest) (*response.CartResponse, error) {
	cartItem, err := s.cartRepo.FindCartItem(req.CartID, req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %v", err)
	}

	variant, err := s.cartRepo.GetProductVariantByID(req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variant: %v", err)
	}

	if req.Quantity == 0 {
		if err := s.cartRepo.DeleteCartItem(cartItem.CartID, cartItem.ProductVariantID); err != nil {
			return nil, fmt.Errorf("failed to delete cart item: %v", err)
		}
	} else {
		if variant.StockQuantity < req.Quantity {
			return nil, fmt.Errorf("insufficient stock: available %d, requested %d", variant.StockQuantity, req.Quantity)
		}
		cartItem.Quantity = req.Quantity
		if err := s.cartRepo.UpdateCartItem(cartItem); err != nil {
			return nil, fmt.Errorf("failed to update cart item: %v", err)
		}
	}

	updatedCart, err := s.cartRepo.GetCartWithItems(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated cart: %v", err)
	}

	return s.calculateCartTotals(updatedCart, nil)
}

func (s *CartService) RemoveItem(req request.RemoveItemRequest) (*response.CartResponse, error) {
	cartItem, err := s.cartRepo.FindCartItem(req.CartID, req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %v", err)
	}

	if err := s.cartRepo.DeleteCartItem(cartItem.CartID, cartItem.ProductVariantID); err != nil {
		return nil, fmt.Errorf("failed to delete cart item: %v", err)
	}

	updatedCart, err := s.cartRepo.GetCartWithItems(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated cart: %v", err)
	}

	return s.calculateCartTotals(updatedCart, nil)
}

func (s *CartService) GetCart(cartID string) (*response.CartResponse, error) {
	cart, err := s.cartRepo.GetCartWithItems(cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	return s.calculateCartTotals(cart, nil)
}

func (s *CartService) ApplyDiscount(req request.ApplyDiscountRequest) (*response.CartResponse, error) {
	cart, err := s.cartRepo.GetCartWithItems(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	discount, err := s.cartRepo.GetDiscountByCode(req.DiscountCode)
	if err != nil {
		return nil, fmt.Errorf("discount not found: %v", err)
	}

	// Validate discount
	if !discount.IsActive {
		return nil, fmt.Errorf("discount is not active")
	}

	now := time.Now()
	if now.Before(discount.StartsAt) {
		return nil, fmt.Errorf("discount has not started yet")
	}

	if now.After(discount.ExpiresAt) {
		return nil, fmt.Errorf("discount has expired")
	}

	if discount.UsageLimit != 0 && discount.UsesCount >= discount.UsageLimit {
		return nil, fmt.Errorf("discount usage limit reached")
	}

	// Calculate current subtotal
	subtotal := 0.0
	for _, item := range cart.CartItems {
		subtotal += float64(item.Quantity) * item.ProductVariant.Price
	}

	if discount.MinimumOrderAmount != 0 && subtotal < discount.MinimumOrderAmount {
		return nil, fmt.Errorf("cart subtotal does not meet minimum order amount for discount")
	}

	return s.calculateCartTotals(cart, discount)
}
