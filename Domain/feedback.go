package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FeedbackCollection = "feedbacks"
)

type Feedback struct {
	ID              primitive.ObjectID `json:"id"`
	UserID          primitive.ObjectID `json:"user_id"`
	RoomID          primitive.ObjectID `json:"room_id"`
	MessageID       primitive.ObjectID `json:"message_id"`
	Question        string             `json:"question"`
	Answer          string             `json:"answer"`
	Strength        []string           `json:"strength"`
	ToImprove       []string           `json:"to_improve"`
	ScorePercentage int                `json:"score_percentage"`
	CreatedAt       int64              `json:"created_at"`
}

type FeedbackRepository interface {
	GetFeedback(c context.Context, roomID string) ([]Feedback, error)
}

type FeedbackUsecase interface {
	GetFeedback(c context.Context, roomID string) ([]Feedback, error)
}
