package ticker

import (
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func updateCounters(client *mongo.Client) error {
	productCount, err := requests.CountDocuments(client, storage.Products)
	if err != nil {
		return err
	}

	userCount, err := requests.CountDocuments(client, storage.Users)
	if err != nil {
		return err
	}

	err = requests.UpdateCounterValue(client, "productId", productCount)
	if err != nil {
		return err
	}

	err = requests.UpdateCounterValue(client, "userId", userCount)
	if err != nil {
		return err
	}

	return nil
}

func CounterTicker(client *mongo.Client) {
	err := updateCounters(client)
	log.Println("err", err)
}
