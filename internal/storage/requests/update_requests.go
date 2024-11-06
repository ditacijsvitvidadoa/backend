package requests

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserCart(client *mongo.Client, userID primitive.ObjectID, productID int) (int64, error) {
	filter := bson.M{"_id": userID}

	update := bson.M{
		"$pull": bson.M{
			"Cart": bson.M{"Id": productID},
		},
	}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func AddProductToCart(client *mongo.Client, userID primitive.ObjectID, update bson.M) (int64, error) {
	filter := bson.M{"_id": userID}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func AddProductToFavourites(client *mongo.Client, userID primitive.ObjectID, productID int) (int64, error) {
	filter := bson.M{"_id": userID}
	update := bson.M{"$addToSet": bson.M{"Favourites": productID}}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func RemoveProductFromFavourites(client *mongo.Client, userID primitive.ObjectID, productID int) (int64, error) {
	filter := bson.M{"_id": userID}
	update := bson.M{"$pull": bson.M{"Favourites": productID}}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func UpdateUserProfileField(client *mongo.Client, userID primitive.ObjectID, fieldPath string, newValue interface{}) error {
	filter := bson.M{"_id": userID}

	update := bson.M{"$set": bson.M{fieldPath: newValue}}

	modifiedCount, err := storage.GeneralUpdate[any](client, storage.Users, filter, update)
	if err != nil {
		return fmt.Errorf("could not update field %s: %v", fieldPath, err)
	}

	if modifiedCount == 0 {
		return fmt.Errorf("user with ID %s not found", userID.Hex())
	}

	fmt.Printf("Successfully updated field '%s' for user: %s\n", fieldPath, userID.Hex())
	return nil
}

func UpdateOrderStatus(client *mongo.Client, userID primitive.ObjectID, newStatus int) (int64, error) {
	filter := bson.M{"_id": userID}

	update := bson.M{"$set": bson.M{"Status": newStatus}}

	return storage.GeneralUpdate[any](client, storage.Orders, filter, update)
}
