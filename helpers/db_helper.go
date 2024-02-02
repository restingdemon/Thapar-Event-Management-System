package helpers

import (
	"context"

	"github.com/restingdemon/thaparEvents/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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