package main

import (
	"time"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCanceled  BookingStatus = "CANCELED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
)

// Booking represents a reservation
type Booking struct {
	ID         int           `json:"id" db:"id"`
	UserID     int           `json:"user_id" db:"user_id"`
	ResourceID int           `json:"resource_id" db:"resource_id"`
	StartTime  time.Time     `json:"start_time" db:"start_time"`
	EndTime    time.Time     `json:"end_time" db:"end_time"`
	Status     BookingStatus `json:"status" db:"status"`
	Notes      string        `json:"notes" db:"notes"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at"`
	CanceledAt *time.Time    `json:"canceled_at,omitempty" db:"canceled_at"`
}

// BookingWithDetails represents a booking with user and resource details
type BookingWithDetails struct {
	Booking
	UserName     string `json:"user_name"`
	UserEmail    string `json:"user_email"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
}

// CreateBookingRequest represents the request to create a booking
type CreateBookingRequest struct {
	ResourceID int       `json:"resource_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
	Notes      string    `json:"notes" validate:"max=500"`
}

// UpdateBookingRequest represents the request to update a booking
type UpdateBookingRequest struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Notes     *string    `json:"notes,omitempty" validate:"omitempty,max=500"`
}

// ListBookingsQuery represents query parameters for listing bookings
type ListBookingsQuery struct {
	UserID     int           `query:"user_id"`
	ResourceID int           `query:"resource_id"`
	Status     BookingStatus `query:"status"`
	StartDate  time.Time     `query:"start_date"`
	EndDate    time.Time     `query:"end_date"`
	Page       int           `query:"page"`
	Size       int           `query:"size"`
}

// BookingConflict represents a booking conflict
type BookingConflict struct {
	ConflictingBookingID int       `json:"conflicting_booking_id"`
	ConflictStartTime    time.Time `json:"conflict_start_time"`
	ConflictEndTime      time.Time `json:"conflict_end_time"`
	Message              string    `json:"message"`
}

// BookingEvent represents an event for the messaging system
type BookingEvent struct {
	Type      string    `json:"type"`
	BookingID int       `json:"booking_id"`
	UserID    int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      Booking   `json:"data"`
}

// BookingEventType defines the types of booking events
type BookingEventType string

const (
	BookingEventCreated   BookingEventType = "booking.created"
	BookingEventUpdated   BookingEventType = "booking.updated"
	BookingEventConfirmed BookingEventType = "booking.confirmed"
	BookingEventCanceled  BookingEventType = "booking.canceled"
)

// AvailabilityCheckRequest represents a request to check availability
type AvailabilityCheckRequest struct {
	ResourceID int       `json:"resource_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
}

// AvailabilityCheckResponse represents the response of availability check
type AvailabilityCheckResponse struct {
	Available bool              `json:"available"`
	Conflicts []BookingConflict `json:"conflicts,omitempty"`
}

// IsValidTransition checks if a status transition is valid
func (b *Booking) IsValidTransition(newStatus BookingStatus) bool {
	switch b.Status {
	case BookingStatusPending:
		return newStatus == BookingStatusConfirmed || newStatus == BookingStatusCanceled
	case BookingStatusConfirmed:
		return newStatus == BookingStatusCanceled || newStatus == BookingStatusCompleted
	case BookingStatusCanceled, BookingStatusCompleted:
		return false // Terminal states
	default:
		return false
	}
}

// CanBeModified checks if a booking can be modified
func (b *Booking) CanBeModified() bool {
	return b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed
}

// Duration returns the duration of the booking
func (b *Booking) Duration() time.Duration {
	return b.EndTime.Sub(b.StartTime)
}

// IsActive checks if the booking is currently active
func (b *Booking) IsActive(now time.Time) bool {
	return b.Status == BookingStatusConfirmed &&
		now.After(b.StartTime) &&
		now.Before(b.EndTime)
}

// IsUpcoming checks if the booking is upcoming
func (b *Booking) IsUpcoming(now time.Time) bool {
	return (b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed) &&
		now.Before(b.StartTime)
}
