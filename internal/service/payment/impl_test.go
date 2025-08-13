package payment

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/service/payment/mocks"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test payment service
func createTestPaymentService(paymentRepo repository.PaymentRepository, cartRepo repository.CartRepository, snapClient SnapClientInterface) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepo,
		cartRepo:     cartRepo,
		snapClient:   snapClient,
		baseURL:      "http://localhost:8080",
		mailer:       nil,
		whatsAppRepo: nil,
	}
}

// Test CreateOrder method
func TestPaymentService_CreateOrder(t *testing.T) {
	t.Run("Success - Create order with valid cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		cart := createTestCartWithItems()
		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
			Notes:                "Test notes",
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockPaymentRepo.EXPECT().GetSeq().Return(int64(1), nil)
		mockPaymentRepo.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any()).Return(nil)

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		// assert.Equal(t, "John Doe", result.CustomerName)
		// assert.Equal(t, "john@example.com", result.CustomerEmail)
		// assert.Equal(t, 250.0, result.Subtotal)
		// assert.Equal(t, 250.0, result.TotalAmount)
		// assert.Equal(t, "pending", result.OrderStatus)
		// assert.Len(t, result.OrderItems, 2)
	})

	t.Run("Error - Cart not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		req := request.CreateOrderRequest{
			CartID:               "invalid-cart",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("invalid-cart").Return(nil, errors.New("cart not found"))

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get cart")

	})

	t.Run("Error - Empty cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		emptyCart := &entity.Cart{
			ID: "cart-123",

			CartItems: []entity.CartItem{},
		}

		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(emptyCart, nil)

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart is empty")

	})

	t.Run("Error - Failed to create order", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		cart := createTestCartWithItems()
		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockPaymentRepo.EXPECT().GetSeq().Return(int64(1), nil)
		mockPaymentRepo.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create order with items")

	})
}

// Test GetOrder method
func TestPaymentService_GetOrder(t *testing.T) {
	t.Run("Success - Get existing order", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		order.OrderItems = []entity.OrderItem{
			{
				ID:               "item-1",
				OrderID:          "order-123",
				ProductVariantID: "variant-1",
				Quantity:         2,
				PriceAtPurchase:  100.0,
				ProductVariant: &entity.ProductVariant{
					ID:    "variant-1",
					Name:  "Test Product",
					Price: 100.0,
				},
			},
		}

		mockPaymentRepo.EXPECT().GetOrderWithItems("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, nil)

		// Act
		result, err := service.GetOrder("order-123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "order-123", result.ID)
		assert.Equal(t, "John Doe", result.CustomerName)
		assert.Len(t, result.OrderItems, 1)
	})

	t.Run("Error - Order not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		mockPaymentRepo.EXPECT().GetOrderWithItems("invalid-order").Return(nil, errors.New("order not found"))

		// Act
		result, err := service.GetOrder("invalid-order")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get order")
	})
}

// Test CreatePayment method
func TestPaymentService_CreatePayment(t *testing.T) {
	t.Run("Success - Create payment with valid order", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(&snap.Response{Token: "test-token", RedirectURL: "http://test.com"}, (*midtrans.Error)(nil))
		mockPaymentRepo.EXPECT().CreatePayment(gomock.Any()).Return(nil)

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "order-123", result.OrderID)
		assert.Equal(t, 250.0, result.Amount)
		assert.Equal(t, entity.PaymentStatusPending, result.Status)
		assert.NotEmpty(t, result.PaymentToken)
		assert.NotEmpty(t, result.PaymentURL)
		assert.NotNil(t, result.ExpiryTime)
	})

	t.Run("Success - Return existing payment", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		existingPayment := createTestPayment()

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(existingPayment, nil)

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingPayment.ID, result.ID)
		assert.Equal(t, existingPayment.OrderID, result.OrderID)
		assert.Equal(t, existingPayment.PaymentToken, result.PaymentToken)
	})

	t.Run("Error - Order not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		mockPaymentRepo.EXPECT().FindOrderByID("invalid-order").Return(nil, errors.New("order not found"))

		// Act
		result, err := service.CreatePayment("invalid-order")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get order")
	})

	t.Run("Error - Midtrans API error", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(nil, &midtrans.Error{Message: "midtrans error"})

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Midtrans returns nil response", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return((*snap.Response)(nil), (*midtrans.Error)(nil))

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Error - Failed to save payment", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(&snap.Response{Token: "test-token", RedirectURL: "http://test.com"}, (*midtrans.Error)(nil))
		mockPaymentRepo.EXPECT().CreatePayment(gomock.Any()).Return(errors.New("database error"))

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create payment")
	})
}

// Additional error handling tests for HandlePaymentNotification
func TestPaymentService_HandlePaymentNotification_ErrorCases(t *testing.T) {
	t.Run("Error - Invalid transaction ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "", // Empty transaction ID
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transaction_id")
	})

	t.Run("Error - Invalid order ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "", // Empty order ID
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order_id")
	})

	t.Run("Error - Payment not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment not found")
	})

	t.Run("Error - Unknown transaction status", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "unknown_status", // Unknown status
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown transaction status")
	})

	t.Run("Error - Database update failure", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(errors.New("database error"))

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update payment and order status")
	})

	t.Run("Error - Order status update failure", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(errors.New("database error"))

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update payment and order status")
	})
}

// Test GetPaymentStatus method
func TestPaymentService_GetPaymentStatus(t *testing.T) {
	t.Run("Success - Get payment status", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		existingPayment := createTestPayment()
		existingPayment.Status = entity.PaymentStatusSuccess
		existingPayment.TransactionID = "txn-123"
		existingPayment.PaymentMethod = entity.PaymentMethodCreditCard
		transactionTime := time.Now()
		existingPayment.TransactionTime = &transactionTime

		mockPaymentRepo.EXPECT().FindPaymentByID("payment-123").Return(existingPayment, nil)

		// Act
		result, err := service.GetPaymentStatus("payment-123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "payment-123", result.ID)
		assert.Equal(t, "order-123", result.OrderID)
		assert.Equal(t, entity.PaymentStatusSuccess, result.Status)
		assert.Equal(t, "txn-123", result.TransactionID)
		assert.Equal(t, string(entity.PaymentMethodCreditCard), result.PaymentMethod)
		assert.Equal(t, 250.0, result.Amount)
	})

	t.Run("Error - Payment not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		mockPaymentRepo.EXPECT().FindPaymentByID("invalid-payment").Return(nil, errors.New("payment not found"))

		// Act
		result, err := service.GetPaymentStatus("invalid-payment")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get payment")
	})
}

// Test HandlePaymentNotification method
func TestPaymentService_HandlePaymentNotification(t *testing.T) {
	t.Run("Success - Handle settlement notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
			GrossAmount:       "250.00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle pending notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "pending",
			PaymentType:       "bank_transfer",
			TransactionTime:   "2023-01-01 12:00:00",
			GrossAmount:       "250.00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "pending").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle failed notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "deny",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
			GrossAmount:       "250.00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "cancelled").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error - Missing transaction ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transaction_id")
	})

	t.Run("Error - Missing order ID", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order_id")
	})

	t.Run("Error - Payment not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "invalid-order",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("invalid-order").Return(nil, errors.New("payment not found"))

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment not found")
	})

	t.Run("Error - Unknown transaction status", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "unknown_status",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown transaction status")
	})

	t.Run("Success - Handle gopay notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "gopay",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle qris notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "qris",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle cstore notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "cstore",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error - Missing transaction status", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transaction_status")
	})

	t.Run("Error - Missing payment type", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid payment_type")
	})

	t.Run("Success - Handle notification with invalid transaction time", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
			TransactionTime:   "invalid-time-format",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err) // Should still succeed even with invalid time format
	})

	t.Run("Success - Handle refund notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "refund",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "refunded").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle capture notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "capture",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle expire notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "expire",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "cancelled").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle cancel notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "cancel",
			PaymentType:       "credit_card",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "cancelled").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Success - Handle shopeepay notification", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "shopeepay",
			TransactionTime:   "2023-01-01 12:00:00",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(nil)

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Error - Failed to update payment", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		payment := createTestPayment()
		notification := request.PaymentNotificationRequest{
			TransactionID:     "txn-123",
			OrderID:           "order-123",
			TransactionStatus: "settlement",
			PaymentType:       "credit_card",
		}

		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(payment, nil)
		mockPaymentRepo.EXPECT().UpdatePaymentAndOrderStatus(gomock.Any(), "order-123", "processing").Return(errors.New("database error"))

		// Act
		err := service.HandlePaymentNotification(notification)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update payment and order status")
	})
}

// Helper function to create test cart with items
func createTestCartWithItems() *entity.Cart {
	return &entity.Cart{
		ID: "cart-123",
		CartItems: []entity.CartItem{
			{
				ID:               "item-1",
				CartID:           "cart-123",
				ProductVariantID: "variant-1",
				Quantity:         2,
				ProductVariant: &entity.ProductVariant{
					ID:    "variant-1",
					Price: 100.0,
					Name:  "Test Product",
				},
			},
			{
				ID:               "item-2",
				CartID:           "cart-123",
				ProductVariantID: "variant-2",
				Quantity:         1,
				ProductVariant: &entity.ProductVariant{
					ID:    "variant-2",
					Price: 50.0,
					Name:  "Test Product 2",
				},
			},
		},
	}
}

// Helper function to create test order
func createTestOrder() *entity.Order {
	return &entity.Order{
		ID:                    "order-123",
		CartID:                "cart-123",
		CustomerName:          "John Doe",
		CustomerEmail:         "john@example.com",
		CustomerPhone:         "+1234567890",
		ShippingStreetAddress: "123 Test St",
		ShippingCity:          "Jakarta",
		ShippingProvince:      "DKI Jakarta",
		ShippingDistrict:      "Kebayoran Baru",
		ShippingPostalCode:    "12190",
		ShippingCourier:       "jne",
		ShippingService:       "REG",
		ShippingCountry:       "Indonesia",
		Subtotal:              250.0,
		DiscountAmount:        0.0,
		ShippingCost:          10000.0,
		TotalAmount:           250.0,
		Currency:              "IDR",
		OrderStatus:           "pending",
		SourceChannel:         "web",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// Helper function to create test payment
func createTestPayment() *entity.Payment {
	expiryTime := time.Now().Add(24 * time.Hour)
	return &entity.Payment{
		ID:           "payment-123",
		OrderID:      "order-123",
		Amount:       250.0,
		Status:       entity.PaymentStatusPending,
		PaymentToken: "test-token",
		PaymentURL:   "https://app.sandbox.midtrans.com/snap/v2/vtweb/test-token",
		ExpiryTime:   &expiryTime,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Test constructor functions
func TestNewPaymentService(t *testing.T) {
	t.Run("Success - Create payment service with dependencies", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)

		// Act
		service := NewPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient, "http://localhost:8080")

		// Assert
		assert.NotNil(t, service)
		assert.Equal(t, mockPaymentRepo, service.paymentRepo)
		assert.Equal(t, mockCartRepo, service.cartRepo)
		assert.Equal(t, mockSnapClient, service.snapClient)
	})
}

func TestNewPaymentServiceWithMidtrans(t *testing.T) {
	t.Run("Success - Create payment service with sandbox Midtrans", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockMailer := mocks.NewMockMailer(ctrl)
		mockWhatsApp := mocks.NewMockWhatsApp(ctrl)
		serverKey := "test-server-key"
		isProduction := false

		// Act
		service := NewPaymentServiceWithMidtrans(mockPaymentRepo, mockCartRepo, serverKey, isProduction, "http://localhost:8080", mockMailer, mockWhatsApp)

		// Assert
		assert.NotNil(t, service)
		assert.Equal(t, mockPaymentRepo, service.paymentRepo)
		assert.Equal(t, mockCartRepo, service.cartRepo)
		assert.NotNil(t, service.snapClient)
	})

	t.Run("Success - Create payment service with production Midtrans", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockMailer := mocks.NewMockMailer(ctrl)
		mockWhatsApp := mocks.NewMockWhatsApp(ctrl)
		serverKey := "test-server-key"
		isProduction := true

		// Act
		service := NewPaymentServiceWithMidtrans(mockPaymentRepo, mockCartRepo, serverKey, isProduction, "http://localhost:8080", mockMailer, mockWhatsApp)

		// Assert
		assert.NotNil(t, service)
		assert.Equal(t, mockPaymentRepo, service.paymentRepo)
		assert.Equal(t, mockCartRepo, service.cartRepo)
		assert.NotNil(t, service.snapClient)
	})
}

// Additional error handling tests
func TestPaymentService_CreateOrder_ErrorCases(t *testing.T) {
	t.Run("Error - Empty cart", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		emptyCart := &entity.Cart{
			ID:        "cart-123",
			CartItems: []entity.CartItem{}, // Empty cart
		}
		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(emptyCart, nil)

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cart is empty")
	})

	t.Run("Error - Cart not found", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(nil, errors.New("cart not found"))

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get cart")
	})

	t.Run("Error - Database failure during order creation", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		cart := createTestCartWithItems()
		req := request.CreateOrderRequest{
			CartID:               "cart-123",
			CustomerName:         "John Doe",
			CustomerEmail:        "john@example.com",
			CustomerPhone:        "+1234567890",
			ShippingAddress:      "123 Test St",
			ShippingCityName:     "Jakarta",
			ShippingProvinceName: "DKI Jakarta",
			ShippingDistrictName: "Kebayoran Baru",
			ShippingPostalCode:   "12190",
			ShippingCourier:      "jne",
			ShippingService:      "REG",
			ShippingCost:         10000,
			TotalWeight:          1000,
		}

		mockCartRepo.EXPECT().GetCartWithItems("cart-123").Return(cart, nil)
		mockPaymentRepo.EXPECT().GetSeq().Return(int64(1), nil)
		mockPaymentRepo.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		// Act
		result, err := service.CreateOrder(req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create order with items")
	})
}

func TestPaymentService_CreatePayment_ErrorCases(t *testing.T) {
	t.Run("Error - Midtrans transaction creation failure", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(nil, &midtrans.Error{
			Message:    "Transaction failed",
			StatusCode: 400,
		})

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create Midtrans transaction")
	})

	t.Run("Error - Midtrans response is nil", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(nil, nil)

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "response is nil")
	})

	t.Run("Error - Empty token in Midtrans response", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(&snap.Response{
			Token:       "",
			RedirectURL: "https://app.sandbox.midtrans.com/snap/v2/vtweb/test-token",
		}, nil)

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "token is empty")
	})

	t.Run("Error - Empty redirect URL in Midtrans response", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("payment not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(&snap.Response{
			Token:       "test-token",
			RedirectURL: "",
		}, nil)

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "redirect URL is empty")
	})

	t.Run("Error - Database failure during payment creation", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
		mockCartRepo := mocks.NewMockCartRepository(ctrl)
		mockSnapClient := mocks.NewMockSnapClientInterface(ctrl)
		service := createTestPaymentService(mockPaymentRepo, mockCartRepo, mockSnapClient)

		order := createTestOrder()
		snapResponse := &snap.Response{
			Token:       "test-token",
			RedirectURL: "https://app.sandbox.midtrans.com/snap/v2/vtweb/test-token",
		}

		mockPaymentRepo.EXPECT().FindOrderByID("order-123").Return(order, nil)
		mockPaymentRepo.EXPECT().FindPaymentByOrderID("order-123").Return(nil, errors.New("not found"))
		mockSnapClient.EXPECT().CreateTransaction(gomock.Any()).Return(snapResponse, nil)
		mockPaymentRepo.EXPECT().CreatePayment(gomock.Any()).Return(errors.New("database error"))

		// Act
		result, err := service.CreatePayment("order-123")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create payment")
	})
}
