package main

import (
	"fmt"
	"time"
)

type ResourceService struct {
	repository ResourceRepository
}

func NewResourceService() *ResourceService {
	return &ResourceService{
		repository: NewResourceRepository(),
	}
}

// Create creates a new resource
func (s *ResourceService) Create(req CreateResourceRequest) (*Resource, error) {
	resource := Resource{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    req.Location,
		Properties:  req.Properties,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := s.repository.Create(&resource); err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	
	return &resource, nil
}

// GetByID retrieves a resource by ID
func (s *ResourceService) GetByID(id int) (*Resource, error) {
	resource, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("resource not found: %w", err)
	}
	
	return resource, nil
}

// List retrieves resources with filtering
func (s *ResourceService) List(query ListResourcesQuery) ([]*Resource, error) {
	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	
	offset := (query.Page - 1) * query.Size
	
	resources, err := s.repository.List(query, query.Size, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	
	return resources, nil
}

// Update updates a resource
func (s *ResourceService) Update(id int, req UpdateResourceRequest) (*Resource, error) {
	resource, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("resource not found: %w", err)
	}
	
	// Update fields if provided
	if req.Name != nil {
		resource.Name = *req.Name
	}
	if req.Type != nil {
		resource.Type = *req.Type
	}
	if req.Description != nil {
		resource.Description = *req.Description
	}
	if req.Capacity != nil {
		resource.Capacity = *req.Capacity
	}
	if req.Location != nil {
		resource.Location = *req.Location
	}
	if req.Properties != nil {
		resource.Properties = req.Properties
	}
	if req.IsActive != nil {
		resource.IsActive = *req.IsActive
	}
	
	resource.UpdatedAt = time.Now()
	
	if err := s.repository.Update(resource); err != nil {
		return nil, fmt.Errorf("failed to update resource: %w", err)
	}
	
	return resource, nil
}

// Delete deletes a resource (soft delete by setting IsActive to false)
func (s *ResourceService) Delete(id int) error {
	resource, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("resource not found: %w", err)
	}
	
	resource.IsActive = false
	resource.UpdatedAt = time.Now()
	
	if err := s.repository.Update(resource); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	
	return nil
}

// GetAvailability checks resource availability for a date range
func (s *ResourceService) GetAvailability(resourceID int, startDate, endDate time.Time) ([]ResourceAvailability, error) {
	// Verify resource exists
	_, err := s.repository.GetByID(resourceID)
	if err != nil {
		return nil, fmt.Errorf("resource not found: %w", err)
	}
	
	// Get availability slots for the resource
	slots, err := s.repository.GetAvailabilitySlots(resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability slots: %w", err)
	}
	
	var availability []ResourceAvailability
	
	// Generate availability for each day in the range
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dayOfWeek := int(date.Weekday())
		
		// Find slots for this day of week
		for _, slot := range slots {
			if slot.DayOfWeek == dayOfWeek && slot.IsActive {
				// Parse time slots
				startTime, _ := time.Parse("15:04", slot.StartTime)
				endTime, _ := time.Parse("15:04", slot.EndTime)
				
				// Combine date with time
				fullStartTime := time.Date(date.Year(), date.Month(), date.Day(),
					startTime.Hour(), startTime.Minute(), 0, 0, date.Location())
				fullEndTime := time.Date(date.Year(), date.Month(), date.Day(),
					endTime.Hour(), endTime.Minute(), 0, 0, date.Location())
				
				availability = append(availability, ResourceAvailability{
					ResourceID: resourceID,
					Date:       date,
					StartTime:  fullStartTime,
					EndTime:    fullEndTime,
					IsBooked:   false, // TODO: Check against bookings
					BookingID:  nil,
				})
			}
		}
	}
	
	return availability, nil
}

// UpdateAvailability updates the availability slots for a resource
func (s *ResourceService) UpdateAvailability(resourceID int, slots []CreateAvailabilitySlotRequest) error {
	// Verify resource exists
	_, err := s.repository.GetByID(resourceID)
	if err != nil {
		return fmt.Errorf("resource not found: %w", err)
	}
	
	// Clear existing slots
	if err := s.repository.ClearAvailabilitySlots(resourceID); err != nil {
		return fmt.Errorf("failed to clear existing slots: %w", err)
	}
	
	// Add new slots
	for _, slotReq := range slots {
		slot := AvailabilitySlot{
			ResourceID: resourceID,
			DayOfWeek:  slotReq.DayOfWeek,
			StartTime:  slotReq.StartTime,
			EndTime:    slotReq.EndTime,
			IsActive:   true,
			CreatedAt:  time.Now(),
		}
		
		if err := s.repository.CreateAvailabilitySlot(&slot); err != nil {
			return fmt.Errorf("failed to create availability slot: %w", err)
		}
	}
	
	return nil
}

// CheckAvailability checks if a resource is available for a specific time period
func (s *ResourceService) CheckAvailability(resourceID int, startTime, endTime time.Time) (bool, error) {
	// Verify resource exists
	_, err := s.repository.GetByID(resourceID)
	if err != nil {
		return false, fmt.Errorf("resource not found: %w", err)
	}
	
	// Get availability for the date
	date := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	availability, err := s.GetAvailability(resourceID, date, date)
	if err != nil {
		return false, err
	}
	
	// Check if any availability slot covers the requested time
	for _, slot := range availability {
		if !slot.IsBooked &&
			(startTime.Equal(slot.StartTime) || startTime.After(slot.StartTime)) &&
			(endTime.Equal(slot.EndTime) || endTime.Before(slot.EndTime)) {
			return true, nil
		}
	}
	
	return false, nil
}
