package app

import (
	"errors"
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
	"strconv"
	"strings"
)

func (a *App) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product, err := collectProductData(r)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	objectId, err := requests.CreateNewProduct(a.client, product)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(objectId)

	sendOk(w)
}

func collectProductData(r *http.Request) (entities.Product, error) {
	productID := utils.GenerateUUID()
	fmt.Println(productID)
	title := r.FormValue("title")
	if title == "" {
		return entities.Product{}, errors.New("title is required")
	}

	articul, err := utils.ParseFormValueAsInt(r.FormValue("articul"))
	if err != nil {
		return entities.Product{}, fmt.Errorf("invalid articul: %v", err)
	}

	code, err := utils.ParseFormValueAsInt(r.FormValue("code"))
	if err != nil {
		return entities.Product{}, fmt.Errorf("invalid code: %v", err)
	}

	description := r.FormValue("description")
	if description == "" {
		return entities.Product{}, errors.New("description is required")
	}

	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil || price <= 0 {
		return entities.Product{}, errors.New("price is required and must be a positive number")
	}

	category := r.FormValue("category")
	if category == "" {
		return entities.Product{}, errors.New("category is required")
	}

	material := r.FormValue("material")
	brand := r.FormValue("brand")
	age := r.FormValue("age")
	inCart := r.FormValue("in_cart") == "false"
	isFavourite := r.FormValue("is_favourite") == "false"

	discount, err := utils.ParseFormValueAsInt(r.FormValue("discount"))
	if err != nil {
		return entities.Product{}, fmt.Errorf("invalid discount: %v", err)
	}
	fmt.Println("discount", discount)

	var sizeInfo *entities.SizeInfo
	if r.FormValue("has_sizes") != "" {
		sizeInfo = &entities.SizeInfo{
			Category: r.FormValue("table.category"),
			HasTable: r.FormValue("has_table") == "false",
		}

		defaultSize := "M"
		sizeValues := r.Form["size_value"]

		sizeInfo.Sizes.DefaultSize = defaultSize
		sizeInfo.Sizes.HasSizes = len(sizeValues) > 0
		sizeInfo.Sizes.SizeValues = sizeValues
	}

	fmt.Println()
	fmt.Printf("sizeInfo %v", sizeInfo)
	fmt.Println()

	var colorInfo *entities.ColorInfo
	colorValues := r.Form["colors"]
	if len(colorValues) > 0 {
		colorInfo = &entities.ColorInfo{Colors: colorValues}
	}

	fmt.Println("colorInfo", colorInfo)

	var characteristics []*entities.Characteristic
	characteristicKeys := r.Form["characteristic_key"]
	characteristicValues := r.Form["characteristic_value"]
	if len(characteristicKeys) == len(characteristicValues) {
		for i := range characteristicKeys {
			key := strings.TrimSpace(characteristicKeys[i])
			value := strings.TrimSpace(characteristicValues[i])

			if key != "" && value != "" {
				characteristics = append(characteristics, &entities.Characteristic{
					Key:   key,
					Value: value,
				})
			}
		}
	}

	imageURLs, err := utils.ParseAndSaveFiles(r, productID)
	if err != nil {
		fmt.Println(err.Error())
		return entities.Product{}, err
	}

	product := entities.Product{
		Id:              productID,
		Title:           title,
		Articul:         articul,
		Code:            code,
		Description:     description,
		Price:           price,
		ImageUrls:       imageURLs,
		Category:        category,
		Material:        material,
		Brand:           brand,
		Age:             age,
		InCart:          inCart,
		IsFavourite:     isFavourite,
		Discount:        discount,
		Characteristics: characteristics,
		Colors:          colorInfo,
		Sizes:           sizeInfo,
	}

	return product, nil
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	cookieValue, err := cookie.GetSessionValue(r, "session")

	var userID string
	var userObjectId primitive.ObjectID
	var cart []entities.CartItem
	var favourites []int32

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
			var productID int32
			if id, ok := products[i]["id"].(int32); ok {
				productID = int32(id)
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

				if utils.CartContains(cart, int32(productID)) {
					product["in_cart"] = true
				}

				if utils.FavouritesContains(favourites, int32(productID)) {
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
