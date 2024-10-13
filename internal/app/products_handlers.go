package app

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/filters"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	cookieValue, err := cookie.GetSessionValue(r, "session")

	var userID string
	var userObjectId primitive.ObjectID
	var cart []entities.CartItem
	var favourites []int

	if err == nil {
		userID, err = cookie.GetUserIDFromCookie(cookieValue)
		if err == nil {
			userObjectId, err = primitive.ObjectIDFromHex(userID)
			if err == nil {
				cart, _ = requests.GetCartByUserID(a.client, userObjectId)
				favourites, _ = requests.GetFavouritesByUserID(a.client, userObjectId)
			}
		}
	}

	filter, err := filters.BuildFilter(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	pageNum, pageSize := filters.GetPaginationParams(r)

	products, err := requests.GetProducts(a.client, filter, &pageNum, &pageSize)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching products: %s", err))
		return
	}

	if userID != "" {
		for i := range products {
			if utils.CartContains(cart, products[i].Id) {
				products[i].InCart = true
			}
			if utils.FavouritesContains(favourites, products[i].Id) {
				products[i].IsFavourite = true
			}
		}
	}

	sendResponse(w, products)
}
