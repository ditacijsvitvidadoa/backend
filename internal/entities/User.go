package entities

type User struct {
	UserID           int           `bson:"Id" json:"id"`
	Password         string        `bson:"Password" json:"password"`
	FullName         FullName      `bson:"FullName" json:"full_name"`
	Prone            string        `bson:"Prone" json:"prone"`
	Email            string        `bson:"Email" json:"email"`
	PostalService    PostalService `bson:"PostalService" json:"postal_service"`
	MarketingConsent bool          `bson:"MarketingConsent" json:"marketing_consent"`
	Cart             []CartItem    `bson:"Cart" json:"cart"`
	Favourites       []int         `bson:"Favourites" json:"favourites"`
}

type FullName struct {
	FirstName  string `bson:"FirstName" json:"first_name"`
	LastName   string `bson:"LastName" json:"last_name"`
	Patronymic string `bson:"Patronymic" json:"patronymic"`
}

type PostalService struct {
	PostalName string `bson:"PostalName" json:"postal_name"`
	Branch     string `bson:"Branch" json:"branch"`
}

type CartItem struct {
	Count int `bson:"Count" json:"count"`
	ID    int `bson:"Id" json:"id"`
}
