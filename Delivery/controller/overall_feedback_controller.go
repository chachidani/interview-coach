package controller

import (
	"net/http"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OverallFeedbackController struct {
	OverallFeedbackUsecase domain.OverallFeedbackUsecase
}

func (c *OverallFeedbackController) CreateOverallFeedback(ctx *gin.Context) {
	var overallFeedback domain.OverallFeedback
	if err := ctx.ShouldBindJSON(&overallFeedback); err != nil {
		ctx.JSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	err := c.OverallFeedbackUsecase.CreateOverallFeedback(ctx, overallFeedback)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	successMessage := "Overall feedback created successfully"
	ctx.JSON(http.StatusCreated, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage})
}

func (c *OverallFeedbackController) GetOverallFeedback(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	overallFeedback, err := c.OverallFeedbackUsecase.GetOverallFeedback(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}	

	successMessage := "Overall feedback retrieved successfully"
	ctx.JSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage, Data: overallFeedback})
}