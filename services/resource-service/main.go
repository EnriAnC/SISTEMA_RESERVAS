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
	// Initialize router
	r := mux.NewRouter()
	
	// Initialize handlers
	resourceHandler := NewResourceHandler()
	
	// Routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/resources", resourceHandler.CreateResource).Methods("POST")
	api.HandleFunc("/resources", resourceHandler.ListResources).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.GetResource).Methods("GET")
	api.HandleFunc("/resources/{id}", resourceHandler.UpdateResource).Methods("PUT")
	api.HandleFunc("/resources/{id}", resourceHandler.DeleteResource).Methods("DELETE")
	api.HandleFunc("/resources/{id}/availability", resourceHandler.GetAvailability).Methods("GET")
	api.HandleFunc("/resources/{id}/availability", resourceHandler.UpdateAvailability).Methods("PUT")
	
	// Server configuration
	server := &http.Server{
		Addr:         ":8002",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		log.Println("Resource Service starting on port 8002...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Resource Service failed to start: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down Resource Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Resource Service forced to shutdown: %v", err)
	}
	
	log.Println("Resource Service stopped")
}
