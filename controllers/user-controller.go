// auth-controller.go
package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	Image           string             `json:"image" bson:"image"`
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
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, fmt.Sprintf("Failed to check user existence: %s", err), http.StatusInternalServerError)
		return
	}

	// If the user doesn't exist, create a new user in the database
	if existingUser == nil {
		if !strings.Contains(user.Email, "@thapar.edu") && (user.Email != "akshay.garg130803@gmail.com" && user.Email != "jiteshkhurana59@gmail.com"){
			http.Error(w, fmt.Sprintf("Not a Thapar user"), http.StatusInternalServerError)
			return
		}
		if(user.Email == "akshay.garg130803@gmail.com" || user.Email == "jiteshkhurana59@gmail.com"){
			user.Role=utils.SuperAdminRole
		}
		err := createUser(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusInternalServerError)
			return
		}
		existingUser, err = helpers.Helper_GetUserByEmail(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find user after create: %s", err), http.StatusInternalServerError)
		}
	} else {
		existingUser.Image = user.Image
		if(existingUser.Email == "akshay.garg130803@gmail.com" || existingUser.Email == "jiteshkhurana59@gmail.com"){
			existingUser.Role=utils.SuperAdminRole
		}
		err = helpers.Helper_UpdateUser(existingUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update user: %s", err), http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT tokens for the user
	token, refreshToken, err := helpers.GenerateAllTokens(existingUser.Email, existingUser.Name, existingUser.Role, existingUser.ID.Hex())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate tokens: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the user data and tokens as a JSON response
	if existingUser.Role == utils.AdminRole {
		society, err := helpers.Helper_GetSocietyByEmail(existingUser.Email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "Society not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get society: %s", err), http.StatusInternalServerError)
			}
			return
		}
		society.Image = existingUser.Image
		err1 := helpers.Helper_UpdateSoc(society)
		if err1 != nil {
			http.Error(w, fmt.Sprintf("Failed to update soc image: %s", err), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"society": map[string]interface{}{
				"_Sid":              society.Soc_ID.Hex(),
				"_Uid":              society.User_ID.Hex(),
				"email":             society.Email,
				"name":              society.Name,
				"year_of_formation": society.YearOfFormation,
				"role":              society.Role,
				"about":             society.About,
				"members":           society.Members,
				"faculty":           society.Faculty,
				"image":             existingUser.Image,
			},
			"token":         token,
			"refresh_token": refreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	} else {
		response := map[string]interface{}{
			"user": map[string]interface{}{
				"_id":    existingUser.ID.Hex(),
				"email":  existingUser.Email,
				"name":   existingUser.Name,
				"phone":  existingUser.Phone,
				"rollno": existingUser.RollNo,
				"branch": existingUser.Branch,
				"batch":  existingUser.Batch,
				"role":   existingUser.Role,
				"image":  existingUser.Image,
			},
			"token":         token,
			"refresh_token": refreshToken,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}

func createUser(googleUser *GoogleUser) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")
	user := models.User{
		Email: googleUser.Email,
		Name:  googleUser.Name,
		Role:  utils.UserRole,
		Image: googleUser.Image,
	}
	if googleUser.Email == "akshay.garg130803@gmail.com" {
		user.Role = utils.SuperAdminRole
	}
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %s", err)
	}

	return nil
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	if email == "" {
		users, err := helpers.Helper_ListAllUsers()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all users: %s", err), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		user, err := helpers.Helper_GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get User: %s", err), http.StatusInternalServerError)
			}
			return
		}

		response, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Failed to marshal user details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the updated user details
	var updatedUser = &models.User{}
	utils.ParseBody(r, updatedUser)

	// Extract email from the context
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

	// Retrieve the user from the database based on the email
	existingUser, err := helpers.Helper_GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("User not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get user: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Convert existingUser to GoogleUser for the update process
	updatedUser = &models.User{
		ID:     existingUser.ID,
		Email:  existingUser.Email,
		Name:   existingUser.Name,
		Phone:  updatedUser.Phone,
		RollNo: updatedUser.RollNo,
		Branch: updatedUser.Branch,
		Batch:  updatedUser.Batch,
		Role:   existingUser.Role,
		Image:  existingUser.Image,
	}

	// Update the user in the database
	err = helpers.Helper_UpdateUser(updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the updated user data as a JSON response
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"_id":    updatedUser.ID.Hex(),
			"email":  updatedUser.Email,
			"name":   updatedUser.Name,
			"phone":  updatedUser.Phone,
			"rollno": updatedUser.RollNo,
			"branch": updatedUser.Branch,
			"batch":  updatedUser.Batch,
			"role":   updatedUser.Role,
			"image":  updatedUser.Image,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserRegistration(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "failed to retrieve email from header", http.StatusInternalServerError)
		return
	}

	// Get event IDs for the provided email
	eventIDs, err := helpers.Helper_GetEventsByUserEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no documents found, return an empty response with status OK
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get event IDs for user: %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch full event details for each event ID
	var events []models.Event
	for _, eventID := range eventIDs {
		event, err := helpers.Helper_GetEventById(eventID.Hex())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get event details: %s", err), http.StatusInternalServerError)
			return
		}
		events = append(events, *event)
	}

	// Marshal response JSON
	jsonResponse, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
