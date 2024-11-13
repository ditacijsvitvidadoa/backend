package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
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
		sendError(w, http.StatusUnauthorized, "Failed to convert user ID from session cookie.")
		return
	}

	UserInfo, err := requests.GetUserByID(a.client, userObjectId)
	if err != nil {
		sendError(w, http.StatusNoContent, "Failed to retrieve user info from storage.")
		return
	}

	if len(UserInfo.Cart) == 0 {
		sendNoContent(w)
		return
	}

	cartItemIDs := make([]string, 0, len(UserInfo.Cart))
	for _, item := range UserInfo.Cart {
		cartItemIDs = append(cartItemIDs, item.ID)
	}

	filters := bson.M{"Id": bson.M{"$in": cartItemIDs}}
	products, err := requests.GetProducts(a.client, filters, nil, nil, nil)
	if err != nil {
		sendError(w, http.StatusNoContent, "Failed to retrieve products from storage.")
		return
	}

	if len(products) == 0 {
		sendResponse(w, []any{})
		return
	}

	productMap := make(map[string]entities.Product)
	for _, product := range products {
		productMap[product.Id] = product
	}

	var cartProducts []entities.CartProduct

	for _, cartItem := range UserInfo.Cart {
		if product, exists := productMap[cartItem.ID]; exists {
			cartProduct := entities.CartProduct{
				Id:       product.Id,
				Articul:  product.Articul,
				Code:     product.Code,
				ImageUrl: product.ImageUrls[0],
				Title:    product.Title,
				Price:    product.Price,
				Discount: product.Discount,
				Count:    cartItem.Count,
			}

			if cartItem.Details != nil {
				if cartItem.Details.Size != "" {
					cartProduct.Size = &cartItem.Details.Size
				}
				if cartItem.Details.Color != "" {
					cartProduct.Color = &cartItem.Details.Color
				}
			}

			cartProducts = append(cartProducts, cartProduct)
		}
	}

	sendResponse(w, cartProducts)
}

func (a *App) deleteCartProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

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
	productID := r.PathValue("id")

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

	details := bson.M{}
	if size := r.URL.Query().Get("size"); size != "" {
		details["Size"] = size
	}
	if color := r.URL.Query().Get("color"); color != "" {
		details["Color"] = color
	}

	update := bson.M{
		"$addToSet": bson.M{
			"Cart": bson.M{
				"Id":      productID,
				"Count":   count,
				"Details": details,
			},
		},
	}

	changeCount, err := requests.AddProductToCart(a.client, userObjectId, update)
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

func (a *App) UpdateCount(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		sendError(w, http.StatusBadRequest, "Dont have new count value")
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

	count, err := strconv.Atoi(countStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Wrong count value; Not be int")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		sendError(w, http.StatusBadRequest, "Dont have id value")
		return
	}
	size := r.URL.Query().Get("size")
	color := r.URL.Query().Get("color")

	filter := bson.M{"Id": id}

	if size != "" {
		filter["Details.Size"] = size
	}
	if color != "" {
		filter["Details.Color"] = color
	}

	wasUpdated, err := requests.UpdateCartProductCount(a.client, userObjectId, filter, count)

	if wasUpdated == 0 {
		sendError(w, http.StatusNoContent, "Nothing was changed")
		return
	}

	sendOk(w)
}
