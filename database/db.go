package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Connect() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		panic(err)
	}
	DB = mongoClient
}

func GetDB() *mongo.Client {
	return DB
}
