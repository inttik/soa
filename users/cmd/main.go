package main

import (
	"log"
	"net/http"
	handler "users/handlers"
	api "users/oas"
)

func main() {
	service := handler.NewService()
	security := handler.NewSecurityHandler()

	srv, err := api.NewServer(&service, &security)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
