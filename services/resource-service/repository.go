package main

import (
	"fmt"
	"strings"
	"sync"
)

// ResourceRepository defines the interface for resource data access
type ResourceRepository interface {
	Create(resource *Resource) error
	GetByID(id int) (*Resource, error)
	Update(resource *Resource) error
	Delete(id int) error
	List(query ListResourcesQuery, limit, offset int) ([]*Resource, error)
	GetAvailabilitySlots(resourceID int) ([]*AvailabilitySlot, error)
	CreateAvailabilitySlot(slot *AvailabilitySlot) error
	ClearAvailabilitySlots(resourceID int) error
}

// InMemoryResourceRepository is a simple in-memory implementation
// TODO: Replace with actual database implementation (PostgreSQL)
type InMemoryResourceRepository struct {
	resources   map[int]*Resource
	slots       map[int]*AvailabilitySlot
	nextResID   int
	nextSlotID  int
	mutex       sync.RWMutex
}

func NewResourceRepository() ResourceRepository {
	return &InMemoryResourceRepository{
		resources:  make(map[int]*Resource),
		slots:      make(map[int]*AvailabilitySlot),
		nextResID:  1,
		nextSlotID: 1,
	}
}

func (r *InMemoryResourceRepository) Create(resource *Resource) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	resource.ID = r.nextResID
	r.nextResID++
	
	r.resources[resource.ID] = resource
	return nil
}

func (r *InMemoryResourceRepository) GetByID(id int) (*Resource, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	resource, exists := r.resources[id]
	if !exists {
		return nil, fmt.Errorf("resource with ID %d not found", id)
	}
	
	if !resource.IsActive {
		return nil, fmt.Errorf("resource with ID %d is inactive", id)
	}
	
	return resource, nil
}

func (r *InMemoryResourceRepository) Update(resource *Resource) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.resources[resource.ID]; !exists {
		return fmt.Errorf("resource with ID %d not found", resource.ID)
	}
	
	r.resources[resource.ID] = resource
	return nil
}

func (r *InMemoryResourceRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	resource, exists := r.resources[id]
	if !exists {
		return fmt.Errorf("resource with ID %d not found", id)
	}
	
	resource.IsActive = false
	return nil
}

func (r *InMemoryResourceRepository) List(query ListResourcesQuery, limit, offset int) ([]*Resource, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var filtered []*Resource
	
	// Filter resources
	for _, resource := range r.resources {
		if !resource.IsActive {
			continue
		}
		
		// Apply filters
		if query.Type != "" && !strings.EqualFold(resource.Type, query.Type) {
			continue
		}
		
		if query.Location != "" && !strings.Contains(strings.ToLower(resource.Location), strings.ToLower(query.Location)) {
			continue
		}
		
		if query.Capacity > 0 && resource.Capacity < query.Capacity {
			continue
		}
		
		filtered = append(filtered, resource)
	}
	
	// Apply pagination
	start := offset
	if start >= len(filtered) {
		return []*Resource{}, nil
	}
	
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	
	return filtered[start:end], nil
}

func (r *InMemoryResourceRepository) GetAvailabilitySlots(resourceID int) ([]*AvailabilitySlot, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var slots []*AvailabilitySlot
	for _, slot := range r.slots {
		if slot.ResourceID == resourceID && slot.IsActive {
			slots = append(slots, slot)
		}
	}
	
	return slots, nil
}

func (r *InMemoryResourceRepository) CreateAvailabilitySlot(slot *AvailabilitySlot) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	slot.ID = r.nextSlotID
	r.nextSlotID++
	
	r.slots[slot.ID] = slot
	return nil
}

func (r *InMemoryResourceRepository) ClearAvailabilitySlots(resourceID int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for id, slot := range r.slots {
		if slot.ResourceID == resourceID {
			delete(r.slots, id)
		}
	}
	
	return nil
}

// TODO: PostgreSQL implementation
/*
type PostgreSQLResourceRepository struct {
	db *sql.DB
}

func NewPostgreSQLResourceRepository(db *sql.DB) ResourceRepository {
	return &PostgreSQLResourceRepository{db: db}
}

func (r *PostgreSQLResourceRepository) Create(resource *Resource) error {
	query := `
		INSERT INTO resources (name, type, description, capacity, location, properties, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	
	err := r.db.QueryRow(
		query,
		resource.Name, resource.Type, resource.Description, resource.Capacity,
		resource.Location, resource.Properties, resource.IsActive,
		resource.CreatedAt, resource.UpdatedAt,
	).Scan(&resource.ID)
	
	return err
}

// ... implement other methods
*/
