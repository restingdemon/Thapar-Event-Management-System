package helpers

import (
	"context"
	"fmt"

	"github.com/restingdemon/thaparEvents/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
			Email:           user.Email,
			Name:            user.Name,
			Phone:           user.Phone,
			RollNo:          user.RollNo,
			Branch:          user.Branch,
			YearOfAdmission: user.YearOfAdmission,
			Role:            user.Role,
		},
	}

	// Update user in the database based on the email
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}

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

func Helper_CreateEvent(event *models.Event) (*mongo.InsertOneResult, error) {
	collection := models.DB.Database("ThaparEventsDb").Collection("event")

	result, err := collection.InsertOne(context.Background(), event)
	if err != nil {
		return result, fmt.Errorf("failed to insert user: %s", err)
	}

	return result, nil
}

// func Helper_GetEventById(Event_ID string) error {
// 	collection := models.DB.Database("ThaparEventsDb").Collection("event")

// 	_, err := collection.InsertOne(context.Background(), event)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert user: %s", err)
// 	}

// 	return nil
// }
