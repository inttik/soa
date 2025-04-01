package main

import (
	"log"
	"net"
	mockstorage "posts/internal/mock_storage"
	"posts/internal/posts_grpc"
	postsservice "posts/internal/posts_service"

	"google.golang.org/grpc"
)

func main() {
	storage := mockstorage.NewMockStorage()
	server := postsservice.NewServer(&storage)

	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	posts_grpc.RegisterPostServiceServer(s, &server)
	log.Println("Start listening at :50001 (real)")
	log.Fatal(s.Serve(lis))
}
