// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// FriendsUserIDGet implements GET /friends/{user_id} operation.
//
// Returns all public friends_id of user.
// If client has private access to user_id (client is
// root or is user_id), then friend aliasas are returned
// as well as friend pairs metadata.
//
// GET /friends/{user_id}
func (UnimplementedHandler) FriendsUserIDGet(ctx context.Context, params FriendsUserIDGetParams) (r FriendsUserIDGetRes, _ error) {
	return r, ht.ErrNotImplemented
}

// FriendsUserIDPost implements POST /friends/{user_id} operation.
//
// Add, update, or remove friend of user.
//
// POST /friends/{user_id}
func (UnimplementedHandler) FriendsUserIDPost(ctx context.Context, req *FriendModify, params FriendsUserIDPostParams) (r FriendsUserIDPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// LoginPost implements POST /login operation.
//
// Checks that you is the user and, if so, returns JWT token.
//
// POST /login
func (UnimplementedHandler) LoginPost(ctx context.Context, req *LoginUserRequest) (r LoginPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ProfileUserIDGet implements GET /profile/{user_id} operation.
//
// Get all user profile information.
//
// GET /profile/{user_id}
func (UnimplementedHandler) ProfileUserIDGet(ctx context.Context, params ProfileUserIDGetParams) (r ProfileUserIDGetRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ProfileUserIDPost implements POST /profile/{user_id} operation.
//
// Check if user profile can be update by current user and, if so, updates it.
//
// POST /profile/{user_id}
func (UnimplementedHandler) ProfileUserIDPost(ctx context.Context, req *ProfileUpdate, params ProfileUserIDPostParams) (r ProfileUserIDPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// RegisterPost implements POST /register operation.
//
// Checks if user can be created and, if so, creates it and returns user_id. Only root can create
// root users.
//
// POST /register
func (UnimplementedHandler) RegisterPost(ctx context.Context, req *CreateUserRequest) (r RegisterPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UserLoginGet implements GET /user/{login} operation.
//
// Returns user id of user with that login, or 404 if there is no such user.
//
// GET /user/{login}
func (UnimplementedHandler) UserLoginGet(ctx context.Context, params UserLoginGetParams) (r UserLoginGetRes, _ error) {
	return r, ht.ErrNotImplemented
}
