package auth_test

import (
	"short-link/internal/auth"
	"short-link/internal/user"
	"short-link/pkg/di"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: u.Email,
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return &user.User{
		Email:    "test@test.com",
		Password: "$2a$10$oEGd4tvKjebd/sZek8YFxOSRfYL2fzdCTXzayFAKq4EGRjK8NXLQe",
	}, nil
}

func TestAuthService_Register(t *testing.T) {
	const initialEmail = "test@example.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(initialEmail, "1", "Test")
	if err != nil {
		if err.Error() != "user already exists" {
			t.Errorf("TestAuthService_Register() error = %v", err)
		}
	}
	if email != initialEmail {
		t.Errorf("TestAuthService_Register() email = %v", email)
	}
}

func TestAuthService_Login(t *testing.T) {
	type fields struct {
		UserRepository di.IUserRepository
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestAuthService_Login",
			fields: fields{
				UserRepository: &MockUserRepository{},
			},
			args: args{
				email:    "test@example.com",
				password: "qwerty",
			},
			want:    "test@example.com",
			wantErr: false,
		},
		{
			name: "TestAuthService_Login2",
			fields: fields{
				UserRepository: &MockUserRepository{},
			},
			args: args{
				email:    "test@example.com",
				password: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := &auth.AuthService{
				UserRepository: tt.fields.UserRepository,
			}
			got, err := authService.Login(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}
