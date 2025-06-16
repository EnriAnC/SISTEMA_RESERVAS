# API Documentation - Sistema de Reservas

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Error Responses
All errors follow a consistent format:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "timestamp": "2024-06-09T10:30:00Z"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

---

## User Management

### Authentication

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user"
  }
}
```

### User Operations

#### Create User
```http
POST /users
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "first_name": "Jane",
  "last_name": "Smith",
  "phone": "+1234567890",
  "role": "user"
}
```

#### Get All Users
```http
GET /users?limit=10&offset=0&role=user
Authorization: Bearer <token>
```

#### Get User by ID
```http
GET /users/{id}
Authorization: Bearer <token>
```

#### Update User
```http
PUT /users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Updated Name",
  "phone": "+0987654321"
}
```

#### Delete User
```http
DELETE /users/{id}
Authorization: Bearer <token>
```

---

## Resource Management

### Resource Operations

#### Get All Resources
```http
GET /resources?type=meeting_room&location=Floor1&available=true&limit=10&offset=0
```

**Query Parameters:**
- `type` - Filter by resource type
- `location` - Filter by location
- `available` - Filter by availability
- `capacity_min` - Minimum capacity
- `capacity_max` - Maximum capacity
- `limit` - Number of results (default: 10)
- `offset` - Pagination offset (default: 0)

**Response:**
```json
{
  "resources": [
    {
      "id": 1,
      "name": "Conference Room A",
      "description": "Large conference room with projector",
      "type": "meeting_room",
      "location": "Floor 1",
      "capacity": 10,
      "price_per_hour": 50.00,
      "amenities": {
        "projector": true,
        "whiteboard": true,
        "video_conference": true
      },
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

#### Create Resource
```http
POST /resources
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Meeting Room",
  "description": "Small meeting space",
  "type": "meeting_room",
  "location": "Floor 2",
  "capacity": 6,
  "price_per_hour": 30.00,
  "amenities": {
    "whiteboard": true,
    "phone": true
  }
}
```

#### Get Resource by ID
```http
GET /resources/{id}
```

#### Update Resource
```http
PUT /resources/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Room Name",
  "capacity": 8,
  "price_per_hour": 35.00
}
```

#### Delete Resource
```http
DELETE /resources/{id}
Authorization: Bearer <token>
```

#### Check Resource Availability
```http
GET /resources/{id}/availability?date=2024-06-10&start_time=09:00&end_time=17:00
```

**Response:**
```json
{
  "resource_id": 1,
  "date": "2024-06-10",
  "availability": [
    {
      "start_time": "09:00",
      "end_time": "10:00",
      "available": true
    },
    {
      "start_time": "10:00",
      "end_time": "12:00",
      "available": false,
      "booking_id": 123
    }
  ]
}
```

---

## Booking Management

### Booking Operations

#### Get All Bookings
```http
GET /bookings?user_id=123&resource_id=456&status=confirmed&start_date=2024-06-01&end_date=2024-06-30&limit=10&offset=0
Authorization: Bearer <token>
```

**Query Parameters:**
- `user_id` - Filter by user
- `resource_id` - Filter by resource
- `status` - Filter by status (pending, confirmed, cancelled, completed)
- `start_date` - Filter bookings from date
- `end_date` - Filter bookings to date
- `limit` - Number of results
- `offset` - Pagination offset

#### Create Booking
```http
POST /bookings
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "notes": "Team meeting for project planning"
}
```

**Response:**
```json
{
  "id": 123,
  "user_id": 456,
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "status": "pending",
  "total_price": 100.00,
  "notes": "Team meeting for project planning",
  "created_at": "2024-06-09T15:30:00Z"
}
```

#### Get Booking by ID
```http
GET /bookings/{id}
Authorization: Bearer <token>
```

#### Update Booking
```http
PUT /bookings/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "start_time": "2024-06-10T11:00:00Z",
  "end_time": "2024-06-10T13:00:00Z",
  "notes": "Updated meeting time"
}
```

#### Cancel Booking
```http
PUT /bookings/{id}/cancel
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Meeting postponed"
}
```

#### Delete Booking
```http
DELETE /bookings/{id}
Authorization: Bearer <token>
```

---

## Notification Management

### Notification Operations

#### Send Notification
```http
POST /notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 123,
  "type": "booking",
  "title": "Booking Confirmed",
  "message": "Your booking for Conference Room A has been confirmed",
  "channel": "email",
  "priority": "high",
  "metadata": {
    "booking_id": 456,
    "resource_name": "Conference Room A"
  }
}
```

#### Get User Notifications
```http
GET /notifications?user_id=123&type=booking&is_read=false&limit=10&offset=0
Authorization: Bearer <token>
```

**Response:**
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
      "created_at": "2024-06-09T10:30:00Z",
      "sent_at": "2024-06-09T10:30:05Z"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

#### Update Notification Status
```http
PUT /notifications/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_read": true
}
```

#### Get Notification Statistics
```http
GET /notifications/stats?user_id=123
Authorization: Bearer <token>
```

**Response:**
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

---

## System Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "services": {
    "user-service": "healthy",
    "resource-service": "healthy",
    "booking-service": "healthy",
    "notification-service": "healthy"
  },
  "timestamp": "2024-06-09T10:30:00Z"
}
```

---

## Data Models

### User Model
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "role": "user",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login": "2024-06-09T08:00:00Z"
}
```

### Resource Model
```json
{
  "id": 1,
  "name": "Conference Room A",
  "description": "Large conference room with projector",
  "type": "meeting_room",
  "location": "Floor 1",
  "capacity": 10,
  "price_per_hour": 50.00,
  "amenities": {
    "projector": true,
    "whiteboard": true,
    "video_conference": true
  },
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Booking Model
```json
{
  "id": 123,
  "user_id": 456,
  "resource_id": 1,
  "start_time": "2024-06-10T10:00:00Z",
  "end_time": "2024-06-10T12:00:00Z",
  "status": "confirmed",
  "total_price": 100.00,
  "notes": "Team meeting",
  "metadata": {},
  "created_at": "2024-06-09T15:30:00Z",
  "updated_at": "2024-06-09T15:30:00Z",
  "cancelled_at": null,
  "cancellation_reason": null
}
```

### Notification Model
```json
{
  "id": 1,
  "user_id": 123,
  "type": "booking",
  "title": "Booking Confirmed",
  "message": "Your booking has been confirmed",
  "channel": "email",
  "priority": "high",
  "status": "sent",
  "is_read": false,
  "metadata": {
    "booking_id": 456
  },
  "created_at": "2024-06-09T10:30:00Z",
  "sent_at": "2024-06-09T10:30:05Z",
  "read_at": null
}
```

---

## Webhook Events

The system publishes events for external integration:

### Booking Events
- `booking.created` - New booking created
- `booking.confirmed` - Booking confirmed
- `booking.cancelled` - Booking cancelled
- `booking.updated` - Booking modified

### User Events
- `user.registered` - New user registered
- `user.updated` - User information updated

### Event Payload Example
```json
{
  "event": "booking.confirmed",
  "timestamp": "2024-06-09T10:30:00Z",
  "data": {
    "booking_id": 123,
    "user_id": 456,
    "resource_id": 1,
    "start_time": "2024-06-10T10:00:00Z",
    "end_time": "2024-06-10T12:00:00Z"
  }
}
```

---

## Rate Limiting

API endpoints are rate limited to prevent abuse:

- **Authentication endpoints**: 5 requests per minute per IP
- **User endpoints**: 100 requests per hour per user
- **Resource endpoints**: 200 requests per hour per user
- **Booking endpoints**: 50 requests per hour per user
- **Notification endpoints**: 100 requests per hour per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1623456789
```

---

## SDK and Examples

### JavaScript/Node.js Example
```javascript
const axios = require('axios');

// Login
const loginResponse = await axios.post('http://localhost:8080/api/v1/auth/login', {
  email: 'user@example.com',
  password: 'password123'
});

const token = loginResponse.data.token;

// Create booking
const bookingResponse = await axios.post(
  'http://localhost:8080/api/v1/bookings',
  {
    resource_id: 1,
    start_time: '2024-06-10T10:00:00Z',
    end_time: '2024-06-10T12:00:00Z',
    notes: 'Team meeting'
  },
  {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  }
);
```

### cURL Examples
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get resources
curl -X GET "http://localhost:8080/api/v1/resources?type=meeting_room" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Create booking
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_id": 1,
    "start_time": "2024-06-10T10:00:00Z",
    "end_time": "2024-06-10T12:00:00Z",
    "notes": "Team meeting"
  }'
```