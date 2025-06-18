# Development Guide

This guide provides instructions for setting up the development environment
and contributing to the Sistema de Reservas project.

## Prerequisites

### Required Software

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)
- **kubectl**: [Install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- **Git**: [Install Git](https://git-scm.com/downloads)

### Optional Tools

- **Postman/Insomnia**: For API testing
- **pgAdmin**: For PostgreSQL database management
- **Lens**: Kubernetes IDE for cluster management

## Development Environment Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd SISTEMA_RESERVAS
```

### 2. Set Up Local Dependencies

```bash
# Install Go dependencies for all services
cd services/user-service && go mod tidy && cd ../..
cd services/resource-service && go mod tidy && cd ../..
cd services/booking-service && go mod tidy && cd ../..
cd services/notification-service && go mod tidy && cd ../..
```

### 3. Start Infrastructure Services

```bash
# Start PostgreSQL, Redis, and monitoring stack
docker-compose up -d postgres redis prometheus grafana
```

### 4. Initialize Database

```bash
# Apply database schema
docker-compose exec postgres psql -U reservas_user -d reservas_db -f /docker-entrypoint-initdb.d/init.sql
```

## Running Services

### Option 1: Run Services Natively

```bash
# Terminal 1 - User Service
cd services/user-service
export PORT=8081
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=reservas_user
export DB_PASSWORD=reservas_pass
export DB_NAME=reservas_db
go run .

# Terminal 2 - Resource Service
cd services/resource-service
export PORT=8082
# ... (same DB config)
go run .

# Terminal 3 - Booking Service
cd services/booking-service
export PORT=8083
# ... (same DB config)
go run .

# Terminal 4 - Notification Service
cd services/notification-service
export PORT=8084
go run .

# Terminal 5 - API Gateway
cd api-gateway
docker run -p 8080:8080 -v $PWD:/etc/krakend/ devopsfaith/krakend run --config /etc/krakend/krakend.json
```

### Option 2: Run with Docker Compose

```bash
# Build and start all services
docker-compose up --build
```

## Testing

### Unit Tests

```bash
# Run tests for all services
./scripts/run-tests.sh

# Or run individually
cd services/user-service && go test ./...
cd services/resource-service && go test ./...
cd services/booking-service && go test ./...
cd services/notification-service && go test ./...
```

### Integration Tests

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags=integration ./tests/...

# Cleanup
docker-compose -f docker-compose.test.yml down
```

### API Testing

Use the provided Postman collection or curl commands:

```bash
# Health checks
curl http://localhost:8080/users/health
curl http://localhost:8080/resources/health
curl http://localhost:8080/bookings/health
curl http://localhost:8080/notifications/health

# Create a user
curl -X POST http://localhost:8080/users/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'
```

## Code Style and Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting
- Use `golint` for linting
- Use `go vet` for code analysis

### Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```Commit
feat(user-service): add user registration endpoint
fix(booking-service): resolve double booking issue
docs(api): update booking endpoints documentation
```

### Code Review Checklist

- [ ] Code follows Go best practices
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] API changes are documented
- [ ] Docker builds successfully
- [ ] No sensitive data in code

## Debugging

### Service Logs

```bash
# Docker Compose logs
docker-compose logs -f user-service
docker-compose logs -f resource-service
docker-compose logs -f booking-service
docker-compose logs -f notification-service

# Kubernetes logs
kubectl logs -f deployment/user-service
kubectl logs -f deployment/resource-service
kubectl logs -f deployment/booking-service
kubectl logs -f deployment/notification-service
```

### Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U reservas_user -d reservas_db

# Common queries
SELECT * FROM users LIMIT 10;
SELECT * FROM resources WHERE available = true;
SELECT * FROM bookings WHERE status = 'confirmed';
```

### Monitoring

- **Prometheus**: <http://localhost:9090>
- **Grafana**: <http://localhost:3000> (admin/admin)
- **Application Metrics**: Each service exposes `/metrics` endpoint

## Adding New Features

### 1. Service Modifications

1. Update models in `models.go`
2. Add business logic in `service.go`
3. Update repository in `repository.go`
4. Add HTTP handlers in `handlers.go`
5. Update routes in `main.go`

### 2. Database Changes

1. Create migration script in `infrastructure/database/migrations/`
2. Update `init.sql` for fresh installations
3. Test migration on development database

### 3. API Gateway Updates

1. Update `api-gateway/krakend.json`
2. Add new routes and backends
3. Update rate limiting if needed

### 4. Documentation Updates

1. Update API documentation in `docs/API.md`
2. Update architecture diagrams if needed
3. Update deployment guides

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port
lsof -i :8081
# Kill process
kill -9 <PID>
```

#### Database Connection Issues

```bash
# Check database status
docker-compose ps postgres
# Restart database
docker-compose restart postgres
```

#### Docker Build Issues

```bash
# Clean Docker cache
docker system prune -f
# Rebuild without cache
docker-compose build --no-cache
```

#### Go Module Issues

```bash
# Clean module cache
go clean -modcache
# Re-download dependencies
go mod download
```

### Performance Issues

1. Check service resource usage: `docker stats`
2. Monitor database connections
3. Check API Gateway metrics
4. Review application logs for errors

## Continuous Integration

### GitHub Actions

The project includes CI/CD workflows:

- **Build**: Compiles all services
- **Test**: Runs unit and integration tests
- **Security**: Scans for vulnerabilities
- **Deploy**: Deploys to staging/production

### Local CI Testing

```bash
# Run CI pipeline locally with act (GitHub Actions local runner)
act push
```

## Contributing

### Pull Request Process

1. Create feature branch: `git checkout -b feature/new-feature`
2. Make changes and add tests
3. Commit using conventional commit format
4. Push branch and create pull request
5. Address review comments
6. Merge after approval

### Issue Reporting

When reporting issues, include:

- Service affected
- Steps to reproduce
- Expected vs actual behavior
- Logs and error messages
- Environment details

## Development Tools

### Recommended VS Code Extensions

- Go extension
- Docker extension
- Kubernetes extension
- REST Client
- GitLens
- Thunder Client (API testing)

### Environment Variables Template

Create `.env` file in each service directory:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=reservas_user
DB_PASSWORD=reservas_pass
DB_NAME=reservas_db

# Service
PORT=8081
ENV=development
LOG_LEVEL=debug

# External services (optional)
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [KrakenD Documentation](https://www.krakend.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Microservices Patterns](https://microservices.io/)

---

For questions or support, please create an issue in the repository or contact the development team.
