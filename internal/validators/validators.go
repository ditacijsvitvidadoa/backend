package validators

import (
	"fmt"
	"net/http"
	"strconv"
)

var (
	allowedCategories = []string{"ForBoys", "ForGirls", "ForInfants", "SoftToys", "BuildingSets", "Bookstore", "Creativity", "ForSchool", "Footswear", "ForSport", "Accessories"}
	allowedAges       = []string{"A", "B", "C", "D"}
	allowedBrands     = []string{"LEGO", "Mattel", "Hasbro", "Fisher-Price"}
	allowedMaterials  = []string{
		"plastic", "wood", "metal", "fabric",
		"rubber", "foam", "silicone",
		"cardboard", "paper", "plush",
		"ceramic", "glass",
	}
	allowedTypes = []string{
		"actionFigures", "dolls", "plushToys",
		"buildingSets", "educationalToys", "vehicles",
		"puzzles", "outdoorToys", "artsAndCrafts",
		"electronicToys", "musicalToys", "boardGames",
	}
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

func ExtractProductID(r *http.Request) (string, error) {
	productID := r.PathValue("id")
	if productID == "" {
		return "", fmt.Errorf("Product ID is empty")
	}

	return productID, nil
}
