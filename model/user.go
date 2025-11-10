package model

// User represents a user entity
type User struct {
	ID    int    `json:"id" gorm:"primaryKey" bson:"_id,omitempty"`
	Name  string `json:"name" gorm:"not null" bson:"name"`
	Email string `json:"email" gorm:"uniqueIndex;not null" bson:"email"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}
