package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/hanifbg/landing_backend/internal/model/static"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func (s *PaymentService) CreateOrder(req request.CreateOrderRequest) (*response.CreateOrderResponse, error) {
	// Get cart with items
	cart, err := s.cartRepo.GetCartWithItems(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	if len(cart.CartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate totals
	subtotal := 0.0
	for _, item := range cart.CartItems {
		subtotal += float64(item.Quantity) * item.ProductVariant.Price
	}

	// Check for discount
	var discountAmount float64 = 0
	var discountCode string

	//crete order number
	nextSeq, err := s.paymentRepo.GetSeq()
	if err != nil {
		return nil, fmt.Errorf("failed to get seq: %v", err)
	}
	orderNumber := fmt.Sprintf("IQB-%d-%05d", time.Now().Year(), nextSeq)

	// Create order
	orderID := uuid.New().String()
	order := &entity.Order{
		ID:                    orderID,
		OrderNumber:           orderNumber,
		CartID:                req.CartID,
		CustomerName:          req.CustomerName,
		CustomerEmail:         req.CustomerEmail,
		CustomerPhone:         req.CustomerPhone,
		ShippingStreetAddress: req.ShippingAddress,
		ShippingCity:          req.ShippingCityName,
		ShippingProvince:      req.ShippingProvinceName,
		ShippingDistrict:      req.ShippingDistrictName,
		ShippingPostalCode:    req.ShippingPostalCode,
		ShippingCourier:       req.ShippingCourier,
		ShippingService:       req.ShippingService,
		ShippingCountry:       "Indonesia",
		Subtotal:              subtotal,
		DiscountAmount:        discountAmount,
		DiscountCodeApplied:   discountCode,
		ShippingCost:          req.ShippingCost,
		TotalAmount:           subtotal - discountAmount + req.ShippingCost,
		Currency:              "IDR",
		OrderStatus:           "pending",
		SourceChannel:         "web",
		Notes:                 req.Notes,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Create order items
	orderItems := make([]entity.OrderItem, 0)
	for _, cartItem := range cart.CartItems {
		orderItem := entity.OrderItem{
			ID:               uuid.New().String(),
			OrderID:          orderID,
			ProductVariant:   cartItem.ProductVariant,
			ProductVariantID: cartItem.ProductVariantID,
			Quantity:         cartItem.Quantity,
			PriceAtPurchase:  cartItem.ProductVariant.Price,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		orderItems = append(orderItems, orderItem)
	}

	// Save order and order items in a single transaction
	if err := s.paymentRepo.CreateOrderWithItems(order, orderItems); err != nil {
		return nil, fmt.Errorf("failed to create order with items: %v", err)
	}

	// Prepare response
	itemResponses := make([]response.OrderItemResponse, 0)
	for _, item := range orderItems {
		itemResponses = append(itemResponses, response.OrderItemResponse{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			PriceAtPurchase:  item.PriceAtPurchase,
		})
	}

	// orderResponse := &response.OrderResponse{
	// 	ID:                  order.ID,
	// 	CartID:              order.CartID,
	// 	CustomerName:        order.CustomerName,
	// 	CustomerEmail:       order.CustomerEmail,
	// 	CustomerPhone:       order.CustomerPhone,
	// 	ShippingAddress:     order.ShippingStreetAddress,
	// 	Subtotal:            order.Subtotal,
	// 	DiscountAmount:      order.DiscountAmount,
	// 	DiscountCodeApplied: order.DiscountCodeApplied,
	// 	ShippingCost:        order.ShippingCost,
	// 	TotalAmount:         order.TotalAmount,
	// 	Currency:            order.Currency,
	// 	OrderStatus:         order.OrderStatus,
	// 	SourceChannel:       order.SourceChannel,
	// 	Notes:               order.Notes,
	// 	OrderItems:          itemResponses,
	// 	CreatedAt:           order.CreatedAt,
	// }

	orderResponse := &response.CreateOrderResponse{
		OrderID:     order.ID,
		OrderNumber: orderNumber,
		TotalAmount: order.TotalAmount,
	}

	//send email
	err = s.mailer.SendOrderConfirmation(order, orderItems)
	if err != nil {
		log.Printf("failed to send order confirmation email: %v", err)
	}
	//send whatsapp
	payload := request.WhatsAppRequest{
		CustomerName:          order.CustomerName,
		OrderNumber:           orderNumber,
		TotalAmount:           order.TotalAmount,
		OrderConfirmationLink: fmt.Sprintf("%s/order-confirmation/%s", s.baseURL, order.ID),
	}

	// Parse the template and execute it with the data
	tmpl, err := template.New("whatsappMessage").Parse(static.WAMessageTemplate)
	if err != nil {
		log.Printf("failed to parse WhatsApp message template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, payload); err != nil {
		log.Printf("failed to execute WhatsApp message template: %v", err)
	}

	// Call the function to send the message if WhatsApp repository is available
	//if s.whatsAppRepo != nil {
	// Format the phone number correctly for WhatsApp
	phoneNumber := s.formatPhoneNumberForWhatsApp(order.CustomerPhone)
	if err := s.whatsAppRepo.SendMessage(phoneNumber, buf.String()); err != nil {
		log.Printf("failed to send WhatsApp message: %v", err)
	}
	//}
	return orderResponse, nil
}

// formatPhoneNumberForWhatsApp converts a regular phone number to WhatsApp format
// It ensures the number starts with "+628" for Indonesia and appends "@s.whatsapp.net"
func (s *PaymentService) formatPhoneNumberForWhatsApp(phoneNumber string) string {
	// Remove any non-digit characters
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	// If the number starts with 0, replace it with 62 (Indonesia country code)
	if strings.HasPrefix(digitsOnly, "0") {
		digitsOnly = "62" + digitsOnly[1:]
	}

	// If the number doesn't start with 62, add it
	if !strings.HasPrefix(digitsOnly, "62") {
		digitsOnly = "62" + digitsOnly
	}

	// Add the plus sign and WhatsApp suffix
	return digitsOnly + "@s.whatsapp.net"
}

func (s *PaymentService) GetOrder(orderID string) (*response.OrderResponse, error) {
	var payment *response.PaymentResponse
	// Get order with items
	order, err := s.paymentRepo.GetOrderWithItems(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	// Prepare response
	itemResponses := make([]response.OrderItemResponse, 0)
	for _, item := range order.OrderItems {
		itemResponses = append(itemResponses, response.OrderItemResponse{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			ProductName:      item.ProductVariant.Name,
			ProductImage:     item.ProductVariant.ImageURL,
			Quantity:         item.Quantity,
			PriceAtPurchase:  item.PriceAtPurchase,
		})
	}

	existingPayment, err := s.paymentRepo.FindPaymentByOrderID(orderID)
	if err != nil {
		fmt.Errorf("failed to get payment: %v", err)
	}

	if existingPayment != nil {
		payment = &response.PaymentResponse{
			ID:            existingPayment.ID,
			OrderID:       existingPayment.OrderID,
			Amount:        existingPayment.Amount,
			Status:        existingPayment.Status,
			PaymentMethod: string(existingPayment.PaymentMethod),
			TransactionID: existingPayment.TransactionID,
			// PaymentToken:  existingPayment.PaymentToken,
			// PaymentURL:    existingPayment.PaymentURL,
			// ExpiryTime:    existingPayment.ExpiryTime,
			CreatedAt: existingPayment.CreatedAt,
		}
	}

	orderResponse := &response.OrderResponse{
		ID:                   order.ID,
		OrderNumber:          order.OrderNumber,
		CartID:               order.CartID,
		CustomerName:         order.CustomerName,
		CustomerEmail:        order.CustomerEmail,
		CustomerPhone:        order.CustomerPhone,
		ShippingAddress:      order.ShippingStreetAddress,
		ShippingCityName:     order.ShippingCity,
		ShippingProvinceName: order.ShippingProvince,
		ShippingDistrictName: order.ShippingDistrict,
		ShippingPostalCode:   order.ShippingPostalCode,
		Subtotal:             order.Subtotal,
		DiscountAmount:       order.DiscountAmount,
		DiscountCodeApplied:  order.DiscountCodeApplied,
		ShippingCost:         order.ShippingCost,
		TotalAmount:          order.TotalAmount,
		Currency:             order.Currency,
		OrderStatus:          order.OrderStatus,
		SourceChannel:        order.SourceChannel,
		Notes:                order.Notes,
		OrderItems:           itemResponses,
		CreatedAt:            order.CreatedAt,
		Payment:              payment,
	}

	return orderResponse, nil
}

func (s *PaymentService) CreatePayment(orderID string) (*response.PaymentResponse, error) {
	// Get order
	order, err := s.paymentRepo.FindOrderByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	// Check if payment already exists
	existingPayment, err := s.paymentRepo.FindPaymentByOrderID(orderID)
	if err == nil && existingPayment != nil {
		// Payment already exists, return it
		return &response.PaymentResponse{
			ID:            existingPayment.ID,
			OrderID:       existingPayment.OrderID,
			Amount:        existingPayment.Amount,
			Status:        existingPayment.Status,
			PaymentMethod: string(existingPayment.PaymentMethod),
			TransactionID: existingPayment.TransactionID,
			PaymentToken:  existingPayment.PaymentToken,
			PaymentURL:    existingPayment.PaymentURL,
			ExpiryTime:    existingPayment.ExpiryTime,
			CreatedAt:     existingPayment.CreatedAt,
		}, nil
	}

	// Create Midtrans transaction
	// Create Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(order.TotalAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.CustomerName,
			Email: order.CustomerEmail,
			Phone: order.CustomerPhone,
			ShipAddr: &midtrans.CustomerAddress{
				Address: order.ShippingStreetAddress,
			},
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeBankTransfer,
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeCreditCard,
		},
		Callbacks: &snap.Callbacks{
			Finish: s.baseURL + "/api/v1/payments/notification",
		},
	}

	// Create transaction
	respSnap, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		// Check if it's a real error (not a typed nil)
		// Convert to string to safely check if it's a meaningful error
		errorStr := fmt.Sprintf("%v", err)
		if errorStr != "<nil>" && errorStr != "" {
			return nil, fmt.Errorf("failed to create Midtrans transaction: %v", err)
		}
	}

	if respSnap == nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: response is nil")
	}

	if respSnap.Token == "" {
		return nil, fmt.Errorf("failed to create Midtrans transaction: token is empty")
	}
	if respSnap.RedirectURL == "" {
		return nil, fmt.Errorf("failed to create Midtrans transaction: redirect URL is empty")
	}

	// Set expiry time (24 hours from now)
	expiryTime := time.Now().Add(24 * time.Hour)

	// Create payment record
	payment := &entity.Payment{
		ID:             uuid.New().String(),
		OrderID:        orderID,
		Amount:         order.TotalAmount,
		Status:         entity.PaymentStatusPending,
		PaymentToken:   respSnap.Token,
		PaymentURL:     respSnap.RedirectURL,
		ExpiryTime:     &expiryTime,
		PaymentDetails: entity.JSONMap{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save payment to database
	if err := s.paymentRepo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// Return payment response
	return &response.PaymentResponse{
		ID:           payment.ID,
		OrderID:      payment.OrderID,
		Amount:       payment.Amount,
		Status:       payment.Status,
		PaymentToken: payment.PaymentToken,
		PaymentURL:   payment.PaymentURL,
		ExpiryTime:   payment.ExpiryTime,
		CreatedAt:    payment.CreatedAt,
	}, nil
}

func (s *PaymentService) GetPaymentStatus(paymentID string) (*response.PaymentStatusResponse, error) {
	// Get payment
	payment, err := s.paymentRepo.FindPaymentByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}

	// Return payment status response
	return &response.PaymentStatusResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Status:          payment.Status,
		TransactionID:   payment.TransactionID,
		TransactionTime: payment.TransactionTime,
		PaymentMethod:   string(payment.PaymentMethod),
		Amount:          payment.Amount,
		UpdatedAt:       payment.UpdatedAt,
	}, nil
}

func (s *PaymentService) HandlePaymentNotification(notification request.PaymentNotificationRequest) error {
	// Extract transaction ID and order ID
	transactionID := notification.TransactionID
	if transactionID == "" {
		return fmt.Errorf("invalid transaction_id")
	}

	orderID := notification.OrderID
	if orderID == "" {
		return fmt.Errorf("invalid order_id")
	}

	// Get payment by order ID
	payment, err := s.paymentRepo.FindPaymentByOrderID(orderID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	// Update transaction ID if not set
	if payment.TransactionID == "" {
		payment.TransactionID = transactionID
	}

	// Extract transaction status
	transactionStatus := notification.TransactionStatus
	if transactionStatus == "" {
		return fmt.Errorf("invalid transaction_status")
	}

	// Extract payment type
	paymentType := notification.PaymentType
	if paymentType == "" {
		return fmt.Errorf("invalid payment_type")
	}

	// Map payment type to PaymentMethod
	switch paymentType {
	case "credit_card":
		payment.PaymentMethod = entity.PaymentMethodCreditCard
	case "bank_transfer":
		payment.PaymentMethod = entity.PaymentMethodBankTransfer
	case "gopay", "shopeepay":
		payment.PaymentMethod = entity.PaymentMethodEWallet
	case "qris":
		payment.PaymentMethod = entity.PaymentMethodQRIS
	case "cstore":
		payment.PaymentMethod = entity.PaymentMethodRetailOutlet
	}

	// Extract transaction time
	transactionTimeStr := notification.TransactionTime
	if transactionTimeStr != "" {
		var transactionTime time.Time
		transactionTime, err = time.Parse("2006-01-02 15:04:05", transactionTimeStr)
		if err == nil {
			payment.TransactionTime = &transactionTime
		}
	}

	// Update payment status based on transaction status
	var paymentStatus entity.PaymentStatus
	var orderStatus string

	switch transactionStatus {
	case "capture", "settlement":
		paymentStatus = entity.PaymentStatusSuccess
		orderStatus = "processing"
	case "pending":
		paymentStatus = entity.PaymentStatusPending
		orderStatus = "pending"
	case "deny", "cancel", "expire":
		paymentStatus = entity.PaymentStatusFailed
		orderStatus = "cancelled"
	case "refund":
		paymentStatus = entity.PaymentStatusRefunded
		orderStatus = "refunded"
	default:
		return fmt.Errorf("unknown transaction status: %s", transactionStatus)
	}

	// Update payment status
	payment.Status = paymentStatus
	payment.UpdatedAt = time.Now()

	// Store full notification data in payment details
	paymentDetails := entity.JSONMap{}
	paymentDetailsBytes, err := json.Marshal(notification)
	if err == nil {
		if err := json.Unmarshal(paymentDetailsBytes, &paymentDetails); err == nil {
			payment.PaymentDetails = paymentDetails
		}
	}

	// Update payment and order status in a single transaction
	if err := s.paymentRepo.UpdatePaymentAndOrderStatus(payment, orderID, orderStatus); err != nil {
		return fmt.Errorf("failed to update payment and order status: %v", err)
	}

	return nil
}
