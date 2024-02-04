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

	if email == "" {
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
	} else {
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
	}
}
