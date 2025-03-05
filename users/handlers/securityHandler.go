package handlers

import (
	"context"
	"log"
	jwttokens "users/internal/jwt_token"
	api "users/oas"

	"github.com/google/uuid"
)

type userMetadataKey string

const (
	USER_METADATA_KEY = userMetadataKey("UserMetadataKey")
)

type userMetadata struct {
	root    bool
	user_id uuid.UUID
}

type securityHandler struct {
	jwt jwttokens.JWTValidator
}

func NewSecurityHandler() (securityHandler, error) {
	jwt, err := jwttokens.NewHandler()
	if err != nil {
		return securityHandler{}, err
	}
	return securityHandler{jwt: jwt}, nil
}

func (*securityHandler) HandleBearerHttpAuthentication(ctx context.Context, operationName api.OperationName, t api.BearerHttpAuthentication) (context.Context, error) {
	log.Println("auth handler!")
	log.Println(t)
	user_metadata := userMetadata{
		root:    true,
		user_id: uuid.New(),
	}
	new_ctx := context.WithValue(ctx, USER_METADATA_KEY, user_metadata)
	return new_ctx, nil
}
