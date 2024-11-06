package entities

type Order struct {
	Status        int         `bson:"Status" json:"status"`
	UserId        string      `bson:"UserId" json:"user_id"`
	FirstName     string      `bson:"FirstName" json:"firstName"`
	LastName      string      `bson:"LastName" json:"lastName"`
	Patronymic    string      `bson:"Patronymic" json:"patronymic"`
	Phone         string      `bson:"Phone" json:"phone"`
	Email         string      `bson:"Email" json:"email"`
	PostalType    string      `bson:"PostalType" json:"postal_type"`
	City          string      `bson:"City" json:"city"`
	ReceivingType string      `bson:"ReceivingType" json:"receiving_type"`
	PostalInfo    interface{} `bson:"PostalInfo" json:"postal_info"`
	Products      []Product   `bson:"Products" json:"products"`
}

type CourierPostalInfo struct {
	Street    string `bson:"street" json:"street"`
	House     string `bson:"house" json:"house"`
	Apartment string `bson:"apartment" json:"apartment"`
	Floor     string `bson:"floor" json:"floor"`
}

type BranchPostalInfo struct {
	Branch string `bson:"branch" json:"branch"`
}
