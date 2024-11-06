package entities

type Product struct {
	Id              string            `bson:"id" json:"id"`
	Articul         int               `bson:"articul" json:"articul"`
	Code            int               `bson:"code" json:"code"`
	ImageUrls       []string          `bson:"image_urls" json:"image_urls"`
	Title           string            `bson:"title" json:"title"`
	Description     string            `bson:"description" json:"description"`
	Price           int               `bson:"price" json:"price"`
	Discount        int               `bson:"discount" json:"discount"`
	Category        string            `bson:"category" json:"category"`
	Material        string            `bson:"material" json:"material"`
	Brand           string            `bson:"brand" json:"brand"`
	Age             string            `bson:"age" json:"age"`
	InCart          bool              `json:"in_cart"`
	IsFavourite     bool              `json:"is_favourite"`
	Count           int               `bson:"count" json:"count"`
	Sizes           *SizeInfo         `bson:"sizes,omitempty" json:"sizes,omitempty"`
	Colors          *ColorInfo        `bson:"colors,omitempty" json:"colors,omitempty"`
	Characteristics []*Characteristic `bson:"characteristics,omitempty" json:"characteristics,omitempty"`
}

type SizeInfo struct {
	Category string `bson:"category" json:"category"`
	HasTable bool   `bson:"has_table" json:"has_table"`
	Sizes    struct {
		DefaultSize string   `bson:"default_size" json:"default_size"`
		HasSizes    bool     `bson:"has_sizes" json:"has_sizes"`
		SizeValues  []string `bson:"size_values" json:"size_values"`
	} `bson:"sizes" json:"sizes"`
}

type ColorInfo struct {
	Colors []string `bson:"colors" json:"colors"`
}

type Characteristic struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}
