package storage

import (
	"context"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Get(client *mongo.Client, collectionName string, filters bson.M, pageNum, pageSize int) ([]entities.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println(pageSize)
	fmt.Println(pageNum)

	collection := client.Database("DutyachiySvitDB").Collection(collectionName)

	var results []entities.Product

	skip := (pageNum - 1) * pageSize

	cursor, err := collection.Find(ctx, filters, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
