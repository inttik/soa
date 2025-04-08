package postgresstorage

import (
	"posts/internal/posts_grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPost(p *postTable) *posts_grpc.Post {
	return &posts_grpc.Post{
		Id:          p.PostId.String(),
		Title:       p.Title,
		Content:     p.Content,
		AuthorId:    p.AuthorId.String(),
		IsPrivate:   p.IsPrivate,
		PublishDate: timestamppb.New(p.PublishDate),
		LastModify:  timestamppb.New(p.LastModify),
		Tags:        p.Tags,
	}
}
