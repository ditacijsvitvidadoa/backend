package entities

type Product struct {
	Id              string            `bson:"Id" json:"id"`
	Articul         int               `bson:"Articul" json:"articul"`
	Code            int               `bson:"Code" json:"code"`
	ImageUrls       []string          `bson:"Image_urls" json:"image_urls"`
	Title           string            `bson:"Title" json:"title"`
	Description     string            `bson:"Description" json:"description"`
	Price           int               `bson:"Price" json:"price"`
	Discount        int               `bson:"Discount" json:"discount"`
	Category        CategoryInfo      `bson:"Category" json:"category"`
	Material        CategoryInfo      `bson:"Material" json:"material"`
	Brand           CategoryInfo      `bson:"Brand" json:"brand"`
	Type            CategoryInfo      `bson:"Type" json:"type"`
	Age             string            `bson:"Age" json:"age"`
	InCart          bool              `bson:"InCart" json:"in_cart"`
	IsFavourite     bool              `bson:"IsFavourite" json:"is_favourite"`
	Count           int               `bson:"Count" json:"count"`
	Sizes           *SizeInfo         `bson:"Sizes,omitempty" json:"sizes,omitempty"`
	Colors          *ColorInfo        `bson:"Colors,omitempty" json:"colors,omitempty"`
	Characteristics []*Characteristic `bson:"Characteristics,omitempty" json:"characteristics,omitempty"`
}

type SizeInfo struct {
	Category string `bson:"Category" json:"category"`
	HasTable bool   `bson:"HasTable" json:"has_table"`
	Sizes    Sizes  `bson:"Sizes" json:"sizes"`
}

type Sizes struct {
	DefaultSize string   `bson:"DefaultSize" json:"default_size"`
	HasSizes    bool     `bson:"HasSizes" json:"has_sizes"`
	SizeValues  []string `bson:"sizeValues" json:"size_values"`
}

type ColorInfo struct {
	DefaultColor string   `bson:"DefaultColor" json:"default_color"`
	Colors       []string `bson:"Colors" json:"colors"`
}

type Characteristic struct {
	Key   string `bson:"Key" json:"key"`
	Value string `bson:"Value" json:"value"`
}
