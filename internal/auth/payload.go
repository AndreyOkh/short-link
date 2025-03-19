package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" schema:"email"`
	Password string `json:"password" validate:"required" schema:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
