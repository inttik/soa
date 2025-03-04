package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	api "users/oas"

	"github.com/google/uuid"
)

type userService struct {
	id int64
}

func NewService() userService {
	return userService{id: 0}
}

func (*userService) FriendsUserIDGet(ctx context.Context, params api.FriendsUserIDGetParams) (api.FriendsUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) FriendsUserIDPost(ctx context.Context, req *api.FriendModify, params api.FriendsUserIDPostParams) (api.FriendsUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) LoginPost(ctx context.Context, req *api.LoginUserRequest) (api.LoginPostRes, error) {
	token := api.JwtToken{}
	user_id := api.UserId{}

	uid := uuid.New()
	value, _ := json.Marshal(uid)
	user_id.UnmarshalJSON(value)

	body := api.LoginPostOK{
		Token:  token,
		UserID: user_id,
	}

	headers := api.LoginPostOKHeaders{}
	headers.SetResponse(body)

	fmt.Println("kek")

	return &headers, nil
}

func (*userService) ProfileUserIDGet(ctx context.Context, params api.ProfileUserIDGetParams) (api.ProfileUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) ProfileUserIDPost(ctx context.Context, req *api.ProfileUpdate, params api.ProfileUserIDPostParams) (api.ProfileUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) RegisterPost(ctx context.Context, req *api.CreateUserRequest) (api.RegisterPostRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) NewError(ctx context.Context, err error) *api.ErrorMessageStatusCode {
	return nil
}
