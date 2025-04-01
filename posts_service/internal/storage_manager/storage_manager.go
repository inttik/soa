package storagemanager

import "posts/internal/posts_grpc"

type StorageManager interface {
	// Create post in DB and returns it. Raise no error
	CreatePost(*posts_grpc.CreatePostRequest) (*posts_grpc.Post, error)
	// Update post in DB and returns it. Raise no error
	UpdatePost(*posts_grpc.PostUpdate) (*posts_grpc.Post, error)
	// Delete post in DB. Raise no error
	DeletePost(*posts_grpc.PostId) error
	// Get post by id in DB. Raise 404, if post not found
	GetPost(*posts_grpc.PostId) (*posts_grpc.Post, error)
	// Get list of posts. Raise no error
	ListPosts(*posts_grpc.ListPostsRequest) (*posts_grpc.ListPostsResponse, error)
}
