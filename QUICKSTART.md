# Quick Start Guide - Sistema de Reservas

This guide will help you get the Sistema de Reservas up and running in minutes.

## Prerequisites

Make sure you have installed:

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)

## Option 1: Docker Compose (Recommended for Testing)

### 1. Clone and Start

```bash
# Clone the repository
git clone <repository-url>
cd SISTEMA_RESERVAS

# Copy environment configuration
cp .env.example .env

# Start all services
docker-compose up -d
```

### 2. Wait for Services to Start

```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

### 3. Verify Installation

```bash
# Test API Gateway
curl http://localhost:8080/users/health
curl http://localhost:8080/resources/health
curl http://localhost:8080/bookings/health
curl http://localhost:8080/notifications/health
```

All health checks should return `200 OK`.

## Option 2: Kubernetes (Production)

### 1. Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f kubernetes/

# Check deployment status
kubectl get pods
kubectl get services
```

### 2. Access Services

```bash
# Port forward API Gateway
kubectl port-forward service/api-gateway 8080:8080

# Or use ingress (if configured)
# kubectl get ingress
```

## Testing the System

### 1. Using Postman

1. Import the collection: `docs/postman-collection.json`
2. Set base URL to `http://localhost:8080`
3. Run the health check requests

### 2. Using curl

```bash
# Create a new user
curl -X POST http://localhost:8080/users/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'

# Login to get JWT token
curl -X POST http://localhost:8080/users/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'

# Use the JWT token from login response for authenticated requests
JWT_TOKEN="your-jwt-token-here"

# Create a resource
curl -X POST http://localhost:8080/resources/api/v1/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "name": "Conference Room A",
    "description": "Large conference room",
    "type": "room",
    "capacity": 20,
    "location": "Building A, Floor 3",
    "price_per_hour": 50.00
  }'

# Make a booking
curl -X POST http://localhost:8080/bookings/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "resource_id": 1,
    "start_time": "2024-12-20T10:00:00Z",
    "end_time": "2024-12-20T12:00:00Z",
    "purpose": "Team meeting"
  }'
```

## Monitoring

### Access Monitoring Dashboards

- **Prometheus**: <http://localhost:9090>
- **Grafana**: <http://localhost:3000> (admin/admin)

### Key Metrics to Monitor

- Service health and uptime
- Response times
- Database connections
- Memory and CPU usage

## Common Issues

### Services Not Starting

```bash
# Check logs
docker-compose logs service-name

# Common fixes
docker-compose down
docker-compose pull
docker-compose up -d
```

### Database Connection Issues

```bash
# Check PostgreSQL
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d
```

### Port Conflicts

```bash
# Check what's using the port
netstat -tulpn | grep :8080

# Kill process or change port in docker-compose.yml
```

## Development

### Running Services Locally

```bash
# Install Go dependencies
cd services/user-service && go mod tidy
cd services/resource-service && go mod tidy
cd services/booking-service && go mod tidy
cd services/notification-service && go mod tidy

# Start infrastructure only
docker-compose up -d postgres redis prometheus grafana

# Run services natively
cd services/user-service && PORT=8081 go run .
cd services/resource-service && PORT=8082 go run .
cd services/booking-service && PORT=8083 go run .
cd services/notification-service && PORT=8084 go run .
```

### Testing

```bash
# Run all tests
./scripts/run-tests.sh

# Run individual service tests
cd services/user-service && go test ./...
```

## Next Steps

1. **Read the Documentation**:
   - [Architecture Guide](docs/ARCHITECTURE.md)
   - [API Documentation](docs/API.md)
   - [Deployment Guide](docs/DEPLOYMENT.md)
   - [Development Guide](docs/DEVELOPMENT.md)

2. **Customize Configuration**:
   - Update `.env` file with your settings
   - Configure external services (email, SMS, etc.)
   - Set up production secrets

3. **Deploy to Production**:
   - Follow the [Deployment Guide](docs/DEPLOYMENT.md)
   - Set up CI/CD pipelines
   - Configure monitoring and alerting

## Getting Help

- Check the [troubleshooting section](docs/DEVELOPMENT.md#troubleshooting) in the development guide
- Review service logs: `docker-compose logs service-name`
- Create an issue in the repository with:
  - Error messages
  - Steps to reproduce
  - Environment details

---

🎉 **Congratulations!** You now have a fully functional microservices-based reservation system running locally.
