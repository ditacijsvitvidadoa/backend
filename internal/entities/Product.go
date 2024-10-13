package entities

type Product struct {
	Id          int    `bson:"id" json:"id"`
	Articul     int    `bson:"articul" json:"articul"`
	Code        int    `bson:"code" json:"code"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Price       int    `bson:"price" json:"price"`
	Discount    int    `bson:"discount" json:"discount"`
	Category    string `bson:"category" json:"category"`
	Material    string `bson:"material" json:"material"`
	Brand       string `bson:"brand" json:"brand"`
	Age         string `bson:"age" json:"age"`
	InCart      bool   `json:"in_cart"`
	IsFavourite bool   `json:"is_favourite"`
	Count       int    `bson:"count" json:"count"`
}
