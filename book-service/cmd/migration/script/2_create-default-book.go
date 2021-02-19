package script

import (
	"context"
	"fmt"
	"log"
	"time"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"book-service/internal/domain"
	"book-service/internal/domain/constant"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		for i := 0; i <= 5; i++ {
			book := domain.Book{
				ID:    primitive.NewObjectID(),
				Title: fmt.Sprintf("Book %d", i),
				Stock: i,
			}
			book.Meta.Create()

			_, err := db.Collection(constant.BookCollection).InsertOne(context.TODO(), &book)
			if err != nil {
				return err
			}

			time.Sleep(1 * time.Millisecond)
		}

		log.Println("success create default book")
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
