package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize repository and service
	repo := NewNotificationRepository()
	service := NewNotificationService(repo)

	// Initialize router
	r := mux.NewRouter()

	// Initialize handlers
	notificationHandler := NewNotificationHandler(service)

	// Routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/notifications", notificationHandler.SendNotification).Methods("POST")
	api.HandleFunc("/notifications", notificationHandler.GetNotifications).Methods("GET")
	api.HandleFunc("/notifications/{id}", notificationHandler.GetNotification).Methods("GET")
	api.HandleFunc("/notifications/{id}/status", notificationHandler.UpdateNotificationStatus).Methods("PUT")
	api.HandleFunc("/notifications/stats", notificationHandler.GetNotificationStats).Methods("GET")

	// Health check endpoint
	r.HandleFunc("/health", notificationHandler.HealthCheck).Methods("GET")

	// Start message consumer in background
	go func() {
		log.Println("Starting notification message consumer...")
		if err := service.StartMessageConsumer(); err != nil {
			log.Printf("Message consumer error: %v", err)
		}
	}()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Server configuration
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Notification Service starting on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Notification Service failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Notification Service...")

	// Stop message consumer
	service.StopMessageConsumer()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Notification Service forced to shutdown: %v", err)
		return
	}

	log.Println("Notification Service stopped")
}
