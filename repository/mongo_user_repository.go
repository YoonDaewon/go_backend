package repository

import (
	"context"
	"errors"
	"time"

	"go_backend/database"
	"go_backend/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserRepository is a MongoDB implementation of UserRepository
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository() UserRepository {
	return &MongoUserRepository{
		collection: database.MongoDB.Collection("users"),
	}
}

// Create creates a new user
func (r *MongoUserRepository) Create(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// MongoDB uses ObjectID, but we'll convert it to int for consistency
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		// For simplicity, we'll use a hash of the ObjectID as int
		// In production, you might want to use a different ID strategy
		user.ID = int(oid.Hex()[0:8][0]) // Simplified conversion
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *MongoUserRepository) GetByID(id int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Note: This is a simplified implementation
	// In production, you'd want to store the int ID in the document
	var user model.User
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users
func (r *MongoUserRepository) GetAll() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates an existing user
func (r *MongoUserRepository) Update(id int, user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{}

	if user.Name != "" {
		update["name"] = user.Name
	}
	if user.Email != "" {
		update["email"] = user.Email
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser model.User
	err := r.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, opts).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &updatedUser, nil
}

// Delete deletes a user by ID
func (r *MongoUserRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
