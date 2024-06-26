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
	Description    string                 `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt      int64                  `json:"createdAt" bson:"createdAt"`
	StartDate      int64                  `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate        int64                  `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Visibility     string                 `json:"visibility" bson:"visibility"`
	Team           string                 `json:"team" bson:"team"`
	MaxTeamMembers int                    `json:"max_team_members,omitempty" bson:"max_team_members,omitempty"`
	MinTeamMembers int                    `json:"min_team_members,omitempty" bson:"min_team_members,omitempty"`
	Additional     map[string]interface{} `json:"additional,omitempty" bson:"additional,omitempty"`
	Parameters     []Parameters           `json:"parameters,omitempty" bson:"parameters,omitempty"`
	EventType      string                 `json:"event_type,omitempty" bson:"event_type,omitempty"`
	EventMode      string                 `json:"event_mode,omitempty" bson:"event_mode,omitempty"`
	Hashtags       []string               `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	SocialMedia    map[string]interface{} `json:"social_media,omitempty" bson:"social_media,omitempty"`
	Prizes         []Prizes               `json:"prizes,omitempty" bson:"prizes,omitempty"`
	Eligibility    string                 `json:"eligibility,omitempty" bson:"eligibility,omitempty"`
	PhotoGallery   []string               `json:"photo_gallery,omitempty" bson:"photo_gallery,omitempty"`
	Venue          string                 `json:"venue,omitempty" bson:"venue,omitempty"`
	Image          string                 `json:"image" bson:"image"`
	Rounds         []Rounds               `json:"rounds" bson:"rounds"`
	Deadlines      []Deadlines            `json:"deadlines" bson:"deadlines"`
	Register       string                 `json:"register" bson:"register"`
	Report         string                 `json:"report" bson:"report"`
}

type Rounds struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}
type Deadlines struct {
	Title       string `json:"title" bson:"title"`
	Date        int64  `json:"date" bson:"date"`
	Description string `json:"description" bson:"description"`
}
type Prizes struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}
type Parameters struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}
