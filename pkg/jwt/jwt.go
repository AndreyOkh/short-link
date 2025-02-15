package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
)

type JWT struct {
	SecretKey string
}

func New(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) CreateToken(email string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	s, err := t.SignedString([]byte(j.SecretKey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return s, nil
}
