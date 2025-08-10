package repository

import (
	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
)

//go:generate mockgen -source=awb_tracking.go -destination=mocks/awb_tracking_repository_mock.go -package=mocks

// AWBTrackingRepository defines the interface for AWB tracking data operations
type AWBTrackingRepository interface {
	// CreateAWBTracking creates a new AWB tracking record
	CreateAWBTracking(awbTracking *entity.AWBTracking) error

	// GetAWBTrackingByOrderID retrieves AWB tracking records by order ID
	GetAWBTrackingByOrderID(orderID uuid.UUID) ([]*entity.AWBTracking, error)

	// GetAWBTrackingByAWBNumber retrieves AWB tracking record by AWB number and courier
	GetAWBTrackingByAWBNumber(awbNumber, courier string) (*entity.AWBTracking, error)

	// UpdateAWBTracking updates an existing AWB tracking record
	UpdateAWBTracking(awbTracking *entity.AWBTracking) error

	// DeleteAWBTracking soft deletes an AWB tracking record
	DeleteAWBTracking(id uuid.UUID) error

	// GetOrderByInvoiceNumber retrieves order by invoice number to validate it exists
	GetOrderByInvoiceNumber(invoiceNumber string) (*entity.Order, error)
}

// AWBTrackingError represents errors from the AWB tracking repository
type AWBTrackingError struct {
	Operation string // Operation that failed
	Err       error  // Original error
}

// Error returns the string representation of the error
func (e *AWBTrackingError) Error() string {
	if e.Err != nil {
		return e.Operation + ": " + e.Err.Error()
	}
	return e.Operation
}

// Unwrap returns the underlying error
func (e *AWBTrackingError) Unwrap() error {
	return e.Err
}
