package main

import (
	"log"
	"net/http"

	handler "users/handlers"
	jwttoken "users/internal/jwt_token"
	mockstorage "users/internal/mock_storage"
	passhandle "users/internal/passhandler"
	"users/oas"
)

func main() {
	err := jwttoken.SetupEnv("secrets/")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := mockstorage.NewMockStorage()
	if err != nil {
		log.Fatal(err)
	}

	adminLogin := oas.LoginString("admin")
	adminPass := oas.PasswordString("admin")
	hashedPass := passhandle.HashPass(adminLogin, adminPass)
	err = storage.MakeRootUser(adminLogin, oas.PasswordString(hashedPass))
	if err != nil {
		log.Fatal(err)
	}

	service, err := handler.NewService(&storage)
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
	mux.Handle("/v1/", http.StripPrefix("/v1", srv))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
