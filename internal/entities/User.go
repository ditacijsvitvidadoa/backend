package entities

type User struct {
	UserID            int               `bson:"Id" json:"id"`
	Password          string            `bson:"Password" json:"password"`
	FullName          FullName          `bson:"FullName" json:"full_name"`
	Phone             string            `bson:"Phone" json:"phone"`
	Email             string            `bson:"Email" json:"email"`
	PostalServiceInfo PostalServiceInfo `bson:"PostalService" json:"postal_service"`
	MarketingConsent  bool              `bson:"MarketingConsent" json:"marketing_consent"`
	Cart              []CartItem        `bson:"Cart" json:"cart"`
	Favourites        []string          `bson:"Favourites" json:"favourites"`
}

type FullName struct {
	FirstName  string `bson:"FirstName" json:"first_name"`
	LastName   string `bson:"LastName" json:"last_name"`
	Patronymic string `bson:"Patronymic" json:"patronymic"`
}

type PostalServiceInfo struct {
	PostalType    string      `bson:"PostalType" json:"postal_type"`
	City          string      `bson:"City" json:"city"`
	CityRef       string      `bson:"CityRef" json:"city_ref"`
	ReceivingType string      `bson:"ReceivingType" json:"receiving_type"`
	PostalInfo    interface{} `bson:"PostalInfo" json:"postal_info"`
}

type CartItem struct {
	Count   int              `bson:"Count" json:"count"`
	ID      string           `bson:"Id" json:"id"`
	Details *CartItemDetails `bson:"Details,omitempty" json:"details,omitempty"`
}

type CartItemDetails struct {
	Size  string `bson:"Size" json:"size"`
	Color string `bson:"Color" json:"color"`
}
