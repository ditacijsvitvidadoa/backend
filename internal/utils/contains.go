package utils

import (
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"log"
)

func CartContains(cart []entities.CartItem, productID, size, color string) bool {
	for _, item := range cart {
		if item.ID == productID {
			log.Printf("Checking item: %+v, size: %s, color: %s", item, size, color)

			if size != "" && item.Details != nil && item.Details.Size != size {
				continue
			}
			if color != "" && item.Details != nil && item.Details.Color != color {
				continue
			}
			return true
		}
	}

	return false
}

func FavouritesContains(favourites []string, productID string) bool {
	for _, id := range favourites {
		if id == productID {
			return true
		}
	}
	return false
}
