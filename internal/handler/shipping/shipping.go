package shipping

import (
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/labstack/echo/v4"
)

func (h *ApiWrapper) GetProvinces(c echo.Context) error {
	var req request.GetProvincesRequest
	req.ID = c.QueryParam("id")

	provinces, err := h.shippingService.GetProvinces(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get provinces",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Provinces retrieved successfully",
		"data":    provinces,
	})
}

func (h *ApiWrapper) GetCities(c echo.Context) error {
	var req request.GetCitiesRequest
	req.ProvinceID = c.Param("province_id")
	req.ID = c.QueryParam("id")

	cities, err := h.shippingService.GetCities(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get cities",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Cities retrieved successfully",
		"data":    cities,
	})
}

func (h *ApiWrapper) CalculateShippingCost(c echo.Context) error {
	var req request.CalculateShippingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"message": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"message": err.Error(),
		})
	}

	costs, err := h.shippingService.CalculateShippingCost(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to calculate shipping cost",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Shipping cost calculated successfully",
		"data":    costs,
	})
}
