package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

type Society struct {
	Soc_ID          primitive.ObjectID `json:"_Sid,omitempty" bson:"_id,omitempty"`
	User_ID         primitive.ObjectID `json:"_Uid,omitempty" bson:"_uid,omitempty"`
	Email           string             `json:"email" bson:"email"`
	Name            string             `json:"name" bson:"name"`
	YearOfFormation string             `json:"year_of_formation" bson:"year_of_formation"`
	Role            string             `json:"role" bson:"role"`
	About           string             `json:"about" bson:"about"`
}
