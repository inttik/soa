package mockstorage

import (
	"posts/internal/posts_grpc"
	"sync"
)

type posts struct {
	data map[string]*posts_grpc.Post
	mx   sync.Mutex
}
