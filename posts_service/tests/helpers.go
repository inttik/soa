package unit_test

import (
	"context"
	mockstorage "posts/internal/mock_storage"
	"posts/internal/posts_grpc"
	postsservice "posts/internal/posts_service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type state struct {
	server posts_grpc.PostServiceServer
	user_a string
	user_b string
}

func (s *state) setup(t *testing.T) {
	storage := mockstorage.NewMockStorage()
	server := postsservice.NewServer(storage)
	s.server = &server
	s.user_a = uuid.NewString()
	s.user_b = uuid.NewString()
}

// CreateHelper для тестирования CreatePost
type CreateHelper struct {
	req    *posts_grpc.CreatePostRequest
	user   string
	root   bool
	exp200 bool
	exp400 bool
	exp403 bool
	exp404 bool
}

func (s *state) applyCreate(h CreateHelper, t *testing.T) *posts_grpc.CreatePostResponse {
	if h.user != "" {
		h.req.Actor = &posts_grpc.Actor{
			UserId: h.user,
			IsRoot: h.root,
		}
	}
	resp, err := s.server.CreatePost(context.Background(), h.req)

	switch {
	case h.exp400:
		assert.Equal(t, posts_grpc.Code_BadRequest, resp.Code)
	case h.exp403:
		assert.Equal(t, posts_grpc.Code_Forbidden, resp.Code)
	case h.exp404:
		assert.Equal(t, posts_grpc.Code_NotFound, resp.Code)
	case h.exp200:
		assert.NoError(t, err)
		assert.Equal(t, posts_grpc.Code_Ok, resp.Code)
		assert.Equal(t, h.req.Title, resp.Post.Title)
		assert.Equal(t, h.req.Content, resp.Post.Content)
		assert.Equal(t, h.req.IsPrivate, resp.Post.IsPrivate)
		assert.Equal(t, h.req.Tags, resp.Post.Tags)
		assert.Equal(t, h.req.Actor.UserId, resp.Post.AuthorId)
		return resp
	default:
		t.Fatal("No expected result specified in helper")
	}
	return nil
}

// UpdateHelper для тестирования UpdatePost
type UpdateHelper struct {
	req    *posts_grpc.UpdatePostRequest
	user   string
	root   bool
	exp200 bool
	exp400 bool
	exp403 bool
	exp404 bool
}

func (s *state) applyUpdate(h UpdateHelper, t *testing.T) *posts_grpc.UpdatePostResponse {
	if h.user != "" {
		h.req.Actor = &posts_grpc.Actor{
			UserId: h.user,
			IsRoot: h.root,
		}
	}
	resp, err := s.server.UpdatePost(context.Background(), h.req)

	var old *posts_grpc.GetPostResponse = nil
	if h.exp200 {
		old, err = s.server.GetPost(context.Background(), &posts_grpc.GetPostRequest{
			Id: h.req.Update.Id,
			Actor: &posts_grpc.Actor{
				UserId: s.user_a,
				IsRoot: true,
			},
		})
		assert.NoError(t, err)
	}

	switch {
	case h.exp400:
		assert.Equal(t, posts_grpc.Code_BadRequest, resp.Code)
	case h.exp403:
		assert.Equal(t, posts_grpc.Code_Forbidden, resp.Code)
	case h.exp404:
		assert.Equal(t, posts_grpc.Code_NotFound, resp.Code)
	case h.exp200:
		assert.NoError(t, err)
		assert.Equal(t, posts_grpc.Code_Ok, resp.Code)
		if h.req.Update.NewTitle != "" {
			assert.Equal(t, h.req.Update.NewTitle, resp.Post.Title)
		} else {
			assert.Equal(t, old.Post.Title, resp.Post.Title)
		}
		if h.req.Update.NewContent != "" {
			assert.Equal(t, h.req.Update.NewContent, resp.Post.Content)
		} else {
			assert.Equal(t, old.Post.Content, resp.Post.Content)
		}
		if h.req.Update.ChangePrivate {
			assert.Equal(t, h.req.Update.IsPrivate, resp.Post.IsPrivate)
		} else {
			assert.Equal(t, old.Post.IsPrivate, resp.Post.IsPrivate)
		}
		for _, tag := range old.Post.Tags {
			if slices.Contains(h.req.Update.RemoveTags, tag) {
				continue
			}
			assert.Contains(t, resp.Post.Tags, tag)
		}
		for _, tag := range h.req.Update.AddTags {
			assert.Contains(t, resp.Post.Tags, tag)
		}
		for _, tag := range h.req.Update.RemoveTags {
			assert.NotContains(t, resp.Post.Tags, tag)
		}
		return resp
	default:
		t.Fatal("No expected result specified in helper")
	}
	return nil
}

// DeleteHelper для тестирования DeletePost
type DeleteHelper struct {
	req    *posts_grpc.DeletePostRequest
	user   string
	root   bool
	exp200 bool
	exp400 bool
	exp403 bool
	exp404 bool
}

func (s *state) applyDelete(h DeleteHelper, t *testing.T) *posts_grpc.DeletePostResponse {
	if h.user != "" {
		h.req.Actor = &posts_grpc.Actor{
			UserId: h.user,
			IsRoot: h.root,
		}
	}
	resp, err := s.server.DeletePost(context.Background(), h.req)

	switch {
	case h.exp400:
		assert.Equal(t, posts_grpc.Code_BadRequest, resp.Code)
	case h.exp403:
		assert.Equal(t, posts_grpc.Code_Forbidden, resp.Code)
	case h.exp404:
		assert.Equal(t, posts_grpc.Code_NotFound, resp.Code)
	case h.exp200:
		assert.NoError(t, err)
		assert.Equal(t, posts_grpc.Code_Ok, resp.Code)
		return resp
	default:
		t.Fatal("No expected result specified in helper")
	}
	return nil
}

// GetHelper для тестирования GetPost
type GetHelper struct {
	req    *posts_grpc.GetPostRequest
	user   string
	root   bool
	exp200 bool
	exp400 bool
	exp403 bool
	exp404 bool
}

func (s *state) applyGet(h GetHelper, t *testing.T) *posts_grpc.GetPostResponse {
	if h.user != "" {
		h.req.Actor = &posts_grpc.Actor{
			UserId: h.user,
			IsRoot: h.root,
		}
	}

	resp, err := s.server.GetPost(context.Background(), h.req)

	switch {
	case h.exp400:
		assert.Equal(t, posts_grpc.Code_BadRequest, resp.Code)
	case h.exp403:
		assert.Equal(t, posts_grpc.Code_Forbidden, resp.Code)
	case h.exp404:
		assert.Equal(t, posts_grpc.Code_NotFound, resp.Code)
	case h.exp200:
		assert.NoError(t, err)
		assert.Equal(t, posts_grpc.Code_Ok, resp.Code)
		return resp
	default:
		t.Fatal("No expected result specified in helper")
	}
	return nil
}

// ListHelper для тестирования ListPosts
type ListHelper struct {
	req          *posts_grpc.ListPostsRequest
	user         string
	root         bool
	exp200       bool
	exp400       bool
	exp403       bool
	exp404       bool
	exp_total    uint32
	exp_from     uint32
	exp_to       uint32
	exp_pub_only bool
}

func (s *state) applyList(h ListHelper, t *testing.T) *posts_grpc.ListPostsResponse {
	if h.user != "" {
		h.req.Actor = &posts_grpc.Actor{
			UserId: h.user,
			IsRoot: h.root,
		}
	}

	resp, err := s.server.ListPosts(context.Background(), h.req)

	switch {
	case h.exp400:
		assert.Equal(t, posts_grpc.Code_BadRequest, resp.Code)
	case h.exp403:
		assert.Equal(t, posts_grpc.Code_Forbidden, resp.Code)
	case h.exp404:
		assert.Equal(t, posts_grpc.Code_NotFound, resp.Code)
	case h.exp200:
		assert.NoError(t, err)
		assert.Equal(t, posts_grpc.Code_Ok, resp.Code)
		assert.Equal(t, h.exp_from, resp.From)
		assert.Equal(t, h.exp_to, resp.To)
		assert.Equal(t, h.exp_total, resp.Total)
		if h.exp_pub_only {
			for _, post := range resp.Posts {
				assert.False(t, post.IsPrivate)
			}
		}
		return resp
	default:
		t.Fatal("No expected result specified in helper")
	}
	return nil
}
