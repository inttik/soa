package project

//go:generate protoc --go_out=./internal/posts_grpc/ --go-grpc_out=./internal/posts_grpc/ posts.proto
