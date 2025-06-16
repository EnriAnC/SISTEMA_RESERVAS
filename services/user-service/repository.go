package main

import (
	"errors"
	"fmt"
	"sync"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	List(limit, offset int) ([]*User, error)
}

// InMemoryUserRepository is a simple in-memory implementation of UserRepository
// TODO: Replace with actual database implementation (PostgreSQL)
type InMemoryUserRepository struct {
	users  map[int]*User
	emails map[string]*User
	nextID int
	mutex  sync.RWMutex
}

func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		emails: make(map[string]*User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if email already exists
	if _, exists := r.emails[user.Email]; exists {
		return errors.New("user with this email already exists")
	}
	
	// Assign ID and store
	user.ID = r.nextID
	r.nextID++
	
	r.users[user.ID] = user
	r.emails[user.Email] = user
	
	return nil
}

func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.emails[email]
	if !exists {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	
	return user, nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	existingUser, exists := r.users[user.ID]
	if !exists {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}
	
	// Update email mapping if email changed
	if existingUser.Email != user.Email {
		delete(r.emails, existingUser.Email)
		r.emails[user.Email] = user
	}
	
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user with ID %d not found", id)
	}
	
	delete(r.users, id)
	delete(r.emails, user.Email)
	
	return nil
}

func (r *InMemoryUserRepository) List(limit, offset int) ([]*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*User, 0, len(r.users))
	count := 0
	skipped := 0
	
	for _, user := range r.users {
		if skipped < offset {
			skipped++
			continue
		}
		
		if count >= limit {
			break
		}
		
		users = append(users, user)
		count++
	}
	
	return users, nil
}

// TODO: PostgreSQL implementation
/*
type PostgreSQLUserRepository struct {
	db *sql.DB
}

func NewPostgreSQLUserRepository(db *sql.DB) UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	
	err := r.db.QueryRow(
		query,
		user.Name, user.Email, user.Password, user.Role, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	
	return err
}

// ... implement other methods
*/
