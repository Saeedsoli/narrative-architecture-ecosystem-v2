// apps/backend/internal/infrastructure/database/mongodb/connection.go

package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig شامل تنظیمات اتصال به MongoDB است.
type MongoConfig struct {
	URI    string
	DBName string
}

// ConnectMongo یک کلاینت MongoDB ایجاد و برمی‌گرداند.
func ConnectMongo(cfg MongoConfig) (*mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to MongoDB: %v", err)
	}

	// بررسی صحت اتصال
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("FATAL: Failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connection established successfully.")
	db := client.Database(cfg.DBName)
	return client, db
}