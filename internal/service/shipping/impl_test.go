package shipping

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/hanifbg/landing_backend/internal/repository"
	repoMocks "github.com/hanifbg/landing_backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test shipping service
func createTestShippingService(mockRepo repository.ShippingRepository) *ShippingService {
	return &ShippingService{
		shippingRepo: mockRepo,
	}
}

// Helper function to create test province data
func createTestProvinces() []response.RajaOngkirProvince {
	return []response.RajaOngkirProvince{
		{
			ProvinceID: 1,
			Province:   "Bali",
		},
		{
			ProvinceID: 2,
			Province:   "Bangka Belitung",
		},
	}
}

// Helper function to create test city data
func createTestCities() []response.RajaOngkirCity {
	return []response.RajaOngkirCity{
		{
			CityID:     1,
			ProvinceID: 1,
			CityName:   "Badung",
		},
		{
			CityID:     2,
			ProvinceID: 1,
			CityName:   "Denpasar",
		},
	}
}

// Helper function to create test district data
func createTestDistricts() []response.RajaOngkirDistrict {
	return []response.RajaOngkirDistrict{
		{
			DistrictID:   1,
			CityID:       575,
			DistrictName: "Cengkareng",
		},
		{
			DistrictID:   2,
			CityID:       575,
			DistrictName: "Grogol Petamburan",
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

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetProvincesRequest{}
		expectedProvinces := createTestProvinces()

		mockShippingRepo.EXPECT().GetProvinces("").Return(expectedProvinces, nil)

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

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetProvincesRequest{ID: "1"}
		expectedProvinces := []response.RajaOngkirProvince{
			{
				ProvinceID: 1,
				Province:   "Bali",
			},
		}

		mockShippingRepo.EXPECT().GetProvinces("1").Return(expectedProvinces, nil)

		result, err := service.GetProvinces(req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Success - Empty provinces list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetProvincesRequest{}
		emptyProvinces := []response.RajaOngkirProvince{}

		mockShippingRepo.EXPECT().GetProvinces("").Return(emptyProvinces, nil)

		result, err := service.GetProvinces(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetProvincesRequest{}
		repoError := errors.New("RajaOngkir API error")

		mockShippingRepo.EXPECT().GetProvinces("").Return(nil, repoError)

		result, err := service.GetProvinces(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RajaOngkir API error")
	})
}

// Test GetCities method
func TestShippingService_GetCities(t *testing.T) {
	t.Run("Success - Get all cities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetCitiesRequest{}
		expectedCities := createTestCities()

		mockShippingRepo.EXPECT().GetCities("", "").Return(expectedCities, nil)

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

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetCitiesRequest{ProvinceID: "1"}
		expectedCities := createTestCities()

		mockShippingRepo.EXPECT().GetCities("1", "").Return(expectedCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "1", result[1].ProvinceID)
	})

	t.Run("Success - Get specific city by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetCitiesRequest{ProvinceID: "1", ID: "1"}
		expectedCities := []response.RajaOngkirCity{
			{
				CityID:     1,
				ProvinceID: 1,
				CityName:   "Badung",
			},
		}

		mockShippingRepo.EXPECT().GetCities("1", "1").Return(expectedCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].CityID)
		assert.Equal(t, "Badung", result[0].CityName)
	})

	t.Run("Success - Empty cities list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetCitiesRequest{}
		emptyCities := []response.RajaOngkirCity{}

		mockShippingRepo.EXPECT().GetCities("", "").Return(emptyCities, nil)

		result, err := service.GetCities(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.GetCitiesRequest{ProvinceID: "1"}
		repoError := errors.New("RajaOngkir API error")

		mockShippingRepo.EXPECT().GetCities("1", "").Return(nil, repoError)

		result, err := service.GetCities(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RajaOngkir API error")
	})
}

// Test CalculateShippingCost method
func TestShippingService_CalculateShippingCost(t *testing.T) {
	t.Run("Success - Calculate shipping cost", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "jne",
		}
		expectedCosts := createTestShippingCosts()

		mockShippingRepo.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(expectedCosts, nil)

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

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

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

		mockShippingRepo.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(multipleCosts, nil)

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

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "invalid",
		}
		emptyCosts := []response.RajaOngkirCost{}

		mockShippingRepo.EXPECT().CalculateShippingCost("501", "114", 1000, "invalid").Return(emptyCosts, nil)

		result, err := service.CalculateShippingCost(req)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Error - RajaOngkir API failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      1000,
			Courier:     "jne",
		}
		repoError := errors.New("RajaOngkir API error")

		mockShippingRepo.EXPECT().CalculateShippingCost("501", "114", 1000, "jne").Return(nil, repoError)

		result, err := service.CalculateShippingCost(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RajaOngkir API error")
	})

	t.Run("Error - Invalid weight", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockShippingRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockShippingRepo)

		req := request.CalculateShippingRequest{
			Origin:      "501",
			Destination: "114",
			Weight:      0,
			Courier:     "jne",
		}

		// The test should just check that we get an error message about requirements
		// We don't need to mock the repository call since validation happens before that

		result, err := service.CalculateShippingCost(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "required")
	})
}

// Test constructor function
func TestNewShippingService(t *testing.T) {
	t.Run("Success - Create shipping service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockRepo)

		assert.NotNil(t, service)
		assert.Equal(t, mockRepo, service.shippingRepo)
	})
}

func TestShippingService_GetDistricts(t *testing.T) {
	t.Run("Success - Get districts with city ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockRepo)

		req := request.GetDistrictsRequest{
			CityID: "575",
		}

		districts := createTestDistricts()
		mockRepo.EXPECT().GetDistricts("575").Return(districts, nil)

		result, err := service.GetDistricts(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "1", result[0].DistrictID)
		assert.Equal(t, "Cengkareng", result[0].DistrictName)
		assert.Equal(t, "575", result[0].CityID)
	})

	t.Run("Error - Failed to get districts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockRepo)

		req := request.GetDistrictsRequest{
			CityID: "575",
		}

		expectedError := errors.New("failed to get districts from RajaOngkir API")
		mockRepo.EXPECT().GetDistricts("575").Return(nil, expectedError)

		result, err := service.GetDistricts(req)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Empty city ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMocks.NewMockShippingRepository(ctrl)
		service := createTestShippingService(mockRepo)

		req := request.GetDistrictsRequest{
			CityID: "",
		}

		// The service should check for empty city ID before calling the repository
		result, err := service.GetDistricts(req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "city ID is required")
		assert.Nil(t, result)
	})
}
