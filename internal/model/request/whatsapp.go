package request

type WhatsAppRequest struct {
	CustomerName          string
	OrderNumber           string
	TotalAmount           float64
	OrderConfirmationLink string
}
