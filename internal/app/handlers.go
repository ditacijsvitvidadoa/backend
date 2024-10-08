package app

import (
	"encoding/json"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/storage"
	"github.com/ditacijsvitvidadoa/backend/internal/validators"
	"log"
	"net/http"
	"strconv"
)

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	age := r.URL.Query().Get("age")
	brand := r.URL.Query().Get("brand")
	material := r.URL.Query().Get("material")
	discount := r.URL.Query().Get("discount")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	productType := r.URL.Query().Get("type")

	filter := make(map[string]interface{})

	if category != "" {
		if !validators.IsValidCategory(category) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid category: %s", category))
			return
		}
		filter["category"] = category
	}

	if age != "" {
		if !validators.IsValidAge(age) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid age: %s", age))
			return
		}
		filter["age"] = age
	}

	if brand != "" {
		if !validators.IsValidBrand(brand) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid brand: %s", brand))
			return
		}
		filter["brand"] = brand
	}

	if material != "" {
		if !validators.IsValidMaterial(material) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid material: %s", material))
			return
		}
		filter["material"] = material
	}

	if productType != "" {
		if !validators.IsValidType(productType) {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid type: %s", productType))
			return
		}
		filter["type"] = productType
	}

	if discount != "" {
		if disc, err := validators.IsValidDiscount(discount); err != nil {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid discount: %s", discount))
			return
		} else {
			filter["discount"] = disc
		}
	}

	if minPrice != "" || maxPrice != "" {
		priceFilter := make(map[string]interface{})
		var minP, maxP int
		var err error

		if minPrice != "" {
			minP, err = strconv.Atoi(minPrice)
			if err != nil {
				sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid minPrice: %s", minPrice))
				return
			}
			priceFilter["$gte"] = minP
		}

		if maxPrice != "" {
			maxP, err = strconv.Atoi(maxPrice)
			if err != nil {
				sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid maxPrice: %s", maxPrice))
				return
			}
			priceFilter["$lte"] = maxP
		}

		if !validators.IsValidPriceRange(minP, maxP) {
			sendError(w, http.StatusBadRequest, "Invalid price range: minPrice is greater than maxPrice or values are negative")
			return
		}

		if len(priceFilter) > 0 {
			filter["price"] = priceFilter
		}
	}

	pageNum := 1
	if page := r.URL.Query().Get("pageNum"); page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing pageNum: %s", err))
			return
		}
		if p < 1 {
			sendError(w, http.StatusBadRequest, "pageNum must be greater than or equal to 1")
			return
		}
		pageNum = p
	}

	pageSize := 10
	if size := r.URL.Query().Get("pageSize"); size != "" {
		s, err := strconv.Atoi(size)
		if err != nil {
			sendError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing pageSize: %s", err))
			return
		}
		if s < 1 {
			sendError(w, http.StatusBadRequest, "pageSize must be greater than or equal to 1")
			return
		}
		pageSize = s
	}

	products, err := storage.Get(a.client, "Products", filter, pageNum, pageSize)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching products: %s", err))
		return
	}

	sendResponse(w, products)
}

func sendError(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"status":"error","message":%s"}`, text)))
}

func sendOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func sendResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("writing response: %s", err)
	}
}
