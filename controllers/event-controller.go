package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	event.Date = time.Now().Unix()
	event.Visibility = false

	result,err := helpers.Helper_CreateEvent(event)
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
