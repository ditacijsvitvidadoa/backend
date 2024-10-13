package utils

import (
	"context"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNextUserID(client *mongo.Client) (int, error) {
	length, err := getLengthFromCollection(client, storage.Users)
	if err != nil {
		return 0, err
	}

	return length + 1, nil
}

func GetNextProductID(client *mongo.Client) (int, error) {
	length, err := getLengthFromCollection(client, storage.Products)
	if err != nil {
		return 0, err
	}

	return length + 1, nil
}

func getLengthFromCollection(client *mongo.Client, collectionName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), storage.MongoDBTimeout)
	defer cancel()

	collection := client.Database(storage.MongoDBName).Collection(collectionName)

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
