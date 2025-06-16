package main

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "EMAIL"
	NotificationTypeSMS   NotificationType = "SMS"
	NotificationTypePush  NotificationType = "PUSH"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "PENDING"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusFailed    NotificationStatus = "FAILED"
	NotificationStatusDelivered NotificationStatus = "DELIVERED"
)

// NotificationTemplate represents a notification template
type NotificationTemplate string

const (
	TemplateBookingCreated   NotificationTemplate = "booking_created"
	TemplateBookingConfirmed NotificationTemplate = "booking_confirmed"
	TemplateBookingCancelled NotificationTemplate = "booking_cancelled"
	TemplateBookingReminder  NotificationTemplate = "booking_reminder"
	TemplateUserWelcome      NotificationTemplate = "user_welcome"
)

// Notification represents a notification record
type Notification struct {
	ID          int                    `json:"id" db:"id"`
	UserID      int                    `json:"user_id" db:"user_id"`
	Type        string                 `json:"type" db:"type"`
	Title       string                 `json:"title" db:"title"`
	Message     string                 `json:"message" db:"message"`
	Channel     string                 `json:"channel" db:"channel"` // email, sms, push, webhook
	Priority    string                 `json:"priority" db:"priority"` // high, normal, low
	Status      string                 `json:"status" db:"status"` // pending, sent, failed
	IsRead      bool                   `json:"is_read" db:"is_read"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"` // Additional data (JSON)
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	SentAt      *time.Time             `json:"sent_at,omitempty" db:"sent_at"`
	ReadAt      *time.Time             `json:"read_at,omitempty" db:"read_at"`
}

// NotificationRequest represents a request to send a notification
type NotificationRequest struct {
	UserID   int                    `json:"user_id" validate:"required"`
	Type     string                 `json:"type" validate:"required"`
	Title    string                 `json:"title" validate:"required,max=200"`
	Message  string                 `json:"message" validate:"required,max=5000"`
	Channel  string                 `json:"channel,omitempty"` // email, sms, push, webhook
	Priority string                 `json:"priority,omitempty"` // high, normal, low
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationStats represents notification statistics for a user
type NotificationStats struct {
	UserID  int `json:"user_id"`
	Total   int `json:"total"`
	Read    int `json:"read"`
	Unread  int `json:"unread"`
	High    int `json:"high"`
	Normal  int `json:"normal"`
	Low     int `json:"low"`
	Sent    int `json:"sent"`
	Failed  int `json:"failed"`
	Pending int `json:"pending"`
}

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	Type      NotificationType   `json:"type" validate:"required,oneof=EMAIL SMS PUSH"`
	Template  NotificationTemplate `json:"template" validate:"required"`
	Recipient string             `json:"recipient" validate:"required"`
	Subject   string             `json:"subject" validate:"required,max=200"`
	Body      string             `json:"body" validate:"required,max=5000"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateStatusRequest represents a request to update notification status
type UpdateStatusRequest struct {
	Status        NotificationStatus `json:"status" validate:"required,oneof=PENDING SENT FAILED DELIVERED"`
	FailureReason *string            `json:"failure_reason,omitempty"`
}

// ListNotificationsQuery represents query parameters for listing notifications
type ListNotificationsQuery struct {
	Type      NotificationType   `query:"type"`
	Status    NotificationStatus `query:"status"`
	Recipient string             `query:"recipient"`
	Template  NotificationTemplate `query:"template"`
	StartDate time.Time          `query:"start_date"`
	EndDate   time.Time          `query:"end_date"`
	Page      int                `query:"page"`
	Size      int                `query:"size"`
}

// BookingEventMessage represents a booking event from the message queue
type BookingEventMessage struct {
	Type      string    `json:"type"`
	BookingID int       `json:"booking_id"`
	UserID    int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		ID         int       `json:"id"`
		UserID     int       `json:"user_id"`
		ResourceID int       `json:"resource_id"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
		Status     string    `json:"status"`
		Notes      string    `json:"notes"`
	} `json:"data"`
}

// UserEventMessage represents a user event from the message queue
type UserEventMessage struct {
	Type      string    `json:"type"`
	UserID    int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"data"`
}

// EmailProvider represents an email service provider
type EmailProvider interface {
	SendEmail(to, subject, body string) error
}

// SMSProvider represents an SMS service provider
type SMSProvider interface {
	SendSMS(to, message string) error
}

// PushProvider represents a push notification provider
type PushProvider interface {
	SendPush(userID string, title, message string) error
}

// TemplateData represents data for template rendering
type TemplateData struct {
	UserName     string    `json:"user_name"`
	UserEmail    string    `json:"user_email"`
	ResourceName string    `json:"resource_name"`
	BookingID    int       `json:"booking_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Notes        string    `json:"notes"`
}

// NotificationBatch represents a batch of notifications to be sent
type NotificationBatch struct {
	ID            int             `json:"id"`
	Notifications []Notification  `json:"notifications"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
}

// IsRetryable checks if a notification can be retried
func (n *Notification) IsRetryable() bool {
	return n.Status == "failed" && n.GetRetryCount() < 3
}

// CanBeSent checks if a notification can be sent
func (n *Notification) CanBeSent() bool {
	return n.Status == "pending" || n.IsRetryable()
}

// GetRetryCount returns retry count from metadata
func (n *Notification) GetRetryCount() int {
	if n.Metadata == nil {
		return 0
	}
	if count, ok := n.Metadata["retry_count"].(float64); ok {
		return int(count)
	}
	return 0
}
