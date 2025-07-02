package product

import "github.com/hanifbg/landing_backend/internal/model/entity"

func (p *ProductService) GetAllProducts(category string) ([]entity.Product, error) {
	var products []entity.Product
	var err error

	if category != "" {
		products, err = p.productRepo.GetAllProductsByCategory(category)
		if err != nil {
			return nil, err
		}
	} else {
		products, err = p.productRepo.GetAllProducts()
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func (p *ProductService) GetProductByID(id string) (*entity.Product, error) {
	product, err := p.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
