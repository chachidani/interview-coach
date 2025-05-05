package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionRoom = "rooms"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Role      string             `bson:"role"`
	Topic     string             `bson:"topic"`
	Messages  []Message          `bson:"messages"`
	CreatedAt int64              `bson:"created_at"`
}

type RoomRequest struct {
	UserID primitive.ObjectID `bson:"user_id"`
	Role   string             `bson:"role"`
	Topic  string             `bson:"topic"`
}

type Message struct {
	Sender    string `bson:"sender"` // "user" or "ai"
	Text      string `bson:"text"`
	VoiceURL  string `bson:"voice_url,omitempty"`
	Timestamp int64  `bson:"timestamp"`
}

type RoomRepository interface {
	CreateRoom(c context.Context, room Room) (string, error)
	GetRoom(c context.Context, roomID primitive.ObjectID) (Room, error)
	GetRoomsWithUserID(c context.Context, userID primitive.ObjectID) ([]Room, error)
	UpdateRoom(c context.Context, roomID primitive.ObjectID, room Room) (Room, error)
	DeleteRoom(c context.Context, roomID primitive.ObjectID) error
	//add message to room
	AddMessageToRoom(c context.Context, roomID primitive.ObjectID, message Message) (Room, error)
}

type RoomUsecase interface {
	CreateRoom(c context.Context, room Room) (string, error)
	GetRoom(c context.Context, roomID primitive.ObjectID) (Room, error)
	GetRoomsWithUserID(c context.Context, userID primitive.ObjectID) ([]Room, error)
	UpdateRoom(c context.Context, roomID primitive.ObjectID, room Room) (Room, error)
	DeleteRoom(c context.Context, roomID primitive.ObjectID) error
	//add message to room
	AddMessageToRoom(c context.Context, roomID primitive.ObjectID, message Message) (Room, error)
}
