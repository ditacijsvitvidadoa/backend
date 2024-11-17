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
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetProducts(client *mongo.Client, filters bson.M, options *options.FindOptions, pageNum, pageSize *int) ([]entities.Product, error) {
	opts := storage.GeneralQueryOptions{
		Filter:   filters,
		Sort:     options,
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

func GetFilters(client *mongo.Client) (*entities.FilterCategory, error) {
	opts := storage.GeneralQueryOptions{
		Filter: bson.M{},
	}

	results, err := storage.GeneralFind[map[string]interface{}](client, "Filters", opts, false)
	if err != nil || len(results) == 0 {
		return nil, err
	}

	var filters entities.FilterCategory
	data, err := bson.Marshal(results[0])
	if err != nil {
		return nil, err
	}

	if err := bson.Unmarshal(data, &filters); err != nil {
		return nil, err
	}

	return &filters, nil
}

func SaveFilters(client *mongo.Client, filters *entities.FilterCategory) error {
	id, err := primitive.ObjectIDFromHex("673a33682cc4a3c5d9a66401")
	if err != nil {
		log.Println("Ошибка при преобразовании _id:", err)
		return err
	}

	filter := bson.M{"_id": id}

	log.Printf("Saving filters: %+v", filters)

	update := bson.D{
		{"$set", bson.D{
			{"categories", filters.Categories},
			{"age", filters.Age},
			{"brand", filters.Brand},
			{"material", filters.Material},
			{"type", filters.Type},
		}},
	}

	result, err := client.Database("DutyachiySvitDB").Collection("Filters").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	log.Printf("Matched %v documents and modified %v documents", result.MatchedCount, result.ModifiedCount)
	return nil
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

func GetPurchaseHistory(client *mongo.Client, userId string) ([]entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), storage.MongoDBTimeout)
	defer cancel()

	// Логирование для проверки userId и фильтра
	fmt.Println("Retrieving purchase history for user:", userId)

	filter := bson.M{"UserId": userId}

	collection := client.Database(storage.MongoDBName).Collection(storage.Orders)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		// Логируем ошибку при выполнении запроса
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []entities.Order

	for cursor.Next(ctx) {
		var order entities.Order
		if err = cursor.Decode(&order); err != nil {
			// Логируем ошибку декодирования
			fmt.Println("Error decoding order:", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = cursor.Err(); err != nil {
		// Логируем ошибку при обходе курсора
		fmt.Println("Cursor error:", err)
		return nil, err
	}

	// Логируем результат
	fmt.Println("Retrieved orders:", orders)

	return orders, nil
}
