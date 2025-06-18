# Copilot Instructions for Sistema de Reservas

## Project Context

This is a cloud-based reservation system built using microservices architecture with Go/Golang.

## Architecture Guidelines

- Follow microservices patterns and principles
- Each service should be independent and scalable
- Use Go best practices and idiomatic code
- Implement proper error handling and logging
- Follow RESTful API design principles
- Use dependency injection patterns
- Implement proper testing strategies (unit, integration)

## Code Style

- Use Go standard formatting (gofmt)
- Follow Go naming conventions
- Use structured logging
- Implement proper HTTP status codes
- Use context for request tracing
- Implement graceful shutdowns

## Security Considerations

- Always validate input data
- Use JWT tokens for authentication
- Implement proper CORS policies
- Use HTTPS/TLS for all communications
- Sanitize database queries to prevent SQL injection
- Implement rate limiting

## Database Patterns

- Use repository pattern for data access
- Implement proper transaction handling
- Use database migrations
- Follow normalization principles

## Testing Guidelines

- Write unit tests for all business logic
- Use table-driven tests when appropriate
- Mock external dependencies
- Implement integration tests for APIs
- Use test containers for database testing

## Documentation

- Update API documentation when adding new endpoints
- Comment complex business logic
- Update README files when adding new features
- Document deployment procedures
