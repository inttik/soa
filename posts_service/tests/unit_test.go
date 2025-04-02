package unit_test

import (
	"posts/internal/posts_grpc"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	s := state{}
	s.setup(t)

	tests := []struct {
		name string
		ch   CreateHelper
	}{
		{
			name: "simple",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     "title",
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "no title",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   s.user_a,
				exp400: true,
			},
		},
		{
			name: "no content",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     "title",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   s.user_a,
				exp400: true,
			},
		},
		{
			name: "no actor",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     "title",
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				exp400: true,
			},
		},
		{
			name: "bad id",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     "title",
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   "my id",
				exp400: true,
			},
		},
		{
			name: "long title",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     strings.Repeat("a", 255),
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "too long title",
			ch: CreateHelper{
				req: &posts_grpc.CreatePostRequest{
					Title:     strings.Repeat("a", 256),
					Content:   "content",
					IsPrivate: true,
					Tags:      []string{"a", "b"},
				},
				user:   s.user_a,
				exp400: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.applyCreate(test.ch, t)
		})
	}
}

func TestUpdate(t *testing.T) {
	s := state{}
	s.setup(t)

	// Сначала создаем пост для тестирования обновления
	createResp := s.applyCreate(CreateHelper{
		req: &posts_grpc.CreatePostRequest{
			Title:     "Original Title",
			Content:   "Original Content",
			IsPrivate: false,
			Tags:      []string{"a", "b"},
		},
		user:   s.user_a,
		exp200: true,
	}, t)
	postID := createResp.GetPost().GetId()

	tests := []struct {
		name string
		uh   UpdateHelper
	}{
		{
			name: "successful update all fields",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:            postID,
						NewTitle:      "New Title",
						NewContent:    "New Content",
						ChangePrivate: true,
						IsPrivate:     true,
						AddTags:       []string{"c"},
						RemoveTags:    []string{"a"},
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "title update",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:       postID,
						NewTitle: "New Title 2",
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "content update",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:         postID,
						NewContent: "New Content 2",
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "add tags",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:      postID,
						AddTags: []string{"c", "d", "e"},
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "remove tags",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:         postID,
						RemoveTags: []string{"b", "e"},
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "change private",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:            postID,
						RemoveTags:    []string{"b", "e"},
						ChangePrivate: true,
						IsPrivate:     false,
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "change private 2",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:            postID,
						RemoveTags:    []string{"b", "e"},
						ChangePrivate: true,
						IsPrivate:     false,
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "no change 1",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id: postID,
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "no change 2",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:        postID,
						IsPrivate: true,
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "no change 3",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:        postID,
						IsPrivate: false,
					},
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "root change",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:       postID,
						NewTitle: "New Title 3",
					},
				},
				user:   s.user_b,
				root:   true,
				exp200: true,
			},
		},
		{
			name: "bad user",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:       postID,
						NewTitle: "Drop",
					},
				},
				user:   s.user_b,
				exp403: true,
			},
		},
		{
			name: "bad post id",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:       s.user_a,
						NewTitle: "Drop",
					},
				},
				user:   s.user_a,
				exp404: true,
			},
		},
		{
			name: "no id",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						NewTitle: "Drop",
					},
				},
				user:   s.user_a,
				exp400: true,
			},
		},
		{
			name: "no user",
			uh: UpdateHelper{
				req: &posts_grpc.UpdatePostRequest{
					Update: &posts_grpc.PostUpdate{
						Id:       postID,
						NewTitle: "Drop",
					},
				},
				exp400: true,
			},
		},
		{
			name: "no update",
			uh: UpdateHelper{
				req:    &posts_grpc.UpdatePostRequest{},
				user:   s.user_a,
				exp400: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.applyUpdate(test.uh, t)
		})
	}
}

func TestDelete(t *testing.T) {
	s := state{}
	s.setup(t)

	postA := s.applyCreate(CreateHelper{
		req: &posts_grpc.CreatePostRequest{
			Title:   "Post by User A",
			Content: "Content",
		},
		user:   s.user_a,
		exp200: true,
	}, t).GetPost()

	postB := s.applyCreate(CreateHelper{
		req: &posts_grpc.CreatePostRequest{
			Title:   "Post by User B",
			Content: "Content",
		},
		user:   s.user_b,
		exp200: true,
	}, t).GetPost()

	tests := []struct {
		name string
		dh   DeleteHelper
	}{
		{
			name: "delete non-existent post",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: uuid.NewString(),
				},
				user:   s.user_a,
				exp404: true,
			},
		},
		{
			name: "empty post id",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: "",
				},
				user:   s.user_a,
				exp400: true,
			},
		},
		{
			name: "invalid post id format",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: "invalid-uuid",
				},
				user:   s.user_a,
				exp400: true,
			},
		},
		{
			name: "delete by non-owner",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: postA.GetId(),
				},
				user:   s.user_b,
				exp403: true,
			},
		},
		{
			name: "missing actor",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: postA.GetId(),
				},
				exp400: true,
			},
		},
		{
			name: "successful delete by owner",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: postA.GetId(),
				},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "delete by root user",
			dh: DeleteHelper{
				req: &posts_grpc.DeletePostRequest{
					Id: postB.GetId(),
				},
				user:   s.user_a,
				root:   true,
				exp200: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.applyDelete(test.dh, t)

			if test.dh.exp200 {
				s.applyGet(GetHelper{
					req:    &posts_grpc.GetPostRequest{Id: test.dh.req.Id},
					user:   s.user_a,
					root:   true,
					exp404: true,
				}, t)
			}
		})
	}
}

func TestGet(t *testing.T) {
	s := state{}
	s.setup(t)

	publicPost := s.applyCreate(CreateHelper{
		req: &posts_grpc.CreatePostRequest{
			Title:     "Public Post",
			Content:   "Public Content",
			IsPrivate: false,
			Tags:      []string{"public"},
		},
		user:   s.user_a,
		exp200: true,
	}, t).GetPost()

	privatePost := s.applyCreate(CreateHelper{
		req: &posts_grpc.CreatePostRequest{
			Title:     "Private Post",
			Content:   "Private Content",
			IsPrivate: true,
			Tags:      []string{"private"},
		},
		user:   s.user_a,
		exp200: true,
	}, t).GetPost()

	tests := []struct {
		name string
		gh   GetHelper
	}{
		{
			name: "get public post: any user",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: publicPost.GetId()},
				user:   s.user_b,
				exp200: true,
			},
		},
		{
			name: "get private post: owner",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: privatePost.GetId()},
				user:   s.user_a,
				exp200: true,
			},
		},
		{
			name: "get private post: root user",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: privatePost.GetId()},
				user:   s.user_b,
				root:   true,
				exp200: true,
			},
		},
		{
			name: "get private post: unauthorized",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: privatePost.GetId()},
				user:   s.user_b,
				exp403: true,
			},
		},
		{
			name: "get non-existent post",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: uuid.NewString()},
				user:   s.user_a,
				exp404: true,
			},
		},
		{
			name: "empty post id",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: ""},
				exp400: true,
			},
		},
		{
			name: "invalid post id format",
			gh: GetHelper{
				req:    &posts_grpc.GetPostRequest{Id: "invalid-uuid"},
				exp400: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := s.applyGet(test.gh, t)

			if test.gh.exp200 {
				assert.NotNil(t, resp.GetPost())
				assert.Equal(t, test.gh.req.GetId(), resp.GetPost().GetId())

				if strings.Contains(test.name, "private") && (test.gh.root == false && test.gh.user != s.user_a) {
					assert.Empty(t, resp.GetPost().GetContent())
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	s := state{}
	s.setup(t)

	for i := range 12 {
		s.applyCreate(CreateHelper{
			req: &posts_grpc.CreatePostRequest{
				Title:     "Public " + strconv.Itoa(i),
				Content:   "Public Content",
				IsPrivate: false,
				Tags:      []string{"public"},
			},
			user:   s.user_a,
			exp200: true,
		}, t)
	}

	for i := range 10 {
		s.applyCreate(CreateHelper{
			req: &posts_grpc.CreatePostRequest{
				Title:     "Public " + strconv.Itoa(i),
				Content:   "Public Content",
				IsPrivate: true,
				Tags:      []string{"public"},
			},
			user:   s.user_b,
			exp200: true,
		}, t)
	}

	for i := range 2 {
		s.applyCreate(CreateHelper{
			req: &posts_grpc.CreatePostRequest{
				Title:     "Private " + strconv.Itoa(i),
				Content:   "Private Content",
				IsPrivate: true,
				Tags:      []string{"public"},
			},
			user:   s.user_a,
			exp200: true,
		}, t)
	}

	tests := []struct {
		name string
		lh   ListHelper
	}{
		{
			name: "a public info",
			lh: ListHelper{
				req:          &posts_grpc.ListPostsRequest{},
				user:         s.user_a,
				exp200:       true,
				exp_from:     0,
				exp_to:       0,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public first",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 1,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     0,
				exp_to:       1,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public third",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      2,
					PageLimit: 1,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     2,
				exp_to:       3,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public chunk 1",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 5,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     0,
				exp_to:       5,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public chunk 2",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      1,
					PageLimit: 5,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     5,
				exp_to:       10,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public chunk 3",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      2,
					PageLimit: 5,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     10,
				exp_to:       12,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a public all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 100,
				},
				user:         s.user_a,
				exp200:       true,
				exp_from:     0,
				exp_to:       12,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "b public all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 100,
				},
				user:         s.user_b,
				exp200:       true,
				exp_from:     0,
				exp_to:       12,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a root public all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 100,
				},
				user:         s.user_a,
				root:         true,
				exp200:       true,
				exp_from:     0,
				exp_to:       12,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "b root public all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 100,
				},
				user:         s.user_b,
				root:         true,
				exp200:       true,
				exp_from:     0,
				exp_to:       12,
				exp_total:    12,
				exp_pub_only: true,
			},
		},
		{
			name: "a private all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       0,
					PageLimit:  100,
					WithHidden: true,
				},
				user:      s.user_a,
				exp200:    true,
				exp_from:  0,
				exp_to:    14,
				exp_total: 14,
			},
		},
		{
			name: "b private all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       0,
					PageLimit:  100,
					WithHidden: true,
				},
				user:      s.user_b,
				exp200:    true,
				exp_from:  0,
				exp_to:    22,
				exp_total: 22,
			},
		},
		{
			name: "b private chunk",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       2,
					PageLimit:  3,
					WithHidden: true,
				},
				user:      s.user_b,
				exp200:    true,
				exp_from:  6,
				exp_to:    9,
				exp_total: 22,
			},
		},
		{
			name: "a root private all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       0,
					PageLimit:  100,
					WithHidden: true,
				},
				user:      s.user_a,
				root:      true,
				exp200:    true,
				exp_from:  0,
				exp_to:    24,
				exp_total: 24,
			},
		},
		{
			name: "b root private all",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       0,
					PageLimit:  100,
					WithHidden: true,
				},
				user:      s.user_b,
				root:      true,
				exp200:    true,
				exp_from:  0,
				exp_to:    24,
				exp_total: 24,
			},
		},
		{
			name: "out of bondaries",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:       5,
					PageLimit:  100,
					WithHidden: true,
				},
				user:      s.user_a,
				root:      true,
				exp200:    true,
				exp_from:  0,
				exp_to:    0,
				exp_total: 24,
			},
		},
		{
			name: "no user",
			lh: ListHelper{
				req: &posts_grpc.ListPostsRequest{
					Page:      0,
					PageLimit: 0,
				},
				exp400: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.applyList(test.lh, t)
		})
	}
}
