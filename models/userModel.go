package models

import (
	"github.com/restingdemon/thaparEvents/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Client

type User struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	RollNo string `json:"rollno" bson:"rollno"`
	Branch string `json:"branch" bson:"branch"`
	Batch string `json:"batch" bson:"batch"`
	Role string `json:"role" bson:"role"`
	Image string `json:"image" bson:"image"`
}

func init() {
	database.Connect()
	DB = database.GetDB()
}