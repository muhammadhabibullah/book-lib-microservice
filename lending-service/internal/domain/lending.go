package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"lending-service/internal/domain/constant"
	"lending-service/pkg/mongodb"
)

type Lending struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	mongodb.Meta `json:"meta" bson:"meta"`
	BookID       primitive.ObjectID     `json:"book_id" bson:"book_id"`
	UserID       primitive.ObjectID     `json:"user_id" bson:"user_id"`
	Status       constant.LendingStatus `json:"status" bson:"status"`
	ReturnDate   time.Time              `json:"return_date" bson:"return_date"`
}

type LendingRepository interface {
	Create(ctx context.Context, lending *Lending) error
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Lending, error)
	FindByID(ctx context.Context, id string) (Lending, error)
	Count(ctx context.Context, filter map[string]interface{}) (int, error)
	Update(ctx context.Context, lending *Lending) error
}
