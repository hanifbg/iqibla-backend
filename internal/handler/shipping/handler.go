package shipping

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ShippingRequest struct {
	Origin struct {
		Kecamatan string `json:"kecamatan"`
		Kota      string `json:"kota"`
	} `json:"origin"`
	Destination struct {
		Kecamatan string `json:"kecamatan"`
		Kota      string `json:"kota"`
	} `json:"destination"`
}

func GetShippingRate(c echo.Context) error {
	var req ShippingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Hardcoded data
	rates := map[string]int{
		"Cempaka Putih-Jakarta Pusat:Menteng-Jakarta Pusat": 15000,
	}

	key := fmt.Sprintf("%s-%s:%s-%s",
		req.Origin.Kecamatan, req.Origin.Kota,
		req.Destination.Kecamatan, req.Destination.Kota,
	)

	rate, exists := rates[key]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Shipping rate not found"})
	}

	return c.JSON(http.StatusOK, map[string]int{"rate": rate})
}
