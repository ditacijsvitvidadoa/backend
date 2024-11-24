package requests

import (
	"context"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateUserCart(client *mongo.Client, userID primitive.ObjectID, productID string) (int64, error) {
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

func AddProductToFavourites(client *mongo.Client, userID primitive.ObjectID, productID string) (int64, error) {
	filter := bson.M{"_id": userID}
	update := bson.M{"$addToSet": bson.M{"Favourites": productID}}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func RemoveProductFromFavourites(client *mongo.Client, userID primitive.ObjectID, productID string) (int64, error) {
	filter := bson.M{"_id": userID}
	update := bson.M{"$pull": bson.M{"Favourites": productID}}

	return storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
}

func UpdateUserProfileField(client *mongo.Client, userID primitive.ObjectID, fieldPath string, newValue interface{}) error {
	filter := bson.M{"_id": userID}

	update := bson.M{"$set": bson.M{fieldPath: newValue}}

	modifiedCount, err := storage.GeneralUpdate[entities.User](client, storage.Users, filter, update)
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

func UpdateCartProductCount(client *mongo.Client, userID primitive.ObjectID, filter bson.M, newCount int) (int64, error) {
	userFilter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"Cart.$[elem].Count": newCount,
		},
	}

	arrayFilters := []interface{}{
		bson.M{"elem.Id": filter["Id"]},
	}

	result, err := client.Database(storage.MongoDBName).Collection(storage.Users).UpdateOne(
		context.Background(),
		userFilter,
		update,
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}),
	)

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func AddOrUpdatePostalServiceInfo(client *mongo.Client, userID primitive.ObjectID, postalServiceInfo entities.PostalServiceInfo) error {
	collectionName := storage.Users
	filter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"PostalService.PostalType":    postalServiceInfo.PostalType,
			"PostalService.City":          postalServiceInfo.City,
			"PostalService.CityRef":       postalServiceInfo.CityRef,
			"PostalService.ReceivingType": postalServiceInfo.ReceivingType,
			"PostalService.PostalInfo":    postalServiceInfo.PostalInfo,
		},
	}

	modifiedCount, err := storage.GeneralUpdate[interface{}](client, collectionName, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update postal service info: %w", err)
	}

	if modifiedCount == 0 {
		newDocument := bson.M{
			"_id":               userID,
			"PostalServiceInfo": postalServiceInfo,
		}

		_, err := storage.GeneralInsert[interface{}](client, collectionName, newDocument)
		if err != nil {
			return fmt.Errorf("failed to insert new postal service info: %w", err)
		}
	}

	return nil
}

func UpdateCounterValue(client *mongo.Client, counterID string, newValue int64) error {
	filter := bson.M{"_id": counterID}
	update := bson.M{"$set": bson.M{"sequence_value": newValue}}

	_, err := storage.GeneralUpdate[interface{}](client, storage.Counters, filter, update)
	return err
}

func UpdateOrCreateProductAnalytics(client *mongo.Client, productId string, field string, increment int) error {
	filter := bson.M{"ProductId": productId}
	existingAnalytics, err := storage.GeneralFind[entities.ProductActivity](client, storage.ProductsActivities, storage.GeneralQueryOptions{
		Filter: filter,
	}, false)
	if err != nil {
		return fmt.Errorf("error fetching product analytics: %w", err)
	}

	if len(existingAnalytics) > 0 {
		update := bson.M{
			"$inc": bson.M{field: increment},
		}
		modifiedCount, err := storage.GeneralUpdate[any](client, storage.ProductsActivities, filter, update)
		if err != nil {
			return fmt.Errorf("error updating analytics: %w", err)
		}
		if modifiedCount == 0 {
			return fmt.Errorf("no documents updated")
		}
		return nil
	}

	newDocument := bson.M{
		"ProductId":   productId,
		"Sales":       0,
		"Clicks":      0,
		"AddedToCart": 0,
		"Favourites":  0,
		field:         increment,
	}

	_, err = storage.GeneralInsert(client, storage.ProductsActivities, newDocument)
	if err != nil {
		return fmt.Errorf("error inserting new analytics document: %w", err)
	}

	return nil
}
