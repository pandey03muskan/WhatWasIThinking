package config

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO_DB *mongo.Client

func MongoDB_Connection() {
	slog.Info("we are in mongoDB connetion function")
	MONGO_URI := os.Getenv("MONGO_URI")
	slog.Info("MongoDB URI", "uri", MONGO_URI)

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		slog.Error("Error while creating mongo client", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		slog.Error("Error while connection", "error", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("Error while pinging DB", "error", err)
	}

	slog.Info("Connected to MongoDB!")
	MONGO_DB = client
}

func GetCollection(collectionName string) *mongo.Collection {
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		slog.Error("DB_NAME is not set in environment variables")
		return nil
	}
	collection := MONGO_DB.Database(DB_NAME).Collection(collectionName)
	slog.Info("collection name is", "collectionName", collectionName)
	return collection
}
