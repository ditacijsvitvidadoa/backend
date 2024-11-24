package entities

type AdminCredentials struct {
	Email    string `bson:"Email"`
	Password string `bson:"Password"`
	Phone    string `bson:"Phone"`
}
