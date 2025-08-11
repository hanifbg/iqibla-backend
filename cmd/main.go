package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hanifbg/landing_backend/config"
	_ "github.com/hanifbg/landing_backend/docs"
	handlerInit "github.com/hanifbg/landing_backend/internal/handler/util"
	repoInit "github.com/hanifbg/landing_backend/internal/repository/util"
	servInit "github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Initialize configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	repo, err := repoInit.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	serv, err := servInit.New(cfg, repo)
	if err != nil {
		fmt.Println("GOT ERROR serv Init", err)
		// Continue to see cache logs even if service init fails
	}

	// Initialize Echo
	e := echo.New()
	e.Logger.SetLevel(4) // INFO level

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Configure CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Ngrok-Skip-Browser-Warning"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	handlerInit.InitHandler(cfg, e, serv)

	// Start server
	serverAddr := "localhost:8081"
	if cfg.AppPort != 0 {
		serverAddr = fmt.Sprintf(":%d", cfg.AppPort)
	}
	go func() {
		fmt.Println("Starting server on", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			e.Logger.Fatal("shutting down the server")
		}
	}()

	log.Printf("Server is running at http://%s", serverAddr)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
