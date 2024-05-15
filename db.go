package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. 🔐")
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("Set your MONGO URI. 🔗")
	}

	clientOptions := options.Client().ApplyURI(uri) // Connect to //MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB! 🎉")
	return client
}
