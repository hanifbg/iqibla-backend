package response

type CartItemResponse struct {
	ID                string                 `json:"id"`
	VariantID         string                 `json:"variant_id"`
	VariantName       string                 `json:"variant_name"`
	VariantPrice      float64                `json:"variant_price"`
	Quantity          int                    `json:"quantity"`
	ImageURL          string                 `json:"image_url"`
	ProductAttributes map[string]interface{} `json:"product_attributes"`
}

type CartResponse struct {
	CartID              string             `json:"cart_id"`
	TotalItems          int                `json:"total_items"`
	SubtotalAmount      float64            `json:"subtotal_amount"`
	DiscountAmount      *float64           `json:"discount_amount,omitempty"`
	DiscountCodeApplied *string            `json:"discount_code_applied,omitempty"`
	Items               []CartItemResponse `json:"items"`
}
