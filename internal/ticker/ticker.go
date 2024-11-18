package ticker

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func GeneralTicker(client *mongo.Client) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("ticker was started")
			if err := updateFiltersIfNecessary(client); err == nil {
				log.Println("Successfully updated filters")
			} else {
				log.Println("Failed to update filters")
			}

			CounterTicker(client)
		}
	}
}
