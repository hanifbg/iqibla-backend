package request

type GetProvincesRequest struct {
	ID string `json:"id,omitempty"`
}

type GetCitiesRequest struct {
	ProvinceID string `json:"province_id,omitempty"`
	ID         string `json:"id,omitempty"`
}

type GetDistrictsRequest struct {
	CityID string `json:"city_id,omitempty"`
}

type CalculateShippingRequest struct {
	Origin      string `json:"origin" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Weight      int    `json:"weight" validate:"required"`
	Courier     string `json:"courier" validate:"required"`
}

// ValidateAWBRequest represents the request to validate and save AWB number
type ValidateAWBRequest struct {
	InvoiceNumber   string  `json:"invoice_number" validate:"required"`                                                                  // Invoice number to link AWB to an order
	AWBNumber       string  `json:"awb_number" validate:"required"`                                                                      // AWB tracking number from courier
	Courier         string  `json:"courier" validate:"required,oneof=jne jnt ninja tiki pos anteraja sicepat sap lion wahana first ide"` // Courier service name
	LastPhoneNumber *string `json:"last_phone_number,omitempty" validate:"omitempty,len=5,numeric"`                                      // Last 5 digits of recipient's phone number (only required for JNE courier)
}
