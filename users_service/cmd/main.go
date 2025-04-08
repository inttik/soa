package main

import (
	"log"
	"net/http"

	handler "users/handlers"
	jwttoken "users/internal/jwt_token"
	passhandle "users/internal/pass_handler"
	postgresstorage "users/internal/postgres_storage"
	"users/oas"
)

func main() {
	err := jwttoken.SetupEnv("secrets/")
	if err != nil {
		log.Fatal(err)
	}

	// storage, err := mockstorage.NewMockStorage()
	storage, err := postgresstorage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	adminLogin := oas.LoginString("admin")
	adminPass := oas.PasswordString("admin")
	hashedPass, err := passhandle.HashPassword(adminPass)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.MakeRootUser(adminLogin, oas.PasswordString(hashedPass))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("User created")

	service, err := handler.NewService(storage)
	if err != nil {
		log.Fatal(err)
	}

	security, err := handler.NewSecurityHandler()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := oas.NewServer(&service, &security)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/users/v1/", http.StripPrefix("/users/v1", srv))

	log.Println("starting server")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
