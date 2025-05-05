package usecases

import (
	"context"
	"time"

	domain "github.com/chachidani/interview-coach-backend/Domain"
)

// signup usecase

type signUpUsecase struct {
	signUpRepository domain.SignUpRepository
	ContextTimeout   time.Duration
}

// GetUser implements domain.SignUpUsecase.
func (s *signUpUsecase) GetUser(c context.Context) ([]domain.User, error) {
	return s.signUpRepository.GetUser(c)
}

// SignUp implements domain.SignUpUsecase.
func (s *signUpUsecase) SignUp(c context.Context, signUpRequest domain.SignUpRequest) (domain.SignUpResponse, error) {
	return s.signUpRepository.SignUp(c, signUpRequest)
}

func NewSignUpUsecase(signUpRepository domain.SignUpRepository, timeout time.Duration) domain.SignUpUsecase {
	return &signUpUsecase{
		signUpRepository: signUpRepository,
		ContextTimeout:   timeout,
	}
}

// login usecase
type loginUsecase struct {
	loginRepository domain.LoginRepository
	ContextTimeout  time.Duration
}

// Login implements domain.LoginUsecase.
func (l *loginUsecase) Login(c context.Context, loginRequest domain.LoginRequest) (domain.LoginResponse, error) {
	return l.loginRepository.Login(c, loginRequest)
}

func NewLoginUsecase(loginRepository domain.LoginRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		loginRepository: loginRepository,
		ContextTimeout:  timeout,
	}
}

// logout usecase
type logoutUsecase struct {
	logoutRepository domain.LogoutRepository
	ContextTimeout   time.Duration
}

// Logout implements domain.LogoutUsecase.
func (l *logoutUsecase) Logout(c context.Context, logoutRequest domain.LogoutRequest) (domain.LogoutResponse, error) {
	panic("unimplemented")
}		

func NewLogoutUsecase(logoutRepository domain.LogoutRepository, timeout time.Duration) domain.LogoutUsecase {
	return &logoutUsecase{
		logoutRepository: logoutRepository,
		ContextTimeout:   timeout,
	}
}



