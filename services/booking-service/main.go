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
	bookingHandler := NewBookingHandler()

	// Routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingHandler.ListBookings).Methods("GET")
	api.HandleFunc("/bookings/{id}", bookingHandler.GetBooking).Methods("GET")
	api.HandleFunc("/bookings/{id}", bookingHandler.UpdateBooking).Methods("PUT")
	api.HandleFunc("/bookings/{id}", bookingHandler.CancelBooking).Methods("DELETE")
	api.HandleFunc("/bookings/{id}/confirm", bookingHandler.ConfirmBooking).Methods("POST")
	api.HandleFunc("/users/{user_id}/bookings", bookingHandler.GetUserBookings).Methods("GET")

	// Server configuration
	server := &http.Server{
		Addr:         ":8003",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("Booking Service starting on port 8003...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Booking Service failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Booking Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Booking Service forced to shutdown: %v", err)
		return
	}

	log.Println("Booking Service stopped")
}
