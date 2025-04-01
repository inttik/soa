package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/uuid"
)

func main() {
	_ = uuid.UUID{}
	postsTarget, _ := url.Parse("http://posts-service:8080")
	postsProxy := httputil.NewSingleHostReverseProxy(postsTarget)

	usersTarget, _ := url.Parse("http://users-service:8080")
	usersProxy := httputil.NewSingleHostReverseProxy(usersTarget)

	http.HandleFunc("/posts/v1/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Forwarding to posts service")
		postsProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/users/v1/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Forwarding to users service")
		usersProxy.ServeHTTP(w, r)
	})

	log.Println("Gateway started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
