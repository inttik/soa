package jwttoken_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	jwttoken "users/internal/jwt_token"
)

const (
	jwtDir = "../../secrets/"
)

func TestToken(t *testing.T) {
	jwttoken.SetupEnv(jwtDir)
	handler, err := jwttoken.NewHandler()
	assert.NoErrorf(t, err, "Handler should be created")

	tests := []struct {
		name     string
		metadata jwttoken.UserMetadata
	}{
		{
			name: "not root user",
			metadata: jwttoken.UserMetadata{
				Root:   false,
				UserId: uuid.New(),
			},
		},
		{
			name: "root user",
			metadata: jwttoken.UserMetadata{
				Root:   true,
				UserId: uuid.UUID{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := test.metadata

			token, err := handler.GenerateJWT(user)
			assert.NoErrorf(t, err, "token should be created")

			actual, err := handler.ReadJWT(token)
			assert.NoErrorf(t, err, "token should be decrypted")

			assert.Equal(t, user, actual)
		})
	}
}
