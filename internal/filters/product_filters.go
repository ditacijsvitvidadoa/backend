package filters

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/validators"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func BuildFilter(r *http.Request) (bson.M, error) {
	filter := bson.M{}

	if err := addCategoryFilter(r, filter); err != nil {
		return nil, err
	}
	if err := addAgeFilter(r, filter); err != nil {
		return nil, err
	}
	if err := addBrandFilter(r, filter); err != nil {
		return nil, err
	}
	if err := addMaterialFilter(r, filter); err != nil {
		return nil, err
	}
	if err := addTypeFilter(r, filter); err != nil {
		return nil, err
	}

	return filter, nil
}

func addCategoryFilter(r *http.Request, filter bson.M) error {
	categoriesParam, err := url.QueryUnescape(r.URL.Query().Get("categories"))
	if err != nil {
		return fmt.Errorf("failed to decode categories: %w", err)
	}

	if categoriesParam != "" {
		categories := strings.Split(categoriesParam, ",")
		var validCategories []string

		for _, category := range categories {
			category = strings.TrimSpace(category)
			if validators.IsValidCategory(category) {
				validCategories = append(validCategories, category)
			} else {
				return fmt.Errorf("invalid category: %s", category)
			}
		}

		if len(validCategories) > 0 {
			filter["Category"] = bson.M{"$in": validCategories}
		}
	}

	return nil
}

func addAgeFilter(r *http.Request, filter bson.M) error {
	ageParam, err := url.QueryUnescape(r.URL.Query().Get("age"))
	if err != nil {
		return fmt.Errorf("failed to decode age: %w", err)
	}

	if ageParam != "" {
		ages := strings.Split(ageParam, ",")
		var validAges []string

		for _, age := range ages {
			age = strings.TrimSpace(age)
			if !validators.IsValidAge(age) {
				return fmt.Errorf("invalid age: %s", age)
			}
			validAges = append(validAges, age)
		}

		filter["age"] = bson.M{"$in": validAges}
	}
	return nil
}

func addBrandFilter(r *http.Request, filter bson.M) error {
	brandParam, err := url.QueryUnescape(r.URL.Query().Get("brand"))
	if err != nil {
		return fmt.Errorf("failed to decode brand: %w", err)
	}

	if brandParam != "" {
		brands := strings.Split(brandParam, ",")
		var validBrands []string

		for _, brand := range brands {
			brand = strings.TrimSpace(brand)
			if !validators.IsValidBrand(brand) {
				return fmt.Errorf("invalid brand: %s", brand)
			}
			validBrands = append(validBrands, brand)
		}

		filter["brand"] = bson.M{"$in": validBrands}
	}
	return nil
}

func addMaterialFilter(r *http.Request, filter bson.M) error {
	materialParam, err := url.QueryUnescape(r.URL.Query().Get("material"))
	if err != nil {
		return fmt.Errorf("failed to decode material: %w", err)
	}

	if materialParam != "" {
		materials := strings.Split(materialParam, ",")
		var validMaterials []string

		for _, material := range materials {
			material = strings.TrimSpace(material)
			if !validators.IsValidMaterial(material) {
				return fmt.Errorf("invalid material: %s", material)
			}
			validMaterials = append(validMaterials, material)
		}

		filter["material"] = bson.M{"$in": validMaterials}
	}
	return nil
}

func addTypeFilter(r *http.Request, filter bson.M) error {
	typeParam, err := url.QueryUnescape(r.URL.Query().Get("type"))
	if err != nil {
		return fmt.Errorf("failed to decode type: %w", err)
	}

	if typeParam != "" {
		types := strings.Split(typeParam, ",")
		var validTypes []string

		for _, productType := range types {
			productType = strings.TrimSpace(productType)
			if !validators.IsValidType(productType) {
				return fmt.Errorf("invalid type: %s", productType)
			}
			validTypes = append(validTypes, productType)
		}

		filter["type"] = bson.M{"$in": validTypes}
	}
	return nil
}

//func addDiscountAndPriceFilters(products []entities.Product, r *http.Request) ([]entities.Product, error) {
//	minPriceStr := r.URL.Query().Get("minPrice")
//	maxPriceStr := r.URL.Query().Get("maxPrice")
//
//	var minPrice, maxPrice int
//	var err error
//
//	if minPriceStr != "" {
//		minPrice, err = strconv.Atoi(minPriceStr)
//		if err != nil {
//			return nil, fmt.Errorf("Invalid minPrice: %s", minPriceStr)
//		}
//	}
//
//	if maxPriceStr != "" {
//		maxPrice, err = strconv.Atoi(maxPriceStr)
//		if err != nil {
//			return nil, fmt.Errorf("Invalid maxPrice: %s", maxPriceStr)
//		}
//	}
//
//	var filteredProducts []entities.Product
//	for _, product := range products {
//		if product.Discount != nil && *product.Discount > 0 {
//			discountedPrice := product.Price - (*product.Discount)
//
//			if (minPriceStr == "" || discountedPrice >= minPrice) &&
//				(maxPriceStr == "" || discountedPrice <= maxPrice) {
//				filteredProducts = append(filteredProducts, product)
//			}
//		} else {
//			// Проверка на обычную цену
//			if (minPriceStr == "" || product.Price >= minPrice) &&
//				(maxPriceStr == "" || product.Price <= maxPrice) {
//				filteredProducts = append(filteredProducts, product)
//			}
//		}
//	}
//
//	return filteredProducts, nil
//}

func GetPaginationParams(r *http.Request) (int, int) {
	pageNum := 1
	if page := r.URL.Query().Get("pageNum"); page != "" {
		p, err := strconv.Atoi(page)
		if err == nil && p > 0 {
			pageNum = p
		}
	}

	pageSize := 24
	if size := r.URL.Query().Get("pageSize"); size != "" {
		s, err := strconv.Atoi(size)
		if err == nil && s > 0 {
			pageSize = s
		}
	}

	return pageNum, pageSize
}
