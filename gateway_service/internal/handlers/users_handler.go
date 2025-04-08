package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type usersHandler struct {
	router     *chi.Mux
	usersProxy *httputil.ReverseProxy
}

func newUserHandler(router *chi.Mux) usersHandler {
	h := usersHandler{
		router: router,
	}
	usersTarget, _ := url.Parse("http://users-service:8080")
	h.usersProxy = httputil.NewSingleHostReverseProxy(usersTarget)
	router.Route("/users/v1/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/register", h.proxyHandler)
		r.Post("/login", h.proxyHandler)
		r.Get("/whoami", h.proxyHandler)
		r.Get("/user/{login}", h.proxyHandler)
		r.Get("/profile/{user_id}", h.proxyHandler)
		r.Post("/profile/{user_id}", h.proxyHandler)
	})

	return h
}

func (h *usersHandler) proxyHandler(w http.ResponseWriter, r *http.Request) {
	h.usersProxy.ServeHTTP(w, r)
}
