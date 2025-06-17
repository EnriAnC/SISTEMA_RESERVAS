package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	service *NotificationService
}

func NewNotificationHandler(service *NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// SendNotification handles sending notifications
// POST /notifications
func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	notification, err := h.service.SendNotification(req)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(notification); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetNotifications retrieves notifications for a user
// GET /notifications?user_id=123&limit=10&offset=0
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	notifications, total, err := h.service.GetUserNotifications(userID, limit, offset)
	if err != nil {
		log.Printf("Error retrieving notifications: %v", err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"notifications": notifications,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ListNotifications is an alias for GetNotifications for backward compatibility
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	h.GetNotifications(w, r)
}

// GetNotification retrieves a specific notification by ID
// GET /notifications/{id}
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := h.service.repo.GetByID(id)
	if err != nil {
		if err.Error() == "notification not found" {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving notification: %v", err)
		http.Error(w, "Failed to retrieve notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notification); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateStatus is an alias for UpdateNotificationStatus
func (h *NotificationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateNotificationStatus(w, r)
}

// UpdateNotificationStatus updates notification read status
// PUT /notifications/{id}/status
func (h *NotificationHandler) UpdateNotificationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	var req struct {
		IsRead bool `json:"is_read"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateNotificationStatus(id, req.IsRead)
	if err != nil {
		if err.Error() == "notification not found" {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating notification status: %v", err)
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetNotificationStats provides notification statistics
// GET /notifications/stats?user_id=123
func (h *NotificationHandler) GetNotificationStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	stats, err := h.service.GetNotificationStats(userID)
	if err != nil {
		log.Printf("Error retrieving notification stats: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// HealthCheck endpoint
// GET /health
func (h *NotificationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "notification-service",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health check response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
