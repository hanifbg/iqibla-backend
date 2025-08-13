package request

// orderEmailItem represents a single item in the order summary for the email
// This aligns with fields used in the mail.html template
// - ProductName
// - Quantity
// - PriceAtPurchase
// We also keep raw values and display formatting simple and let template show raw numbers for now
// Consider adding proper currency formatting if needed later

type OrderEmailItem struct {
	ProductName     string
	Quantity        int
	PriceAtPurchase string
}

// orderEmailData represents all data needed by the mail.html template
// Fields map to template variables like .CustomerName, .OrderNumber, etc.

type OrderEmailData struct {
	CustomerName          string
	OrderNumber           string
	OrderItems            []OrderEmailItem
	SubtotalAmount        float64
	ShippingCost          float64
	TotalAmount           string
	OrderConfirmationLink string
}
