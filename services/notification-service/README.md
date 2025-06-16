# Notification Service

The Notification Service handles all notification delivery across multiple channels including email, SMS, push notifications, and webhooks. It processes events from other services and manages notification preferences and delivery status.

## Features

- **Multi-channel delivery**: Email, SMS, push notifications, webhooks
- **Event-driven architecture**: Processes events from other microservices
- **Notification history**: Tracks all sent notifications with read status
- **Priority handling**: High, normal, low priority notifications
- **Delivery status tracking**: Pending, sent, failed status
- **User notification preferences**: Channel preferences and settings
- **Statistics and analytics**: Notification delivery metrics

## API Endpoints

### Send Notification
```http
POST /notifications
Content-Type: application/json

{
  "user_id": 123,
  "type": "booking",
  "title": "Booking Confirmed",
  "message": "Your booking has been confirmed",
  "channel": "email",
  "priority": "high",
  "metadata": {
    "booking_id": 456,
    "resource_name": "Conference Room A"
  }
}
```

### Get User Notifications
```http
GET /notifications?user_id=123&limit=10&offset=0
```

Response:
```json
{
  "notifications": [
    {
      "id": 1,
      "user_id": 123,
      "type": "booking",
      "title": "Booking Confirmed",
      "message": "Your booking has been confirmed",
      "channel": "email",
      "priority": "high",
      "is_read": false,
      "status": "sent",
      "created_at": "2024-01-15T10:30:00Z",
      "sent_at": "2024-01-15T10:30:05Z"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

### Update Notification Status
```http
PUT /notifications/{id}/status
Content-Type: application/json

{
  "is_read": true
}
```

### Get Notification Statistics
```http
GET /notifications/stats?user_id=123
```

Response:
```json
{
  "user_id": 123,
  "total": 50,
  "read": 30,
  "unread": 20,
  "high": 10,
  "normal": 35,
  "low": 5,
  "sent": 45,
  "failed": 2,
  "pending": 3
}
```

### Health Check
```http
GET /health
```

## Notification Types

- **booking**: Booking confirmations, cancellations, modifications
- **reminder**: Upcoming booking reminders
- **welcome**: Welcome messages for new users
- **system**: System maintenance and updates
- **promotional**: Marketing and promotional content

## Notification Channels

- **email**: Email notifications (default)
- **sms**: SMS text messages
- **push**: Mobile push notifications
- **webhook**: HTTP webhook callbacks

## Priority Levels

- **high**: Urgent notifications (booking confirmations, cancellations)
- **normal**: Standard notifications (reminders, updates)
- **low**: Non-urgent notifications (promotional content)

## Event Processing

The service listens for events from other microservices:

### Booking Events
- `booking_confirmed`: Sent when a booking is confirmed
- `booking_cancelled`: Sent when a booking is cancelled
- `booking_reminder`: Sent before booking start time

### User Events
- `user_registered`: Sent when a new user registers

## Configuration

Set the following environment variables:

```bash
# Server configuration
PORT=8084
ENV=development

# Email configuration (for production)
EMAIL_SERVICE_URL=https://api.sendgrid.com
EMAIL_API_KEY=your_sendgrid_api_key

# SMS configuration (for production)
SMS_SERVICE_URL=https://api.twilio.com
SMS_ACCOUNT_SID=your_twilio_sid
SMS_AUTH_TOKEN=your_twilio_token

# Push notification configuration (for production)
FCM_SERVER_KEY=your_firebase_server_key
APNS_KEY_ID=your_apple_key_id
APNS_TEAM_ID=your_apple_team_id
```

## Docker Usage

Build the image:
```bash
docker build -t notification-service .
```

Run the container:
```bash
docker run -p 8084:8084 \
  -e PORT=8084 \
  -e ENV=production \
  notification-service
```

## Development

Install dependencies:
```bash
go mod tidy
```

Run the service:
```bash
go run .
```

Run tests:
```bash
go test ./...
```

## Architecture

The service follows a clean architecture pattern:

- **Handlers**: HTTP request handling and response formatting
- **Service**: Business logic and notification processing
- **Repository**: Data persistence and retrieval
- **Models**: Data structures and validation

## Integration

### With Other Services

The notification service integrates with:

- **User Service**: Retrieves user preferences and contact information
- **Booking Service**: Receives booking events for notifications
- **Resource Service**: Gets resource information for notifications

### External Services

For production deployment, integrate with:

- **Email**: SendGrid, AWS SES, Mailgun
- **SMS**: Twilio, AWS SNS, Nexmo
- **Push**: Firebase FCM, Apple APNS
- **Webhook**: Custom HTTP endpoints

## Monitoring

Key metrics to monitor:

- Notification delivery rates by channel
- Failed notification counts
- Response times for notification sending
- User engagement with notifications
- Channel preference distribution

## Future Enhancements

- [ ] Template-based notifications
- [ ] User notification preferences management
- [ ] Scheduled notifications
- [ ] Notification batching and throttling
- [ ] A/B testing for notification content
- [ ] Rich media support (images, attachments)
- [ ] Delivery receipt tracking
- [ ] Notification analytics dashboard
