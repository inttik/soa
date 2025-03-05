package jwttokens

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
	root    bool
	user_id uuid.UUID
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
	jwtPublicFile  = "secrets/signature.pub"
	jwtPrivateFile = "secrets/signature.pem"
)

func NewHandler() (jwtHandler, error) {
	private_key, err := os.ReadFile(jwtPrivateFile)
	if err != nil {
		return jwtHandler{}, err
	}
	public_key, err := os.ReadFile(jwtPublicFile)
	if err != nil {
		return jwtHandler{}, err
	}
	jwt_private, err := jwt.ParseRSAPrivateKeyFromPEM(private_key)
	if err != nil {
		return jwtHandler{}, err
	}
	jwt_public, err := jwt.ParseRSAPublicKeyFromPEM(public_key)
	if err != nil {
		return jwtHandler{}, err
	}
	return jwtHandler{jwtPublic: jwt_public, jwtPrivate: jwt_private}, nil
}

func (h jwtHandler) GenerateJWT(metadata UserMetadata) (string, error) {
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 48).Unix(),
		"iat":     time.Now(),
		"iss":     "user server",
		"user_id": metadata.user_id,
		"root":    metadata.root,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(h.jwtPublic)
	if err != nil {
		log.Fatal("Fail to sign token")
		return "", err
	}
	return token, nil

}

func (h jwtHandler) ReadJWT(token string) (UserMetadata, error) {
	parsed_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("bad signing method")
		}
		return h.jwtPublic, nil
	})
	if err != nil {
		return UserMetadata{}, err
	}
	if !parsed_token.Valid {
		return UserMetadata{}, errors.New("token is not valid")
	}
	claims, ok := parsed_token.Claims.(jwt.MapClaims)
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

	user_id, ok := claims["user_id"].(uuid.UUID)
	if !ok {
		return UserMetadata{}, errors.New("token has no user_id field")
	}

	root, ok := claims["root"].(bool)
	if !ok {
		return UserMetadata{}, errors.New("token has no root field")
	}
	return UserMetadata{user_id: user_id, root: root}, nil
}
