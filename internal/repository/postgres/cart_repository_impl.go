package postgres

import (
	"github.com/hanifbg/landing_backend/internal/model/entity"
)

func (r *RepoDatabase) FindCartByID(cartID string) (*entity.Cart, error) {
	var cart entity.Cart
	result := r.DB.Preload("CartItems.ProductVariant").First(&cart, "id = ?", cartID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cart, nil
}

func (r *RepoDatabase) CreateCart(cart *entity.Cart) error {
	return r.DB.Create(cart).Error
}

func (r *RepoDatabase) CreateCartItem(item *entity.CartItem) error {
	return r.DB.Create(item).Error
}

func (r *RepoDatabase) FindCartItem(cartID, variantID string) (*entity.CartItem, error) {
	var item entity.CartItem
	result := r.DB.Where("cart_id = ? AND product_variant_id = ?", cartID, variantID).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *RepoDatabase) UpdateCartItem(item *entity.CartItem) error {
	return r.DB.Save(item).Error
}

func (r *RepoDatabase) DeleteCartItem(cartID, variantID string) error {
	return r.DB.Where("cart_id = ? AND product_variant_id = ?", cartID, variantID).Delete(&entity.CartItem{}).Error
}

func (r *RepoDatabase) GetCartItemsByCartID(cartID string) ([]entity.CartItem, error) {
	var items []entity.CartItem
	result := r.DB.Preload("ProductVariant").Where("cart_id = ?", cartID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *RepoDatabase) GetDiscountByCode(code string) (*entity.Discount, error) {
	var discount entity.Discount
	result := r.DB.Where("code = ? AND is_active = true", code).First(&discount)
	if result.Error != nil {
		return nil, result.Error
	}
	return &discount, nil
}

func (r *RepoDatabase) GetProductVariantByID(variantID string) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	result := r.DB.Where("id = ? AND is_active = true", variantID).First(&variant)
	if result.Error != nil {
		return nil, result.Error
	}
	return &variant, nil
}

func (r *RepoDatabase) GetCartWithItems(cartID string) (*entity.Cart, error) {
	var cart entity.Cart

	result := r.DB.
		Preload("CartItems").
		Preload("CartItems.ProductVariant").
		Where("id = ?", cartID).
		First(&cart)

	if result.Error != nil {
		return nil, result.Error
	}

	return &cart, nil
}
