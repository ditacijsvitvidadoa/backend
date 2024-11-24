package entities

type ProductActivity struct {
	ProductId   string `bson:"ProductId" json:"product_id"`
	Sales       int    `bson:"Sales" json:"sales"`
	Clicks      int    `bson:"Clicks" json:"clicked"`
	AddedToCart int    `bson:"AddedToCart" json:"added_to_cart"`
	Favourites  int    `bson:"Favourites" json:"favourites"`
}
