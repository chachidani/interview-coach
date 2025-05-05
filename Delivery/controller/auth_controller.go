package controller

import (
	"net/http"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/config"
	"github.com/gin-gonic/gin"
)

type SignUpController struct {
	SignUpUsecase domain.SignUpUsecase
}

type LoginController struct {
	LoginUsecase domain.LoginUsecase
}

type LogoutController struct {
	LogoutUsecase domain.LogoutUsecase
}


func (uc *SignUpController) SignUp(c *gin.Context) {
	var signUpRequest domain.SignUpRequest
	if err := c.ShouldBindJSON(&signUpRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	signUpResponse, err := uc.SignUpUsecase.SignUp(c, signUpRequest)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusCreated, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &signUpResponse.Message})
}


func (uc *LoginController) Login(c *gin.Context) {
	var loginRequest domain.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	loginResponse, err := uc.LoginUsecase.Login(c, loginRequest)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusOK, loginResponse)	
}

func (uc *SignUpController) GetUser(c *gin.Context) {
	var users []domain.User
	users, err := uc.SignUpUsecase.GetUser(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}


func (uc *LogoutController) Logout(c *gin.Context) {
	var logoutRequest domain.LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	logoutResponse, err := uc.LogoutUsecase.Logout(c, logoutRequest)	
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &logoutResponse.Message})
}


