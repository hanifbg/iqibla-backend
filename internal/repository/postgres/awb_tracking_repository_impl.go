package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/repository"
	"gorm.io/gorm"
)

// AWBTrackingRepositoryImpl implements the AWBTrackingRepository interface
type AWBTrackingRepositoryImpl struct {
	db *gorm.DB
}

// NewAWBTrackingRepository creates a new instance of AWBTrackingRepositoryImpl
func NewAWBTrackingRepository(db *gorm.DB) repository.AWBTrackingRepository {
	return &AWBTrackingRepositoryImpl{
		db: db,
	}
}

// CreateAWBTracking creates a new AWB tracking record
func (r *AWBTrackingRepositoryImpl) CreateAWBTracking(awbTracking *entity.AWBTracking) error {
	if err := r.db.Create(awbTracking).Error; err != nil {
		return &repository.AWBTrackingError{
			Operation: "CreateAWBTracking",
			Err:       err,
		}
	}
	return nil
}

// GetAWBTrackingByOrderID retrieves AWB tracking records by order ID
func (r *AWBTrackingRepositoryImpl) GetAWBTrackingByOrderID(orderID uuid.UUID) ([]*entity.AWBTracking, error) {
	var awbTrackings []*entity.AWBTracking
	if err := r.db.Where("order_id = ? AND deleted_at IS NULL", orderID).Find(&awbTrackings).Error; err != nil {
		return nil, &repository.AWBTrackingError{
			Operation: "GetAWBTrackingByOrderID",
			Err:       err,
		}
	}
	return awbTrackings, nil
}

// GetAWBTrackingByAWBNumber retrieves AWB tracking record by AWB number and courier
func (r *AWBTrackingRepositoryImpl) GetAWBTrackingByAWBNumber(awbNumber, courier string) (*entity.AWBTracking, error) {
	var awbTracking entity.AWBTracking
	err := r.db.Where("awb_number = ? AND courier = ? AND deleted_at IS NULL", awbNumber, courier).First(&awbTracking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found, not an error
		}
		return nil, &repository.AWBTrackingError{
			Operation: "GetAWBTrackingByAWBNumber",
			Err:       err,
		}
	}
	return &awbTracking, nil
}

// UpdateAWBTracking updates an existing AWB tracking record
func (r *AWBTrackingRepositoryImpl) UpdateAWBTracking(awbTracking *entity.AWBTracking) error {
	if err := r.db.Save(awbTracking).Error; err != nil {
		return &repository.AWBTrackingError{
			Operation: "UpdateAWBTracking",
			Err:       err,
		}
	}
	return nil
}

// DeleteAWBTracking soft deletes an AWB tracking record
func (r *AWBTrackingRepositoryImpl) DeleteAWBTracking(id uuid.UUID) error {
	if err := r.db.Delete(&entity.AWBTracking{}, id).Error; err != nil {
		return &repository.AWBTrackingError{
			Operation: "DeleteAWBTracking",
			Err:       err,
		}
	}
	return nil
}

// GetOrderByInvoiceNumber retrieves order by invoice number to validate it exists
func (r *AWBTrackingRepositoryImpl) GetOrderByInvoiceNumber(invoiceNumber string) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Where("order_number = ? AND deleted_at IS NULL", invoiceNumber).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &repository.AWBTrackingError{
				Operation: "GetOrderByInvoiceNumber",
				Err:       fmt.Errorf("order with invoice number '%s' not found", invoiceNumber),
			}
		}
		return nil, &repository.AWBTrackingError{
			Operation: "GetOrderByInvoiceNumber",
			Err:       err,
		}
	}
	return &order, nil
}
