package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"book-service/pkg/mongodb"
)

type Book struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	mongodb.Meta `json:"meta" bson:"meta"`
	Title        string `json:"title" bson:"title"`
	Stock        int    `json:"stock" bson:"stock"`
}

type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Book, error)
	Count(ctx context.Context, filter map[string]interface{}) (int, error)
	FindByID(ctx context.Context, id string) (Book, error)
	FindByTitle(ctx context.Context, title string) (Book, error)
	Update(ctx context.Context, book *Book) error
	UpdateStock(ctx context.Context, book *Book, newStock int) error
	Delete(ctx context.Context, book *Book) error
}
