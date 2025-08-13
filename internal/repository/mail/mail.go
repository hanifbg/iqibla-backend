package mail

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"

	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"github.com/hanifbg/landing_backend/internal/model/request"
	"gopkg.in/gomail.v2"
)

func (m *Mailer) Send(from, to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	return m.Dialer.DialAndSend(mailer)
}

// SendOrderConfirmation builds the email body from template and sends it to the customer
func (m *Mailer) SendOrderConfirmation(order *entity.Order, items []entity.OrderItem) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}
	if order.CustomerEmail == "" {
		return fmt.Errorf("customer email is empty")
	}

	// Build items list
	itemsEmail := make([]request.OrderEmailItem, 0, len(order.OrderItems))
	for _, it := range items {
		name := ""
		if it.ProductVariant != nil {
			name = it.ProductVariant.Name
		}
		fmt.Println("ProductName:", name)
		fmt.Println("Quantity:", it.Quantity)
		fmt.Println("PriceAtPurchase:", it.PriceAtPurchase)
		itemsEmail = append(itemsEmail, request.OrderEmailItem{
			ProductName:     name,
			Quantity:        it.Quantity,
			PriceAtPurchase: it.PriceAtPurchase,
		})
	}

	// Build data for template
	data := request.OrderEmailData{
		CustomerName:          order.CustomerName,
		OrderNumber:           order.OrderNumber,
		OrderItems:            itemsEmail,
		SubtotalAmount:        order.Subtotal,
		ShippingCost:          order.ShippingCost,
		TotalAmount:           order.TotalAmount,
		OrderConfirmationLink: buildOrderLink(order.ID),
	}

	// Load and execute HTML template
	tplPath := resolveTemplatePath()
	tpl, err := htmltemplate.New("mail").ParseFiles(tplPath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var buf bytes.Buffer
	// Execute template by name (file base name)
	if err := tpl.ExecuteTemplate(&buf, filepath.Base(tplPath), data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Determine sender from config or env
	from := getSMTPFrom()
	subject := fmt.Sprintf("Order Confirmation #%s", order.OrderNumber)

	fmt.Println("from:", from)
	return m.Send(from, order.CustomerEmail, subject, buf.String())
}

// resolveTemplatePath returns absolute path to mail.html template
func resolveTemplatePath() string {
	// Prefer absolute path relative to project root when running in repo
	defaultPath := filepath.Join("internal", "model", "static", "mail.html")
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath
	}
	// Try with working dir adjustments
	wd, err := os.Getwd()
	if err == nil {
		p := filepath.Join(wd, "internal", "model", "static", "mail.html")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Fallback to just filename (ParseFiles can search rel path)
	return defaultPath
}

// buildOrderLink creates a link for the order confirmation page
func buildOrderLink(orderID string) string {
	cfg, err := config.GetConfig()
	base := ""
	if err == nil && cfg != nil && cfg.BaseURL != "" {
		base = cfg.BaseURL
	}
	if base == "" {
		base = "http://localhost:8080"
	}
	// Assuming frontend route to view order details
	return fmt.Sprintf("%s/order-confirmation/%s", base, orderID)
}

func getSMTPFrom() string {
	cfg, err := config.GetConfig()
	if err == nil && cfg != nil && cfg.SMTPFrom != "" {
		return cfg.SMTPFrom
	}
	// fallback to username if from is not set
	if err == nil && cfg != nil && cfg.SMTPUsername != "" {
		return cfg.SMTPUsername
	}
	return "no-reply@example.com"
}
