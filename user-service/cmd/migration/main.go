package main

import (
	"log"

	"github.com/joho/godotenv"
	migrate "github.com/xakep666/mongo-migrate"

	_ "user-service/cmd/migration/script" // migration script
	"user-service/pkg/mongodb"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	db := mongodb.GetDatabase()
	migrate.SetDatabase(db)

	if err := migrate.Up(migrate.AllAvailable); err != nil {
		log.Fatal(err)
	}
}
