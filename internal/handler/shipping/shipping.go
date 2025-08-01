package shipping

import (
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/labstack/echo/v4"
)

// GetProvinces godoc
// @Summary Get provinces
// @Description Get a list of all provinces or specific province by ID
// @Tags shipping
// @Produce json
// @Param id query string false "Province ID (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/provinces [get]
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

// GetCities godoc
// @Summary Get cities
// @Description Get cities by province ID and/or city ID
// @Tags shipping
// @Produce json
// @Param province_id path string true "Province ID"
// @Param id query string false "City ID (optional)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/cities/{province_id} [get]
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

// GetDistricts godoc
// @Summary Get districts
// @Description Get districts by city ID
// @Tags shipping
// @Produce json
// @Param city_id path string true "City ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/districts/{city_id} [get]
func (h *ApiWrapper) GetDistricts(c echo.Context) error {
	var req request.GetDistrictsRequest
	req.CityID = c.Param("city_id")

	districts, err := h.shippingService.GetDistricts(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get districts",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Districts retrieved successfully",
		"data":    districts,
	})
}

// CalculateShippingCost godoc
// @Summary Calculate shipping cost
// @Description Calculate shipping cost based on origin, destination, weight, and courier
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body request.CalculateShippingRequest true "Calculate shipping cost request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/shipping/cost [post]
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
