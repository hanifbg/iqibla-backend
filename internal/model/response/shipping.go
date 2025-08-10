package response

type ProvinceResponse struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type CityResponse struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	CityName   string `json:"city_name"`
}

type ShippingCostResponse struct {
	Service     string  `json:"service"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	ETD         string  `json:"etd"`
}

// ValidateAWBResponse represents the response for AWB validation
type ValidateAWBResponse struct {
	ID            string `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	AWBNumber     string `json:"awb_number"`
	Courier       string `json:"courier"`
	IsValidated   bool   `json:"is_validated"`
	Message       string `json:"message"`
}

// RajaOngkirTrackingResponse represents the response from RajaOngkir tracking API
type RajaOngkirTrackingResponse struct {
	Meta RajaOngkirTrackingMeta `json:"meta"`
	Data interface{}            `json:"data"`
}

// RajaOngkirTrackingMeta represents the meta information in tracking response
type RajaOngkirTrackingMeta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type KomerceProvinceResponse struct {
	Meta Meta                 `json:"meta"`
	Data []RajaOngkirProvince `json:"data"`
}

type KomerceCityResponse struct {
	Meta Meta             `json:"meta"`
	Data []RajaOngkirCity `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type RajaOngkirResponse struct {
	RajaOngkir struct {
		Query   interface{} `json:"query"`
		Status  Status      `json:"status"`
		Results interface{} `json:"results"`
	} `json:"rajaongkir"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type RajaOngkirProvince struct {
	ProvinceID int    `json:"id"`
	Province   string `json:"name"`
}

type RajaOngkirCity struct {
	CityID     int    `json:"id"`
	ProvinceID int    `json:"-"` // Not included in API response anymore, will be set from URL
	Province   string `json:"-"` // Not included in API response anymore
	Type       string `json:"-"` // Not included in API response anymore
	CityName   string `json:"name"`
	PostalCode string `json:"-"` // Not included in API response anymore
}

type DistrictResponse struct {
	DistrictID   string `json:"district_id"`
	CityID       string `json:"city_id"`
	DistrictName string `json:"district_name"`
}

type RajaOngkirDistrict struct {
	DistrictID   int    `json:"id"`
	CityID       int    `json:"-"` // Not included in API response anymore
	City         string `json:"-"` // Not included in API response anymore
	DistrictName string `json:"name"`
	Type         string `json:"-"` // Not included in API response anymore
}

type RajaOngkirCost struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	ETD         string `json:"etd"`
}
