package filters

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func BuildFilter(r *http.Request) (bson.M, error) {
	filter := bson.M{}

	if err := buildSearchFilter(r, filter); err != nil {
		return nil, err
	}
	if err := addDiscountAndPriceFilters(r, filter); err != nil {
		return nil, err
	}
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

func buildSearchFilter(r *http.Request, filter bson.M) error {
	search := r.URL.Query().Get("search")

	if search == "" {
		return nil
	}

	searchFilter := bson.M{
		"$or": []bson.M{
			{"Title": bson.M{"$regex": search, "$options": "i"}},
			{"Description": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	if existing, ok := filter["$or"]; ok {
		filter["$or"] = append(existing.([]bson.M), searchFilter["$or"].([]bson.M)...)
	} else {
		filter["$or"] = searchFilter["$or"]
	}

	return nil
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

func addDiscountAndPriceFilters(r *http.Request, filter bson.M) error {
	minPriceStr := r.URL.Query().Get("minPrice")
	maxPriceStr := r.URL.Query().Get("maxPrice")

	fmt.Printf("minPriceStr: %s, maxPriceStr: %s\n", minPriceStr, maxPriceStr)

	// Переменная для фильтра по цене
	var priceFilters []bson.M
	if minPriceStr != "" {
		minPrice, err := strconv.Atoi(minPriceStr)
		if err != nil {
			return fmt.Errorf("invalid minPrice: %s", minPriceStr)
		}
		fmt.Printf("Parsed minPrice: %d\n", minPrice)
		priceFilters = append(priceFilters, bson.M{
			"$or": []interface{}{
				bson.M{"Discount": bson.M{"$gte": minPrice, "$ne": 0}},
				bson.M{"Price": bson.M{"$gte": minPrice}},
			},
		})
	}

	if maxPriceStr != "" {
		maxPrice, err := strconv.Atoi(maxPriceStr)
		if err != nil {
			return fmt.Errorf("invalid maxPrice: %s", maxPriceStr)
		}
		fmt.Printf("Parsed maxPrice: %d\n", maxPrice)
		priceFilters = append(priceFilters, bson.M{
			"$or": []interface{}{
				bson.M{"Discount": bson.M{"$lte": maxPrice, "$ne": 0}},
				bson.M{"Price": bson.M{"$lte": maxPrice}},
			},
		})
	}

	if len(priceFilters) > 0 {
		filter["$and"] = priceFilters
	}

	return nil
}

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

func AddSortOrderFilter(r *http.Request, filter bson.M, options *options.FindOptions) error {
	sortOrderParam := r.URL.Query().Get("sortOrder")

	if sortOrderParam != "" {
		if !validators.IsValidSortOrder(sortOrderParam) {
			return fmt.Errorf("invalid sortOrder: %s", sortOrderParam)
		}

		var sortOrder bson.D
		switch sortOrderParam {
		case "ascending":
			sortOrder = bson.D{{"Discount", -1}, {"Price", 1}}
		case "descending":
			sortOrder = bson.D{{"Discount", -1}, {"Price", -1}}
		case "popular":
			sortOrder = bson.D{}
		default:
			sortOrder = bson.D{}
		}

		fmt.Printf("sortOrder: %+v\n", sortOrder)
		options.SetSort(sortOrder)
		fmt.Printf("options after setting sort: %+v\n", options)
	}

	return nil
}
