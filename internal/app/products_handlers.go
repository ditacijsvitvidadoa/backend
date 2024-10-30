package app

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/filters"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"github.com/ditacijsvitvidadoa/backend/internal/validators"
	"go.mongodb.org/mongo-driver/bson"
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

	fmt.Println(filter)

	allProducts, err := requests.GetProducts(a.client, filter, nil, nil)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching all products: %s", err))
		return
	}

	details := buildProductDetails(allProducts)

	pageNum, pageSize := filters.GetPaginationParams(r)
	products, err := requests.GetProducts(a.client, filter, &pageNum, &pageSize)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching paginated products: %s", err))
		return
	}

	if userID != "" {
		for i := range products {
			var productID int
			if id, ok := products[i]["id"].(int32); ok {
				productID = int(id)
			}

			if utils.CartContains(cart, productID) {
				products[i]["in_cart"] = true
			} else {
				products[i]["in_cart"] = false
			}

			if utils.FavouritesContains(favourites, productID) {
				products[i]["is_favourite"] = true
			} else {
				products[i]["is_favourite"] = false
			}
		}
	}

	response := map[string]interface{}{
		"products": products,
		"details":  details,
	}

	sendResponse(w, response)
}

func (a *App) getProductByID(w http.ResponseWriter, r *http.Request) {
	productID, err := validators.ExtractProductID(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	filter := bson.M{"id": productID}

	products, err := requests.GetProducts(a.client, filter, nil, nil)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching product: %s", err))
		return
	}

	if len(products) == 0 {
		sendError(w, http.StatusNotFound, "Product not found")
		return
	}

	product := products[0]

	product["in_cart"] = false
	product["is_favourite"] = false

	cookieValue, err := cookie.GetSessionValue(r, "session")
	if err == nil {
		userID, err := cookie.GetUserIDFromCookie(cookieValue)
		if err == nil {
			userObjectId, err := primitive.ObjectIDFromHex(userID)
			if err == nil {
				cart, _ := requests.GetCartByUserID(a.client, userObjectId)
				favourites, _ := requests.GetFavouritesByUserID(a.client, userObjectId)

				if utils.CartContains(cart, int(productID)) {
					product["in_cart"] = true
				}

				if utils.FavouritesContains(favourites, int(productID)) {
					product["is_favourite"] = true
				}
			}
		}
	}

	sendResponse(w, product)
}

func buildProductDetails(products []map[string]interface{}) map[string]interface{} {
	if len(products) == 0 {
		return map[string]interface{}{
			"total_count":       0,
			"min_price_product": nil,
			"max_price_product": nil,
		}
	}

	var minProduct, maxProduct map[string]interface{}
	minPrice := int(^uint(0) >> 1)
	maxPrice := 0

	for _, product := range products {
		currentPrice := getEffectivePrice(product)

		if currentPrice < minPrice {
			minPrice = currentPrice
			minProduct = product
		}
		if currentPrice > maxPrice {
			maxPrice = currentPrice
			maxProduct = product
		}
	}

	return map[string]interface{}{
		"total_count": len(products),
		"min_price_product": map[string]interface{}{
			"id":    minProduct["id"],
			"price": minPrice,
		},
		"max_price_product": map[string]interface{}{
			"id":    maxProduct["id"],
			"price": maxPrice,
		},
	}
}

func getEffectivePrice(product map[string]interface{}) int {
	if discount, ok := product["discount"].(int32); ok && discount > 0 {
		return int(discount)
	}
	if price, ok := product["price"].(int32); ok {
		return int(price)
	}
	return 0
}
