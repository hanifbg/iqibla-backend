package request

import "strings"

type SendMessage struct {
	ChatID          int64  `json:"chat_id"`
	MessageThreadID int64  `json:"message_thread_id"`
	Text            string `json:"text"`
	ParseMode       string `json:"parse_mode"`
}

// ShippingAddressData holds the granular data for the shipping address.
type TelegramShippingAddressData struct {
	ShippingStreetAddress string
	ShippingCity          string
	ShippingProvince      string
	ShippingDistrict      string
	ShippingPostalCode    string
	ShippingCountry       string
}

type TelegramRequest struct {
	OrderID         string
	OrderNumber     string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	TotalAmount     float64
	ShippingAddress string
	ShippingCourier string
	ShippingService string
	OrderItems      []struct {
		ProductName     string
		Quantity        int
		PriceAtPurchase float64
	}
}

// FormatShippingAddress combines granular address data into a single formatted string.
func FormatShippingAddress(data TelegramShippingAddressData) string {
	var parts []string
	if data.ShippingStreetAddress != "" {
		parts = append(parts, data.ShippingStreetAddress)
	}
	if data.ShippingDistrict != "" {
		parts = append(parts, data.ShippingDistrict)
	}
	if data.ShippingCity != "" {
		parts = append(parts, data.ShippingCity)
	}
	if data.ShippingProvince != "" {
		parts = append(parts, data.ShippingProvince)
	}
	if data.ShippingPostalCode != "" {
		parts = append(parts, data.ShippingPostalCode)
	}
	if data.ShippingCountry != "" {
		parts = append(parts, data.ShippingCountry)
	}
	return strings.Join(parts, ", ")
}
