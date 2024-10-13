package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func (a *App) getCartProducts(w http.ResponseWriter, r *http.Request) {
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

	UserInfo, err := requests.GetUserByID(a.client, userObjectId)
	if err != nil {
		sendError(w, http.StatusNoContent, "Failed to retrieve user info from storage.")
		return
	}

	if UserInfo.Cart == nil {
		sendError(w, http.StatusOK, "Cart is empty.")
		return
	}

	cart := UserInfo.Cart
	var cartItemIDs []int
	cartCountMap := make(map[int]int)

	for _, item := range cart {
		cartItemIDs = append(cartItemIDs, item.ID)
		cartCountMap[item.ID] = item.Count
	}

	filters := bson.M{"id": bson.M{"$in": cartItemIDs}}

	products, err := requests.GetProducts(a.client, filters, nil, nil)
	if err != nil {
		sendError(w, http.StatusNoContent, "Failed to retrieve products from storage.")
		return
	}

	for i := range products {
		if count, exists := cartCountMap[products[i].Id]; exists {
			products[i].Count = count
		} else {
			products[i].Count = 0
		}
	}

	sendResponse(w, products)
}

func (a *App) deleteCartProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID.")
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

	changeCount, err := requests.UpdateUserCart(a.client, userObjectId, productID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to remove product from cart.")
		return
	}

	if changeCount == 0 {
		sendError(w, http.StatusInternalServerError, "Nothing was changed")
		return
	}

	sendOk(w)
}

func (a *App) addCartProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID.")
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

	countStr := r.URL.Query().Get("count")
	count := 1

	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid count value.")
			return
		}
	}

	changeCount, err := requests.AddProductToCart(a.client, userObjectId, productID, count)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to add product to cart.")
		return
	}

	if changeCount == 0 {
		sendError(w, http.StatusInternalServerError, "Nothing was changed")
		return
	}

	sendOk(w)
}
