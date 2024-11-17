package ticker

import (
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func addUniqueCategory(categories []entities.CategoryInfo, category entities.CategoryInfo) []entities.CategoryInfo {
	for _, item := range categories {
		if item.Value == category.Value {
			return categories
		}
	}
	return append(categories, category)
}

func updateFilters(products []entities.Product, filters *entities.FilterCategory) {
	for _, product := range products {
		filters.Categories.Items = addUniqueCategory(filters.Categories.Items, product.Category)

		filters.Brand.Items = addUniqueCategory(filters.Brand.Items, product.Brand)

		filters.Material.Items = addUniqueCategory(filters.Material.Items, product.Material)

		filters.Type.Items = addUniqueCategory(filters.Type.Items, product.Type)
	}
}

func updateFiltersIfNecessary(client *mongo.Client) error {
	currentFilters, err := requests.GetFilters(client)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	log.Println(currentFilters)

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

	log.Println(products)

	initialFilters := *currentFilters
	updateFilters(products, currentFilters)

	log.Println(equalFilters(&initialFilters, currentFilters))

	if !equalFilters(&initialFilters, currentFilters) {
		if err = requests.SaveFilters(client, currentFilters); err != nil {
			return err
		}
	}

	return nil
}

func equalFilters(a, b *entities.FilterCategory) bool {
	return len(a.Categories.Items) == len(b.Categories.Items) &&
		len(a.Brand.Items) == len(b.Brand.Items) &&
		len(a.Material.Items) == len(b.Material.Items) &&
		len(a.Type.Items) == len(b.Type.Items)
}
