package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AWBTracking represents an AWB (Air Way Bill) tracking record in the system
type AWBTracking struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	OrderID         uuid.UUID     `json:"order_id" db:"order_id"`
	AWBNumber       string        `json:"awb_number" db:"awb_number"`
	Courier         string        `json:"courier" db:"courier"`
	LastPhoneNumber *string       `json:"last_phone_number,omitempty" db:"last_phone_number"`
	IsValidated     bool          `json:"is_validated" db:"is_validated"`
	TrackingData    *TrackingData `json:"tracking_data,omitempty" db:"tracking_data"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
}

// TableName specifies the table name for GORM
func (AWBTracking) TableName() string {
	return "awb_tracking"
}

// TrackingData represents the detailed tracking information from RajaOngkir API
type TrackingData struct {
	Delivered      bool            `json:"delivered"`
	Summary        TrackingSummary `json:"summary"`
	Details        TrackingDetails `json:"details"`
	DeliveryStatus DeliveryStatus  `json:"delivery_status"`
	Manifest       []Manifest      `json:"manifest"`
}

// TrackingSummary represents summary information of the shipment
type TrackingSummary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

// TrackingDetails represents detailed information of the shipment
type TrackingDetails struct {
	WaybillNumber    string `json:"waybill_number"`
	WaybillDate      string `json:"waybill_date"`
	WaybillTime      string `json:"waybill_time"`
	Weight           string `json:"weight"`
	Origin           string `json:"origin"`
	Destination      string `json:"destination"`
	ShipperName      string `json:"shipper_name"`
	ShipperAddress1  string `json:"shipper_address1"`
	ShipperAddress2  string `json:"shipper_address2"`
	ShipperAddress3  string `json:"shipper_address3"`
	ShipperCity      string `json:"shipper_city"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverAddress1 string `json:"receiver_address1"`
	ReceiverAddress2 string `json:"receiver_address2"`
	ReceiverAddress3 string `json:"receiver_address3"`
	ReceiverCity     string `json:"receiver_city"`
}

// DeliveryStatus represents the delivery status information
type DeliveryStatus struct {
	Status      string `json:"status"`
	PODReceiver string `json:"pod_receiver"`
	PODDate     string `json:"pod_date"`
	PODTime     string `json:"pod_time"`
}

// Manifest represents a single tracking event in the shipment journey
type Manifest struct {
	ManifestCode        string `json:"manifest_code"`
	ManifestDescription string `json:"manifest_description"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name"`
}

// Scan implements the sql.Scanner interface for TrackingData
func (td *TrackingData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), td)
	}

	return json.Unmarshal(bytes, td)
}

// Value implements the driver.Valuer interface for TrackingData
func (td TrackingData) Value() (driver.Value, error) {
	return json.Marshal(td)
}

// ValidCouriers returns a list of valid courier codes
func ValidCouriers() []string {
	return []string{
		"jne", "jnt", "ninja", "tiki", "pos",
		"anteraja", "sicepat", "sap", "lion",
		"wahana", "first", "ide",
	}
}

// IsValidCourier checks if the given courier code is valid
func IsValidCourier(courier string) bool {
	validCouriers := ValidCouriers()
	for _, validCourier := range validCouriers {
		if courier == validCourier {
			return true
		}
	}
	return false
}
