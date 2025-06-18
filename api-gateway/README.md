# API Gateway

The API Gateway serves as the single entry point for all client requests to the reservation system microservices. It's built using KrakenD, a high-performance API Gateway that provides routing, load balancing, authentication, and other essential features.

## Features

- **Unified API**: Single entry point for all microservices
- **Load Balancing**: Distributes requests across service instances  
- **CORS Support**: Cross-origin resource sharing configuration
- **Rate Limiting**: Request throttling and protection
- **Request/Response Transformation**: Data manipulation capabilities
- **Health Checks**: Service health monitoring
- **Metrics Collection**: Performance and usage analytics
- **Caching**: Response caching for improved performance

## Configuration

The gateway is configured through `krakend.json` which defines:

- **Endpoints**: Client-facing API endpoints
- **Backends**: Internal microservice endpoints
- **Routing Rules**: How requests are forwarded
- **Middleware**: Additional processing features

## API Endpoints

All endpoints are prefixed with `/api/v1`:

### User Management

- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `POST /api/v1/auth/login` - User authentication

### Resource Management

- `GET /api/v1/resources` - List resources
- `POST /api/v1/resources` - Create resource
- `GET /api/v1/resources/{id}` - Get resource by ID
- `PUT /api/v1/resources/{id}` - Update resource
- `DELETE /api/v1/resources/{id}` - Delete resource
- `GET /api/v1/resources/{id}/availability` - Check availability

### Booking Management

- `GET /api/v1/bookings` - List bookings
- `POST /api/v1/bookings` - Create booking
- `GET /api/v1/bookings/{id}` - Get booking by ID
- `PUT /api/v1/bookings/{id}` - Update booking
- `DELETE /api/v1/bookings/{id}` - Delete booking
- `PUT /api/v1/bookings/{id}/cancel` - Cancel booking

### Notification Management

- `GET /api/v1/notifications` - List notifications
- `POST /api/v1/notifications` - Send notification
- `PUT /api/v1/notifications/{id}/status` - Update notification status
- `GET /api/v1/notifications/stats` - Get notification statistics

### System

- `GET /health` - System health check

## Service Discovery

The gateway uses static service discovery with the following service hosts:

- **user-service**: `http://user-service:8081`
- **resource-service**: `http://resource-service:8082`
- **booking-service**: `http://booking-service:8083`
- **notification-service**: `http://notification-service:8084`

## Security

### CORS Configuration

- Allows all origins (configure for production)
- Supports all HTTP methods
- Allows common headers (Authorization, Content-Type)
- Credentials support disabled by default

### Authentication

- JWT token validation (configured per endpoint)
- Bearer token support
- Token forwarding to backend services

## Monitoring

### Metrics

- Endpoint performance metrics
- Backend response times
- Error rates and status codes
- Request volume statistics
- Available on port 8090

### Logging

- Structured JSON logging
- Request/response logging
- Error tracking
- Debug mode support

## Docker Usage

Build the image:

```bash
docker build -t api-gateway .
```

Run the container:

```bash
docker run -p 8080:8080 api-gateway
```

## Development

Start KrakenD with configuration:

```bash
krakend run -d -c krakend.json
```

Validate configuration:

```bash
krakend check -d -c krakend.json
```

## Configuration Examples

### Adding Authentication

```json
{
  "endpoint": "/api/v1/protected-endpoint",
  "extra_config": {
    "auth/validator": {
      "alg": "HS256",
      "key": "your-secret-key",
      "issuer": "your-issuer",
      "audience": ["your-audience"]
    }
  }
}
```

### Adding Rate Limiting

```json
{
  "endpoint": "/api/v1/limited-endpoint",
  "extra_config": {
    "qos/ratelimit/router": {
      "max_rate": 100,
      "capacity": 100
    }
  }
}
```

### Adding Response Caching

```json
{
  "endpoint": "/api/v1/cached-endpoint",
  "extra_config": {
    "qos/http-cache": {
      "ttl": "300s"
    }
  }
}
```

## Production Considerations

1. **Security**:
   - Configure specific CORS origins
   - Enable JWT authentication
   - Set up rate limiting
   - Use HTTPS/TLS

2. **Performance**:
   - Enable response caching
   - Configure connection pooling
   - Set appropriate timeouts
   - Use circuit breakers

3. **Monitoring**:
   - Set up metrics collection
   - Configure log aggregation
   - Enable health checks
   - Set up alerting

4. **High Availability**:
   - Deploy multiple gateway instances
   - Use load balancer in front
   - Configure service mesh
   - Set up failover strategies

## Troubleshooting

### Common Issues

1. **Service Unreachable**:
   - Check service hostnames in configuration
   - Verify network connectivity
   - Check service health endpoints

2. **CORS Errors**:
   - Update allowed origins
   - Check allowed methods and headers
   - Verify preflight requests

3. **Authentication Failures**:
   - Verify JWT configuration
   - Check token format and claims
   - Validate signing keys

4. **Performance Issues**:
   - Enable caching where appropriate
   - Optimize timeout values
   - Check backend service performance
   - Monitor gateway metrics
