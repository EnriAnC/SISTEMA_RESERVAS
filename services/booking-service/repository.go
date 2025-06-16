package main

import (
	"fmt"
	"sync"
	"time"
)

// BookingRepository defines the interface for booking data access
type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id int) (*Booking, error)
	Update(booking *Booking) error
	Delete(id int) error
	List(query ListBookingsQuery, limit, offset int) ([]*Booking, error)
	GetConflictingBookings(resourceID int, startTime, endTime time.Time) ([]*Booking, error)
	GetByUserID(userID int, limit, offset int) ([]*Booking, error)
	GetByResourceID(resourceID int, limit, offset int) ([]*Booking, error)
}

// InMemoryBookingRepository is a simple in-memory implementation
// TODO: Replace with actual database implementation (PostgreSQL)
type InMemoryBookingRepository struct {
	bookings map[int]*Booking
	nextID   int
	mutex    sync.RWMutex
}

func NewBookingRepository() BookingRepository {
	return &InMemoryBookingRepository{
		bookings: make(map[int]*Booking),
		nextID:   1,
	}
}

func (r *InMemoryBookingRepository) Create(booking *Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	booking.ID = r.nextID
	r.nextID++
	
	r.bookings[booking.ID] = booking
	return nil
}

func (r *InMemoryBookingRepository) GetByID(id int) (*Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	booking, exists := r.bookings[id]
	if !exists {
		return nil, fmt.Errorf("booking with ID %d not found", id)
	}
	
	return booking, nil
}

func (r *InMemoryBookingRepository) Update(booking *Booking) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.bookings[booking.ID]; !exists {
		return fmt.Errorf("booking with ID %d not found", booking.ID)
	}
	
	r.bookings[booking.ID] = booking
	return nil
}

func (r *InMemoryBookingRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.bookings[id]; !exists {
		return fmt.Errorf("booking with ID %d not found", id)
	}
	
	delete(r.bookings, id)
	return nil
}

func (r *InMemoryBookingRepository) List(query ListBookingsQuery, limit, offset int) ([]*Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var filtered []*Booking
	
	// Filter bookings
	for _, booking := range r.bookings {
		// Apply filters
		if query.UserID > 0 && booking.UserID != query.UserID {
			continue
		}
		
		if query.ResourceID > 0 && booking.ResourceID != query.ResourceID {
			continue
		}
		
		if query.Status != "" && booking.Status != query.Status {
			continue
		}
		
		if !query.StartDate.IsZero() && booking.StartTime.Before(query.StartDate) {
			continue
		}
		
		if !query.EndDate.IsZero() && booking.EndTime.After(query.EndDate.AddDate(0, 0, 1)) {
			continue
		}
		
		filtered = append(filtered, booking)
	}
	
	// Apply pagination
	start := offset
	if start >= len(filtered) {
		return []*Booking{}, nil
	}
	
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	
	return filtered[start:end], nil
}

func (r *InMemoryBookingRepository) GetConflictingBookings(resourceID int, startTime, endTime time.Time) ([]*Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var conflicts []*Booking
	
	for _, booking := range r.bookings {
		// Skip cancelled bookings
		if booking.Status == BookingStatusCancelled {
			continue
		}
		
		// Check if booking is for the same resource
		if booking.ResourceID != resourceID {
			continue
		}
		
		// Check for time overlap
		if r.timeOverlaps(booking.StartTime, booking.EndTime, startTime, endTime) {
			conflicts = append(conflicts, booking)
		}
	}
	
	return conflicts, nil
}

func (r *InMemoryBookingRepository) GetByUserID(userID int, limit, offset int) ([]*Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var userBookings []*Booking
	count := 0
	skipped := 0
	
	for _, booking := range r.bookings {
		if booking.UserID != userID {
			continue
		}
		
		if skipped < offset {
			skipped++
			continue
		}
		
		if count >= limit {
			break
		}
		
		userBookings = append(userBookings, booking)
		count++
	}
	
	return userBookings, nil
}

func (r *InMemoryBookingRepository) GetByResourceID(resourceID int, limit, offset int) ([]*Booking, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var resourceBookings []*Booking
	count := 0
	skipped := 0
	
	for _, booking := range r.bookings {
		if booking.ResourceID != resourceID {
			continue
		}
		
		if skipped < offset {
			skipped++
			continue
		}
		
		if count >= limit {
			break
		}
		
		resourceBookings = append(resourceBookings, booking)
		count++
	}
	
	return resourceBookings, nil
}

// timeOverlaps checks if two time periods overlap
func (r *InMemoryBookingRepository) timeOverlaps(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

// TODO: PostgreSQL implementation
/*
type PostgreSQLBookingRepository struct {
	db *sql.DB
}

func NewPostgreSQLBookingRepository(db *sql.DB) BookingRepository {
	return &PostgreSQLBookingRepository{db: db}
}

func (r *PostgreSQLBookingRepository) Create(booking *Booking) error {
	query := `
		INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`
	
	err := r.db.QueryRow(
		query,
		booking.UserID, booking.ResourceID, booking.StartTime, booking.EndTime,
		booking.Status, booking.Notes, booking.CreatedAt, booking.UpdatedAt,
	).Scan(&booking.ID)
	
	return err
}

func (r *PostgreSQLBookingRepository) GetConflictingBookings(resourceID int, startTime, endTime time.Time) ([]*Booking, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, notes, created_at, updated_at
		FROM bookings 
		WHERE resource_id = $1 
		AND status != 'CANCELLED'
		AND start_time < $3 
		AND end_time > $2`
	
	rows, err := r.db.Query(query, resourceID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var bookings []*Booking
	for rows.Next() {
		booking := &Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.ResourceID,
			&booking.StartTime, &booking.EndTime, &booking.Status,
			&booking.Notes, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	
	return bookings, nil
}

// ... implement other methods
*/
