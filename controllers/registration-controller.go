package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

//register participant ----> for user only

func CreateRegistration(w http.ResponseWriter, r *http.Request) {
	var regisDetails = &models.Registration{}
	utils.ParseBody(r, regisDetails)

	vars := mux.Vars(r)
	eventId, ok := vars["eventId"]
	if !ok {
		http.Error(w, "failed to retrieve event Id from header", http.StatusInternalServerError)
		return
	}

	existingEvent, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Event not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}

	for key := range existingEvent.Parameters {
		if _, ok := regisDetails.Parameters[key]; !ok {
			http.Error(w, fmt.Sprintf("missing parameter '%s' in registration", key), http.StatusBadRequest)
			return
		}
	}

	if existingEvent.Team != regisDetails.Team {
		http.Error(w, fmt.Sprintf("Invaid type of Registration"), http.StatusNotFound)
		return
	}

	if regisDetails.Team == true {
		if regisDetails.TeamName == "" || len(regisDetails.TeamEmails) == 0 {
			http.Error(w, fmt.Sprintf("Team Name or Team Members are empty"), http.StatusBadRequest)
			return
		}

		if len(regisDetails.TeamEmails) > existingEvent.MaxTeamMembers || len(regisDetails.TeamEmails) < existingEvent.MinTeamMembers {
			http.Error(w, fmt.Sprintf("Exceeded maximum number of team members allowed"), http.StatusBadRequest)
			return
		}

		for _, teamEmail := range regisDetails.TeamEmails {
			isRegistered, err := helpers.Helper_IsTeamMemberRegisteredForEvent(existingEvent.Event_ID, teamEmail)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to check team member registration: %v", err), http.StatusInternalServerError)
				return
			}
			if isRegistered {
				http.Error(w, fmt.Sprintf("Participant with email %s is already registered for the event", teamEmail), http.StatusBadRequest)
				return
			}
		}

	} else {
		existingRegistration, err := helpers.Helper_GetRegistrationByEmailAndEvent(regisDetails.Email, existingEvent.Event_ID)
		if err == nil && existingRegistration != nil {
			http.Error(w, fmt.Sprintf("Participant with email %s is already registered for the event", regisDetails.Email), http.StatusBadRequest)
			return
		}
	}

	regisDetails.Event_ID = existingEvent.Event_ID
	regisDetails.Soc_ID = existingEvent.Soc_ID
	regisDetails.Soc_Email = existingEvent.Soc_Email

	result, err := helpers.Helper_CreateRegistration(regisDetails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register participant: %v", err), http.StatusInternalServerError)
		return
	}
	regisDetails.Participant_ID = result.InsertedID.(primitive.ObjectID)

	response, err := json.Marshal(regisDetails)
	if err != nil {
		http.Error(w, "Failed to marshal registration details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//get registration ----> for user only to check its registration within team emails with particular eventId, if no team event then just check email

//getAllRegistrations ----> for admin and superadmin
