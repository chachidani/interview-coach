package controller

import (
	"net/http"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/config"
	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	FeedbackUsecase domain.FeedbackUsecase
}



func (c *FeedbackController) GetFeedback(ctx *gin.Context) {
	roomID := ctx.Param("room_id")

	feedback, err := c.FeedbackUsecase.GetFeedback(ctx, roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	successMessage := "Feedback retrieved successfully"
	ctx.JSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage, Data: feedback})
}
