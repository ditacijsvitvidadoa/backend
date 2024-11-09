package requests

import (
	"context"
	"errors"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	password2 "github.com/vzglad-smerti/password_hash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProducts(client *mongo.Client, filters bson.M, pageNum, pageSize *int) ([]entities.Product, error) {
	opts := storage.GeneralQueryOptions{
		Filter:   filters,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
	return storage.GeneralFind[entities.Product](client, storage.Products, opts, false)
}

func GetAll(client *mongo.Client, CollectionName string) ([]map[string]interface{}, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{},
	}

	return storage.GeneralFind[map[string]interface{}](client, CollectionName, opts, false)
}

func LogInAccount(client *mongo.Client, email, password string) (primitive.ObjectID, error) {
	filter := bson.M{"Email": email}

	var result struct {
		ID       primitive.ObjectID `bson:"_id"`
		Password string             `bson:"Password"`
	}

	err := client.Database(storage.MongoDBName).Collection(storage.Users).FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, fmt.Errorf("user not found")
		}
		return primitive.NilObjectID, fmt.Errorf("error querying user: %v", err)
	}

	isValid, err := password2.Verify(result.Password, password)
	if err != nil {
		fmt.Printf("Error verifying password: %v\n", err)
		return primitive.NilObjectID, err
	}

	if !isValid {
		return primitive.NilObjectID, fmt.Errorf("invalid password")
	}

	return result.ID, nil
}

func GetUserByID(client *mongo.Client, userID primitive.ObjectID) (entities.User, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{"_id": userID},
	}

	results, err := storage.GeneralFind[entities.User](client, storage.Users, opts, true)
	if err != nil {
		return entities.User{}, err
	}

	if len(results) == 0 {
		return entities.User{}, errors.New("user not found")
	}

	return results[0], nil
}

func GetCartByUserID(client *mongo.Client, userID primitive.ObjectID) ([]entities.CartItem, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{"_id": userID},
	}

	results, err := storage.GeneralFind[entities.User](client, storage.Users, opts, true)
	if err != nil {
		return nil, err
	}

	return results[0].Cart, nil
}

func GetFavouritesByUserID(client *mongo.Client, userID primitive.ObjectID) ([]string, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{"_id": userID},
	}

	results, err := storage.GeneralFind[entities.User](client, storage.Users, opts, true)
	if err != nil {
		return nil, err
	}

	return results[0].Favourites, nil
}

func IsEmailExists(client *mongo.Client, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), storage.MongoDBTimeout)
	defer cancel()

	collection := client.Database(storage.MongoDBName).Collection(storage.Users)

	var existingUser entities.User
	err := collection.FindOne(ctx, bson.M{"Email": email}).Decode(&existingUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
