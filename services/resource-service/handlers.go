package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ResourceHandler struct {
	resourceService *ResourceService
}

func NewResourceHandler() *ResourceHandler {
	return &ResourceHandler{
		resourceService: NewResourceService(),
	}
}

// CreateResource handles POST /api/v1/resources
func (h *ResourceHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.Create(req)
	if err != nil {
		log.Printf("Error creating resource: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListResources handles GET /api/v1/resources
func (h *ResourceHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	query := h.parseListResourcesQuery(r)

	resources, err := h.resourceService.List(query)
	if err != nil {
		log.Printf("Error listing resources: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resources); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// parseListResourcesQuery parses query parameters for listing resources
func (h *ResourceHandler) parseListResourcesQuery(r *http.Request) ListResourcesQuery {
	query := ListResourcesQuery{}

	if resourceType := r.URL.Query().Get("type"); resourceType != "" {
		query.Type = resourceType
	}
	if location := r.URL.Query().Get("location"); location != "" {
		query.Location = location
	}
	if capacity := r.URL.Query().Get("min_capacity"); capacity != "" {
		if minCapacity, err := strconv.Atoi(capacity); err == nil {
			query.Capacity = minCapacity
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

// GetResource handles GET /api/v1/resources/{id}
func (h *ResourceHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid resource ID: %v", err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.GetByID(id)
	if err != nil {
		log.Printf("Error getting resource: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateResource handles PUT /api/v1/resources/{id}
func (h *ResourceHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid resource ID: %v", err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	var req UpdateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.Update(id, req)
	if err != nil {
		log.Printf("Error updating resource: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteResource handles DELETE /api/v1/resources/{id}
func (h *ResourceHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid resource ID: %v", err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	if err := h.resourceService.Delete(id); err != nil {
		log.Printf("Error deleting resource: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAvailability handles GET /api/v1/resources/{id}/availability
func (h *ResourceHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid resource ID: %v", err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	startDate, endDate, err := h.parseDateRange(r)
	if err != nil {
		log.Printf("Error parsing date range: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	availability, err := h.resourceService.GetAvailability(id, startDate, endDate)
	if err != nil {
		log.Printf("Error getting availability: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(availability); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// parseDateRange parses start_date and end_date query parameters
func (h *ResourceHandler) parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format (YYYY-MM-DD)")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format (YYYY-MM-DD)")
	}

	return startDate, endDate, nil
}

// UpdateAvailability handles PUT /api/v1/resources/{id}/availability
func (h *ResourceHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid resource ID: %v", err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	var slots []CreateAvailabilitySlotRequest
	if err := json.NewDecoder(r.Body).Decode(&slots); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.resourceService.UpdateAvailability(id, slots); err != nil {
		log.Printf("Error updating availability: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
