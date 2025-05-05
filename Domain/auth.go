package domain

import (
	"context"

)

const (
	CollectionUser = "users"
)

// Signup

type SignUpRequest struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type SignUpRepository interface {
	SignUp(c context.Context, signUpRequest SignUpRequest) (SignUpResponse, error)
	GetUser(c context.Context) ([]User, error)
}

type SignUpUsecase interface {
	SignUp(c context.Context, signUpRequest SignUpRequest) (SignUpResponse, error)
	GetUser(c context.Context) ([]User, error)
}

// Login

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginRepository interface {
	Login(c context.Context, loginRequest LoginRequest) (LoginResponse, error)
}

type LoginUsecase interface {
	Login(c context.Context, loginRequest LoginRequest) (LoginResponse, error)
}

// Logout

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type LogoutRepository interface {
	Logout(c context.Context, logoutRequest LogoutRequest) (LogoutResponse, error)
}

type LogoutUsecase interface {
	Logout(c context.Context, logoutRequest LogoutRequest) (LogoutResponse, error)
}
