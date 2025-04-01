package postsservice

import (
	"context"
	grpc "posts/internal/posts_grpc"
	storagemanager "posts/internal/storage_manager"

	"github.com/google/uuid"
)

type server struct {
	grpc.UnimplementedPostServiceServer
	manager storagemanager.StorageManager
}

func NewServer(manager storagemanager.StorageManager) server {
	return server{manager: manager}
}

var (
	CREATE_BAD_REQUEST = grpc.CreatePostResponse{Code: grpc.Code_BadRequest}
	UPDATE_BAD_REQUEST = grpc.UpdatePostResponse{Code: grpc.Code_BadRequest}
	DELETE_BAD_REQUEST = grpc.DeletePostResponse{Code: grpc.Code_BadRequest}
	GET_BAD_REQUEST    = grpc.GetPostResponse{Code: grpc.Code_BadRequest}
	LIST_BAD_REQUEST   = grpc.ListPostsResponse{Code: grpc.Code_BadRequest}
)

func check_uuid(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func (s *server) CreatePost(_ context.Context, req *grpc.CreatePostRequest) (*grpc.CreatePostResponse, error) {
	if req.Title == "" {
		return &CREATE_BAD_REQUEST, nil
	}
	if len(req.Title) > 255 {
		return &CREATE_BAD_REQUEST, nil
	}
	if req.Content == "" {
		return &CREATE_BAD_REQUEST, nil
	}
	if req.Actor == nil {
		return &CREATE_BAD_REQUEST, nil
	}
	if req.Actor.UserId == "" {
		return &CREATE_BAD_REQUEST, nil
	}
	if check_uuid(req.Actor.UserId) != nil {
		return &CREATE_BAD_REQUEST, nil
	}

	post, err := s.manager.CreatePost(req)
	if err != nil {
		return nil, err
	}
	return &grpc.CreatePostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func (s *server) UpdatePost(_ context.Context, req *grpc.UpdatePostRequest) (*grpc.UpdatePostResponse, error) {
	if req.Update == nil {
		return &UPDATE_BAD_REQUEST, nil
	}
	if req.Update.Id == nil {
		return &UPDATE_BAD_REQUEST, nil
	}
	if req.Update.Id.Id == "" {
		return &UPDATE_BAD_REQUEST, nil
	}
	if req.Actor == nil {
		return &UPDATE_BAD_REQUEST, nil
	}
	if req.Actor.UserId == "" {
		return &UPDATE_BAD_REQUEST, nil
	}
	if check_uuid(req.Update.Id.Id) != nil {
		return &UPDATE_BAD_REQUEST, nil
	}
	if check_uuid(req.Actor.UserId) != nil {
		return &UPDATE_BAD_REQUEST, nil
	}

	post, err := s.manager.GetPost(req.Update.Id)
	if err != nil {
		return &grpc.UpdatePostResponse{Code: grpc.Code_NotFound}, nil
	}
	if post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.UpdatePostResponse{Code: grpc.Code_Forbidden}, nil
	}
	post, err = s.manager.UpdatePost(req.Update)
	if err != nil {
		return nil, err
	}
	return &grpc.UpdatePostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func (s *server) DeletePost(_ context.Context, req *grpc.DeletePostRequest) (*grpc.DeletePostResponse, error) {
	if req.Id == nil {
		return &DELETE_BAD_REQUEST, nil
	}
	if req.Id.Id == "" {
		return &DELETE_BAD_REQUEST, nil
	}
	if req.Actor == nil {
		return &DELETE_BAD_REQUEST, nil
	}
	if req.Actor.UserId == "" {
		return &DELETE_BAD_REQUEST, nil
	}
	if check_uuid(req.Id.Id) != nil {
		return &DELETE_BAD_REQUEST, nil
	}
	if check_uuid(req.Actor.UserId) != nil {
		return &DELETE_BAD_REQUEST, nil
	}

	post, err := s.manager.GetPost(req.Id)
	if err != nil {
		return &grpc.DeletePostResponse{Code: grpc.Code_NotFound}, nil
	}
	if post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.DeletePostResponse{Code: grpc.Code_Forbidden}, nil
	}
	err = s.manager.DeletePost(post.Id)
	if err != nil {
		return nil, err
	}
	return &grpc.DeletePostResponse{Code: grpc.Code_Ok}, nil
}

func (s *server) GetPost(_ context.Context, req *grpc.GetPostRequest) (*grpc.GetPostResponse, error) {
	if req.Id == nil {
		return &GET_BAD_REQUEST, nil
	}
	if req.Id.Id == "" {
		return &GET_BAD_REQUEST, nil
	}
	if req.Actor == nil {
		return &GET_BAD_REQUEST, nil
	}
	if req.Actor.UserId == "" {
		return &GET_BAD_REQUEST, nil
	}
	if check_uuid(req.Id.Id) != nil {
		return &GET_BAD_REQUEST, nil
	}
	if check_uuid(req.Actor.UserId) != nil {
		return &GET_BAD_REQUEST, nil
	}

	post, err := s.manager.GetPost(req.Id)
	if err != nil {
		return &grpc.GetPostResponse{Code: grpc.Code_NotFound}, nil
	}
	if post.IsPrivate && post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.GetPostResponse{Code: grpc.Code_Forbidden}, nil
	}
	return &grpc.GetPostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func (s *server) ListPosts(_ context.Context, req *grpc.ListPostsRequest) (*grpc.ListPostsResponse, error) {
	if req.Actor == nil {
		return &LIST_BAD_REQUEST, nil
	}
	if req.Actor.UserId == "" {
		return &LIST_BAD_REQUEST, nil
	}
	if check_uuid(req.Actor.UserId) != nil {
		return &LIST_BAD_REQUEST, nil
	}

	resp, err := s.manager.ListPosts(req)
	if err != nil {
		return nil, err
	}
	resp.Code = grpc.Code_Ok
	return resp, nil
}
