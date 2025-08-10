package shipping

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

func (s *ShippingService) GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error) {
	// Validate input if necessary

	// Call repository
	provinces, err := s.ShippingRepo.GetProvinces(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces: %w", err)
	}

	// Transform to service response format
	var result []response.ProvinceResponse
	for _, province := range provinces {
		result = append(result, response.ProvinceResponse{
			ProvinceID: strconv.Itoa(province.ProvinceID),
			Province:   province.Province,
		})
	}

	return result, nil
}

func (s *ShippingService) GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error) {
	// Validate input if necessary

	// Call repository
	cities, err := s.ShippingRepo.GetCities(req.ProvinceID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities: %w", err)
	}

	// Transform to service response format
	var result []response.CityResponse
	for _, city := range cities {
		result = append(result, response.CityResponse{
			CityID:     strconv.Itoa(city.CityID),
			ProvinceID: strconv.Itoa(city.ProvinceID),
			CityName:   city.CityName,
		})
	}

	return result, nil
}

func (s *ShippingService) GetDistricts(req request.GetDistrictsRequest) ([]response.DistrictResponse, error) {
	// Validate input
	if req.CityID == "" {
		return nil, fmt.Errorf("city ID is required")
	}

	// Call repository
	districts, err := s.ShippingRepo.GetDistricts(req.CityID)
	if err != nil {
		return nil, err
	}

	var result []response.DistrictResponse
	for _, district := range districts {
		result = append(result, response.DistrictResponse{
			DistrictID:   strconv.Itoa(district.DistrictID),
			CityID:       req.CityID, // Use the city ID from the request parameter
			DistrictName: district.DistrictName,
		})
	}

	return result, nil
}

func (s *ShippingService) CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error) {
	// Validate input
	if req.Origin == "" || req.Destination == "" || req.Weight <= 0 || req.Courier == "" {
		return nil, fmt.Errorf("origin, destination, weight, and courier are required")
	}

	// Call repository
	costs, err := s.ShippingRepo.CalculateShippingCost(req.Origin, req.Destination, req.Weight, req.Courier)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping cost: %w", err)
	}

	var result []response.ShippingCostResponse
	for _, cost := range costs {
		result = append(result, response.ShippingCostResponse{
			Service:     cost.Service,
			Description: cost.Description,
			Cost:        float64(cost.Cost),
			ETD:         cost.ETD,
		})
	}

	return result, nil
}

// ValidateAndSaveAWB validates AWB number with RajaOngkir and saves it to database
func (s *ShippingService) ValidateAndSaveAWB(req request.ValidateAWBRequest) (*response.ValidateAWBResponse, error) {
	// Step 1: Validate that the invoice number exists
	order, err := s.AWBTrackingRepo.GetOrderByInvoiceNumber(req.InvoiceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to validate invoice number: %w", err)
	}

	// Step 2: Convert order ID to UUID
	orderID, err := uuid.Parse(order.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	// Step 3: Check if AWB already exists for this order and courier
	existingAWB, err := s.AWBTrackingRepo.GetAWBTrackingByAWBNumber(req.AWBNumber, req.Courier)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing AWB: %w", err)
	}
	if existingAWB != nil {
		return nil, fmt.Errorf("AWB number '%s' for courier '%s' already exists", req.AWBNumber, req.Courier)
	}

	// Step 4: Validate AWB with RajaOngkir API
	trackingData, err := s.ShippingRepo.ValidateAWB(req.AWBNumber, req.Courier, req.LastPhoneNumber)
	if err != nil {
		return &response.ValidateAWBResponse{
			AWBNumber:     req.AWBNumber,
			Courier:       req.Courier,
			InvoiceNumber: req.InvoiceNumber,
			IsValidated:   false,
			Message:       "Invalid AWB number",
		}, nil // Return successful response with validation failure
	}

	// Step 5: Extract tracking data from API response
	var parsedTrackingData *entity.TrackingData
	if trackingData != nil && trackingData.Data != nil {
		// Parse the tracking data - this depends on the actual structure from RajaOngkir
		trackingBytes, _ := json.Marshal(trackingData.Data)
		var trackingInfo entity.TrackingData
		if json.Unmarshal(trackingBytes, &trackingInfo) == nil {
			parsedTrackingData = &trackingInfo
		}
	}

	// Step 6: Create AWB tracking record
	awbTracking := &entity.AWBTracking{
		ID:              uuid.New(),
		OrderID:         orderID,
		AWBNumber:       req.AWBNumber,
		Courier:         req.Courier,
		LastPhoneNumber: req.LastPhoneNumber,
		IsValidated:     true,
		TrackingData:    parsedTrackingData,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Step 7: Save to database
	if err := s.AWBTrackingRepo.CreateAWBTracking(awbTracking); err != nil {
		return nil, fmt.Errorf("failed to save AWB tracking: %w", err)
	}

	// Step 8: Return success response
	return &response.ValidateAWBResponse{
		ID:            awbTracking.ID.String(),
		InvoiceNumber: req.InvoiceNumber,
		AWBNumber:     req.AWBNumber,
		Courier:       req.Courier,
		IsValidated:   true,
		Message:       "AWB number validated and saved successfully",
	}, nil
}
