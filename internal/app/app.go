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

func (a *App) GetRouter() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /api/get-products", a.getProducts)
	r.HandleFunc("GET /api/get-product/{id}", a.getProductByID)
	r.HandleFunc("GET /api/check-auth", a.checkAuthentication)
	r.HandleFunc("GET /api/get-cart-products", a.getCartProducts)

	r.HandleFunc("POST /api/login", a.logIn)
	r.HandleFunc("GET /api/user-account", a.getProfileInfo)
	r.HandleFunc("POST /api/create-account", a.createUserAccount)
	r.HandleFunc("POST /api/logout", a.logout)

	r.HandleFunc("POST /api/account-update/firstname", a.updateFirstName)
	r.HandleFunc("POST /api/account-update/lastname", a.updateLastName)
	r.HandleFunc("POST /api/account-update/patronymic", a.updatePatronymic)
	r.HandleFunc("POST /api/account-update/phone", a.updatePhoneNumber)
	r.HandleFunc("POST /api/account-update/email", a.updateEmail)
	r.HandleFunc("POST /api/account-update/password", a.updatePassword)

	r.HandleFunc("DELETE /api/delete-cart-product/{id}", a.deleteCartProduct)
	r.HandleFunc("PUT /api/add-product-to-cart/{id}", a.addCartProduct)

	r.HandleFunc("PUT /api/add-favourite-product/{id}", a.addFavouriteProduct)
	r.HandleFunc("DELETE /api/delete-favourite-product/{id}", a.deleteFavouriteProduct)

	return corsMiddleware(r)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
