package entity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAWBTracking_TrackingDataScanAndValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected *TrackingData
		hasError bool
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: &TrackingData{},
			hasError: false,
		},
		{
			name:  "valid JSON bytes",
			input: []byte(`{"delivered": true, "summary": {"courier_code": "J&T"}}`),
			expected: &TrackingData{
				Delivered: true,
				Summary: TrackingSummary{
					CourierCode: "J&T",
				},
			},
			hasError: false,
		},
		{
			name:  "valid JSON string",
			input: `{"delivered": false, "summary": {"courier_code": "JNE"}}`,
			expected: &TrackingData{
				Delivered: false,
				Summary: TrackingSummary{
					CourierCode: "JNE",
				},
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var td TrackingData
			err := td.Scan(tc.input)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.input != nil {
					assert.Equal(t, tc.expected.Delivered, td.Delivered)
					assert.Equal(t, tc.expected.Summary.CourierCode, td.Summary.CourierCode)
				}
			}
		})
	}
}

func TestAWBTracking_TrackingDataValue(t *testing.T) {
	td := TrackingData{
		Delivered: true,
		Summary: TrackingSummary{
			CourierCode:   "J&T",
			CourierName:   "J&T Express",
			WaybillNumber: "JX123456789",
			Status:        "DELIVERED",
		},
		Details: TrackingDetails{
			WaybillNumber: "JX123456789",
			Weight:        "1000",
			Origin:        "JAKARTA",
			Destination:   "BANDUNG",
		},
	}

	value, err := td.Value()
	assert.NoError(t, err)

	// Verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.True(t, result["delivered"].(bool))
}

func TestValidCouriers(t *testing.T) {
	couriers := ValidCouriers()

	expectedCouriers := []string{
		"jne", "jnt", "ninja", "tiki", "pos",
		"anteraja", "sicepat", "sap", "lion",
		"wahana", "first", "ide",
	}

	assert.Equal(t, len(expectedCouriers), len(couriers))
	for _, expected := range expectedCouriers {
		assert.Contains(t, couriers, expected)
	}
}

func TestIsValidCourier(t *testing.T) {
	testCases := []struct {
		courier  string
		expected bool
	}{
		{"jne", true},
		{"jnt", true},
		{"ninja", true},
		{"tiki", true},
		{"pos", true},
		{"anteraja", true},
		{"sicepat", true},
		{"sap", true},
		{"lion", true},
		{"wahana", true},
		{"first", true},
		{"ide", true},
		{"invalid", false},
		{"JNE", false}, // case sensitive
		{"", false},
		{"unknown", false},
	}

	for _, tc := range testCases {
		t.Run(tc.courier, func(t *testing.T) {
			result := IsValidCourier(tc.courier)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAWBTracking_EntityCreation(t *testing.T) {
	orderID := uuid.New()
	awbID := uuid.New()
	phoneNumber := "12345"
	now := time.Now()

	awb := &AWBTracking{
		ID:              awbID,
		OrderID:         orderID,
		AWBNumber:       "JX123456789",
		Courier:         "jnt",
		LastPhoneNumber: &phoneNumber,
		IsValidated:     true,
		TrackingData: &TrackingData{
			Delivered: true,
			Summary: TrackingSummary{
				CourierCode:   "J&T",
				CourierName:   "J&T Express",
				WaybillNumber: "JX123456789",
				Status:        "DELIVERED",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Verify all fields are set correctly
	assert.Equal(t, awbID, awb.ID)
	assert.Equal(t, orderID, awb.OrderID)
	assert.Equal(t, "JX123456789", awb.AWBNumber)
	assert.Equal(t, "jnt", awb.Courier)
	assert.Equal(t, &phoneNumber, awb.LastPhoneNumber)
	assert.True(t, awb.IsValidated)
	assert.NotNil(t, awb.TrackingData)
	assert.True(t, awb.TrackingData.Delivered)
	assert.Equal(t, "J&T", awb.TrackingData.Summary.CourierCode)
	assert.Equal(t, now, awb.CreatedAt)
	assert.Equal(t, now, awb.UpdatedAt)
	assert.Nil(t, awb.DeletedAt)
}

func TestTrackingData_CompleteStructure(t *testing.T) {
	manifest := []Manifest{
		{
			ManifestCode:        "100",
			ManifestDescription: "Package received",
			ManifestDate:        "2025-01-01",
			ManifestTime:        "10:00:00",
			CityName:            "JAKARTA",
		},
		{
			ManifestCode:        "200",
			ManifestDescription: "Package delivered",
			ManifestDate:        "2025-01-02",
			ManifestTime:        "15:30:00",
			CityName:            "BANDUNG",
		},
	}

	td := &TrackingData{
		Delivered: true,
		Summary: TrackingSummary{
			CourierCode:   "J&T",
			CourierName:   "J&T Express",
			WaybillNumber: "JX123456789",
			ServiceCode:   "REG",
			WaybillDate:   "2025-01-01",
			ShipperName:   "Test Shipper",
			ReceiverName:  "Test Receiver",
			Origin:        "JAKARTA",
			Destination:   "BANDUNG",
			Status:        "DELIVERED",
		},
		Details: TrackingDetails{
			WaybillNumber:    "JX123456789",
			WaybillDate:      "2025-01-01",
			WaybillTime:      "10:00:00",
			Weight:           "1000",
			Origin:           "JAKARTA",
			Destination:      "BANDUNG",
			ShipperName:      "Test Shipper",
			ShipperAddress1:  "Jl. Test Street No. 1",
			ShipperCity:      "JAKARTA",
			ReceiverName:     "Test Receiver",
			ReceiverAddress1: "Jl. Receiver Street No. 2",
			ReceiverCity:     "BANDUNG",
		},
		DeliveryStatus: DeliveryStatus{
			Status:      "DELIVERED",
			PODReceiver: "Test Receiver",
			PODDate:     "2025-01-02",
			PODTime:     "15:30:00",
		},
		Manifest: manifest,
	}

	// Test Value method
	value, err := td.Value()
	require.NoError(t, err)

	// Test Scan method
	var scannedTD TrackingData
	err = scannedTD.Scan(value)
	require.NoError(t, err)

	// Verify all data is preserved
	assert.Equal(t, td.Delivered, scannedTD.Delivered)
	assert.Equal(t, td.Summary.CourierCode, scannedTD.Summary.CourierCode)
	assert.Equal(t, td.Details.WaybillNumber, scannedTD.Details.WaybillNumber)
	assert.Equal(t, td.DeliveryStatus.Status, scannedTD.DeliveryStatus.Status)
	assert.Len(t, scannedTD.Manifest, 2)
	assert.Equal(t, td.Manifest[0].ManifestCode, scannedTD.Manifest[0].ManifestCode)
	assert.Equal(t, td.Manifest[1].ManifestDescription, scannedTD.Manifest[1].ManifestDescription)
}
