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

func CreatePostBad(text string) (*grpc.CreatePostResponse, error) {
	resp := grpc.CreatePostResponse{
		Code:  grpc.Code_BadRequest,
		Error: text,
	}
	return &resp, nil
}

func (s *server) CreatePost(_ context.Context, req *grpc.CreatePostRequest) (*grpc.CreatePostResponse, error) {
	if req.Title == "" {
		return CreatePostBad("Field 'title' not found")
	}
	if len(req.Title) > 255 {
		return CreatePostBad("Title too long")
	}
	if req.Content == "" {
		return CreatePostBad("Field 'content' not found")
	}
	if req.Actor == nil {
		return CreatePostBad("Actor not found")
	}
	if req.Actor.UserId == "" {
		return CreatePostBad("User id not found")
	}
	if check_uuid(req.Actor.UserId) != nil {
		return CreatePostBad("User id is not uuid")
	}

	post, err := s.manager.CreatePost(req)
	if err != nil {
		return nil, err
	}
	return &grpc.CreatePostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func UpdatePostBad(text string) (*grpc.UpdatePostResponse, error) {
	resp := grpc.UpdatePostResponse{
		Code:  grpc.Code_BadRequest,
		Error: text,
	}
	return &resp, nil
}

func (s *server) UpdatePost(_ context.Context, req *grpc.UpdatePostRequest) (*grpc.UpdatePostResponse, error) {
	if req.Update == nil {
		return UpdatePostBad("No update")
	}
	if req.Update.Id == "" {
		return UpdatePostBad("Post id not found")
	}
	if req.Actor == nil {
		return UpdatePostBad("Actor not found")
	}
	if req.Actor.UserId == "" {
		return UpdatePostBad("User id not found")
	}
	if check_uuid(req.Update.Id) != nil {
		return UpdatePostBad("Post id is not UUID")
	}
	if check_uuid(req.Actor.UserId) != nil {
		return UpdatePostBad("User id is not UUID")
	}

	post, err := s.manager.GetPost(req.Update.Id)
	if err != nil {
		return &grpc.UpdatePostResponse{
			Code:  grpc.Code_NotFound,
			Error: "Post not found",
		}, nil
	}
	if post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.UpdatePostResponse{
			Code:  grpc.Code_Forbidden,
			Error: "User can't modify the post",
		}, nil
	}
	post, err = s.manager.UpdatePost(req.Update)
	if err != nil {
		return nil, err
	}
	return &grpc.UpdatePostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func DeletePostBad(text string) (*grpc.DeletePostResponse, error) {
	resp := grpc.DeletePostResponse{
		Code:  grpc.Code_BadRequest,
		Error: text,
	}
	return &resp, nil
}

func (s *server) DeletePost(_ context.Context, req *grpc.DeletePostRequest) (*grpc.DeletePostResponse, error) {
	if req.Id == "" {
		return DeletePostBad("Post id not found")
	}
	if req.Actor == nil {
		return DeletePostBad("Actor not found")
	}
	if req.Actor.UserId == "" {
		return DeletePostBad("User id not found")
	}
	if check_uuid(req.Id) != nil {
		return DeletePostBad("Post id is not UUID")
	}
	if check_uuid(req.Actor.UserId) != nil {
		return DeletePostBad("User id is not UUID")
	}

	post, err := s.manager.GetPost(req.Id)
	if err != nil {
		return &grpc.DeletePostResponse{
			Code:  grpc.Code_NotFound,
			Error: "Post not found",
		}, nil
	}
	if post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.DeletePostResponse{
			Code:  grpc.Code_Forbidden,
			Error: "User can't delete post",
		}, nil
	}
	err = s.manager.DeletePost(post.Id)
	if err != nil {
		return nil, err
	}
	return &grpc.DeletePostResponse{Code: grpc.Code_Ok}, nil
}

func GetPostBad(text string) (*grpc.GetPostResponse, error) {
	resp := grpc.GetPostResponse{
		Code:  grpc.Code_BadRequest,
		Error: text,
	}
	return &resp, nil
}

func (s *server) GetPost(_ context.Context, req *grpc.GetPostRequest) (*grpc.GetPostResponse, error) {
	if req.Id == "" {
		return GetPostBad("Post id not found")
	}
	if req.Actor == nil {
		return GetPostBad("Actor not found")
	}
	if req.Actor.UserId == "" {
		return GetPostBad("User id not found")
	}
	if check_uuid(req.Id) != nil {
		return GetPostBad("Post id is not UUID")
	}
	if check_uuid(req.Actor.UserId) != nil {
		return GetPostBad("User id is not UUID")
	}

	post, err := s.manager.GetPost(req.Id)
	if err != nil {
		return &grpc.GetPostResponse{
			Code:  grpc.Code_NotFound,
			Error: "Post not found",
		}, nil
	}
	if post.IsPrivate && post.AuthorId != req.Actor.UserId && !req.Actor.IsRoot {
		return &grpc.GetPostResponse{
			Code:  grpc.Code_Forbidden,
			Error: "User can't get that post",
		}, nil
	}
	return &grpc.GetPostResponse{Code: grpc.Code_Ok, Post: post}, nil
}

func ListPostBad(text string) (*grpc.ListPostsResponse, error) {
	resp := grpc.ListPostsResponse{
		Code:  grpc.Code_BadRequest,
		Error: text,
	}
	return &resp, nil
}

func (s *server) ListPosts(_ context.Context, req *grpc.ListPostsRequest) (*grpc.ListPostsResponse, error) {
	if req.Actor == nil {
		return ListPostBad("Actor not found")
	}
	if req.Actor.UserId == "" {
		return ListPostBad("User id not found")
	}
	if check_uuid(req.Actor.UserId) != nil {
		return ListPostBad("User id is not UUID")
	}

	resp, err := s.manager.ListPosts(req)
	if err != nil {
		return nil, err
	}
	resp.Code = grpc.Code_Ok
	return resp, nil
}
