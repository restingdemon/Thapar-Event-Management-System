package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

// var DB *mongo.Client


type Society struct {

	Soc_ID    primitive.ObjectID `json:"_Sid,omitempty" bson:"_Sid,omitempty"`
	User_ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	YearOfFormation string `json:"year_of_formation" bson:"year_of_formation"`
	Role string `json:"role" bson:"role"`
}