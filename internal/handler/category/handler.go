package category

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *ApiWrapper) GetCategoryBySlug(c echo.Context) error {
	slug := c.Param("slug")

	category, err := h.CategoryService.GetCategoryBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch category",
		})
	}

	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Category not found",
		})
	}

	return c.JSON(http.StatusOK, category)
}
