package postgresstorage

import (
	"errors"
	"posts/internal/posts_grpc"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNoPost = errors.New("post not found")
)

func (ps *postgresStorage) CreatePost(req *posts_grpc.CreatePostRequest) (*posts_grpc.Post, error) {
	newPost := postTable{
		PostId:      uuid.New(),
		AuthorId:    uuid.MustParse(req.Actor.UserId),
		Title:       req.Title,
		Content:     req.Content,
		IsPrivate:   req.IsPrivate,
		Tags:        nil,
		PublishDate: time.Now(),
		LastModify:  time.Now(),
	}

	err := ps.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&newPost)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return toPost(&newPost), nil
}

func (ps *postgresStorage) UpdatePost(req *posts_grpc.PostUpdate) (*posts_grpc.Post, error) {
	post := postTable{}
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("post_id = ?", req.Id).Limit(1).Find(&post)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoPost
		}

		if req.NewTitle != "" {
			post.Title = req.NewTitle
		}
		if req.NewContent != "" {
			post.Content = req.NewContent
		}
		if req.ChangePrivate {
			post.IsPrivate = req.IsPrivate
		}

		cur_tags := make(map[string]struct{})
		for _, str := range post.Tags {
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

		result = tx.Save(&post)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoPost
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return toPost(&post), nil
}

func (ps *postgresStorage) DeletePost(id string) error {
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("post_id = ?", id).Delete(&postTable{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoPost
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (ps *postgresStorage) GetPost(id string) (*posts_grpc.Post, error) {
	post := postTable{}
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("post_id = ?", id).Limit(1).Find(&post)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNoPost
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return toPost(&post), nil
}

// Quite unoptimal
func (ps *postgresStorage) ListPosts(req *posts_grpc.ListPostsRequest) (*posts_grpc.ListPostsResponse, error) {
	var posts []postTable

	result := ps.db.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	var index = uint32(0)
	var from = req.PageLimit * req.Page
	var to = req.PageLimit * (req.Page + 1)

	resp := posts_grpc.ListPostsResponse{}
	resp.From = 0
	resp.To = 0
	resp.Posts = make([]*posts_grpc.Post, 0, req.PageLimit)
	resp.Total = 0

	for _, post := range posts {
		if post.IsPrivate {
			if !req.Actor.IsRoot && req.Actor.UserId != post.AuthorId.String() {
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
			resp.Posts = append(resp.Posts, toPost(&post))
			resp.To = index + 1
		}
		index += 1
	}

	return &resp, nil
}
