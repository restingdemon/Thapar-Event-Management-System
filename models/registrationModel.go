package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Registration struct {
	Participant_ID primitive.ObjectID `json:"_Pid,omitempty" bson:"_id,omitempty"`
	Event_ID       primitive.ObjectID `json:"_Eid,omitempty" bson:"_eid,omitempty"`
	Soc_ID         primitive.ObjectID `json:"_Sid" bson:"_sid"`
	Soc_Email      string             `json:"_semail" bson:"_semail"`

	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	RollNo  string `json:"rollno" bson:"rollno"`
	PhoneNo string `json:"phoneno" bson:"phoneno"`

	Team       string   `json:"team" bson:"team"`
	TeamName   string   `json:"team_name,omitempty" bson:"team_name,omitempty"`
	TeamEmails []string `json:"team_emails,omitempty" bson:"team_emails,omitempty"`

	Parameters []Parameters `json:"parameters,omitempty" bson:"parameters,omitempty"`
}
