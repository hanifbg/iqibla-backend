package shipping

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/hanifbg/landing_backend/internal/service/shipping/mocks"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test shipping service
func createTestShippingService(rajaOngkirClient RajaOngkirClientInterface) *ShippingService {
	return &ShippingService{
		rajaOngkirClient: rajaOngkirClient,
	}
}

// Helper function to create test province data
func createTestProvinces() []response.RajaOngkirProvince {
	return []response.RajaOngkirProvince{
		{
			ProvinceID: "1",
			Province:   "Bali",
		},
		{
			ProvinceID: "2",
			Province:   "Bangka Belitung",
		},
	}
}

// Helper function to create test city data
func createTestCities() []response.RajaOngkirCity {
	return []response.RajaOngkirCity{
		{
			CityID:     "1",
			ProvinceID: "1",
			Province:   "Bali",
			Type:       "Kabupaten",
			CityName:   "Badung",
			PostalCode: "80351",
		},
		{
			CityID:     "2",
			ProvinceID: "1",
			Province:   "Bali",
			Type:       "Kota",
			CityName:   "Denpasar",
			PostalCode: "80111",
		},
	}
}

// Helper function to create test shipping cost data
func createTestShippingCosts() []response.RajaOngkirCost {
	return []response.RajaOngkirCost{
		{
			Code: "jne",
			Name: "Jalur Nugraha Ekakurir (JNE)",
			Costs: []struct {
				Service     string `json:"service"`
				Description string `json:"description"`
				Cost        []struct {
					Value int    `json:"value"`
					ETD   string `json:"etd"`
					Note  string `json:"note"`
				} `json:"cost"`
			}{
				{
					Service:     "REG",
					Description: "Layanan Reguler",
					Cost: []struct {
						Value int    `json:"value"`
						ETD   string `json:"etd"`
						Note  string `json:"note"`
					}{
						{
							Value: 15000,
							ETD:   "1-2",
							Note:  "",
						},
					},
				},
				{
					Service:     "OKE",
					Description: "Ongkos Kirim Ekonomis",
					Cost: []struct {
						Value int    `json:"value"`
						ETD   string `json:"etd"`
						Note  string `json:"note"`
					}{
						{
							Value: 12000,
							ETD:   "2-3",
							Note:  "",
						},
					},
				},
			},
		},
	}
}

// Test GetProvinces method
func TestShippingService_GetProvinces(t *testing.T) {
	t.Run("Success - Get all provinces", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetProvincesRequest{}
		expectedProvinces := createTestProvinces()

		mockRajaOngkirClient.EXPECT().GetProvinces("").Return(expectedProvinces, nil)

		result, err := service.GetProvinces(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
		assert.Equal(t, "2", result[1].ProvinceID)
		assert.Equal(t, "Bangka Belitung", result[1].Province)
	})

	t.Run("Success - Get specific province by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetProvincesRequest{ID: "1"}
		expectedProvinces := []response.RajaOngkirProvince{
			{
				ProvinceID: "1",
				Province:   "Bali",
			},
		}

		mockRajaOngkirClient.EXPECT().GetProvinces("1").Return(expectedProvinces, nil)

		result, err := service.GetProvinces(req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Success - Empty provinces list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetProvincesRequest{}
		emptyProvinces := []response.RajaOngkirProvince{}

		mockRajaOngkirClient.EXPECT().GetProvinces("").Return(emptyProvinces, nil)

		result, err := service.GetProvinces(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetProvincesRequest{}
		expectedError := errors.New("RajaOngkir API error")

		mockRajaOngkirClient.EXPECT().GetProvinces("").Return(nil, expectedError)

		result, err := service.GetProvinces(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

// Test GetCities method
func TestShippingService_GetCities(t *testing.T) {
	t.Run("Success - Get all cities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetCitiesRequest{}
		expectedCities := createTestCities()

		mockRajaOngkirClient.EXPECT().GetCities("", "").Return(expectedCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "1", result[0].CityID)
		assert.Equal(t, "Badung", result[0].CityName)
		assert.Equal(t, "2", result[1].CityID)
		assert.Equal(t, "Denpasar", result[1].CityName)
	})

	t.Run("Success - Get cities by province ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetCitiesRequest{ProvinceID: "1"}
		expectedCities := createTestCities()

		mockRajaOngkirClient.EXPECT().GetCities("1", "").Return(expectedCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "1", result[1].ProvinceID)
	})

	t.Run("Success - Get specific city by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetCitiesRequest{ProvinceID: "1", ID: "1"}
		expectedCities := []response.RajaOngkirCity{
			{
				CityID:     "1",
				ProvinceID: "1",
				Province:   "Bali",
				Type:       "Kabupaten",
				CityName:   "Badung",
				PostalCode: "80351",
			},
		}

		mockRajaOngkirClient.EXPECT().GetCities("1", "1").Return(expectedCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].CityID)
		assert.Equal(t, "Badung", result[0].CityName)
	})

	t.Run("Success - Empty cities list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetCitiesRequest{}
		emptyCities := []response.RajaOngkirCity{}

		mockRajaOngkirClient.EXPECT().GetCities("", "").Return(emptyCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.GetCitiesRequest{ProvinceID: "1"}
		expectedError := errors.New("RajaOngkir API error")

		mockRajaOngkirClient.EXPECT().GetCities("1", "").Return(nil, expectedError)

		result, err := service.GetCities(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

// Test CalculateShippingCost method
func TestShippingService_CalculateShippingCost(t *testing.T) {
	t.Run("Success - Calculate shipping cost", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "jne",
		}
		expectedCosts := createTestShippingCosts()

		mockRajaOngkirClient.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(expectedCosts, nil)

		result, err := service.CalculateShippingCost(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2) // Two services: REG and OKE
		assert.Equal(t, "REG", result[0].Service)
		assert.Equal(t, "Layanan Reguler", result[0].Description)
		assert.Equal(t, float64(15000), result[0].Cost)
		assert.Equal(t, "1-2", result[0].ETD)
		assert.Equal(t, "OKE", result[1].Service)
		assert.Equal(t, "Ongkos Kirim Ekonomis", result[1].Description)
		assert.Equal(t, float64(12000), result[1].Cost)
		assert.Equal(t, "2-3", result[1].ETD)
	})

	t.Run("Success - Multiple couriers with multiple services", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "jne",
		}

		// Create test data with multiple cost values for one service
		multipleCosts := []response.RajaOngkirCost{
			{
				Code: "jne",
				Name: "Jalur Nugraha Ekakurir (JNE)",
				Costs: []struct {
					Service     string `json:"service"`
					Description string `json:"description"`
					Cost        []struct {
						Value int    `json:"value"`
						ETD   string `json:"etd"`
						Note  string `json:"note"`
					} `json:"cost"`
				}{
					{
						Service:     "REG",
						Description: "Layanan Reguler",
						Cost: []struct {
							Value int    `json:"value"`
							ETD   string `json:"etd"`
							Note  string `json:"note"`
						}{
							{
								Value: 15000,
								ETD:   "1-2",
								Note:  "",
							},
							{
								Value: 18000,
								ETD:   "1-1",
								Note:  "Express",
							},
						},
					},
				},
			},
		}

		mockRajaOngkirClient.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(multipleCosts, nil)

		result, err := service.CalculateShippingCost(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2) // Two cost values for REG service
		assert.Equal(t, "REG", result[0].Service)
		assert.Equal(t, float64(15000), result[0].Cost)
		assert.Equal(t, "REG", result[1].Service)
		assert.Equal(t, float64(18000), result[1].Cost)
	})

	t.Run("Success - Empty shipping costs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "invalid",
		}
		emptyCosts := []response.RajaOngkirCost{}

		mockRajaOngkirClient.EXPECT().CalculateShippingCost("501", "114", 1000, "invalid").Return(emptyCosts, nil)

		result, err := service.CalculateShippingCost(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "jne",
		}
		expectedError := errors.New("RajaOngkir API error")

		mockRajaOngkirClient.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(nil, expectedError)

		result, err := service.CalculateShippingCost(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("Error - Invalid weight", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      0,
			Courier:     "jne",
		}
		expectedError := errors.New("invalid weight")

		mockRajaOngkirClient.EXPECT().CalculateShippingCost("501", "114", 0, "jne").Return(nil, expectedError)

		result, err := service.CalculateShippingCost(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

// Test constructor function
func TestNewShippingService(t *testing.T) {
	t.Run("Success - Create shipping service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRajaOngkirClient := mocks.NewMockRajaOngkirClientInterface(ctrl)
		service := createTestShippingService(mockRajaOngkirClient)

		assert.NotNil(t, service)
		assert.Equal(t, mockRajaOngkirClient, service.rajaOngkirClient)
	})
}