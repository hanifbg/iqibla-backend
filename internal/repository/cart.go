package repository

import "github.com/hanifbg/landing_backend/internal/model/entity"

type CartRepository interface {
	// Cart operations
	FindCartByID(cartID string) (*entity.Cart, error)
	CreateCart(cart *entity.Cart) error

	// Cart item operations
	CreateCartItem(item *entity.CartItem) error
	FindCartItem(cartID, variantID string) (*entity.CartItem, error)
	UpdateCartItem(item *entity.CartItem) error
	DeleteCartItem(cartID, variantID string) error
	GetCartItemsByCartID(cartID string) ([]entity.CartItem, error)
	GetCartWithItems(cartID string) (*entity.Cart, error)

	// Related operations
	GetDiscountByCode(code string) (*entity.Discount, error)
	GetProductVariantByID(variantID string) (*entity.ProductVariant, error)
}
