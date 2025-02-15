package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"short-link/internal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (authService *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := authService.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	_, err = authService.UserRepository.Create(&user.User{
		Email:    email,
		Password: string(bcryptPassword),
		Name:     name,
	})
	if err != nil {
		return "", err
	} else {
		return email, nil
	}
}

func (authService *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := authService.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return email, nil
}
