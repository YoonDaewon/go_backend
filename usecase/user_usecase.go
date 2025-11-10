package usecase

import (
	"errors"

	"go_backend/model"
	"go_backend/repository"
)

// UserUsecase handles user business logic
type UserUsecase interface {
	CreateUser(req *model.CreateUserRequest) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	UpdateUser(id int, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(id int) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (u *userUsecase) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
	// Business logic validation
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	return u.userRepo.Create(user)
}

// GetUserByID retrieves a user by ID
func (u *userUsecase) GetUserByID(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return u.userRepo.GetByID(id)
}

// GetAllUsers retrieves all users
func (u *userUsecase) GetAllUsers() ([]*model.User, error) {
	return u.userRepo.GetAll()
}

// UpdateUser updates an existing user
func (u *userUsecase) UpdateUser(id int, req *model.UpdateUserRequest) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	return u.userRepo.Update(id, user)
}

// DeleteUser deletes a user by ID
func (u *userUsecase) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	return u.userRepo.Delete(id)
}

