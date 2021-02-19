package script

import (
	"context"
	"log"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"

	"user-service/internal/domain/constant"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		err := db.CreateCollection(context.TODO(), constant.UserCollection)
		if err != nil {
			return err
		}

		log.Println("success create user collection")
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
