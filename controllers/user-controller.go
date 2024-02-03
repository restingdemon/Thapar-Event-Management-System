// auth-controller.go
package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoogleUser struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email           string             `json:"email" bson:"email"`
	Name            string             `json:"name" bson:"name"`
	Phone           string             `json:"phone" bson:"phone"`
	RollNo          string             `json:"rollno" bson:"rollno"`
	Branch          string             `json:"branch" bson:"branch"`
	YearOfAdmission string             `json:"year_of_admission" bson:"year_of_admission"`
	Role            string             `json:"role" bson:"role"`
	Token           string             `json:"token"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the token
	var user = &GoogleUser{}
	utils.ParseBody(r, user)

	if user.Email == "" {
		http.Error(w, fmt.Sprintf("No email provided"), http.StatusBadRequest)
		return
	}

	if !utils.IsloginValid(user.Email, user.Token) {
		http.Error(w, fmt.Sprintf("User token not valid"), http.StatusBadRequest)
		return
	}
	// Check if the user already exists in the database based on their email
	existingUser, err := helpers.Helper_GetUserByEmail(user.Email)
	if err!= nil && !errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, fmt.Sprintf("Failed to check user existence: %s", err), http.StatusInternalServerError)
		return
	}

	// If the user doesn't exist, create a new user in the database
	if existingUser == nil {
		err := createUser(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusInternalServerError)
			return
		}
		existingUser, err = helpers.Helper_GetUserByEmail(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find user after create: %s", err), http.StatusInternalServerError)
		}
	}

	// Generate JWT tokens for the user
	token, refreshToken, err := helpers.GenerateAllTokens(existingUser.Email, existingUser.Name, existingUser.Role, existingUser.ID.Hex())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate tokens: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the user data and tokens as a JSON response
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"_id":               existingUser.ID.Hex(),
			"email":             existingUser.Email,
			"name":              existingUser.Name,
			"phone":             existingUser.Phone,
			"rollno":            existingUser.RollNo,
			"branch":            existingUser.Branch,
			"year_of_admission": existingUser.YearOfAdmission,
			"role":              existingUser.Role,
		},
		"token":         token,
		"refresh_token": refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUser(googleUser *GoogleUser) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	user := models.User{
		Email: googleUser.Email,
		Name:  googleUser.Name,
		Role:  "user",
	}

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %s", err)
	}

	return nil
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	// Retrieve the user from the database based on the ID
	user, err := helpers.Helper_GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("User not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get user: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Return the user data as a JSON response
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"_id":               user.ID.Hex(),
			"email":             user.Email,
			"name":              user.Name,
			"phone":             user.Phone,
			"rollno":            user.RollNo,
			"branch":            user.Branch,
			"year_of_admission": user.YearOfAdmission,
			"role":              user.Role,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
