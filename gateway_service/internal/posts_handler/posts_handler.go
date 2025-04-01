package posthandler

import (
	"gateway/internal/posts_grpc"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostsHandler struct {
	router *chi.Mux
	conn   *grpc.ClientConn
	client posts_grpc.PostServiceClient

	usersUrl    string
	usersClient *http.Client
}

func NewPostsHandler(router *chi.Mux) PostsHandler {
	h := PostsHandler{}
	h.router = router

	conn, err := grpc.NewClient("posts-service:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	h.conn = conn
	h.client = posts_grpc.NewPostServiceClient(h.conn)
	h.usersUrl = "http://users-service:8080/users/v1"
	h.usersClient = &http.Client{
		Timeout: 2 * time.Second,
	}

	router.Route("/posts/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(h.authMiddleware())
		r.Post("/create_post", h.createPost)
		r.Post("/posts/{id}", h.updatePost)
		r.Delete("/posts/{id}", h.deletePost)
		r.Get("/posts/{id}", h.getPost)
		r.Get("/feed", h.getFeed)
	})

	return h
}

type userInfo struct {
	UserId string `json:"user_id"`
	IsRoot bool   `json:"is_root"`
}

// TODO hands

func (h *PostsHandler) createPost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}
	log.Println(userInfo)
}

func (h *PostsHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}
	log.Println(userInfo)
}

func (h *PostsHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}
	log.Println(userInfo)
}

func (h *PostsHandler) getPost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}
	log.Println(userInfo)
}

func (h *PostsHandler) getFeed(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}
	log.Println(userInfo)
}
