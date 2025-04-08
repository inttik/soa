package handlers

import (
	posthandler "gateway/internal/posts_handler"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type gatewayServer struct {
	router *chi.Mux
	users  usersHandler
	posts  posthandler.PostsHandler
	server *http.Server
}

func NewGatewayServer() gatewayServer {
	router := chi.NewRouter()
	users := newUserHandler(router)
	posts := posthandler.NewPostsHandler(router)

	return gatewayServer{
		router: router,
		users:  users,
		posts:  posts,
	}
}

func (s *gatewayServer) Start(port string) {
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	log.Println("starting server at port ", port)
	s.server.ListenAndServe()
}
