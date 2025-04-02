package mockstorage

import (
	"errors"
	"posts/internal/posts_grpc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockStorage struct {
	p posts
}

func NewMockStorage() mockStorage {
	return mockStorage{posts{data: make(map[string]*posts_grpc.Post)}}
}

func (s *mockStorage) CreatePost(req *posts_grpc.CreatePostRequest) (*posts_grpc.Post, error) {
	new_uuid := uuid.New().String()

	post := posts_grpc.Post{
		Id:          new_uuid,
		Title:       req.Title,
		Content:     req.Content,
		AuthorId:    req.Actor.UserId,
		IsPrivate:   req.IsPrivate,
		PublishDate: timestamppb.Now(),
		LastModify:  timestamppb.Now(),
		Tags:        req.Tags,
	}

	s.p.mx.Lock()
	defer s.p.mx.Unlock()

	s.p.data[new_uuid] = &post

	return &post, nil
}

func (s *mockStorage) UpdatePost(req *posts_grpc.PostUpdate) (*posts_grpc.Post, error) {
	s.p.mx.Lock()
	defer s.p.mx.Unlock()

	post := s.p.data[req.GetId()]
	if req.GetNewTitle() != "" {
		post.Title = req.GetNewTitle()
	}
	if req.GetNewContent() != "" {
		post.Content = req.GetNewContent()
	}
	if req.GetChangePrivate() == true {
		post.IsPrivate = req.GetIsPrivate()
	}

	cur_tags := make(map[string]struct{})
	for _, str := range post.GetTags() {
		cur_tags[str] = struct{}{}
	}
	for _, str := range req.GetRemoveTags() {
		delete(cur_tags, str)
	}
	for _, str := range req.GetAddTags() {
		cur_tags[str] = struct{}{}
	}

	tags := make([]string, 0, len(cur_tags))
	for tag := range cur_tags {
		tags = append(tags, tag)
	}
	post.Tags = tags

	return post, nil
}

func (s *mockStorage) DeletePost(id string) error {
	s.p.mx.Lock()
	defer s.p.mx.Unlock()

	delete(s.p.data, id)

	return nil
}

func (s *mockStorage) GetPost(id string) (*posts_grpc.Post, error) {
	s.p.mx.Lock()
	defer s.p.mx.Unlock()

	post, ok := s.p.data[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *mockStorage) ListPosts(req *posts_grpc.ListPostsRequest) (*posts_grpc.ListPostsResponse, error) {
	s.p.mx.Lock()
	defer s.p.mx.Unlock()

	var index = uint32(0)
	var from = req.PageLimit * req.Page
	var to = req.PageLimit * (req.Page + 1)

	resp := posts_grpc.ListPostsResponse{}
	resp.From = 0
	resp.To = 0
	resp.Posts = make([]*posts_grpc.Post, 0, req.PageLimit)
	resp.Total = 0

	for _, post := range s.p.data {
		if post.IsPrivate {
			if !req.Actor.IsRoot && req.Actor.UserId != post.AuthorId {
				continue
			}
			if !req.WithHidden {
				continue
			}
		}
		resp.Total += 1
		if index >= from && index < to {
			if len(resp.Posts) == 0 {
				resp.From = index
			}
			resp.Posts = append(resp.Posts, post)
			resp.To = index + 1
		}
		index += 1
	}
	return &resp, nil
}
