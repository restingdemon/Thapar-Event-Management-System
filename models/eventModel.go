package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Event_ID       primitive.ObjectID     `json:"_Eid,omitempty" bson:"_id,omitempty"`
	Soc_ID         primitive.ObjectID     `json:"_Sid" bson:"_sid"`
	User_ID        primitive.ObjectID     `json:"_Uid" bson:"_uid"`
	Soc_Email      string                 `json:"email" bson:"email"`
	Soc_Name       string                 `json:"soc_name" bson:"soc_name"`
	Title          string                 `json:"title" bson:"title"`
	Description    string                 `json:"description" bson:"description"`
	Date           int64                  `json:"date" bson:"date"`
	Visibility     bool                   `json:"visibility" bson:"visibility"`
	Team           bool                   `json:"team" bson:"team"`
	MaxTeamMembers int                    `json:"max_team_members,omitempty" bson:"max_team_members,omitempty"`
	MinTeamMembers int                    `json:"min_team_members,omitempty" bson:"min_team_members,omitempty"`
	Additional     map[string]interface{} `json:"additional,omitempty" bson:"additional,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
	EventType      string                 `json:"event_type" bson:"event_type"`
	EventMode      string                 `json:"event_mode" bson:"event_mode"`
	Hashtags       []string               `json:"hashtags" bson:"hashtags"`
	SocialMedia    map[string]interface{} `json:"social_media,omitempty" bson:"social_media,omitempty"`
	Prizes         map[string]interface{} `json:"prizes" bson:"prizes"`
	Eligibility    string                 `json:"eligibility" bson:"eligibility"`
	PhotoGallery   []string               `json:"photo_gallery,omitempty" bson:"photo_gallery,omitempty"`
}
