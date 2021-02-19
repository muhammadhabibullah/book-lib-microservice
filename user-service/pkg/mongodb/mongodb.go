package mongodb

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client      *mongo.Client
	database    *mongo.Database
	connectOnce sync.Once
)

func connectDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	clientOptions = clientOptions.SetMaxPoolSize(50)

	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	database = client.Database(os.Getenv("MONGODB_DATABASE"))
}

func GetClient() *mongo.Client {
	connectOnce.Do(connectDatabase)
	return client
}

func GetDatabase() *mongo.Database {
	connectOnce.Do(connectDatabase)
	return database
}
