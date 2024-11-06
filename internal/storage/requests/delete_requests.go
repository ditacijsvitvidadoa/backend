package requests

import (
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteByObjectID(client *mongo.Client, collectionName string, id primitive.ObjectID) (int64, error) {
	filter := bson.M{"_id": id}
	return storage.GeneralDelete[any](client, collectionName, filter)
}

func ArchiveOrder(client *mongo.Client, id primitive.ObjectID) (int64, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{"_id": id},
	}
	order, err := storage.GeneralFind[map[string]interface{}](client, storage.Orders, opts, true)
	if err != nil {
		return 0, err
	}

	if len(order) == 0 {
		return 0, mongo.ErrNoDocuments
	}

	deletedCount, err := storage.GeneralDelete[map[string]interface{}](client, storage.Orders, opts.Filter)
	if err != nil {
		return 0, err
	}

	_, err = storage.GeneralInsert(client, storage.ArchiveOrders, order[0])
	if err != nil {
		_, rollbackErr := storage.GeneralInsert(client, storage.Orders, order[0])
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}

	return deletedCount, nil
}

func RefreshOrder(client *mongo.Client, id primitive.ObjectID) (int64, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{"_id": id},
	}
	order, err := storage.GeneralFind[map[string]interface{}](client, storage.ArchiveOrders, opts, true)
	if err != nil {
		return 0, err
	}

	if len(order) == 0 {
		return 0, mongo.ErrNoDocuments
	}

	deletedCount, err := storage.GeneralDelete[map[string]interface{}](client, storage.ArchiveOrders, opts.Filter)
	if err != nil {
		return 0, err
	}

	_, err = storage.GeneralInsert(client, storage.Orders, order[0])
	if err != nil {
		_, rollbackErr := storage.GeneralInsert(client, storage.ArchiveOrders, order[0])
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}

	return deletedCount, nil
}
