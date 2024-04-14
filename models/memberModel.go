package models

import (
	//"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

type Members struct {
	Name      string `json:"name" bson:"name"`
	Position  string `json:"position" bson:"position"`
	Email     string `json:"email" bson:"email"`
	Instagram string `json:"instagram" bson:"instagram"`
	LinkedIn  string `json:"linkedin" bson:"linkedin"`
}

type Faculty struct {
	Name      string `json:"name" bson:"name"`
	Position  string `json:"position" bson:"position"`
	Email     string `json:"email" bson:"email"`
	LinkedIn  string `json:"linkedin" bson:"linkedin"`
}
