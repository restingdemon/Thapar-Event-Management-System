// auth-controller.go
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/restingdemon/thaparEvents/helpers"
	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the token
	var tokenData map[string]string
	utils.ParseBody(r, &tokenData)

	accessToken := tokenData["accessToken"]
	if accessToken == "" {
		http.Error(w, "Missing access token", http.StatusBadRequest)
		return
	}

	// Fetch user info from Google using the access token
	googleUser, err := fetchGoogleUserInfo(accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch user info from Google: %s", err), http.StatusInternalServerError)
		return
	}

	// Check if the user already exists in the database based on their email
	existingUser, err := getUserByEmail(googleUser.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check user existence: %s", err), http.StatusInternalServerError)
		return
	}

	// If the user doesn't exist, create a new user in the database
	if existingUser == nil {
		err := createUser(googleUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusInternalServerError)
			return
		}
		existingUser, err = getUserByEmail(googleUser.Email)
		if err != nil {
			http.Error(w,fmt.Sprintf("Failed to find user after create: %s", err), http.StatusInternalServerError)
		}
	}

	// Generate JWT tokens for the user
	token, refreshToken, err := helpers.GenerateAllTokens(googleUser.Email, googleUser.Name, "user", existingUser.ID.Hex())
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
		},
		"token":         token,
		"refresh_token": refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchGoogleUserInfo(accessToken string) (*GoogleUser, error) {
	// Make an HTTP request to Google User Info API
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info from Google: %s", resp.Status)
	}

	// Decode the response body into a GoogleUser struct
	var googleUser GoogleUser
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func getUserByEmail(email string) (*models.User, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	filter := bson.M{"email": email}
	user := &models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func createUser(googleUser *GoogleUser) error {
	collection := models.DB.Database("your-database").Collection("users")

	objectID, err := primitive.ObjectIDFromHex(googleUser.ID)
	if err != nil {
		return fmt.Errorf("failed to convert ID to ObjectID: %s", err)
	}
	user := models.User{
		ID:              objectID,
		Email:           googleUser.Email,
		Name:            googleUser.Name,
	}

	_, err = collection.InsertOne(context.TODO(), user)
	return err
}
