package handlers

import (
	"context"
	"log"

	jwttoken "users/internal/jwt_token"
	api "users/oas"
)

type userMetadataKey string

const (
	UserMetadataKey = userMetadataKey("UserMetadataKey")
)

type securityHandler struct {
	jwt jwttoken.JWTValidator
}

func NewSecurityHandler() (securityHandler, error) {
	jwt, err := jwttoken.NewHandler()
	if err != nil {
		return securityHandler{}, err
	}
	return securityHandler{jwt: &jwt}, nil
}

func (h *securityHandler) HandleBearerHttpAuthentication(ctx context.Context, operationName api.OperationName, t api.BearerHttpAuthentication) (context.Context, error) {
	metadata, err := h.jwt.ReadJWT(t.Token)
	if err != nil {
		log.Println(err)
		return ctx, nil
	}
	newCtx := context.WithValue(ctx, UserMetadataKey, metadata)
	return newCtx, nil
}
