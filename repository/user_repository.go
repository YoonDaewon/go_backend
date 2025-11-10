package repository

import (
	"errors"
	"sync"

	"go_backend/model"
)

// UserRepository handles user data operations
type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(id int, user *model.User) (*model.User, error)
	Delete(id int) error
}

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	users map[int]*model.User
	mu    sync.RWMutex
	idSeq int
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make(map[int]*model.User),
		idSeq: 1,
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.idSeq
	r.idSeq++
	r.users[user.ID] = user

	return user, nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id int) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetAll retrieves all users
func (r *InMemoryUserRepository) GetAll() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(id int, user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingUser, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}

	return existingUser, nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

