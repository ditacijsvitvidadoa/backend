package app

import "net/http"

func GetRouter() *http.ServeMux {
	r := http.NewServeMux()

	return r
}
