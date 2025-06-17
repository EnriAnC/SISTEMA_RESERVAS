package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type BookingHandler struct {
	bookingService *BookingService
}

func NewBookingHandler() *BookingHandler {
	return &BookingHandler{
		bookingService: NewBookingService(),
	}
}

// CreateBooking handles POST /api/v1/bookings
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user ID from JWT token
	userID := 1 // Placeholder

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	if req.StartTime.Before(time.Now()) {
		http.Error(w, "Cannot create booking in the past", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.Create(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
	}
}

// parseListBookingsQuery extracts and validates query parameters for listing bookings
func parseListBookingsQuery(r *http.Request) ListBookingsQuery {
	query := ListBookingsQuery{}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		if id, err := strconv.Atoi(userID); err == nil {
			query.UserID = id
		}
	}

	if resourceID := r.URL.Query().Get("resource_id"); resourceID != "" {
		if id, err := strconv.Atoi(resourceID); err == nil {
			query.ResourceID = id
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		query.Status = BookingStatus(status)
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			query.StartDate = date
		}
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			query.EndDate = date
		}
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = p
		}
	}

	if size := r.URL.Query().Get("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			query.Size = s
		}
	}

	return query
}

// ListBookings handles GET /api/v1/bookings
func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	query := parseListBookingsQuery(r)

	bookings, err := h.bookingService.List(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("Error encoding bookings response: %v", err)
	}
}

// GetBooking handles GET /api/v1/bookings/{id}
func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}
	booking, err := h.bookingService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
	}
}

// UpdateBooking handles PUT /api/v1/bookings/{id}
func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var req UpdateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.StartTime != nil && req.EndTime != nil {
		if req.StartTime.After(*req.EndTime) || req.StartTime.Equal(*req.EndTime) {
			http.Error(w, "End time must be after start time", http.StatusBadRequest)
			return
		}
	}
	booking, err := h.bookingService.Update(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
	}
}

// CancelBooking handles DELETE /api/v1/bookings/{id}
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.Cancel(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ConfirmBooking handles POST /api/v1/bookings/{id}/confirm
func (h *BookingHandler) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.Confirm(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
	}
}

// GetUserBookings handles GET /api/v1/users/{user_id}/bookings
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := ListBookingsQuery{
		UserID: userID,
	}

	// Parse additional query parameters
	if status := r.URL.Query().Get("status"); status != "" {
		query.Status = BookingStatus(status)
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = p
		}
	}

	if size := r.URL.Query().Get("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			query.Size = s
		}
	}

	bookings, err := h.bookingService.List(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("Error encoding bookings response: %v", err)
	}
}

// CheckAvailability handles POST /api/v1/bookings/check-availability
func (h *BookingHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	var req AvailabilityCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.bookingService.CheckAvailability(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding availability response: %v", err)
	}
}
