package main

import (
	"IT/DutyachiySvit/backend/internal/app"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
		fmt.Println("port -> 8080")
	}

	router := app.GetRouter()

	err := http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println(err)
	}
}
