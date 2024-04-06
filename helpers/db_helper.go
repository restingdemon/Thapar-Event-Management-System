package helpers

import (
	"context"
	"fmt"

	"github.com/restingdemon/thaparEvents/models"
	"github.com/restingdemon/thaparEvents/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// **********USER*************************
func Helper_GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	filter := bson.M{"_id": userID}
	user := &models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Helper_GetUserByEmail(email string) (*models.User, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	filter := bson.M{"email": email}
	user := &models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Helper_ListAllUsers() ([]models.User, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %s", err)
	}

	return users, nil
}

func Helper_UpdateUser(user *models.User) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("users")

	update := bson.M{
		"$set": models.User{
			Email:  user.Email,
			Name:   user.Name,
			Phone:  user.Phone,
			RollNo: user.RollNo,
			Branch: user.Branch,
			Batch:  user.Batch,
			Role:   user.Role,
			Image:  user.Image,
		},
	}

	// Update user in the database based on the email
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

// **********SOCIETY******************
func Helper_CreateSociety(societydetails *models.Society) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("society")

	_, err := collection.InsertOne(context.Background(), societydetails)
	if err != nil {
		return fmt.Errorf("failed to insert user: %s", err)
	}

	return nil
}

func Helper_GetSocietyByEmail(email string) (*models.Society, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("society")

	filter := bson.M{"email": email}
	society := &models.Society{}
	err := collection.FindOne(context.TODO(), filter).Decode(society)
	if err != nil {
		return nil, err
	}

	return society, nil
}

func Helper_ListAllSocieties() ([]models.Society, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("society")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var societies []models.Society
	if err := cursor.All(context.TODO(), &societies); err != nil {
		return nil, fmt.Errorf("failed to decode societies: %s", err)
	}

	return societies, nil
}

func Helper_GetSocietyById(societyId string) (*models.Society, error) {
    collection := models.DB.Database("ThaparEventsDb").Collection("society")
    objID, err := primitive.ObjectIDFromHex(societyId)

    if err != nil {
        return nil, err
    }

    filter := bson.M{"_id": objID}
    society := &models.Society{}
	
    err = collection.FindOne(context.Background(), filter).Decode(society)
    if err != nil {
        return nil, err
    }

    return society, nil
}

func Helper_UpdateSoc(soc *models.Society) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("society")

	update := bson.M{
		"$set": models.Society{
			Soc_ID:          soc.Soc_ID,
			User_ID:         soc.User_ID,
			Email:           soc.Email,
			Name:            soc.Name,
			YearOfFormation: soc.YearOfFormation,
			Role:            soc.Role,
			About:           soc.About,
		},
	}

	// Update user in the database based on the email
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": soc.Email}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

// ***********EVENT***************
func Helper_CreateEvent(event *models.Event) (*mongo.InsertOneResult, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")

	result, err := collection.InsertOne(context.Background(), event)
	if err != nil {
		return result, fmt.Errorf("failed to insert user: %s", err)
	}

	return result, nil
}

func Helper_GetEventById(Event_ID string) (*models.Event, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")

	objID, err := primitive.ObjectIDFromHex(Event_ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	event := &models.Event{}
	err = collection.FindOne(context.TODO(), filter).Decode(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func Helper_UpdateEvent(event *models.Event) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")

	update := bson.M{
		"$set": models.Event{
			Event_ID:       event.Event_ID,
			Soc_ID:         event.Soc_ID,
			User_ID:        event.User_ID,
			Soc_Email:      event.Soc_Email,
			Title:          event.Title,
			Description:    event.Description,
			Date:           event.Date,
			Additional:     event.Additional,
			Parameters:     event.Parameters,
			Team:           event.Team,
			MaxTeamMembers: event.MaxTeamMembers,
			MinTeamMembers: event.MinTeamMembers,
			Visibility:     event.Visibility,
			Soc_Name:       event.Soc_Name,
			EventType:      event.EventType,
			EventMode:      event.EventMode,
			Hashtags:       event.Hashtags,
			SocialMedia:    event.SocialMedia,
			Prizes:         event.Prizes,
			Eligibility:    event.Eligibility,
		},
	}

	// Update user in the database based on the email
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": event.Event_ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

func Helper_GetAllEvents(event_type string) ([]models.Event, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")
	filter := bson.M{"visibility": true}
	if event_type != "" {
		filter = bson.M{"visibility": true, "event_type": event_type}
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []models.Event
	if err := cursor.All(context.TODO(), &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %s", err)
	}

	return events, nil
}

func Helper_DeleteEvent(Event_ID string) error {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")

	objID, err := primitive.ObjectIDFromHex(Event_ID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	return nil
}

// ******REGISTRATIONS************
func Helper_CreateRegistration(registration *models.Registration) (*mongo.InsertOneResult, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("registrations")

	result, err := collection.InsertOne(context.Background(), registration)
	if err != nil {
		return result, fmt.Errorf("failed to register user: %s", err)
	}

	return result, nil
}

func Helper_GetRegistrationByEmailAndEvent(email string, eventID primitive.ObjectID) (*models.Registration, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("registrations")

	filter := bson.M{"email": email, "_eid": eventID}
	registration := &models.Registration{}
	err := collection.FindOne(context.TODO(), filter).Decode(registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no registration found
		}
		return nil, err
	}

	return registration, nil
}

func Helper_IsTeamMemberRegisteredForEvent(eventId primitive.ObjectID, teamEmail string) (*models.Registration, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("registrations")
	filter := bson.M{
		"_eid":        eventId,
		"team_emails": bson.M{"$in": []string{teamEmail}},
	}
	registration := &models.Registration{}
	err := collection.FindOne(context.TODO(), filter).Decode(registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no registration found
		}
		return nil, err
	}

	return registration, nil
}

func Helper_GetAllRegistrations(userType string, eventID, Soc_ID primitive.ObjectID) ([]models.Registration, error) {
	var registrations []models.Registration
	collection := models.DB.Database("ThaparEventsDb").Collection("registrations")

	filter := bson.M{"_eid": eventID}

	// If the user is an admin filtering by society ID
	if userType == utils.AdminRole {
		filter["_sid"] = Soc_ID
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var registration models.Registration
		if err := cursor.Decode(&registration); err != nil {
			return nil, err
		}
		registrations = append(registrations, registration)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}
