package main

import (
	"log"
	"net"
	postgresstorage "posts/internal/postgres_storage"
	"posts/internal/posts_grpc"
	postsservice "posts/internal/posts_service"

	"google.golang.org/grpc"
)

func main() {
	// storage := mockstorage.NewMockStorage()
	storage, err := postgresstorage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := postsservice.NewServer(storage)

	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	posts_grpc.RegisterPostServiceServer(s, &server)
	log.Println("Start listening at :50001 (real)")
	log.Fatal(s.Serve(lis))
}
