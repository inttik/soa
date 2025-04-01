package main

import (
	"gateway/internal/handlers"

	"github.com/google/uuid"
)

func main() {
	_ = uuid.UUID{}

	server := handlers.NewGatewayServer()
	server.Start("8080")
}
