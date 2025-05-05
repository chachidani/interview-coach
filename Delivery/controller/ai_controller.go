package controller

import (
	"net/http"


	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiUsecase domain.GeminiUsecase
}

func (c *AIController) GenerateResponse(ctx *gin.Context) {
	request := domain.GeminiRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})	
		return
	}

	response, err := c.aiUsecase.GenerateResponse(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})	
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": response})
}


