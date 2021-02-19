package script

import (
	"context"
	"log"
	"os"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"user-service/internal/domain"
	"user-service/internal/domain/constant"
	"user-service/pkg/password"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {

		emailRole := make(map[string]string)
		emailRole[os.Getenv("ADMIN_EMAIL")] = constant.AdminRole
		emailRole[os.Getenv("LIBRARIAN_EMAIL")] = constant.LibrarianRole
		emailRole[os.Getenv("MEMBER_EMAIL")] = constant.MemberRole
		defaultPass := os.Getenv("DEFAULT_USER_PASSWORD")

		for email, role := range emailRole {
			hashedPassword, _ := password.Hash(defaultPass)
			user := domain.User{
				ID:             primitive.NewObjectID(),
				Email:          email,
				HashedPassword: hashedPassword,
				Role:           role,
			}
			user.Meta.Create()
			_, err := db.Collection(constant.UserCollection).InsertOne(context.TODO(), &user)
			if err != nil {
				return err
			}
		}

		log.Println("success create default users")
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
