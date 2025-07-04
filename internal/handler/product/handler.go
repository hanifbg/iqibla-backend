package product

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAllProducts godoc
// @Summary Get all active products
// @Description Retrieves all active products with their variants, optionally filtered by category
// @Tags products
// @Accept json
// @Produce json
// @Param category query string false "Filter products by category"
// @Success 200 {array} entity.Product
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [get]
func (h *ApiWrapper) GetAllProducts(c echo.Context) error {
	// Get the optional category query parameter
	category := c.QueryParam("category")

	// Get all products from the service
	products, err := h.ProductService.GetAllProducts(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	// If category is specified, filter the products
	// if category != "" {
	// 	filteredProducts := make([]entity.Product, 0)
	// 	for _, product := range products {
	// 		if product.Category == category {
	// 			filteredProducts = append(filteredProducts, product)
	// 		}
	// 	}
	// 	products = filteredProducts
	// }

	return c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieves a single product by its ID with variants
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [get]
func (h *ApiWrapper) GetProductByID(c echo.Context) error {
	id := c.Param("id")

	product, err := h.ProductService.GetProductByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product",
		})
	}

	if product == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}

	return c.JSON(http.StatusOK, product)
}
