package request

type AddItemRequest struct {
	CartID    string `json:"cart_id,omitempty"`
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateItemRequest struct {
	CartID    string `json:"cart_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=0"`
}

type RemoveItemRequest struct {
	CartID    string `json:"cart_id" binding:"required"`
	VariantID string `json:"variant_id" binding:"required"`
}

type ApplyDiscountRequest struct {
	CartID      string `json:"cart_id" binding:"required"`
	DiscountCode string `json:"discount_code" binding:"required"`
}