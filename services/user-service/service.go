package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository UserRepository
	jwtSecret  string
}

func NewUserService() *UserService {
	return &UserService{
		repository: NewUserRepository(),
		jwtSecret:  "your-secret-key", // TODO: Move to environment variable
	}
}

// Create creates a new user
func (s *UserService) Create(req CreateUserRequest) (*UserProfile, error) {
	// Check if user already exists
	existingUser, _ := s.repository.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := s.repository.Create(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	profile := user.ToProfile()
	return &profile, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int) (*UserProfile, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	profile := user.ToProfile()
	return &profile, nil
}

// Update updates a user
func (s *UserService) Update(id int, req UpdateUserRequest) (*UserProfile, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, _ := s.repository.GetByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already taken by another user")
		}
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	
	user.UpdatedAt = time.Now()
	
	if err := s.repository.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	
	profile := user.ToProfile()
	return &profile, nil
}

// Delete deletes a user
func (s *UserService) Delete(id int) error {
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	
	if !user.IsActive {
		return "", errors.New("user account is disabled")
	}
	
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	
	// Generate JWT token (stub - implement actual JWT generation)
	token := fmt.Sprintf("jwt_token_for_user_%d_%d", user.ID, time.Now().Unix())
	
	return token, nil
}

// RefreshToken refreshes a JWT token
func (s *UserService) RefreshToken(token string) (string, error) {
	// TODO: Implement actual JWT validation and refresh
	// This is a stub implementation
	
	// Validate existing token
	if token == "" {
		return "", errors.New("invalid token")
	}
	
	// Generate new token (stub)
	newToken := fmt.Sprintf("refreshed_%s_%d", token, time.Now().Unix())
	
	return newToken, nil
}

// ValidateToken validates a JWT token and returns user ID
func (s *UserService) ValidateToken(token string) (int, error) {
	// TODO: Implement actual JWT validation
	// This is a stub implementation
	
	if token == "" {
		return 0, errors.New("invalid token")
	}
	
	// Stub: extract user ID from token (in real implementation, parse JWT)
	return 1, nil
}
