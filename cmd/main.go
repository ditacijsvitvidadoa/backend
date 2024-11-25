package main

import (
	"context"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/app"
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/ditacijsvitvidadoa/backend/internal/mongo_conn"
	"github.com/ditacijsvitvidadoa/backend/internal/ticker"
	"log"
	"net/http"
	"os"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	fmt.Println("port ->", port)

	cache, err := cash.RedisConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	client, err := mongo_conn.MongoConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	a := app.NewApp(client, cache)
	router := a.GetRouter()

	go ticker.GeneralTicker(client)

	fmt.Println("Listening on port:", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
