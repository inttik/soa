package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"users/oas"

	"github.com/google/uuid"
)

type userService struct {
	id int64
}

func NewService() userService {
	return userService{id: 0}
}

func (*userService) FriendsUserIDGet(ctx context.Context, params oas.FriendsUserIDGetParams) (oas.FriendsUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) FriendsUserIDPost(ctx context.Context, req *oas.FriendModify, params oas.FriendsUserIDPostParams) (oas.FriendsUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) LoginPost(ctx context.Context, req *oas.LoginUserRequest) (oas.LoginPostRes, error) {
	token := oas.JwtToken{}
	user_id := oas.UserId{}

	user_metadata := ctx.Value(USER_METADATA_KEY)
	if user_metadata == nil {
		fmt.Println("Not logined")
	} else {
		user_id = oas.UserId(user_metadata.(userMetadata).user_id)
		fmt.Println("Logined as", user_id)
	}

	body := oas.LoginPostOK{
		Token:  token,
		UserID: user_id,
	}

	headers := oas.LoginPostOKHeaders{}
	headers.SetResponse(body)

	return &headers, nil
}

func (*userService) ProfileUserIDGet(ctx context.Context, params oas.ProfileUserIDGetParams) (oas.ProfileUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) ProfileUserIDPost(ctx context.Context, req *oas.ProfileUpdate, params oas.ProfileUserIDPostParams) (oas.ProfileUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) RegisterPost(ctx context.Context, req *oas.CreateUserRequest) (oas.RegisterPostRes, error) {
	log.Println("Register handler")
	user_metadata := ctx.Value(USER_METADATA_KEY)
	if user_metadata == nil {
		fmt.Println("Not logined")
	} else {
		user_metadata := user_metadata.(userMetadata)
		user_id := oas.UserId(user_metadata.user_id)
		root_str := ""
		if !user_metadata.root {
			root_str = "(not root)"
		} else {
			root_str = "(root)"
		}
		fmt.Println("Logined as", uuid.UUID(user_id).String(), root_str)
	}
	return nil, errors.New("not implemented")
}
