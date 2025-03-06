package main

import (
	"log"
	"net/http"

	handler "users/handlers"
	mockstorage "users/internal/mock_storage"
	passhandle "users/internal/pass_handle"
	"users/oas"
)

func main() {
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
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
