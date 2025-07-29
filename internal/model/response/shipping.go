package response

type ProvinceResponse struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type CityResponse struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type ShippingCostResponse struct {
	Service     string  `json:"service"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	ETD         string  `json:"etd"`
}

type KomerceProvinceResponse struct {
	Meta Meta   `json:"meta"`
	Data []RajaOngkirProvince `json:"data"`
}

type KomerceCityResponse struct {
	Meta Meta   `json:"meta"`
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
	ProvinceID int `json:"id"`
	Province   string `json:"name"`
}

type RajaOngkirCity struct {
	CityID     int `json:"id"`
	ProvinceID int `json:"province_id"`
	Province   string `json:"province_name"`
	Type       string `json:"type"`
	CityName   string `json:"name"`
	PostalCode string `json:"postal_code"`
}

type RajaOngkirCost struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Costs []struct {
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost []struct {
			Value int    `json:"value"`
			ETD   string `json:"etd"`
			Note  string `json:"note"`
		} `json:"cost"`
	} `json:"costs"`
}