package handlers

import (
	"context"
	"errors"
	"log"
	"strings"

	jwttoken "users/internal/jwt_token"
	pass "users/internal/pass_handler"
	storage "users/internal/storage_manager"
	"users/oas"

	"github.com/google/uuid"
)

type userService struct {
	jwt jwttoken.JWTValidator
	sm  storage.StorageManager
}

func NewService(sm storage.StorageManager) (userService, error) {
	jwt, err := jwttoken.NewHandler()
	if err != nil {
		return userService{}, err
	}
	return userService{jwt: &jwt, sm: sm}, nil
}

func (s *userService) RegisterPost(ctx context.Context, req *oas.CreateUserRequest) (oas.RegisterPostRes, error) {
	userMetadata := ctx.Value(UserMetadataKey)

	newRoot := req.Root.Or(false)
	if newRoot {
		allowed := false
		switch u := userMetadata.(type) {
		default:
			allowed = false
		case jwttoken.UserMetadata:
			allowed = u.Root
		}

		if !allowed {
			return &oas.RegisterPostForbidden{Data: strings.NewReader("to create root user you should be root")}, nil
		}
	}

	hashed_pass, err := pass.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = oas.PasswordString(hashed_pass)

	newUUID, err := s.sm.CreateUser(req)

	if err != nil {
		return &oas.RegisterPostBadRequest{Data: strings.NewReader(err.Error())}, nil
	}

	resp := oas.UserId(newUUID)
	return &resp, nil
}

func (s *userService) LoginPost(ctx context.Context, req *oas.LoginUserRequest) (oas.LoginPostRes, error) {
	currentUUID, err := s.sm.GetUserId(req.Login)
	if err != nil {
		return &oas.LoginPostNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	correctPass, err := s.sm.GetPassword(oas.UserId(currentUUID))
	if err != nil {
		return &oas.LoginPostNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	if !pass.ComparePassword(string(req.Password), string(correctPass)) {
		return &oas.LoginPostBadRequest{Data: strings.NewReader("bad password")}, nil
	}

	profile, err := s.sm.GetProfile(oas.UserId(currentUUID))
	if err != nil {
		return &oas.LoginPostNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	isRoot, ok := profile.Root.Get()
	if !ok {
		log.Fatal("INVARIANT BREAK: storage mannager doesn't returned root flag")
		return nil, errors.New("INVARIANT BREAK: storage mannager doesn't returned root flag")
	}

	metadata := jwttoken.UserMetadata{
		Root:   bool(isRoot),
		UserId: uuid.UUID(currentUUID),
	}
	token, err := s.jwt.GenerateJWT(metadata)
	if err != nil {
		return nil, err
	}

	resp := oas.LoginPostOK{
		Token:  oas.JwtToken(token),
		UserID: oas.UserId(currentUUID),
	}

	return &resp, nil
}

func (s *userService) WhoamiGet(ctx context.Context) (oas.WhoamiGetRes, error) {
	userMetadata := ctx.Value(UserMetadataKey)
	switch u := userMetadata.(type) {
	default:
		return &oas.WhoamiGetForbidden{Data: strings.NewReader("unauthorized")}, nil
	case jwttoken.UserMetadata:
		return &oas.WhoamiGetOK{
			UserID: oas.UserId(u.UserId),
			IsRoot: oas.RootFlag(u.Root),
		}, nil
	}
}

func (s *userService) UserLoginGet(ctx context.Context, params oas.UserLoginGetParams) (oas.UserLoginGetRes, error) {
	id, err := s.sm.GetUserId(params.Login)
	if err != nil {
		return &oas.UserLoginGetNotFound{Data: strings.NewReader(err.Error())}, nil
	}
	user_id := oas.UserId(id)
	return &user_id, nil
}

func (s *userService) ProfileUserIDGet(ctx context.Context, params oas.ProfileUserIDGetParams) (oas.ProfileUserIDGetRes, error) {
	userMetadata := ctx.Value(UserMetadataKey)

	adv_access := false
	switch u := userMetadata.(type) {
	default:
		adv_access = false
	case jwttoken.UserMetadata:
		adv_access = u.Root || (u.UserId == uuid.UUID(params.UserID))
	}

	profile, err := s.sm.GetProfile(params.UserID)
	if err != nil {
		return &oas.ProfileUserIDGetNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	if !adv_access {
		profile.Root.Reset()
		profile.RegDate.Reset()
		profile.LastModify.Reset()
	}

	return &profile, nil
}

func (s *userService) ProfileUserIDPost(ctx context.Context, req *oas.ProfileUpdate, params oas.ProfileUserIDPostParams) (oas.ProfileUserIDPostRes, error) {
	userMetadata := ctx.Value(UserMetadataKey)

	auth := false
	adv_access := false
	switch u := userMetadata.(type) {
	default:
		auth = false
		adv_access = false
	case jwttoken.UserMetadata:
		auth = true
		adv_access = u.Root || (u.UserId == uuid.UUID(params.UserID))
	}
	if !auth {
		return &oas.ProfileUserIDPostUnauthorized{Data: strings.NewReader("user not authorized")}, nil
	}
	if !adv_access {
		return &oas.ProfileUserIDPostForbidden{Data: strings.NewReader("user has no access")}, nil
	}

	err := s.sm.UpdateProfile(params.UserID, req)
	if err != nil {
		return &oas.ProfileUserIDPostNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	profile, err := s.sm.GetProfile(params.UserID)
	if err != nil {
		return &oas.ProfileUserIDPostNotFound{Data: strings.NewReader(err.Error())}, nil
	}

	return &profile, nil
}

func (s *userService) FriendsUserIDGet(ctx context.Context, params oas.FriendsUserIDGetParams) (oas.FriendsUserIDGetRes, error) {
	return nil, errors.New("not implemented")
}

func (s *userService) FriendsUserIDPost(ctx context.Context, req *oas.FriendModify, params oas.FriendsUserIDPostParams) (oas.FriendsUserIDPostRes, error) {
	return nil, errors.New("not implemented")
}
