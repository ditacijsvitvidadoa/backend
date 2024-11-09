package entities

type CartProduct struct {
	Id       string  `bson:"Id" json:"id"`
	Articul  int     `bson:"Articul" json:"articul"`
	Code     int     `bson:"Code" json:"code"`
	ImageUrl string  `bson:"Image_url" json:"image_url"`
	Title    string  `bson:"Title" json:"title"`
	Price    int     `bson:"Price" json:"price"`
	Discount int     `bson:"Discount" json:"discount"`
	Count    int     `bson:"Count" json:"count"`
	Size     *string `bson:"Size,omitempty" json:"size,omitempty"`
	Color    *string `bson:"Color,omitempty" json:"color,omitempty"`
}
