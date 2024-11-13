package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (a *App) GetFavouritesProducts(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	productsIds, err := requests.GetFavouritesByUserID(a.client, userObjectId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve products IDs from session cookie.")
		return
	}

	filter := bson.M{
		"Id": bson.M{
			"$in": productsIds,
		},
	}

	products, err := requests.GetProducts(a.client, filter, nil, nil, nil)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve products.")
		return
	}

	sendResponse(w, products)
}

func (a *App) addFavouriteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

	if productID == "" {
		sendError(w, http.StatusNotFound, "Product ID is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	changeCount, err := requests.AddProductToFavourites(a.client, userObjectId, productID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to add product to favorites.")
		return
	}

	if changeCount == 0 {
		sendError(w, http.StatusInternalServerError, "Nothing was changed")
		return
	}

	sendOk(w)
}

func (a *App) deleteFavouriteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

	if productID == "" {
		sendError(w, http.StatusNotFound, "Product ID is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	changeCount, err := requests.RemoveProductFromFavourites(a.client, userObjectId, productID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to remove product from favorites.")
		return
	}

	if changeCount == 0 {
		sendError(w, http.StatusInternalServerError, "Nothing was changed")
		return
	}

	sendOk(w)
}
