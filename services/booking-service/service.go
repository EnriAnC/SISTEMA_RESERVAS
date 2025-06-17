package main

import (
	"fmt"
	"time"
)

type BookingService struct {
	repository BookingRepository
	// TODO: Add HTTP clients for User and Resource services
}

func NewBookingService() *BookingService {
	return &BookingService{
		repository: NewBookingRepository(),
	}
}

// Create creates a new booking after validating availability
func (s *BookingService) Create(userID int, req CreateBookingRequest) (*Booking, error) {
	// Check resource availability (TODO: Call Resource Service)
	conflicts, err := s.repository.GetConflictingBookings(req.ResourceID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check conflicts: %w", err)
	}

	if len(conflicts) > 0 {
		return nil, fmt.Errorf("resource is not available for the selected time slot")
	}

	// Create booking
	booking := Booking{
		UserID:     userID,
		ResourceID: req.ResourceID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     BookingStatusPending,
		Notes:      req.Notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repository.Create(&booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// TODO: Publish booking created event
	s.publishEvent(BookingEventCreated, &booking)

	return &booking, nil
}

// GetByID retrieves a booking by ID
func (s *BookingService) GetByID(id int) (*BookingWithDetails, error) {
	booking, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	// TODO: Enrich with user and resource details from other services
	bookingWithDetails := &BookingWithDetails{
		Booking:      *booking,
		UserName:     "User Name",        // TODO: Get from User Service
		UserEmail:    "user@example.com", // TODO: Get from User Service
		ResourceName: "Resource Name",    // TODO: Get from Resource Service
		ResourceType: "room",             // TODO: Get from Resource Service
	}

	return bookingWithDetails, nil
}

// List retrieves bookings with filtering
func (s *BookingService) List(query ListBookingsQuery) ([]*BookingWithDetails, error) {
	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}

	offset := (query.Page - 1) * query.Size

	bookings, err := s.repository.List(query, query.Size, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}

	// TODO: Enrich with user and resource details
	var enrichedBookings []*BookingWithDetails
	for _, booking := range bookings {
		enrichedBookings = append(enrichedBookings, &BookingWithDetails{
			Booking:      *booking,
			UserName:     "User Name",
			UserEmail:    "user@example.com",
			ResourceName: "Resource Name",
			ResourceType: "room",
		})
	}

	return enrichedBookings, nil
}

// Update updates a booking
func (s *BookingService) Update(id int, req UpdateBookingRequest) (*Booking, error) {
	booking, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	if !booking.CanBeModified() {
		return nil, fmt.Errorf("booking cannot be modified in its current state: %s", booking.Status)
	}

	// Check for conflicts if time is being changed
	if req.StartTime != nil || req.EndTime != nil {
		startTime := booking.StartTime
		endTime := booking.EndTime

		if req.StartTime != nil {
			startTime = *req.StartTime
		}
		if req.EndTime != nil {
			endTime = *req.EndTime
		}

		conflicts, err := s.repository.GetConflictingBookings(booking.ResourceID, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to check conflicts: %w", err)
		}

		// Filter out the current booking from conflicts
		var actualConflicts []*Booking
		for _, conflict := range conflicts {
			if conflict.ID != booking.ID {
				actualConflicts = append(actualConflicts, conflict)
			}
		}

		if len(actualConflicts) > 0 {
			return nil, fmt.Errorf("resource is not available for the selected time slot")
		}

		booking.StartTime = startTime
		booking.EndTime = endTime
	}

	if req.Notes != nil {
		booking.Notes = *req.Notes
	}

	booking.UpdatedAt = time.Now()

	if err := s.repository.Update(booking); err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	// TODO: Publish booking updated event
	s.publishEvent(BookingEventUpdated, booking)

	return booking, nil
}

// Cancel cancels a booking
func (s *BookingService) Cancel(id int) error {
	booking, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if !booking.IsValidTransition(BookingStatusCanceled) {
		return fmt.Errorf("booking cannot be canceled in its current state: %s", booking.Status)
	}

	booking.Status = BookingStatusCanceled
	booking.UpdatedAt = time.Now()
	now := time.Now()
	booking.CanceledAt = &now

	if err := s.repository.Update(booking); err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// TODO: Publish booking canceled event
	s.publishEvent(BookingEventCanceled, booking)

	return nil
}

// Confirm confirms a booking
func (s *BookingService) Confirm(id int) (*Booking, error) {
	booking, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	if !booking.IsValidTransition(BookingStatusConfirmed) {
		return nil, fmt.Errorf("booking cannot be confirmed in its current state: %s", booking.Status)
	}

	booking.Status = BookingStatusConfirmed
	booking.UpdatedAt = time.Now()

	if err := s.repository.Update(booking); err != nil {
		return nil, fmt.Errorf("failed to confirm booking: %w", err)
	}

	// TODO: Publish booking confirmed event
	s.publishEvent(BookingEventConfirmed, booking)

	return booking, nil
}

// CheckAvailability checks if a resource is available for booking
func (s *BookingService) CheckAvailability(req AvailabilityCheckRequest) (*AvailabilityCheckResponse, error) {
	conflicts, err := s.repository.GetConflictingBookings(req.ResourceID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	response := &AvailabilityCheckResponse{
		Available: len(conflicts) == 0,
	}

	if len(conflicts) > 0 {
		for _, conflict := range conflicts {
			response.Conflicts = append(response.Conflicts, BookingConflict{
				ConflictingBookingID: conflict.ID,
				ConflictStartTime:    conflict.StartTime,
				ConflictEndTime:      conflict.EndTime,
				Message:              fmt.Sprintf("Booking #%d conflicts with requested time", conflict.ID),
			})
		}
	}

	return response, nil
}

// GetUpcomingBookings gets upcoming bookings for a user
func (s *BookingService) GetUpcomingBookings(userID int) ([]*Booking, error) {
	now := time.Now()
	query := ListBookingsQuery{
		UserID:    userID,
		StartDate: now,
		EndDate:   now.AddDate(0, 1, 0), // Next month
		Status:    BookingStatusConfirmed,
	}

	bookings, err := s.repository.List(query, 50, 0) // Get up to 50 upcoming bookings
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming bookings: %w", err)
	}

	return bookings, nil
}

// publishEvent publishes a booking event (stub implementation)
// TODO: Implement actual message publishing to RabbitMQ
func (s *BookingService) publishEvent(eventType BookingEventType, booking *Booking) {
	event := BookingEvent{
		Type:      string(eventType),
		BookingID: booking.ID,
		UserID:    booking.UserID,
		Timestamp: time.Now(),
		Data:      *booking,
	}

	// TODO: Publish to message queue
	fmt.Printf("Publishing event: %s for booking %d\n", event.Type, event.BookingID)
}
