package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type App struct {
	client *mongo.Client
	cash   *cash.RedisClient
}

func NewApp(client *mongo.Client, cash *cash.RedisClient) *App {
	app := &App{
		client: client,
		cash:   cash,
	}

	return app
}

func (a *App) GetRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /api/get-products", a.getProducts)
	r.HandleFunc("GET /api/check-auth", a.checkAuthentication)
	r.HandleFunc("GET /api/get-cart-products", a.getCartProducts)

	r.HandleFunc("GET /api/login", a.logIn)
	r.HandleFunc("GET /api/user/account", a.getProfileInfo)
	r.HandleFunc("POST /api/create-account", a.createUserAccount)
	r.HandleFunc("POST /api/logout", a.logout)

	r.HandleFunc("PUT /api/account-update/firstname", a.updateFirstName)
	r.HandleFunc("PUT /api/account-update/lastname", a.updateLastName)
	r.HandleFunc("PUT /api/account-update/patronymic", a.updatePatronymic)
	r.HandleFunc("PUT /api/account-update/phone", a.updatePhoneNumber)
	r.HandleFunc("PUT /api/account-update/email", a.updateEmail)
	r.HandleFunc("PUT /api/account-update/password", a.updatePassword)

	r.HandleFunc("DELETE /api/delete-cart-product/{id}", a.deleteCartProduct)
	r.HandleFunc("POST /api/add-product-to-cart/{id}", a.addCartProduct)

	r.HandleFunc("POST /api/add-favourite-product/{id}", a.addFavouriteProduct)
	r.HandleFunc("DELETE /api/delete-favourite-product/{id}", a.deleteFavouriteProduct)

	return r
}
