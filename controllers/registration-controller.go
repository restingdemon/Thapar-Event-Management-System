package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
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
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "failed to retrieve event Id from header", http.StatusInternalServerError)
		return
	}
	regisDetails.Email = email
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

	for key := range regisDetails.Parameters {
		if _, ok := existingEvent.Parameters[key]; !ok {
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
			if err == nil && isRegistered != nil {
				http.Error(w, fmt.Sprintf("Participant with email %s is already registered for the event", regisDetails.Email), http.StatusBadRequest)
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

	err = sendRegistrationConfirmationEmail(regisDetails, existingEvent.Title)
	if err != nil {
		// Log error but don't halt the registration process
		fmt.Printf("Failed to send email notification: %v\n", err)
	}

	response, err := json.Marshal(regisDetails)
	if err != nil {
		http.Error(w, "Failed to marshal registration details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// check registration ----> for user only to check its registration within team emails with particular eventId, if no team event then just check email
func CheckRegistration(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventId, ok := vars["eventId"]
	if !ok {
		http.Error(w, "failed to retrieve event Id from header", http.StatusInternalServerError)
		return
	}
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "failed to retrieve email from header", http.StatusInternalServerError)
		return
	}
	existingEvent, err := helpers.Helper_GetEventById(eventId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprint("Event not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Event: %s", err), http.StatusInternalServerError)
		}
		return
	}
	type Response struct {
		Message              bool                 `json:"message"`
		ExistingRegistration *models.Registration `json:"existing_registration,omitempty"`
	}

	response := &Response{}

	if existingEvent.Team {
		isRegistered, err := helpers.Helper_IsTeamMemberRegisteredForEvent(existingEvent.Event_ID, email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check team member registration: %v", err), http.StatusInternalServerError)
			return
		} else if err == nil && isRegistered != nil {
			response.Message = true
			response.ExistingRegistration = isRegistered
		} else {
			response.Message = false
		}

	} else {
		existingRegistration, err := helpers.Helper_GetRegistrationByEmailAndEvent(email, existingEvent.Event_ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check team member registration: %v", err), http.StatusInternalServerError)
			return
		} else if err == nil && existingRegistration != nil {
			response.Message = true
			response.ExistingRegistration = existingRegistration
		} else {
			response.Message = false
		}
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResponse)
}

// getAllRegistrations ----> for admin and superadmin
func GetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["eventId"]
	if !ok {
		http.Error(w, "Event ID not provided in the request", http.StatusBadRequest)
		return
	}
	userType, ok := r.Context().Value("userType").(string)
	if !ok {
		http.Error(w, "User type not available in the request context", http.StatusInternalServerError)
		return
	}

	var socID primitive.ObjectID
	if userType == utils.AdminRole {
		// Fetch society ID only for admin role
		var ok bool
		socID, ok = r.Context().Value("societyID").(primitive.ObjectID)
		if !ok {
			http.Error(w, "Society ID not available in the request context for admin", http.StatusInternalServerError)
			return
		}
	}

	objectEventID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		http.Error(w, "Invalid Event ID provided", http.StatusBadRequest)
		return
	}
	registrations, err := helpers.Helper_GetAllRegistrations(userType, objectEventID, socID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get registrations: %s", err), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(registrations)
	if err != nil {
		http.Error(w, "Failed to marshal registration details", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func sendRegistrationConfirmationEmail(regisDetails *models.Registration, eventName string) error {

	d := gomail.NewDialer("smtp.gmail.com", 587, "thapar.events.ajak@gmail.com", "Ajak@123")

	var recipientEmails []string
	if regisDetails.Team {
		recipientEmails = regisDetails.TeamEmails
	} else {
		recipientEmails = []string{regisDetails.Email}
	}
	for _, recipientEmail := range recipientEmails {
		m := gomail.NewMessage()
		m.SetHeader("From", "thapar.events.ajak@gmail.com")
		m.SetHeader("To", recipientEmail)
		m.SetHeader("Subject", "Registration Confirmation for "+eventName)
		m.SetBody("text/html", "Dear Participant,<br><br>Thank you for registering for "+eventName+".<br><br>We look forward to seeing you at the event.<br><br>Best Regards,<br>AJAK")

		// Send the email
		if err := d.DialAndSend(m); err != nil {
			fmt.Printf("Failed to send email notification to %s: %v\n", recipientEmail, err)
			return err
		}
	}
	return nil
}
