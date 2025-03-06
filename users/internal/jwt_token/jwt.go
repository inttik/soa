package jwttoken

import (
	"crypto/rsa"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserMetadata struct {
	Root   bool
	UserId uuid.UUID
}

type JWTValidator interface {
	GenerateJWT(metadata UserMetadata) (string, error)
	ReadJWT(token string) (UserMetadata, error)
}

type jwtHandler struct {
	jwtPublic  *rsa.PublicKey
	jwtPrivate *rsa.PrivateKey
}

const (
	EnvJwtPublic  = "JWT_PUBLIC"
	EnvJwtPrivate = "JWT_PRIVATE"
)

func NewHandler() (jwtHandler, error) {
	privateKey, exists := os.LookupEnv(EnvJwtPrivate)
	if !exists {
		return jwtHandler{}, errors.New("no env variable for jwt privater")
	}
	publicKey, exists := os.LookupEnv(EnvJwtPublic)
	if !exists {
		return jwtHandler{}, errors.New("no env variable for jwt public")
	}
	jwtPrivate, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return jwtHandler{}, err
	}
	jwtPublic, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return jwtHandler{}, err
	}
	return jwtHandler{jwtPublic: jwtPublic, jwtPrivate: jwtPrivate}, nil
}

func SetupEnv(path string) error {
	privateKey, err := os.ReadFile(path + "/signature.pem")
	if err != nil {
		_, exists := os.LookupEnv(EnvJwtPrivate)
		if !exists {
			return err
		}
	}
	err = os.Setenv(EnvJwtPrivate, string(privateKey))
	if err != nil {
		return err
	}

	publicKey, err := os.ReadFile(path + "/signature.pub")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvJwtPublic, string(publicKey))
	if err != nil {
		return err
	}
	return nil
}

func (h *jwtHandler) GenerateJWT(metadata UserMetadata) (string, error) {
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 48).Unix(),
		"iat":     time.Now(),
		"iss":     "user server",
		"user_id": metadata.UserId.String(),
		"root":    metadata.Root,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(h.jwtPrivate)
	if err != nil {
		log.Fatal("Fail to sign token: ", err)
		return "", err
	}
	return token, nil
}

func (h *jwtHandler) ReadJWT(token string) (UserMetadata, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("bad signing method")
		}
		return h.jwtPublic, nil
	})
	if err != nil {
		return UserMetadata{}, err
	}
	if !parsedToken.Valid {
		return UserMetadata{}, errors.New("token is not valid")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return UserMetadata{}, errors.New("token has no jwt.MapClaims")
	}
	expDate, err := claims.GetExpirationTime()
	if err != nil {
		return UserMetadata{}, err
	}
	if time.Now().After(expDate.Time) {
		return UserMetadata{}, errors.New("outdated jwt token")
	}

	userIdString, ok := claims["user_id"].(string)
	if !ok {
		return UserMetadata{}, errors.New("token has no user_id field")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return UserMetadata{}, err
	}

	root, ok := claims["root"].(bool)
	if !ok {
		return UserMetadata{}, errors.New("token has no root field")
	}
	return UserMetadata{UserId: userId, Root: root}, nil
}
