package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongoDBName    = "DutyachiySvitDB"
	MongoDBTimeout = 10 * time.Second

	Users         = "Users"
	Products      = "Products"
	Counters      = "Counters"
	Orders        = "Orders"
	ArchiveOrders = "ArchiveOrders"
)

type GeneralQueryOptions struct {
	Filter     bson.M
	Projection bson.M
	PageNum    *int
	PageSize   *int
	Sort       *options.FindOptions
}

func GeneralFind[T any](client *mongo.Client, collectionName string, opts GeneralQueryOptions, singleResult bool) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()

	collection := client.Database(MongoDBName).Collection(collectionName)

	var results []T

	if singleResult {
		var result T
		err := collection.FindOne(ctx, opts.Filter, options.FindOne().SetProjection(opts.Projection)).Decode(&result)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, errors.New("document not found")
			}
			return nil, err
		}
		return []T{result}, nil
	}

	findOptions := options.Find().SetProjection(opts.Projection)

	if opts.Sort != nil {
		findOptions.SetSort(opts.Sort.Sort)
	}

	// Pagination handling
	if opts.PageNum != nil && opts.PageSize != nil {
		skip := int64((*opts.PageNum - 1) * *opts.PageSize)
		limit := int64(*opts.PageSize)
		findOptions.SetSkip(skip).SetLimit(limit)
	}

	cursor, err := collection.Find(ctx, opts.Filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GeneralUpdate[T any](client *mongo.Client, collectionName string, filter bson.M, update bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()

	collection := client.Database(MongoDBName).Collection(collectionName)

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func GeneralDelete[T any](client *mongo.Client, collectionName string, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()

	collection := client.Database(MongoDBName).Collection(collectionName)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func GeneralInsert[T any](client *mongo.Client, collectionName string, document T) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()

	collection := client.Database(MongoDBName).Collection(collectionName)

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
