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

const (
	// DefaultPort is the default port for the user service
	DefaultPort = ":8001"
	// ShutdownTimeout is the timeout for graceful shutdown
	ShutdownTimeout = 30 * time.Second
	// ServerReadTimeout is the read timeout for HTTP server
	ServerReadTimeout = 15 * time.Second
	// ServerWriteTimeout is the write timeout for HTTP server
	ServerWriteTimeout = 15 * time.Second
	// ServerIdleTimeout is the idle timeout for HTTP server
	ServerIdleTimeout = 60 * time.Second
)

func main() {
	log.Println("Starting User Service...")

	// Initialize router
	r := mux.NewRouter()

	// Initialize handlers
	userHandler := NewUserHandler()

	// Setup routes
	setupRoutes(r, userHandler)

	// Create server with proper configuration
	server := &http.Server{
		Addr:         getPort(),
		Handler:      r,
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("User Service starting on port %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

// setupRoutes configures all the routes for the user service
func setupRoutes(r *mux.Router, userHandler *UserHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	// User CRUD operations
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Authentication endpoints
	api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", userHandler.RefreshToken).Methods("POST")

	// Health check endpoint
	api.HandleFunc("/health", healthCheck).Methods("GET")
}

// healthCheck provides a simple health check endpoint
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"healthy","service":"user-service"}`))
	if err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

// getPort returns the port from environment variable or default port
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return DefaultPort
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down User Service...")
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return
	}

	log.Println("User Service stopped")
}
