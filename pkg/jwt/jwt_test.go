package jwt_test

import (
	"short-link/pkg/jwt"
	"testing"
)

func TestJWT_CreateToken(t *testing.T) {
	const email = "test@example.com"
	jwtService := jwt.New("12")
	res, err := jwtService.CreateToken(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Error(err)
	}
	isValid, data := jwtService.ParseToken(res)
	if !isValid {
		t.Fatal("invalid token")
	}
	if data.Email != email {
		t.Fatal("invalid email")
	}
	t.Log(data.Email)
}
