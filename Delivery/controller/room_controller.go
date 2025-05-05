package controller

import (
	"net/http"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomController struct {
	RoomUsecase domain.RoomUsecase
}

func (uc *RoomController) CreateRoom(c *gin.Context) {
	var room domain.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	roomResponse, err := uc.RoomUsecase.CreateRoom(c, room)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusCreated, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &roomResponse})
}

func (uc *RoomController) GetRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: "room ID is required", SuccessResponse: false})
		return
	}

	roomResponse, err := uc.RoomUsecase.GetRoom(c, roomID)
	if err != nil {
		if err.Error() == "room not found" {
			c.IndentedJSON(http.StatusNotFound, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, config.ResponseData{
		Error:           false,
		SuccessResponse: true,
		SuccessMessage:  nil,
		Data:            roomResponse,
	})
}

func (uc *RoomController) GetRoomsWithUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: "user ID is required", SuccessResponse: false})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	rooms, err := uc.RoomUsecase.GetRoomsWithUserID(c, objectID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusOK, config.ResponseData{
		Error:           false,
		SuccessResponse: true,
		Data:            rooms,
	})
}

func (uc *RoomController) UpdateRoom(c *gin.Context) {
	roomID := c.Param("id")
	var room domain.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	roomResponse, err := uc.RoomUsecase.UpdateRoom(c, roomID, room)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	successMessage := "Room updated successfully"
	c.IndentedJSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage, Data: roomResponse})
}

func (uc *RoomController) DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")
	err := uc.RoomUsecase.DeleteRoom(c, roomID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	successMessage := "Room deleted successfully"
	c.IndentedJSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage})
}

func (uc *RoomController) AddMessageToRoom(c *gin.Context) {
	roomID := c.Param("id")
	var message domain.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	roomResponse, err := uc.RoomUsecase.AddMessageToRoom(c, roomID, message)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}
	successMessage := "Message added to room successfully"

	c.IndentedJSON(http.StatusOK, config.ResponseData{
		Error:           false,
		SuccessResponse: true,
		SuccessMessage:  &successMessage,
		Data:            roomResponse,
	})
}
