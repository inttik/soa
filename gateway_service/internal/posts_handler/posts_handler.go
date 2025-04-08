package posthandler

import (
	"encoding/json"
	"gateway/internal/posts_grpc"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostsHandler struct {
	router    *chi.Mux
	conn      *grpc.ClientConn
	rpcClient posts_grpc.PostServiceClient

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
	h.rpcClient = posts_grpc.NewPostServiceClient(h.conn)
	h.usersUrl = "http://users-service:8080/users/v1"
	h.usersClient = &http.Client{
		Timeout: 2 * time.Second,
	}

	router.Route("/posts/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(h.authMiddleware())
		r.Post("/create_post", h.createPost)
		r.Post("/posts/{post_id}", h.updatePost)
		r.Delete("/posts/{post_id}", h.deletePost)
		r.Get("/posts/{post_id}", h.getPost)
		r.Get("/feed", h.getFeed)
	})

	return h
}

type userInfo struct {
	UserId string `json:"user_id"`
	IsRoot bool   `json:"is_root"`
}

func (h *PostsHandler) createPost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}

	var request posts_grpc.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Request error. "+err.Error(), http.StatusBadRequest)
		return
	}

	request.Actor = &posts_grpc.Actor{}
	request.Actor.UserId = userInfo.UserId
	request.Actor.IsRoot = userInfo.IsRoot

	resp, err := h.rpcClient.CreatePost(r.Context(), &request)
	if err != nil {
		http.Error(w, "Server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch resp.Code {
	case posts_grpc.Code_Bad:
		http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
		return
	case posts_grpc.Code_BadRequest:
		http.Error(w, "Bad request. "+resp.Error, http.StatusBadRequest)
		return
	case posts_grpc.Code_Forbidden:
		http.Error(w, "Forbidden. "+resp.Error, http.StatusForbidden)
		return
	case posts_grpc.Code_NotFound:
		http.Error(w, "Not found. "+resp.Error, http.StatusNotFound)
		return
	case posts_grpc.Code_Ok:
		w.Header().Set("Content-type", "application/json")
		resp_post := getPostResponse(resp.Post)
		if err := json.NewEncoder(w).Encode(resp_post); err != nil {
			http.Error(w, "Json encode error. "+err.Error(), http.StatusInternalServerError)
		}
		return

	}

	http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
	return
}

func (h *PostsHandler) updatePost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}

	var request posts_grpc.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request.Update); err != nil {
		http.Error(w, "Request error. "+err.Error(), http.StatusBadRequest)
		return
	}

	if request.Update == nil {
		request.Update = &posts_grpc.PostUpdate{}
	}
	request.Update.Id = chi.URLParam(r, "post_id")

	request.Actor = &posts_grpc.Actor{}
	request.Actor.UserId = userInfo.UserId
	request.Actor.IsRoot = userInfo.IsRoot

	resp, err := h.rpcClient.UpdatePost(r.Context(), &request)
	if err != nil {
		http.Error(w, "Server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch resp.Code {
	case posts_grpc.Code_Bad:
		http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
		return
	case posts_grpc.Code_BadRequest:
		http.Error(w, "Bad request. "+resp.Error, http.StatusBadRequest)
		return
	case posts_grpc.Code_Forbidden:
		http.Error(w, "Forbidden. "+resp.Error, http.StatusForbidden)
		return
	case posts_grpc.Code_NotFound:
		http.Error(w, "Not found. "+resp.Error, http.StatusNotFound)
		return
	case posts_grpc.Code_Ok:
		w.Header().Set("Content-type", "application/json")
		resp_post := getPostResponse(resp.Post)
		if err := json.NewEncoder(w).Encode(resp_post); err != nil {
			http.Error(w, "Json encode error. "+err.Error(), http.StatusInternalServerError)
		}
		return

	}

	http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
	return
}

func (h *PostsHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}

	var request posts_grpc.DeletePostRequest
	request.Id = chi.URLParam(r, "post_id")

	request.Actor = &posts_grpc.Actor{}
	request.Actor.UserId = userInfo.UserId
	request.Actor.IsRoot = userInfo.IsRoot

	resp, err := h.rpcClient.DeletePost(r.Context(), &request)
	if err != nil {
		http.Error(w, "Server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch resp.Code {
	case posts_grpc.Code_Bad:
		http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
		return
	case posts_grpc.Code_BadRequest:
		http.Error(w, "Bad request. "+resp.Error, http.StatusBadRequest)
		return
	case posts_grpc.Code_Forbidden:
		http.Error(w, "Forbidden. "+resp.Error, http.StatusForbidden)
		return
	case posts_grpc.Code_NotFound:
		http.Error(w, "Not found. "+resp.Error, http.StatusNotFound)
		return
	case posts_grpc.Code_Ok:
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte("Post " + request.Id + " deleted."))
		return

	}

	http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
	return
}

func (h *PostsHandler) getPost(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}

	var request posts_grpc.GetPostRequest

	request.Id = chi.URLParam(r, "post_id")

	request.Actor = &posts_grpc.Actor{}
	request.Actor.UserId = userInfo.UserId
	request.Actor.IsRoot = userInfo.IsRoot

	resp, err := h.rpcClient.GetPost(r.Context(), &request)
	if err != nil {
		http.Error(w, "Server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch resp.Code {
	case posts_grpc.Code_Bad:
		http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
		return
	case posts_grpc.Code_BadRequest:
		http.Error(w, "Bad request. "+resp.Error, http.StatusBadRequest)
		return
	case posts_grpc.Code_Forbidden:
		http.Error(w, "Forbidden. "+resp.Error, http.StatusForbidden)
		return
	case posts_grpc.Code_NotFound:
		http.Error(w, "Not found. "+resp.Error, http.StatusNotFound)
		return
	case posts_grpc.Code_Ok:
		w.Header().Set("Content-type", "application/json")
		resp_post := getPostResponse(resp.Post)
		if err := json.NewEncoder(w).Encode(resp_post); err != nil {
			http.Error(w, "Json encode error. "+err.Error(), http.StatusInternalServerError)
		}
		return

	}

	http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
	return
}

func (h *PostsHandler) getFeed(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("userInfo").(*userInfo)
	if !ok {
		log.Fatal("Has no user info")
	}

	var request posts_grpc.ListPostsRequest

	if r.URL.Query().Has("page") {
		u64, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 32)
		if err != nil {
			http.Error(w, "Bad parameter page. "+err.Error(), http.StatusBadRequest)
			return
		}
		request.Page = uint32(u64)
	} else {
		request.Page = 0
	}
	if r.URL.Query().Has("page_limit") {
		u64, err := strconv.ParseUint(r.URL.Query().Get("page_limit"), 10, 32)
		if err != nil {
			http.Error(w, "Bad parameter page_limit. "+err.Error(), http.StatusBadRequest)
			return
		}
		request.PageLimit = uint32(u64)
	} else {
		request.PageLimit = 10
	}
	if r.URL.Query().Has("with_hidden") {
		request.WithHidden = true
	} else {
		request.WithHidden = false
	}

	request.Actor = &posts_grpc.Actor{}
	request.Actor.UserId = userInfo.UserId
	request.Actor.IsRoot = userInfo.IsRoot

	resp, err := h.rpcClient.ListPosts(r.Context(), &request)
	if err != nil {
		http.Error(w, "Server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch resp.Code {
	case posts_grpc.Code_Bad:
		http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
		return
	case posts_grpc.Code_BadRequest:
		http.Error(w, "Bad request. "+resp.Error, http.StatusBadRequest)
		return
	case posts_grpc.Code_Forbidden:
		http.Error(w, "Forbidden. "+resp.Error, http.StatusForbidden)
		return
	case posts_grpc.Code_NotFound:
		http.Error(w, "Not found. "+resp.Error, http.StatusNotFound)
		return
	case posts_grpc.Code_Ok:
		w.Header().Set("Content-type", "application/json")
		resp_post := getFeedResponse(resp)
		if err := json.NewEncoder(w).Encode(resp_post); err != nil {
			http.Error(w, "Json encode error. "+err.Error(), http.StatusInternalServerError)
		}
		return

	}

	http.Error(w, "Bad return code. "+resp.Error, http.StatusInternalServerError)
	return
}
