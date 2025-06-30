package request

type GetProvincesRequest struct {
	ID string `json:"id,omitempty"`
}

type GetCitiesRequest struct {
	ProvinceID string `json:"province_id,omitempty"`
	ID         string `json:"id,omitempty"`
}

type CalculateShippingRequest struct {
	Origin      string `json:"origin" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Weight      int    `json:"weight" validate:"required"`
	Courier     string `json:"courier" validate:"required"`
}