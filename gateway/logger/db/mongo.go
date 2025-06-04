/* Mongo Database

MongoDB pool singleton

Joaquin Badillo
2025-06-04
*/

package db

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	once       sync.Once
	connectErr error
)

func GetMongoClient() (*mongo.Client, error) {
	once.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = "mongodb://db:27017" // Default fallback
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		opts := options.Client().ApplyURI(uri)
		client, connectErr = mongo.Connect(ctx, opts)
		if connectErr != nil {
			return
		}

		// Verify connection
		connectErr = client.Ping(ctx, nil)
		if connectErr != nil {
			log.Printf("Failed to ping MongoDB: %v", connectErr)
		} else {
			log.Println("MongoDB connection established")
		}
	})

	return client, connectErr
}

func CloseMongoClient(ctx context.Context) error {
	if client == nil {
		return nil // Nothing to close
	}
	return client.Disconnect(ctx)
}
