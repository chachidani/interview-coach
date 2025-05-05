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

	c.IndentedJSON(http.StatusCreated, roomResponse)
}


func (uc *RoomController) GetRoom(c *gin.Context) {
	

	roomID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	roomResponse, err := uc.RoomUsecase.GetRoom(c, objectID)
	if err != nil {
		if err.Error() == "room not found" {
			c.IndentedJSON(http.StatusNotFound, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, roomResponse)
}	


func (uc *RoomController) GetRoomsWithUserID(c *gin.Context) {
	userID := c.Param("id")
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

	c.IndentedJSON(http.StatusOK, rooms)
}


func (uc *RoomController) UpdateRoom(c *gin.Context) {
	roomID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})	
		return
	}

	var room domain.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}		

	roomResponse, err := uc.RoomUsecase.UpdateRoom(c, objectID, room)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	c.IndentedJSON(http.StatusOK, roomResponse)
}	


func (uc *RoomController) DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})	
		return
	}

	err = uc.RoomUsecase.DeleteRoom(c, objectID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})
		return
	}

	successMessage := "Room deleted successfully"
	c.IndentedJSON(http.StatusOK, config.ResponseData{Error: false, SuccessResponse: true, SuccessMessage: &successMessage})
}


func (uc *RoomController) AddMessageToRoom(c *gin.Context) {
	roomID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})	
		return
	}

	var message domain.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})	
		return
	}

	roomResponse, err := uc.RoomUsecase.AddMessageToRoom(c, objectID, message)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, config.ResponseData{Error: true, ErrorMessage: err.Error(), SuccessResponse: false})	
		return
	}

	c.IndentedJSON(http.StatusOK, roomResponse)
}



