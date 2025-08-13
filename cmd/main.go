package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-context-engineering-template/config"
	httpDelivery "backend-context-engineering-template/internal/delivery/http"
	"backend-context-engineering-template/internal/delivery/http/handlers"
	"backend-context-engineering-template/internal/repository/postgres"
	"backend-context-engineering-template/internal/usecase"
	"backend-context-engineering-template/pkg/database"
	"backend-context-engineering-template/pkg/logger"
)

func main() {
	cfg := config.Load()

	appLogger := logger.New(cfg.Log.Level)
	appLogger.Info("Starting application...")

	dbConfig := database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig, appLogger)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			appLogger.WithError(err).Error("Failed to close database connection")
		}
	}()

	productRepo := postgres.NewProductRepository(db, appLogger)
	productUseCase := usecase.NewProductUseCase(productRepo, appLogger)
	productHandler := handlers.NewProductHandler(productUseCase, appLogger)

	router := httpDelivery.SetupRouter(productHandler, appLogger)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Addr, cfg.HTTP.Port),
		Handler: router,
	}

	go func() {
		appLogger.WithField("addr", server.Addr).Info("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.WithError(err).Fatal("Server forced to shutdown")
	}

	appLogger.Info("Server exited")
}
