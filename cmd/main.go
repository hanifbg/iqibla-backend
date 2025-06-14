package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hanifbg/landing_backend/config"
	_ "github.com/hanifbg/landing_backend/docs"
	handlerInit "github.com/hanifbg/landing_backend/internal/handler/util"
	repoInit "github.com/hanifbg/landing_backend/internal/repository/util"
	servInit "github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

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
	}

	// Initialize Echo
	e := echo.New()
	e.Logger.SetLevel(4) // INFO level

	// Initialize handlers
	handlerInit.InitHandler(cfg, e, serv)

	// Start server
	serverAddr := "localhost:8080"
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
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
