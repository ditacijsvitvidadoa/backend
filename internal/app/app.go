package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/gorilla/mux"
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
	r := mux.NewRouter()

	r.HandleFunc("/api/get-products", a.getProducts).Methods("GET")
	r.HandleFunc("/api/get-product/{id}", a.getProductByID).Methods("GET")
	r.HandleFunc("/api/check-auth", a.checkAuthentication).Methods("GET")
	r.HandleFunc("/api/get-cart-products", a.getCartProducts).Methods("GET")
	r.HandleFunc("/api/create-product", a.CreateProduct).Methods("POST")
	r.HandleFunc("/api/get-products-filter", a.getProductsFilter).Methods("GET")
	r.HandleFunc("/api/update-product-analytics", a.updateProductAnalytics).Methods("PUT")

	r.HandleFunc("/api/login", a.logIn).Methods("POST")
	r.HandleFunc("/api/user-account", a.getProfileInfo).Methods("GET")
	r.HandleFunc("/api/create-account", a.createUserAccount).Methods("POST")
	r.HandleFunc("/api/logout", a.logout).Methods("POST")

	r.HandleFunc("/api/account-update/firstname", a.updateFirstName).Methods("POST")
	r.HandleFunc("/api/account-update/lastname", a.updateLastName).Methods("POST")
	r.HandleFunc("/api/account-update/patronymic", a.updatePatronymic).Methods("POST")
	r.HandleFunc("/api/account-update/phone", a.updatePhoneNumber).Methods("POST")
	r.HandleFunc("/api/account-update/email", a.updateEmail).Methods("POST")
	r.HandleFunc("/api/account-update/password", a.updatePassword).Methods("POST")
	r.HandleFunc("/api/update-postal-info", a.addOrUpdatePostalServiceInfo).Methods("POST")
	r.HandleFunc("/api/marketing-consent", a.updateMarketingConsent).Methods("POST")

	r.HandleFunc("/api/get-purchases-history", a.PurchasesHistory).Methods("GET")

	r.HandleFunc("/api/send-to-support", a.sendToSupport).Methods("POST")

	r.HandleFunc("/api/delete-cart-product/{id}", a.deleteCartProduct).Methods("DELETE")
	r.HandleFunc("/api/add-product-to-cart/{id}", a.addCartProduct).Methods("PUT")

	r.HandleFunc("/api/update-cart-product-count", a.UpdateCount).Methods("PUT")

	r.HandleFunc("/api/get-favoutires-products", a.GetFavouritesProducts).Methods("GET")
	r.HandleFunc("/api/add-favourite-product/{id}", a.addFavouriteProduct).Methods("PUT")
	r.HandleFunc("/api/delete-favourite-product/{id}", a.deleteFavouriteProduct).Methods("DELETE")

	r.HandleFunc("/api/add-order", a.AddOrder).Methods("PUT")
	r.HandleFunc("/api/delete-order/{id}", a.DeleteOrder).Methods("DELETE")
	r.HandleFunc("/api/get-orders", a.GetOrders).Methods("GET")
	r.HandleFunc("/api/change-order-status/{id}", a.ChangeOrderStatus).Methods("PUT")
	r.HandleFunc("/api/archive-order/{id}", a.ArchiveOrder).Methods("DELETE")

	r.HandleFunc("/api/get-archive-orders", a.GetArchiveOrders).Methods("GET")
	r.HandleFunc("/api/refresh-archive-order/{id}", a.RefreshArchiveOrder).Methods("DELETE")

	r.HandleFunc("/api/get-cities", a.GetAllCities).Methods("GET")
	r.HandleFunc("/api/get-postals/{city_ref}", a.GetPostalsFromCity).Methods("GET")

	r.HandleFunc("/api/login-admin-panel", a.loginAdminPanel).Methods("POST")
	r.HandleFunc("/api/check-admin-auth", a.checkAdminSession).Methods("GET")

	staticDir := "/app/static/frontend"
	productsDir := "/app/static/products"

	fileServerFrontend := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/frontend/").Handler(http.StripPrefix("/static/frontend/", fileServerFrontend))

	fileServerProducts := http.FileServer(http.Dir(productsDir))
	r.PathPrefix("/static/products/").Handler(http.StripPrefix("/static/products/", fileServerProducts))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	return r
}
