package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoInstance *mongo.Database
	mongoOnce    sync.Once
)

func ConnectMongo() {
	GetMongo()
}

func GetMongo() *mongo.Database {
	mongoOnce.Do(func() {
		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			log.Print("MONGO_URI not set, MongoDB disabled")
			return
		}

		dbName := os.Getenv("MONGO_DATABASE")
		if dbName == "" {
			dbName = "amigo"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return
		}

		if err := client.Ping(ctx, nil); err != nil {
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}

		log.Printf("connected to MongoDB (%s)", dbName)
		mongoInstance = client.Database(dbName)
	})
	return mongoInstance
}
