package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
)

type JWTData struct {
	Email string
}

type JWT struct {
	SecretKey string
}

func New(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) CreateToken(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})
	s, err := t.SignedString([]byte(j.SecretKey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return s, nil
}

func (j *JWT) ParseToken(tokenString string) (bool, *JWTData) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return false, nil
	}
	email, ok := t.Claims.(jwt.MapClaims)["email"]
	if !ok {
		return false, nil
	} else {
		return t.Valid, &JWTData{
			Email: email.(string),
		}
	}

}
