package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterSociety(w http.ResponseWriter, r *http.Request) {

	var societyDetails = &models.Society{}
	utils.ParseBody(r, societyDetails)

	if societyDetails.Email == "" || societyDetails.Role == "" {
		http.Error(w, fmt.Sprintf("No email or role provided"), http.StatusBadRequest)
		return
	}
	existingSoc, err1 := helpers.Helper_GetSocietyByEmail(societyDetails.Email)
	if err1 != nil {
		if errors.Is(err1, mongo.ErrNoDocuments) {
		} else {
			http.Error(w, fmt.Sprintf("error is %v", err1), http.StatusBadRequest)
			return
		}
	}
	if err1 != mongo.ErrNoDocuments && existingSoc != nil {
		http.Error(w, fmt.Sprintf("Soc already exists"), http.StatusBadRequest)
		return
	}
	societyDetails.Role = utils.AdminRole
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := models.DB.StartSession()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start session: %v", err), http.StatusInternalServerError)
		return
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start transaction: %v", err), http.StatusInternalServerError)
		return
	}

	user, err := helpers.Helper_GetUserByEmail(societyDetails.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			session.AbortTransaction(ctx)
			http.Error(w, fmt.Sprintf("User not found"), http.StatusNotFound)
		} else {
			session.AbortTransaction(ctx)
			http.Error(w, fmt.Sprintf("Failed to get user: %s", err), http.StatusInternalServerError)
		}
		return
	}

	user.Role = societyDetails.Role
	err = helpers.Helper_UpdateUser(user)
	if err != nil {
		session.AbortTransaction(ctx)
		http.Error(w, fmt.Sprintf("Failed to update user role: %v", err), http.StatusInternalServerError)
		return
	}

	societyDetails.User_ID = user.ID
	societyDetails.Name = user.Name
	err = helpers.Helper_CreateSociety(societyDetails)
	if err != nil {
		session.AbortTransaction(ctx)
		http.Error(w, fmt.Sprintf("Failed to create society: %v", err), http.StatusInternalServerError)
		return
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Society registered successfully"))
}

func GetSocietyDetails(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	email := queryParams.Get("email")
	society_id := queryParams.Get("societyId")
	if email != "" {
		society, err := helpers.Helper_GetSocietyByEmail(email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "Society not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get society: %s", err), http.StatusInternalServerError)
			}
			return
		}
		response, err := json.Marshal(society)
		if err != nil {
			http.Error(w, "Failed to marshal society details", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else if society_id != "" {
		society, err := helpers.Helper_GetSocietyById(society_id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "Society not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get society: %s", err), http.StatusInternalServerError)
			}
			return
		}
		response, err := json.Marshal(society)
		if err != nil {
			http.Error(w, "Failed to marshal society details", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		societies, err := helpers.Helper_ListAllSocieties()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all societies: %s", err), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(societies)
		if err != nil {
			http.Error(w, "Failed to marshal societies", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func UpdateSociety(w http.ResponseWriter, r *http.Request) {
	var updatedSociety = &models.Society{}
	utils.ParseBody(r, updatedSociety)
	emailValue := r.Context().Value("email")
	if emailValue == nil {
		http.Error(w, "Email not found in context", http.StatusInternalServerError)
		return
	}
	roleValue := r.Context().Value("role")
	if roleValue == nil {
		http.Error(w, "Role not found in context", http.StatusInternalServerError)
		return
	}

	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	role, ok := roleValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}

	if role == utils.AdminRole || role == utils.SuperAdminRole {
		existingSoc, err := helpers.Helper_GetSocietyByEmail(email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, fmt.Sprintf("Soc not found"), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get Soc: %s", err), http.StatusInternalServerError)
			}
			return
		}
		updatedSociety = &models.Society{
			Soc_ID:          existingSoc.Soc_ID,
			User_ID:         existingSoc.User_ID,
			Email:           existingSoc.Email,
			Role:            existingSoc.Role,
			Image:           existingSoc.Image,
			Name:            updatedSociety.Name,
			YearOfFormation: updatedSociety.YearOfFormation,
			About:           updatedSociety.About,
			Members:         updatedSociety.Members,
			Faculty:         updatedSociety.Faculty,
		}
		err = helpers.Helper_UpdateSoc(updatedSociety)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update soc: %s", err), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(updatedSociety)
		if err != nil {
			http.Error(w, "Failed to marshal society details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func GetSocEvents(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	soc_email := queryParams.Get("soc_email")
	if soc_email == "" {
		http.Error(w, fmt.Sprintln("society not provided"), http.StatusNotFound)
		return
	}
	events, err := helpers.Helper_GetSocEvents(soc_email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintln("Not events found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get events: %s", err), http.StatusInternalServerError)
		}
		return
	}

	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, fmt.Sprintln("Unable to marshal events into json"), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func GetSocDashboard(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	if email == "" {
		http.Error(w, "Email not found in context", http.StatusInternalServerError)
		return
	}

	_, err := helpers.Helper_GetSocietyByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("Soc not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get Soc: %s", err), http.StatusInternalServerError)
		}
		return
	}
	type Response struct {
		TotalEvents    int64 `json:"totalEvents"`
		UpcomingEvents int64 `json:"upcomingEvents"`
		TeamMembers    int64 `json:"teamMembers"`
	}

	totalEvents, upcomingEvents, members, ok := helpers.Helper_GetSocDashboard(email)
	if ok != nil {
		http.Error(w, fmt.Sprintf("Failed to get dashboard"), http.StatusInternalServerError)
		return
	}
	resp := Response{
		TotalEvents:    totalEvents,
		UpcomingEvents: upcomingEvents,
		TeamMembers:    members,
	}

	output, err1 := json.Marshal(resp)
	if err1 != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
