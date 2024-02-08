package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Event_ID    primitive.ObjectID     `json:"_Eid,omitempty" bson:"_id,omitempty"`
	Soc_ID      primitive.ObjectID     `json:"_Sid" bson:"_sid"`
	User_ID     primitive.ObjectID     `json:"_Uid" bson:"_uid"`
	Soc_Email   string                 `json:"email" bson:"email"`
	Title       string                 `json:"title" bson:"title"`
	Description string                 `json:"description" bson:"description"`
	Date        int64                  `json:"date" bson:"date"`
	Visibility  bool                   `json:"visibility" bson:"visibility"`
	Additional  map[string]interface{} `json:"additional,omitempty" bson:"additional,omitempty"`
}
