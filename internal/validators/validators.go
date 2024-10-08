package validators

import (
	"strconv"
)

var (
	allowedCategories = []string{"ForBoys", "ForGirls", "ForInfants", "SoftToys", "BuildingSets", "Bookstore", "Creativity", "ForSchool", "Footswear", "ForSport", "Accessories"}
	allowedAges       = []string{"A", "B", "C", "D"}
	allowedBrands     = []string{"brand1", "brand2", "brand3"}
	allowedMaterials  = []string{"plastic", "wood", "metal"}
	allowedTypes      = []string{"type1", "type2", "type3"}
)

func IsValidCategory(category string) bool {
	for _, c := range allowedCategories {
		if category == c {
			return true
		}
	}
	return false
}

func IsValidAge(age string) bool {
	for _, a := range allowedAges {
		if age == a {
			return true
		}
	}
	return false
}

func IsValidBrand(brand string) bool {
	for _, b := range allowedBrands {
		if brand == b {
			return true
		}
	}
	return false
}

func IsValidMaterial(material string) bool {
	for _, m := range allowedMaterials {
		if material == m {
			return true
		}
	}
	return false
}

func IsValidType(productType string) bool {
	for _, t := range allowedTypes {
		if productType == t {
			return true
		}
	}
	return false
}

func IsValidDiscount(discount string) (bool, error) {
	return strconv.ParseBool(discount)
}

func IsValidPriceRange(from, to int) bool {
	return from >= 0 && to >= from
}
