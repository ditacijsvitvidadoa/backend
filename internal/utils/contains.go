package utils

import "github.com/ditacijsvitvidadoa/backend/internal/entities"

func CartContains(cart []entities.CartItem, productID int32) bool {
	for _, item := range cart {
		if item.ID == productID {
			return true
		}
	}
	return false
}

func FavouritesContains(favourites []int32, productID int32) bool {
	for _, id := range favourites {
		if id == productID {
			return true
		}
	}
	return false
}
