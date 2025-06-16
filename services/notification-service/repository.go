package main

import (
	"errors"
	"sync"
)

type NotificationRepository struct {
	notifications map[int]*Notification
	mu           sync.RWMutex
	nextID       int
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		notifications: make(map[int]*Notification),
		nextID:       1,
	}
}

func (r *NotificationRepository) GetNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	id := r.nextID
	r.nextID++
	return id
}

func (r *NotificationRepository) Create(notification *Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.notifications[notification.ID] = notification
	return nil
}

func (r *NotificationRepository) GetByID(id int) (*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	notification, exists := r.notifications[id]
	if !exists {
		return nil, errors.New("notification not found")
	}
	
	return notification, nil
}

func (r *NotificationRepository) GetByUserID(userID, limit, offset int) ([]*Notification, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var userNotifications []*Notification
	
	// Filter notifications by user ID
	for _, notification := range r.notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	
	total := len(userNotifications)
	
	// Apply pagination
	start := offset
	if start > total {
		start = total
	}
	
	end := start + limit
	if end > total {
		end = total
	}
	
	if start >= total {
		return []*Notification{}, total, nil
	}
	
	// Sort by creation time (newest first)
	for i := 0; i < len(userNotifications)-1; i++ {
		for j := i + 1; j < len(userNotifications); j++ {
			if userNotifications[i].CreatedAt.Before(userNotifications[j].CreatedAt) {
				userNotifications[i], userNotifications[j] = userNotifications[j], userNotifications[i]
			}
		}
	}
	
	return userNotifications[start:end], total, nil
}

func (r *NotificationRepository) GetAll() ([]*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	notifications := make([]*Notification, 0, len(r.notifications))
	for _, notification := range r.notifications {
		notifications = append(notifications, notification)
	}
	
	return notifications, nil
}

func (r *NotificationRepository) Update(notification *Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.notifications[notification.ID]; !exists {
		return errors.New("notification not found")
	}
	
	r.notifications[notification.ID] = notification
	return nil
}

func (r *NotificationRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.notifications[id]; !exists {
		return errors.New("notification not found")
	}
	
	delete(r.notifications, id)
	return nil
}

func (r *NotificationRepository) GetByType(notificationType string) ([]*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var typeNotifications []*Notification
	
	for _, notification := range r.notifications {
		if notification.Type == notificationType {
			typeNotifications = append(typeNotifications, notification)
		}
	}
	
	return typeNotifications, nil
}

func (r *NotificationRepository) GetUnreadByUserID(userID int) ([]*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var unreadNotifications []*Notification
	
	for _, notification := range r.notifications {
		if notification.UserID == userID && !notification.IsRead {
			unreadNotifications = append(unreadNotifications, notification)
		}
	}
	
	// Sort by creation time (newest first)
	for i := 0; i < len(unreadNotifications)-1; i++ {
		for j := i + 1; j < len(unreadNotifications); j++ {
			if unreadNotifications[i].CreatedAt.Before(unreadNotifications[j].CreatedAt) {
				unreadNotifications[i], unreadNotifications[j] = unreadNotifications[j], unreadNotifications[i]
			}
		}
	}
	
	return unreadNotifications, nil
}

// TODO: Replace with PostgreSQL implementation
/*
Example PostgreSQL implementation:

func (r *NotificationRepository) Create(notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, channel, priority, metadata, is_read, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	
	err := r.db.QueryRow(query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Channel,
		notification.Priority,
		notification.Metadata,
		notification.IsRead,
		notification.Status,
		notification.CreatedAt,
	).Scan(&notification.ID)
	
	return err
}
*/
