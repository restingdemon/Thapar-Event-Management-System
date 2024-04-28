package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"sort"
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

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		http.Error(w, "Failed to retrieve email from url or no email in url", http.StatusInternalServerError)
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
	event.CreatedAt = time.Now().Unix()
	event.Visibility = "false"
	event.Team = "false"
	event.Register = "false"
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
	if updatedEvent.Title != "" {
		existingEvent.Title = updatedEvent.Title
	}
	if updatedEvent.Description != "" {
		existingEvent.Description = updatedEvent.Description
	}
	if updatedEvent.StartDate != 0 {
		existingEvent.StartDate = updatedEvent.StartDate
	}
	if updatedEvent.EndDate != 0 {
		existingEvent.EndDate = updatedEvent.EndDate
	}
	if updatedEvent.Additional != nil {
		existingEvent.Additional = updatedEvent.Additional
	}
	if updatedEvent.Parameters != nil {
		existingEvent.Parameters = updatedEvent.Parameters
	}
	if updatedEvent.Team != "" {
		existingEvent.Team = updatedEvent.Team
	}
	if updatedEvent.MaxTeamMembers != 0 {
		existingEvent.MaxTeamMembers = updatedEvent.MinTeamMembers
	}
	if updatedEvent.MinTeamMembers != 0 {
		existingEvent.MinTeamMembers = updatedEvent.MinTeamMembers
	}
	if updatedEvent.EventType != "" {
		existingEvent.EventType = updatedEvent.EventType
	}
	if updatedEvent.EventMode != "" {
		existingEvent.EventMode = updatedEvent.EventMode
	}
	if updatedEvent.Hashtags != nil {
		existingEvent.Hashtags = updatedEvent.Hashtags
	}
	if updatedEvent.SocialMedia != nil {
		existingEvent.SocialMedia = updatedEvent.SocialMedia
	}
	if updatedEvent.Prizes != nil {
		existingEvent.Prizes = updatedEvent.Prizes
	}
	if updatedEvent.Eligibility != "" {
		existingEvent.Eligibility = updatedEvent.Eligibility
	}
	if updatedEvent.Venue != "" {
		existingEvent.Venue = updatedEvent.Venue
	}
	if updatedEvent.Visibility != "" {
		existingEvent.Visibility = updatedEvent.Visibility
	}
	if updatedEvent.Rounds != nil {
		existingEvent.Rounds = updatedEvent.Rounds
	}
	if updatedEvent.Deadlines != nil {
		existingEvent.Deadlines = updatedEvent.Deadlines
	}
	if updatedEvent.Register != "" {
		existingEvent.Register = updatedEvent.Register
	}
	err = helpers.Helper_UpdateEvent(existingEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(existingEvent)
	if err != nil {
		http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Event         models.Event `json:"event" bson:"event"`
		Registartions int64        `json:"registrations" bson:"registrations"`
	}
	queryParams := r.URL.Query()
	eventType := queryParams.Get("event_type")
	event_id := queryParams.Get("eventId")

	if event_id != "" {
		event, err := helpers.Helper_GetEventById(event_id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get event: %s", err), http.StatusInternalServerError)
			}
			return
		}
		registrations, err := helpers.Helper_GetEventDashboard(utils.AdminRole, event.Event_ID, event.Soc_ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get registrations: %s", err), http.StatusInternalServerError)
			return
		}
		var response res
		response.Event = *event
		response.Registartions = registrations
		responses, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responses)
		return
	}
	events, err := helpers.Helper_GetAllEvents(eventType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}
	// Sort events by start date, with the nearest event from today shown first
	sort.Slice(events, func(i, j int) bool {
		today := time.Now().Unix()
		return math.Abs(float64(events[i].StartDate-today)) < math.Abs(float64(events[j].StartDate-today))
	})

	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Failed to marshal events details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func GetNotVisibleEvents(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Event         models.Event `json:"event" bson:"event"`
		Registartions int64        `json:"registrations" bson:"registrations"`
	}
	queryParams := r.URL.Query()
	eventType := queryParams.Get("event_type")
	event_id := queryParams.Get("eventId")

	if event_id != "" {
		event, err := helpers.Helper_GetEventById(event_id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get event: %s", err), http.StatusInternalServerError)
			}
			return
		}
		registrations, err := helpers.Helper_GetEventDashboard(utils.AdminRole, event.Event_ID, event.Soc_ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get registrations: %s", err), http.StatusInternalServerError)
			return
		}
		var response res
		response.Event = *event
		response.Registartions = registrations
		responses, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal event details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responses)
		return
	}
	events, err := helpers.Helper_GetNotVisibleEvents(eventType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get events: %v", err), http.StatusInternalServerError)
		return
	}
	// Sort events by start date, with the nearest event from today shown first
	sort.Slice(events, func(i, j int) bool {
		today := time.Now().Unix()
		return math.Abs(float64(events[i].StartDate-today)) < math.Abs(float64(events[j].StartDate-today))
	})

	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Failed to marshal events details", http.StatusInternalServerError)
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
		CreatedAt:      time.Now().Unix(),
		StartDate:      existingEvent.StartDate,
		EndDate:        existingEvent.EndDate,
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
		PhotoGallery:   existingEvent.PhotoGallery,
		Image:          existingEvent.Image,
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

func UploadPhotos(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get event ID from request params
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	// Retrieve event by ID to ensure it exists
	event, err := helpers.Helper_GetEventById(eventID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		return
	}
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

	if event.Soc_Email != email {
		http.Error(w, "You can only add photos in your own society event", http.StatusForbidden)
		return
	}

	files := r.MultipartForm.File["photos"]
	var photoURLs []string

	// Loop through each file and upload to Cloudinary
	for _, fileHeader := range files {

		// Check file size
		if fileHeader.Size > (10 * 1024 * 1024) { // 10 MB in bytes
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "File size exceeds the limit of 10 MB", http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Upload file to Cloudinary
		uploadResult, err := helpers.UploadToS3(r.Context(), file, "photos")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload photo: %s", err), http.StatusInternalServerError)
			return
		}

		// Append the uploaded photo URL to the photoURLs slice
		photoURLs = append(photoURLs, uploadResult)
	}

	// Append uploaded photo URLs to the event's photo gallery
	event.PhotoGallery = append(event.PhotoGallery, photoURLs...)

	// Update event in the database with the new photo gallery
	err = helpers.Helper_UpdateEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Photos uploaded successfully"))
}

func DeletePhoto(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventID := vars["eventId"]
	queryParams := r.URL.Query()
	photoURL := queryParams.Get("photoURL")

	event, err := helpers.Helper_GetEventById(eventID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		return
	}
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

	if event.Soc_Email != email {
		http.Error(w, "You can only delete photos in your own society event", http.StatusForbidden)
		return
	}

	// Check if photo URL exists in the event's photo gallery
	var photoIndex = -1
	for i, url := range event.PhotoGallery {
		if url == photoURL {
			photoIndex = i
			break
		}
	}

	if photoIndex == -1 {
		http.Error(w, "Photo not found in the photo gallery", http.StatusNotFound)
		return
	}

	// Remove the photo URL from the event's photo gallery
	event.PhotoGallery = append(event.PhotoGallery[:photoIndex], event.PhotoGallery[photoIndex+1:]...)

	// Update event in the database with the modified photo gallery
	err = helpers.Helper_UpdateEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo deleted successfully"))
}

func UploadPoster(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get event ID from request params
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	// Retrieve event by ID to ensure it exists
	event, err := helpers.Helper_GetEventById(eventID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Event not found: %s", err), http.StatusNotFound)
		return
	}
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

	if event.Soc_Email != email {
		http.Error(w, "You can only add photos in your own society event", http.StatusForbidden)
		return
	}
	var files []*multipart.FileHeader
	var fileType string
	if r.MultipartForm.File["report"] != nil {
		files = r.MultipartForm.File["report"]
		fileType = "report"
	} else if r.MultipartForm.File["photos"] != nil {
		files = r.MultipartForm.File["photos"]
		fileType = "photos"
	} else {
		http.Error(w, "No files found in request", http.StatusBadRequest)
		return
	}

	// Loop through each file and upload to Cloudinary
	for _, fileHeader := range files {

		// Check file size
		if fileHeader.Size > (10 * 1024 * 1024) { // 10 MB in bytes
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "File size exceeds the limit of 10 MB", http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Upload file to Cloudinary
		uploadResult, err := helpers.UploadToS3(r.Context(), file, eventID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload photo: %s", err), http.StatusInternalServerError)
			return
		}

		if fileType == "photos" {
			event.Image = uploadResult
		} else if fileType == "report" {
			event.Report = uploadResult
		}

	}
	err = helpers.Helper_UpdateEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Photos uploaded successfully"))
}

func GetEventDashboard(w http.ResponseWriter, r *http.Request) {
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
	type Response struct {
		TotalRegistrations int64 `json:"totalregistrations"`
	}
	registrations, err := helpers.Helper_GetEventDashboard(userType, objectEventID, socID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get registrations: %s", err), http.StatusInternalServerError)
		return
	}
	resp := Response{
		TotalRegistrations: registrations,
	}
	response, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal registration", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
