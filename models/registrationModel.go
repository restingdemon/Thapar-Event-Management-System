package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Registration struct {
	Participant_ID primitive.ObjectID     `json:"_Pid,omitempty" bson:"_id,omitempty"`
	Event_ID       primitive.ObjectID     `json:"_Eid,omitempty" bson:"_eid,omitempty"`
	Soc_ID         primitive.ObjectID     `json:"_Sid" bson:"_sid"`
	Soc_Email      string                 `json:"email" bson:"email"`
	
	Name           string                 `json:"title" bson:"title"`
	Email          string                 `json:"description" bson:"description"`
	RollNo         int64                  `json:"date" bson:"date"`
	PhoneNo        bool                   `json:"visibility" bson:"visibility"`
	
	Team           bool                   `json:"team" bson:"team"`
	TeamName       string                 `json:"team_name" bson:"team_name"`
	TeamEmails     []string               `json:"team_emails" bson:"team_emails"`
	
	Parameters     map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
}
