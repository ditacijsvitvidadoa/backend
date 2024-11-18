package ticker

import (
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func addUniqueCategory(categories []entities.CategoryInfo, category entities.CategoryInfo) []entities.CategoryInfo {
	if category.Value == "" || category.LabelUA == "" {
		return categories
	}

	for _, item := range categories {
		if item.Value == category.Value {
			return categories
		}
	}
	return append(categories, category)
}

func updateFilters(products []entities.Product, filters *entities.FilterCategory) {
	for _, product := range products {
		filters.Brand.Items = addUniqueCategory(filters.Brand.Items, product.Brand)
		filters.Material.Items = addUniqueCategory(filters.Material.Items, product.Material)
		filters.Type.Items = addUniqueCategory(filters.Type.Items, product.Type)
	}

	filters.Brand.Items = cleanUnusedFilterItems(filters.Brand.Items, products, func(p entities.Product) entities.CategoryInfo { return p.Brand })
	filters.Material.Items = cleanUnusedFilterItems(filters.Material.Items, products, func(p entities.Product) entities.CategoryInfo { return p.Material })
	filters.Type.Items = cleanUnusedFilterItems(filters.Type.Items, products, func(p entities.Product) entities.CategoryInfo { return p.Type })
}

func cleanUnusedFilterItems(filterItems []entities.CategoryInfo, products []entities.Product, extractFunc func(entities.Product) entities.CategoryInfo) []entities.CategoryInfo {
	var validItems []entities.CategoryInfo

	for _, filterItem := range filterItems {
		found := false
		for _, product := range products {
			if extractFunc(product) == filterItem {
				found = true
				break
			}
		}
		if found {
			validItems = append(validItems, filterItem)
		}
	}

	return validItems
}

func updateFiltersIfNecessary(client *mongo.Client) error {
	currentFilters, err := requests.GetFilters(client)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	log.Println("Current Filters:", currentFilters)

	if currentFilters == nil {
		currentFilters = &entities.FilterCategory{
			Categories: entities.Filter{Title: "Категорії"},
			Age:        entities.Filter{Title: "Вік"},
			Brand:      entities.Filter{Title: "Бренд"},
			Material:   entities.Filter{Title: "Матеріали"},
			Type:       entities.Filter{Title: "Типи"},
		}
	}

	products, err := requests.GetProducts(client, bson.M{}, nil, nil, nil)
	if err != nil {
		return err
	}

	updateFilters(products, currentFilters)

	log.Println("Updated Filters:", currentFilters)

	if err = requests.SaveFilters(client, currentFilters); err != nil {
		return err
	}

	return nil
}
