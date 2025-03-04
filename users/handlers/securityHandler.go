package handlers

import (
	"context"
	"errors"
	api "users/oas"
)

type securityHandler struct {
}

func NewSecurityHandler() securityHandler {
	return securityHandler{}
}

func (*securityHandler) HandleBearerHttpAuthentication(ctx context.Context, operationName api.OperationName, t api.BearerHttpAuthentication) (context.Context, error) {
	return nil, errors.New("not implemented")
}
