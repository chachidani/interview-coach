package usecases

import (
	"context"
	"time"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomUsecase struct {
	roomRepository domain.RoomRepository
	ContextTimeout time.Duration
}

// AddMessageToRoom implements domain.RoomUsecase.
func (r *roomUsecase) AddMessageToRoom(c context.Context, roomID string, message domain.Message) (domain.Room, error) {
	return r.roomRepository.AddMessageToRoom(c, roomID, message)
}

// CreateRoom implements domain.RoomUsecase.
func (r *roomUsecase) CreateRoom(c context.Context, room domain.Room) (string, error) {
	return r.roomRepository.CreateRoom(c, room)
}

// DeleteRoom implements domain.RoomUsecase.
func (r *roomUsecase) DeleteRoom(c context.Context, roomID string) error {
	return r.roomRepository.DeleteRoom(c, roomID)
}

// GetRoom implements domain.RoomUsecase.
func (r *roomUsecase) GetRoom(c context.Context, roomID string) (domain.Room, error) {
	return r.roomRepository.GetRoom(c, roomID)
}

// GetRoomsWithUserID implements domain.RoomUsecase.
func (r *roomUsecase) GetRoomsWithUserID(c context.Context, userID primitive.ObjectID) ([]domain.Room, error) {
	return r.roomRepository.GetRoomsWithUserID(c, userID)
}

// UpdateRoom implements domain.RoomUsecase.
func (r *roomUsecase) UpdateRoom(c context.Context, roomID string, room domain.Room) (domain.Room, error) {
	return r.roomRepository.UpdateRoom(c, roomID, room)
}

func NewRoomUsecase(roomRepository domain.RoomRepository, timeout time.Duration) domain.RoomUsecase {
	return &roomUsecase{
		roomRepository: roomRepository,
		ContextTimeout: timeout,
	}
}
