package cart

import (
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/labstack/echo/v4"
)

// AddItem godoc
// @Summary Add item to cart
// @Description Adds a product variant to the cart with specified quantity
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.AddItemRequest true "Add item request"
// @Success 200 {object} response.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cart/add [post]
func (h *ApiWrapper) AddItem(c echo.Context) error {
	var req request.AddItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if req.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Quantity must be greater than 0"})
	}

	response, err := h.cartService.AddItem(req)
	if err != nil {
		switch err.Error() {
		case "failed to find cart":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart not found"})
		case "failed to get product variant":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product variant not found"})
		case "insufficient stock":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateItemQuantity godoc
// @Summary Update cart item quantity
// @Description Updates the quantity of a product variant in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.UpdateItemRequest true "Update quantity request"
// @Success 200 {object} response.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cart/update-quantity [post]
func (h *ApiWrapper) UpdateItemQuantity(c echo.Context) error {
	var req request.UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if req.Quantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Quantity cannot be negative"})
	}

	response, err := h.cartService.UpdateItemQuantity(req)
	if err != nil {
		switch err.Error() {
		case "cart item not found":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart item not found"})
		case "failed to get product variant":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product variant not found"})
		case "insufficient stock":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Removes a product variant from the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.RemoveItemRequest true "Remove item request"
// @Success 200 {object} response.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cart/remove [post]
func (h *ApiWrapper) RemoveItem(c echo.Context) error {
	var req request.RemoveItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	response, err := h.cartService.RemoveItem(req)
	if err != nil {
		if err.Error() == "cart item not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart item not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, response)
}

// GetCart godoc
// @Summary Get cart details
// @Description Retrieves cart details including items and totals
// @Tags cart
// @Accept json
// @Produce json
// @Param cart_id path string true "Cart ID"
// @Success 200 {object} response.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cart/{cart_id} [get]
func (h *ApiWrapper) GetCart(c echo.Context) error {
	cartID := c.Param("cart_id")
	if cartID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cart ID is required"})
	}

	response, err := h.cartService.GetCart(cartID)
	if err != nil {
		if err.Error() == "failed to get cart" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, response)
}

// ApplyDiscount godoc
// @Summary Apply discount to cart
// @Description Applies a discount code to the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body request.ApplyDiscountRequest true "Apply discount request"
// @Success 200 {object} response.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cart/apply-discount [post]
func (h *ApiWrapper) ApplyDiscount(c echo.Context) error {
	var req request.ApplyDiscountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	response, err := h.cartService.ApplyDiscount(req)
	if err != nil {
		switch err.Error() {
		case "failed to get cart":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cart not found"})
		case "discount not found":
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Discount code not found"})
		case "discount is not active":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Discount is not active"})
		case "discount has not started yet":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Discount has not started yet"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, response)
}
