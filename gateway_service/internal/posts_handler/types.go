package posthandler

import (
	"gateway/internal/posts_grpc"
	"time"
)

type postResponse struct {
	PostId      string    `json:"post_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
	AuthorId    string    `json:"author_id"`
	PublishDate time.Time `json:"publish_date"`
	LastMoidfy  time.Time `json:"last_modify"`
}

func getPostResponse(post *posts_grpc.Post) postResponse {
	return postResponse{
		PostId:      post.Id,
		Title:       post.Title,
		Content:     post.Content,
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
		AuthorId:    post.AuthorId,
		PublishDate: post.PublishDate.AsTime(),
		LastMoidfy:  post.LastModify.AsTime(),
	}
}

type feedResponse struct {
	Posts []postResponse `json:"posts"`
	From  uint32         `json:"from"`
	To    uint32         `json:"to"`
	Total uint32         `json:"total"`
}

func getFeedResponse(posts *posts_grpc.ListPostsResponse) feedResponse {
	var resp feedResponse

	resp.Posts = make([]postResponse, 0, len(posts.Posts))
	for _, p := range posts.Posts {
		resp.Posts = append(resp.Posts, getPostResponse(p))
	}

	resp.From = posts.From
	resp.To = posts.To
	resp.Total = posts.Total

	return resp
}
