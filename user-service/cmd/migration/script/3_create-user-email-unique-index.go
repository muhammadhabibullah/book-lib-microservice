package script

import (
	"context"
	"log"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"user-service/internal/domain/constant"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		opt := options.Index().SetName(constant.UserEmailUniqueIndex).
			SetUnique(true)
		keys := bson.D{{"email", 1}}
		model := mongo.IndexModel{Keys: keys, Options: opt}

		idx, err := db.Collection(constant.UserCollection).Indexes().
			CreateOne(context.TODO(), model)
		if err != nil {
			return err
		}

		log.Printf("success create %s\n", idx)
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
