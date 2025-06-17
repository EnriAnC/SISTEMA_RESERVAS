package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type NotificationService struct {
	repo *NotificationRepository
}

func NewNotificationService(repo *NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// SendNotification processes and sends a notification
func (s *NotificationService) SendNotification(req NotificationRequest) (*Notification, error) {
	// Validate request
	if req.UserID == 0 {
		return nil, errors.New("user_id is required")
	}
	if req.Type == "" {
		return nil, errors.New("notification type is required")
	}
	if req.Message == "" {
		return nil, errors.New("message is required")
	}

	// Create notification
	notification := &Notification{
		ID:        s.repo.GetNextID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		Channel:   req.Channel,
		Priority:  req.Priority,
		Metadata:  req.Metadata,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	// Set default values
	if notification.Channel == "" {
		notification.Channel = ChannelEmail
	}
	if notification.Priority == "" {
		notification.Priority = PriorityNormal
	}

	// Save notification
	if err := s.repo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notification based on channel
	if err := s.sendToChannel(notification); err != nil {
		log.Printf("Failed to send notification via %s: %v", notification.Channel, err)
		// Mark as failed but don't return error - notification is still stored
		notification.Status = StatusFailed
		if updateErr := s.repo.Update(notification); updateErr != nil {
			log.Printf("Failed to update notification status: %v", updateErr)
			return nil, fmt.Errorf("failed to update notification status: %w", updateErr)
		}
	} else {
		notification.Status = StatusSent
		notification.SentAt = &notification.CreatedAt
		if updateErr := s.repo.Update(notification); updateErr != nil {
			log.Printf("Failed to update notification status: %v", updateErr)
			return nil, fmt.Errorf("failed to update notification status: %w", updateErr)
		}
	}

	return notification, nil
}

// GetUserNotifications retrieves notifications for a specific user
func (s *NotificationService) GetUserNotifications(userID, limit, offset int) ([]*Notification, int, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// UpdateNotificationStatus updates the read status of a notification
func (s *NotificationService) UpdateNotificationStatus(id int, isRead bool) error {
	notification, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if notification == nil {
		return errors.New("notification not found")
	}

	notification.IsRead = isRead
	if isRead && notification.ReadAt == nil {
		now := time.Now()
		notification.ReadAt = &now
	}

	return s.repo.Update(notification)
}

// GetNotificationStats provides statistics about user notifications
func (s *NotificationService) GetNotificationStats(userID int) (*NotificationStats, error) {
	notifications, _, err := s.repo.GetByUserID(userID, 1000, 0) // Get all for stats
	if err != nil {
		return nil, err
	}

	stats := &NotificationStats{
		UserID: userID,
		Total:  len(notifications),
	}

	for _, notification := range notifications {
		if notification.IsRead {
			stats.Read++
		} else {
			stats.Unread++
		}

		switch notification.Priority {
		case PriorityHigh:
			stats.High++
		case PriorityNormal:
			stats.Normal++
		case PriorityLow:
			stats.Low++
		}

		switch notification.Status {
		case StatusSent:
			stats.Sent++
		case StatusFailed:
			stats.Failed++
		default:
			stats.Pending++
		}
	}

	return stats, nil
}

// sendToChannel simulates sending notification via different channels
func (s *NotificationService) sendToChannel(notification *Notification) error {
	log.Printf("Sending notification %d via %s to user %d",
		notification.ID, notification.Channel, notification.UserID)

	switch notification.Channel {
	case ChannelEmail:
		return s.sendEmail(notification)
	case "sms":
		return s.sendSMS(notification)
	case "push":
		return s.sendPushNotification(notification)
	case "webhook":
		return s.sendWebhook(notification)
	default:
		return fmt.Errorf("unsupported channel: %s", notification.Channel)
	}
}

// sendEmail simulates email sending
func (s *NotificationService) sendEmail(notification *Notification) error {
	log.Printf("[EMAIL] To User %d: %s - %s",
		notification.UserID, notification.Title, notification.Message)

	// TODO: Integrate with email service (SendGrid, AWS SES, etc.)
	// For now, just simulate success
	time.Sleep(100 * time.Millisecond) // Simulate network delay
	return nil
}

// sendSMS simulates SMS sending
func (s *NotificationService) sendSMS(notification *Notification) error {
	log.Printf("[SMS] To User %d: %s",
		notification.UserID, notification.Message)

	// TODO: Integrate with SMS service (Twilio, AWS SNS, etc.)
	time.Sleep(200 * time.Millisecond) // Simulate network delay
	return nil
}

// sendPushNotification simulates push notification sending
func (s *NotificationService) sendPushNotification(notification *Notification) error {
	log.Printf("[PUSH] To User %d: %s - %s",
		notification.UserID, notification.Title, notification.Message)

	// TODO: Integrate with push service (Firebase FCM, Apple APNS, etc.)
	time.Sleep(150 * time.Millisecond) // Simulate network delay
	return nil
}

// sendWebhook simulates webhook sending
func (s *NotificationService) sendWebhook(notification *Notification) error {
	log.Printf("[WEBHOOK] To User %d: %s",
		notification.UserID, notification.Message)

	// TODO: Send HTTP POST to configured webhook URL
	payload := map[string]interface{}{
		"user_id":  notification.UserID,
		"type":     notification.Type,
		"title":    notification.Title,
		"message":  notification.Message,
		"metadata": notification.Metadata,
	}

	payloadJSON, _ := json.Marshal(payload)
	log.Printf("[WEBHOOK] Payload: %s", string(payloadJSON))

	time.Sleep(300 * time.Millisecond) // Simulate network delay
	return nil
}

// ProcessEvents handles incoming events from other services
func (s *NotificationService) ProcessEvents(event map[string]interface{}) error {
	eventType, ok := event["type"].(string)
	if !ok {
		return errors.New("event type is required")
	}

	userID, ok := event["user_id"].(float64)
	if !ok {
		return errors.New("user_id is required")
	}

	log.Printf("Processing event: %s for user %d", eventType, int(userID))

	// Create notification based on event type
	var req NotificationRequest
	req.UserID = int(userID)

	switch eventType {
	case "booking_confirmed":
		req.Type = "booking"
		req.Title = "Booking Confirmed"
		req.Message = "Your booking has been confirmed successfully"
		req.Priority = PriorityHigh
		req.Channel = ChannelEmail
	case "booking_canceled":
		req.Type = "booking"
		req.Title = "Booking Canceled"
		req.Message = "Your booking has been canceled"
		req.Priority = PriorityHigh
		req.Channel = ChannelEmail
	case "booking_reminder":
		req.Type = "reminder"
		req.Title = "Booking Reminder"
		req.Message = "Your booking is starting soon"
		req.Priority = PriorityNormal
		req.Channel = "push"
	case "user_registered":
		req.Type = "welcome"
		req.Title = "Welcome!"
		req.Message = "Welcome to our reservation system"
		req.Priority = PriorityNormal
		req.Channel = ChannelEmail
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	req.Metadata = event
	_, err := s.SendNotification(req)
	return err
}

// StartMessageConsumer starts the message consumer for processing events
func (s *NotificationService) StartMessageConsumer() error {
	log.Println("Message consumer started (mock implementation)")
	// TODO: Implement real message queue consumer (RabbitMQ, Kafka, etc.)
	return nil
}

// StopMessageConsumer stops the message consumer
func (s *NotificationService) StopMessageConsumer() {
	log.Println("Message consumer stopped")
	// TODO: Implement graceful shutdown of message consumer
}
