package handler

import (
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/labstack/echo/v4"
)

type ShippingHandler struct {
	shippingService service.ShippingService
}

func NewShippingHandler(shippingService service.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
	}
}

func (h *ShippingHandler) GetProvinces(c echo.Context) error {
	var req request.GetProvincesRequest
	req.ID = c.QueryParam("id")

	provinces, err := h.shippingService.GetProvinces(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to get provinces",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Provinces retrieved successfully",
		"data": provinces,
	})
}

func (h *ShippingHandler) GetCities(c echo.Context) error {
	var req request.GetCitiesRequest
	req.ProvinceID = c.QueryParam("province_id")
	req.ID = c.QueryParam("id")

	cities, err := h.shippingService.GetCities(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to get cities",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Cities retrieved successfully",
		"data": cities,
	})
}

func (h *ShippingHandler) CalculateShippingCost(c echo.Context) error {
	var req request.CalculateShippingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request",
			"message": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"message": err.Error(),
		})
	}

	costs, err := h.shippingService.CalculateShippingCost(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to calculate shipping cost",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Shipping cost calculated successfully",
		"data": costs,
	})
}