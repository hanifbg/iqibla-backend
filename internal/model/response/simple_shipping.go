package response

// SimpleRajaOngkirCity represents the city response from RajaOngkir API in simplified format
type SimpleRajaOngkirCity struct {
	CityID int    `json:"id"`
	Name   string `json:"name"`
}

// KomerceSimpleCityResponse represents the response from the Komerce City API with simplified city format
type KomerceSimpleCityResponse struct {
	Meta Meta                   `json:"meta"`
	Data []SimpleRajaOngkirCity `json:"data"`
}
