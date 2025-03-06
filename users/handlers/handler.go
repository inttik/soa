package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"

	jwttoken "users/internal/jwt_token"
	"users/oas"

	"github.com/google/uuid"
)

type userService struct {
	jwt jwttoken.JWTValidator
}

func NewService() (userService, error) {
	jwt, err := jwttoken.NewHandler()
	if err != nil {
		return userService{}, err
	}
	return userService{jwt: jwt}, nil
}

func (*userService) FriendsUserIDGet(ctx context.Context, params oas.FriendsUserIDGetParams) (oas.FriendsUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (*userService) FriendsUserIDPost(ctx context.Context, req *oas.FriendModify, params oas.FriendsUserIDPostParams) (oas.FriendsUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}

func (s *userService) LoginPost(ctx context.Context, req *oas.LoginUserRequest) (oas.LoginPostRes, error) {
	userId := oas.UserId(uuid.New())

	metadata := jwttoken.UserMetadata{
		Root:   false,
		UserId: uuid.UUID(userId),
	}

	token, err := s.jwt.GenerateJWT(metadata)
	if err != nil {
		return nil, err
	}

	log.Println(token)

	body := oas.LoginPostOK{
		Token:  oas.JwtToken(token),
		UserID: userId,
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
	userMetadata := ctx.Value(UserMetadataKey)
	if userMetadata == nil {
		fmt.Println("Not logined")
	} else {
		userMetadata := userMetadata.(jwttoken.UserMetadata)
		userId := oas.UserId(userMetadata.UserId)
		rootStr := ""
		if !userMetadata.Root {
			rootStr = "(not root)"
		} else {
			rootStr = "(root)"
		}
		fmt.Println("Logined as", uuid.UUID(userId).String(), rootStr)
	}
	return nil, errors.New("not implemented")
}
