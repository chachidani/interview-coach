package repository

import (
	"context"
	"fmt"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type signUpRepository struct {
	database        mongo.Database
	collection      string
	passwordService *middleware.PasswordService
}

// GetUser implements domain.SignUpRepository.
func (s *signUpRepository) GetUser(c context.Context) ([]domain.User, error) {
	collection := s.database.Collection(s.collection)
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []domain.User
	if err := cursor.All(c, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// SignUp implements domain.SignUpRepository.
func (s *signUpRepository) SignUp(c context.Context, signUpRequest domain.SignUpRequest) (domain.SignUpResponse, error) {
	collection := s.database.Collection(s.collection)

	var existingUser domain.User
	err := collection.FindOne(c, bson.M{"email": signUpRequest.Email}).Decode(&existingUser)
	if err == nil {
		return domain.SignUpResponse{}, fmt.Errorf("user with email %s already exists", signUpRequest.Email)
	}

	if err != mongo.ErrNoDocuments {
		return domain.SignUpResponse{}, err
	}

	hashedPassword, err := s.passwordService.HashPassword(signUpRequest.Password)
	if err != nil {
		return domain.SignUpResponse{}, err
	}

	user := domain.User{
		ID:       primitive.NewObjectID(),
		Username: signUpRequest.Username,
		Email:    signUpRequest.Email,
		Password: hashedPassword,
		Rooms:    []string{},
	}

	_, err = collection.InsertOne(c, user)
	if err != nil {
		return domain.SignUpResponse{}, err
	}

	return domain.SignUpResponse{
		Message: "User created successfully",
	}, nil
}

func NewSignUpRepository(database mongo.Database, collection string, passwordService *middleware.PasswordService) domain.SignUpRepository {
	return &signUpRepository{
		database:        database,
		collection:      collection,
		passwordService: passwordService,
	}
}

// login repository
type loginRepository struct {
	database        mongo.Database
	collection      string
	passwordService *middleware.PasswordService
	jwtService      *middleware.JWTService
}

// Login implements domain.LoginRepository.
func (l *loginRepository) Login(c context.Context, loginRequest domain.LoginRequest) (domain.LoginResponse, error) {
	collection := l.database.Collection(l.collection)

	var user domain.User
	err := collection.FindOne(c, bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.LoginResponse{}, fmt.Errorf("user not found")
		}
		return domain.LoginResponse{}, err
	}

	if err := l.passwordService.VerifyPassword(user.Password, loginRequest.Password); err != nil {
		return domain.LoginResponse{}, fmt.Errorf("invalid password")
	}

	token, err := l.jwtService.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("failed to generate token")
	}

	return domain.LoginResponse{
		Message: "Login successful",
		Token:   token,
	}, nil

}

func NewLoginRepository(database mongo.Database, collection string, passwordService *middleware.PasswordService, jwtService *middleware.JWTService) domain.LoginRepository {
	return &loginRepository{
		database:        database,
		collection:      collection,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// logout repository
type logoutRepository struct {
	database   mongo.Database
	collection string
}

func NewLogoutRepository(database mongo.Database, collection string) domain.LogoutRepository {
	return &logoutRepository{
		database:   database,
		collection: collection,
	}
}

// Logout implements domain.LogoutRepository.
func (l *logoutRepository) Logout(c context.Context, logoutRequest domain.LogoutRequest) (domain.LogoutResponse, error) {
	panic("unimplemented")
}
