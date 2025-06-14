package product

import "github.com/hanifbg/landing_backend/internal/model/entity"

func (p *ProductService) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product

	products, err := p.productRepo.GetAllProducts()
	if err != nil {
		return nil, err
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
