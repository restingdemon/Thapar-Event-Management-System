package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Connect() {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://ajak:xu2pgtDcfiat6YSq@thaparevent.xaxfle4.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	DB = mongoClient
}

func GetDB() *mongo.Client{
	return DB
}
