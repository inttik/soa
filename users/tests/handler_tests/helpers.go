package handlers_test

import (
	"context"
	"errors"
	"testing"
	"users/handlers"
	jwttoken "users/internal/jwt_token"
	mockstorage "users/internal/mock_storage"
	"users/internal/passhandler"
	"users/oas"

	"github.com/stretchr/testify/assert"
)

const (
	jwtDir = "../../secrets/"
)

type state struct {
	service    oas.Handler
	security   oas.SecurityHandler
	adminToken oas.BearerHttpAuthentication
}

func (s *state) setup(t *testing.T) {
	err := jwttoken.SetupEnv(jwtDir)
	assert.NoErrorf(t, err, "Env should be setted up")

	storage, err := mockstorage.NewMockStorage()
	assert.NoErrorf(t, err, "Mock storage should be created")

	adminLogin := oas.LoginString("admin")
	adminPass := oas.PasswordString("admin")
	hashedPass := passhandler.HashPass(adminLogin, adminPass)
	err = storage.MakeRootUser(adminLogin, oas.PasswordString(hashedPass))
	assert.NoErrorf(t, err, "Mock storage should have admin user")

	serv, err := handlers.NewService(&storage)
	assert.NoErrorf(t, err, "Service should be created")
	s.service = &serv

	secr, err := handlers.NewSecurityHandler()
	assert.NoErrorf(t, err, "Security handler should be created")
	s.security = &secr

	lh := loginHelper{
		req: oas.LoginUserRequest{
			Login:    "admin",
			Password: "admin",
		},
		expect200: true,
	}

	resp, err := s.applyLogin(lh, t)
	assert.NoErrorf(t, err, "Admin should be able to login")
	s.adminToken.Token = string(resp.Token)
}

type registerHelper struct {
	req       oas.CreateUserRequest
	admin     bool
	expect201 bool
	expect400 bool
	expect403 bool
}

func (s *state) applyRegister(req registerHelper, t *testing.T) (*oas.UserId, error) {
	err := req.req.Validate()
	if err != nil {
		assert.Equal(t, false, req.expect201)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect403)
		return nil, err
	}

	ctx := context.Background()

	if req.admin {
		ctx, err = s.security.HandleBearerHttpAuthentication(ctx, oas.RegisterPostOperation, s.adminToken)
		assert.NoErrorf(t, err, "Security service raise no error")
	}

	ret, err := s.service.RegisterPost(ctx, &req.req)
	assert.NoErrorf(t, err, "service raise no errors")
	switch u := ret.(type) {
	case *oas.UserId:
		assert.Equal(t, true, req.expect201)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, false, req.expect403)
		return u, nil
	case *oas.RegisterPostBadRequest:
		assert.Equal(t, false, req.expect201)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect403)
		return nil, errors.New("400")
	case *oas.RegisterPostForbidden:
		assert.Equal(t, false, req.expect201)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, true, req.expect403)
		return nil, errors.New("403")
	}
	t.Fatal("Unreachable")
	return nil, errors.New("500")
}

type loginHelper struct {
	req       oas.LoginUserRequest
	expect200 bool
	expect400 bool
	expect404 bool
}

func (s *state) applyLogin(req loginHelper, t *testing.T) (*oas.LoginPostOK, error) {
	err := req.req.Validate()
	if err != nil {
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect404)
		return nil, err
	}

	ctx := context.Background()

	ret, err := s.service.LoginPost(ctx, &req.req)
	assert.NoErrorf(t, err, "service raise no errors")
	switch u := ret.(type) {
	case *oas.LoginPostOK:
		assert.Equal(t, true, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, false, req.expect404)
		return u, nil
	case *oas.LoginPostBadRequest:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect404)
		return nil, errors.New("400")
	case *oas.LoginPostNotFound:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, true, req.expect404)
		return nil, errors.New("404")
	}
	t.Fatal("Unreachable")
	return nil, errors.New("500")
}

type userGetHelp struct {
	req       oas.LoginString
	expect200 bool
	expect404 bool
}

func (s *state) applyUserGet(req userGetHelp, t *testing.T) (*oas.UserId, error) {
	request := oas.UserLoginGetParams{
		Login: req.req,
	}
	ctx := context.Background()

	ret, err := s.service.UserLoginGet(ctx, request)
	assert.NoErrorf(t, err, "Server raises no error")
	switch u := ret.(type) {
	case *oas.UserId:
		assert.Equal(t, true, req.expect200)
		assert.Equal(t, false, req.expect404)
		return u, nil
	case *oas.UserLoginGetNotFound:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect404)
		return nil, errors.New("404")
	}
	t.Fatal("Unreachable")
	return nil, errors.New("500")
}

type profileGetHelper struct {
	req       oas.UserId
	logined   *oas.LoginUserRequest
	admin     bool
	expect200 bool
	expect404 bool
}

func (s *state) applyProfileGet(req profileGetHelper, t *testing.T) (*oas.ProfileInfo, error) {
	request := oas.ProfileUserIDGetParams{
		UserID: req.req,
	}

	ctx := context.Background()

	if req.logined != nil {
		resp, err := s.applyLogin(loginHelper{
			req: oas.LoginUserRequest{
				Login:    req.logined.Login,
				Password: req.logined.Password,
			},
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "Should be loggined")

		err = resp.Token.Validate()
		assert.NoErrorf(t, err, "Token should be correct")
		token := oas.BearerHttpAuthentication{
			Token: string(resp.Token),
		}
		ctx, err = s.security.HandleBearerHttpAuthentication(ctx, oas.ProfileUserIDGetOperation, token)
		assert.NoErrorf(t, err, "Security service raise no error")
	}
	if req.admin {
		var err error
		ctx, err = s.security.HandleBearerHttpAuthentication(ctx, oas.ProfileUserIDGetOperation, s.adminToken)
		assert.NoErrorf(t, err, "Security service raise no error")
	}

	ret, err := s.service.ProfileUserIDGet(ctx, request)
	assert.NoErrorf(t, err, "Service raise no errors")
	switch u := ret.(type) {
	case *oas.ProfileInfo:
		assert.Equal(t, true, req.expect200)
		assert.Equal(t, false, req.expect404)
		return u, nil
	case *oas.ProfileUserIDGetNotFound:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect404)
		return nil, errors.New("404")
	}
	t.Fatal("Unreachable")
	return nil, errors.New("500")
}

type profilePostHelper struct {
	target    *oas.UserId
	req       oas.ProfileUpdate
	register  *registerHelper
	logined   *oas.LoginUserRequest
	admin     bool
	expect200 bool
	expect400 bool
	expect401 bool
	expect403 bool
	expect404 bool
}

func (s *state) applyProfilePost(req profilePostHelper, t *testing.T) (*oas.ProfileInfo, error) {
	params := oas.ProfileUserIDPostParams{}

	if req.register != nil {
		uid, err := s.applyRegister(*req.register, t)
		assert.NoErrorf(t, err, "successfully registrated")
		params.UserID = *uid
	} else {
		params.UserID = *req.target
	}

	err := req.req.Validate()
	if err != nil {
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect401)
		assert.Equal(t, false, req.expect403)
		assert.Equal(t, false, req.expect404)
		return nil, errors.New("400")
	}

	ctx := context.Background()

	if req.logined != nil {
		resp, err := s.applyLogin(loginHelper{
			req: oas.LoginUserRequest{
				Login:    req.logined.Login,
				Password: req.logined.Password,
			},
			expect200: true,
		}, t)
		assert.NoErrorf(t, err, "Should be loggined")

		err = resp.Token.Validate()
		assert.NoErrorf(t, err, "Token should be correct")
		token := oas.BearerHttpAuthentication{
			Token: string(resp.Token),
		}
		ctx, err = s.security.HandleBearerHttpAuthentication(ctx, oas.ProfileUserIDGetOperation, token)
		assert.NoErrorf(t, err, "Security service raise no error")
	}
	if req.admin {
		var err error
		ctx, err = s.security.HandleBearerHttpAuthentication(ctx, oas.ProfileUserIDGetOperation, s.adminToken)
		assert.NoErrorf(t, err, "Security service raise no error")
	}

	ret, err := s.service.ProfileUserIDPost(ctx, &req.req, params)
	assert.NoErrorf(t, err, "service raise no errors")
	switch u := ret.(type) {
	case *oas.ProfileInfo:
		assert.Equal(t, true, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, false, req.expect401)
		assert.Equal(t, false, req.expect403)
		assert.Equal(t, false, req.expect404)
		return u, nil
	case *oas.ProfileUserIDPostBadRequest:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, true, req.expect400)
		assert.Equal(t, false, req.expect401)
		assert.Equal(t, false, req.expect403)
		assert.Equal(t, false, req.expect404)
		return nil, errors.New("400")
	case *oas.ProfileUserIDPostUnauthorized:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, true, req.expect401)
		assert.Equal(t, false, req.expect403)
		assert.Equal(t, false, req.expect404)
		return nil, errors.New("401")
	case *oas.ProfileUserIDPostForbidden:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, false, req.expect401)
		assert.Equal(t, true, req.expect403)
		assert.Equal(t, false, req.expect404)
		return nil, errors.New("403")
	case *oas.ProfileUserIDPostNotFound:
		assert.Equal(t, false, req.expect200)
		assert.Equal(t, false, req.expect400)
		assert.Equal(t, false, req.expect401)
		assert.Equal(t, false, req.expect403)
		assert.Equal(t, true, req.expect404)
		return nil, errors.New("404")
	}
	t.Fatal("Unreachable")
	return nil, errors.New("500")
}
