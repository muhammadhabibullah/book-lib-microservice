package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"user-service/pkg/mongodb"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	mongodb.Meta   `json:"meta" bson:"meta"`
	Email          string `json:"email" bson:"email"`
	HashedPassword string `json:"hashed_password" bson:"hashed_password"`
	Role           string `json:"role" bson:"role"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Fetch(ctx context.Context, filter map[string]interface{}) ([]User, error)
	Count(ctx context.Context, filter map[string]interface{}) (int, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
}
