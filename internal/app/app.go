package app

import (
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type App struct {
	client *mongo.Client
}

func NewApp(client *mongo.Client) *App {
	app := &App{
		client: client,
	}

	return app
}

func (a *App) GetRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/api/get-products", a.getProducts)

	return r
}
