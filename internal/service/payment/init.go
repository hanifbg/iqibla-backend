package payment

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// SnapClientInterface defines the interface for Midtrans Snap client
type SnapClientInterface interface {
	CreateTransaction(req *snap.Request) (*snap.Response, *midtrans.Error)
}

type PaymentService struct {
	paymentRepo  repository.PaymentRepository
	cartRepo     repository.CartRepository
	snapClient   SnapClientInterface
	baseURL      string
	mailer       repository.Mailer
	whatsAppRepo repository.WhatsApp
}

// New creates a PaymentService following the same pattern as other services
func New(cfg *config.AppConfig, repo *util.RepoWrapper) *PaymentService {
	return NewPaymentServiceWithMidtrans(repo.PaymentRepo, repo.CartRepo,
		cfg.MidtransServerKey, cfg.IsProduction, cfg.BaseURL,
		repo.MailRepo, repo.WhatsAppRepo)
}

func NewPaymentService(paymentRepo repository.PaymentRepository, cartRepo repository.CartRepository, snapClient SnapClientInterface, baseURL string) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		cartRepo:    cartRepo,
		snapClient:  snapClient,
		baseURL:     baseURL,
	}
}

// NewPaymentServiceWithMidtrans creates a PaymentService with a real Midtrans client
func NewPaymentServiceWithMidtrans(paymentRepo repository.PaymentRepository, cartRepo repository.CartRepository,
	midtransServerKey string, isProduction bool, baseURL string,
	mailer repository.Mailer, whatsapp repository.WhatsApp) *PaymentService {
	// Initialize Midtrans client
	// Set environment based on isProduction flag
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	// Correct way to initialize Snap client according to Midtrans Go SDK documentation
	var snapClient snap.Client
	snapClient.New(midtransServerKey, midtrans.EnvironmentType(env))

	return &PaymentService{
		paymentRepo:  paymentRepo,
		cartRepo:     cartRepo,
		snapClient:   &snapClient,
		baseURL:      baseURL,
		mailer:       mailer,
		whatsAppRepo: whatsapp,
	}
}
