package db

import (
	"context"
	"time"

	"github.com/JustBrowsing/query-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient(ctx context.Context, cfg config.MongoDBConfig) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(uint64(cfg.PoolSize)).
		SetConnectTimeout(cfg.Timeout * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Ping the primary to verify that the connection is established
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// NewMongoDBDatabase creates a new MongoDB database
func NewMongoDBDatabase(client *mongo.Client, cfg config.MongoDBConfig) *mongo.Database {
	return client.Database(cfg.Database)
}