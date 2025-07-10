package category

import "github.com/hanifbg/landing_backend/internal/model/entity"

func (s *CategoryService) GetCategoryBySlug(slug string) (*entity.Category, error) {
	category, err := s.categoryRepo.GetCategoryBySlug(slug)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, nil
	}

	products, err := s.productRepo.GetAllProductsByCategorySlug(slug)
	if err != nil {
		return nil, err
	}

	category.Products = products

	return category, nil
}
