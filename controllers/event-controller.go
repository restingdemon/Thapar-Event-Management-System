package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event = &models.Event{}
	utils.ParseBody(r, event)

	emailValue := r.Context().Value("email")
	if emailValue == nil {
		http.Error(w, "Email not found in context", http.StatusInternalServerError)
		return
	}

	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}

	soc_details, err := helpers.Helper_GetSocietyByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Society not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get society: %s", err), http.StatusInternalServerError)
		}
		return
	}

	event.Soc_ID = soc_details.Soc_ID
	event.User_ID = soc_details.User_ID
	event.Soc_Email = soc_details.Email
	event.Soc_Name = soc_details.Name
	event.Date = time.Now().Unix()
	event.Visibility = false

	result, err := helpers.Helper_CreateEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
		return
	}
	event.Event_ID = result.InsertedID.(primitive.ObjectID)
	response, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//event update

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var updatedEvent = &models.Event{}
	utils.ParseBody(r, updatedEvent)
	eventIdVal := r.Context().Value("eventId")
	if eventIdVal == nil {
		http.Error(w, "Event Id not found in context", http.StatusInternalServerError)
		return
	}
	roleValue := r.Context().Value("role")
	if roleValue == nil {
		http.Error(w, "Role not found in context", http.StatusInternalServerError)
		return
	}

	eventId, ok := eventIdVal.(string)
	if !ok {
		http.Error(w, "Failed to retrieve Event Id from context", http.StatusInternalServerError)
		return
	}
	role, ok := roleValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve role from context", http.StatusInternalServerError)
		return
	}

	existingEvent, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}
	if role == utils.AdminRole {
		emailVal := r.Context().Value("email")
		if emailVal == nil {
			http.Error(w, "Email not found in context", http.StatusInternalServerError)
			return
		}

		email, ok := emailVal.(string)
		if !ok {
			http.Error(w, "Failed to retrieve Email from context", http.StatusInternalServerError)
			return
		}

		if existingEvent.Soc_Email != email {
			http.Error(w, "You can only update event in your own society", http.StatusForbidden)
			return
		}

	}

	updatedEvent = &models.Event{
		Soc_ID:         existingEvent.Soc_ID,
		Event_ID:       existingEvent.Event_ID,
		User_ID:        existingEvent.User_ID,
		Soc_Email:      existingEvent.Soc_Email,
		Soc_Name:       existingEvent.Soc_Name,
		Visibility:     existingEvent.Visibility,
		Title:          updatedEvent.Title,
		Description:    updatedEvent.Description,
		Date:           updatedEvent.Date,
		Additional:     updatedEvent.Additional,
		Parameters:     updatedEvent.Parameters,
		Team:           updatedEvent.Team,
		MaxTeamMembers: updatedEvent.MaxTeamMembers,
		MinTeamMembers: updatedEvent.MinTeamMembers,
		EventType:      updatedEvent.EventType,
		EventMode:      updatedEvent.EventMode,
		Hashtags:       updatedEvent.Hashtags,
		SocialMedia:    updatedEvent.SocialMedia,
		Prizes:         updatedEvent.Prizes,
		Eligibility:    updatedEvent.Eligibility,
	}
	err = helpers.Helper_UpdateEvent(updatedEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedEvent)
	if err != nil {
		http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	event_type := queryParams.Get("event_type")
	events, err := helpers.Helper_GetAllEvents(event_type)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Failed to marshal events details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetEventById(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["eventId"]
	if eventId == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	event, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}

	response, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateVisibility(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventId, ok := vars["eventId"]
	var updatedEvent = &models.Event{}
	utils.ParseBody(r, updatedEvent)
	if !ok {
		http.Error(w, fmt.Sprintf("Failed to get events: %v", ok), http.StatusBadRequest)
		return
	}
	existingEvent, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}

	updatedEvent = &models.Event{
		Soc_ID:         existingEvent.Soc_ID,
		Event_ID:       existingEvent.Event_ID,
		User_ID:        existingEvent.User_ID,
		Soc_Email:      existingEvent.Soc_Email,
		Title:          existingEvent.Title,
		Description:    existingEvent.Description,
		Date:           existingEvent.Date,
		Additional:     existingEvent.Additional,
		Parameters:     existingEvent.Parameters,
		Team:           existingEvent.Team,
		MaxTeamMembers: existingEvent.MaxTeamMembers,
		MinTeamMembers: existingEvent.MinTeamMembers,
		Soc_Name:       existingEvent.Soc_Name,
		EventType:      existingEvent.EventType,
		EventMode:      existingEvent.EventMode,
		Hashtags:       existingEvent.Hashtags,
		SocialMedia:    existingEvent.SocialMedia,
		Prizes:         existingEvent.Prizes,
		Eligibility:    existingEvent.Eligibility,
		Visibility:     updatedEvent.Visibility,
	}
	err = helpers.Helper_UpdateEvent(updatedEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Visibility set"))
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventIdVal := r.Context().Value("eventId")
	if eventIdVal == nil {
		http.Error(w, "Event Id not found in context", http.StatusInternalServerError)
		return
	}
	roleValue := r.Context().Value("role")
	if roleValue == nil {
		http.Error(w, "Role not found in context", http.StatusInternalServerError)
		return
	}

	eventId, ok := eventIdVal.(string)
	if !ok {
		http.Error(w, "Failed to retrieve Event Id from context", http.StatusInternalServerError)
		return
	}
	role, ok := roleValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve role from context", http.StatusInternalServerError)
		return
	}

	existingEvent, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}
	if role == utils.AdminRole {
		emailVal := r.Context().Value("email")
		if emailVal == nil {
			http.Error(w, "Email not found in context", http.StatusInternalServerError)
			return
		}

		email, ok := emailVal.(string)
		if !ok {
			http.Error(w, "Failed to retrieve Email from context", http.StatusInternalServerError)
			return
		}

		if existingEvent.Soc_Email != email {
			http.Error(w, "You can only delete event in your own society", http.StatusForbidden)
			return
		}

	}

	err = helpers.Helper_DeleteEvent(eventId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete event: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event Deleted Successfully"))
}
