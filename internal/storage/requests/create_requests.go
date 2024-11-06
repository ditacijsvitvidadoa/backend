package requests

import (
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateNewUser(client *mongo.Client, newUser entities.User) (primitive.ObjectID, error) {
	return storage.GeneralInsert(client, storage.Users, newUser)
}

func CreateNewOrder(client *mongo.Client, order entities.Order) (primitive.ObjectID, error) {
	return storage.GeneralInsert(client, storage.Orders, order)
}

func CreateNewProduct(client *mongo.Client, product entities.Product) (primitive.ObjectID, error) {
	return storage.GeneralInsert(client, storage.Products, product)
}
