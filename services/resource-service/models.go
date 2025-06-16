package main

import (
	"time"
)

// Resource represents a bookable resource (room, equipment, etc.)
type Resource struct {
	ID          int                    `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Type        string                 `json:"type" db:"type"` // "room", "equipment", "vehicle", etc.
	Description string                 `json:"description" db:"description"`
	Capacity    int                    `json:"capacity" db:"capacity"`
	Location    string                 `json:"location" db:"location"`
	Properties  map[string]interface{} `json:"properties" db:"properties"` // Flexible properties (JSON)
	IsActive    bool                   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// AvailabilitySlot represents time slots when a resource is available
type AvailabilitySlot struct {
	ID         int       `json:"id" db:"id"`
	ResourceID int       `json:"resource_id" db:"resource_id"`
	DayOfWeek  int       `json:"day_of_week" db:"day_of_week"` // 0=Sunday, 1=Monday, ..., 6=Saturday
	StartTime  string    `json:"start_time" db:"start_time"`   // HH:MM format
	EndTime    string    `json:"end_time" db:"end_time"`       // HH:MM format
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ResourceAvailability represents the availability status of a resource for a specific date/time
type ResourceAvailability struct {
	ResourceID int       `json:"resource_id"`
	Date       time.Time `json:"date"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	IsBooked   bool      `json:"is_booked"`
	BookingID  *int      `json:"booking_id,omitempty"`
}

// CreateResourceRequest represents the request to create a new resource
type CreateResourceRequest struct {
	Name        string                 `json:"name" validate:"required,min=2,max=100"`
	Type        string                 `json:"type" validate:"required,oneof=room equipment vehicle space"`
	Description string                 `json:"description" validate:"max=500"`
	Capacity    int                    `json:"capacity" validate:"required,min=1"`
	Location    string                 `json:"location" validate:"required,max=200"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// UpdateResourceRequest represents the request to update a resource
type UpdateResourceRequest struct {
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Type        *string                `json:"type,omitempty" validate:"omitempty,oneof=room equipment vehicle space"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Capacity    *int                   `json:"capacity,omitempty" validate:"omitempty,min=1"`
	Location    *string                `json:"location,omitempty" validate:"omitempty,max=200"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// CreateAvailabilitySlotRequest represents the request to create availability slot
type CreateAvailabilitySlotRequest struct {
	DayOfWeek int    `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

// ListResourcesQuery represents query parameters for listing resources
type ListResourcesQuery struct {
	Type     string `query:"type"`
	Location string `query:"location"`
	Capacity int    `query:"min_capacity"`
	Page     int    `query:"page"`
	Size     int    `query:"size"`
}

// AvailabilityQuery represents query parameters for checking availability
type AvailabilityQuery struct {
	StartDate time.Time `query:"start_date" validate:"required"`
	EndDate   time.Time `query:"end_date" validate:"required"`
	StartTime string    `query:"start_time"`
	EndTime   string    `query:"end_time"`
}
