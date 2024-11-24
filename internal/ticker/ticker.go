package ticker

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func GeneralTicker(client *mongo.Client) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("ticker was started")
			if err := updateFiltersIfNecessary(client); err == nil {
			} else {
			}

			CounterTicker(client)
		}
	}
}
