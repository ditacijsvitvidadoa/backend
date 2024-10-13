package main

import (
	"context"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/app"
	cash2 "github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/ditacijsvitvidadoa/backend/internal/mongo_conn"
	"log"
	"net/http"
	"os"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	fmt.Println("port -> ", port)

	cash, err := cash2.RedisConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer cash.Close()

	client, err := mongo_conn.MongoConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	a := app.NewApp(client, cash)
	router := a.GetRouter()

	fmt.Println("Listening on port:", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Println(err)
	}
}
