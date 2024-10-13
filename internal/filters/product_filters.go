package filters

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/validators"
	"net/http"
	"strconv"
)

func BuildFilter(r *http.Request) (map[string]interface{}, error) {
	filter := make(map[string]interface{})

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
	if err := addDiscountAndPriceFilters(r, filter); err != nil {
		return nil, err
	}

	return filter, nil
}

func addCategoryFilter(r *http.Request, filter map[string]interface{}) error {
	category := r.URL.Query().Get("category")
	if category != "" {
		if !validators.IsValidCategory(category) {
			return fmt.Errorf("Invalid category: %s", category)
		}
		filter["category"] = category
	}
	return nil
}

func addAgeFilter(r *http.Request, filter map[string]interface{}) error {
	age := r.URL.Query().Get("age")
	if age != "" {
		if !validators.IsValidAge(age) {
			return fmt.Errorf("Invalid age: %s", age)
		}
		filter["age"] = age
	}
	return nil
}

func addBrandFilter(r *http.Request, filter map[string]interface{}) error {
	brand := r.URL.Query().Get("brand")
	if brand != "" {
		if !validators.IsValidBrand(brand) {
			return fmt.Errorf("Invalid brand: %s", brand)
		}
		filter["brand"] = brand
	}
	return nil
}

func addMaterialFilter(r *http.Request, filter map[string]interface{}) error {
	material := r.URL.Query().Get("material")
	if material != "" {
		if !validators.IsValidMaterial(material) {
			return fmt.Errorf("Invalid material: %s", material)
		}
		filter["material"] = material
	}
	return nil
}

func addTypeFilter(r *http.Request, filter map[string]interface{}) error {
	productType := r.URL.Query().Get("type")
	if productType != "" {
		if !validators.IsValidType(productType) {
			return fmt.Errorf("Invalid type: %s", productType)
		}
		filter["type"] = productType
	}
	return nil
}

func addDiscountAndPriceFilters(r *http.Request, filter map[string]interface{}) error {
	// Добавить логику обработки скидок и цен
	// ...
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

	pageSize := 10
	if size := r.URL.Query().Get("pageSize"); size != "" {
		s, err := strconv.Atoi(size)
		if err == nil && s > 0 {
			pageSize = s
		}
	}

	return pageNum, pageSize
}
