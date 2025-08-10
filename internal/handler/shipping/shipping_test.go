package shipping

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockShippingService struct {
	mock.Mock
}

func (m *MockShippingService) GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error) {
	args := m.Called(req)
	return args.Get(0).([]response.ProvinceResponse), args.Error(1)
}

func (m *MockShippingService) GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error) {
	args := m.Called(req)
	return args.Get(0).([]response.CityResponse), args.Error(1)
}

func (m *MockShippingService) GetDistricts(req request.GetDistrictsRequest) ([]response.DistrictResponse, error) {
	args := m.Called(req)
	return args.Get(0).([]response.DistrictResponse), args.Error(1)
}

func (m *MockShippingService) CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error) {
	args := m.Called(req)
	return args.Get(0).([]response.ShippingCostResponse), args.Error(1)
}

func (m *MockShippingService) ValidateAndSaveAWB(req request.ValidateAWBRequest) (*response.ValidateAWBResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.ValidateAWBResponse), args.Error(1)
}

// ValidatorMock mocks the validator functionality
type ValidatorMock struct{}

func (v *ValidatorMock) Validate(i interface{}) error {
	if req, ok := i.(*request.ValidateAWBRequest); ok {
		if req.InvoiceNumber == "INVALID" {
			return errors.New("validation failed")
		}
	}
	return nil
}

func TestValidateAWB_Success(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &ValidatorMock{}

	req := request.ValidateAWBRequest{
		InvoiceNumber: "INV001",
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/shipping/awb/validate", bytes.NewReader(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockService := new(MockShippingService)
	h := &ApiWrapper{shippingService: mockService}

	responseData := &response.ValidateAWBResponse{
		ID:            "123e4567-e89b-12d3-a456-426614174000",
		InvoiceNumber: "INV001",
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
		IsValidated:   true,
		Message:       "AWB number validated and saved successfully",
	}

	mockService.On("ValidateAndSaveAWB", req).Return(responseData, nil)

	// Test
	err := h.ValidateAWB(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseMap)

	assert.Equal(t, "AWB validated and saved successfully", responseMap["message"])

	// Verify service was called
	mockService.AssertExpectations(t)
}

func TestValidateAWB_InvalidRequest(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &ValidatorMock{}

	// Invalid JSON
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/shipping/awb/validate", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockService := new(MockShippingService)
	h := &ApiWrapper{shippingService: mockService}

	// Test
	err := h.ValidateAWB(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var responseMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseMap)

	assert.Equal(t, "Invalid request", responseMap["error"])
}

func TestValidateAWB_ValidationFailed(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &ValidatorMock{}

	req := request.ValidateAWBRequest{
		InvoiceNumber: "INVALID", // This will cause validation to fail
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/shipping/awb/validate", bytes.NewReader(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockService := new(MockShippingService)
	h := &ApiWrapper{shippingService: mockService}

	// Test
	err := h.ValidateAWB(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var responseMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseMap)

	assert.Equal(t, "Validation failed", responseMap["error"])
}

func TestValidateAWB_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &ValidatorMock{}

	req := request.ValidateAWBRequest{
		InvoiceNumber: "INV001",
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/shipping/awb/validate", bytes.NewReader(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockService := new(MockShippingService)
	h := &ApiWrapper{shippingService: mockService}

	mockService.On("ValidateAndSaveAWB", req).Return(nil, errors.New("service error"))

	// Test
	err := h.ValidateAWB(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var responseMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseMap)

	assert.Equal(t, "Failed to validate AWB", responseMap["error"])

	// Verify service was called
	mockService.AssertExpectations(t)
}

func TestValidateAWB_AWBInvalid(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &ValidatorMock{}

	req := request.ValidateAWBRequest{
		InvoiceNumber: "INV001",
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/shipping/awb/validate", bytes.NewReader(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	mockService := new(MockShippingService)
	h := &ApiWrapper{shippingService: mockService}

	responseData := &response.ValidateAWBResponse{
		InvoiceNumber: "INV001",
		AWBNumber:     "JX123456789",
		Courier:       "jnt",
		IsValidated:   false,
		Message:       "Invalid AWB number",
	}

	mockService.On("ValidateAndSaveAWB", req).Return(responseData, nil)

	// Test
	err := h.ValidateAWB(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var responseMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseMap)

	assert.Equal(t, "Invalid AWB number", responseMap["error"])
	assert.Equal(t, responseData.Message, responseMap["message"])

	// Verify service was called
	mockService.AssertExpectations(t)
}
