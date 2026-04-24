package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO_DB *mongo.Client

func MongoDB_Connection() {
	fmt.Println("we are in mongoDB connetion function")
	MONGO_URI := os.Getenv("MONGO_URI")
	fmt.Println("MongoDB URI:", MONGO_URI)

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		fmt.Println("err while creating mongo client", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("error while connection :", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("error while ping to DB :", err)
	}

	fmt.Println("Connected to MongoDB!")
	MONGO_DB = client
}

func GetCollection(collectionName string) *mongo.Collection {
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		fmt.Println("DB_NAME is not set in environment variables")
		return nil
	}
	collection := MONGO_DB.Database(DB_NAME).Collection(collectionName)
	fmt.Println("collection name is ", collectionName)
	return collection
}
