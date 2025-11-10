package repository

import (
	"errors"

	"go_backend/database"
	"go_backend/model"

	"gorm.io/gorm"
)

// PostgresUserRepository is a PostgreSQL implementation of UserRepository
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository() UserRepository {
	return &PostgresUserRepository{
		db: database.PostgresDB,
	}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users
func (r *PostgresUserRepository) GetAll() ([]*model.User, error) {
	var users []*model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(id int, user *model.User) (*model.User, error) {
	var existingUser model.User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}

	if err := r.db.Model(&existingUser).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}

// Delete deletes a user by ID
func (r *PostgresUserRepository) Delete(id int) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
