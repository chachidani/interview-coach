package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionOverallFeedback = "overall_feedbacks"
)

type OverallFeedback struct {
	ID              primitive.ObjectID `json:"id"`
	UserID          primitive.ObjectID `json:"user_id"`
	Strength        []string           `json:"strength"`
	Improvement     []string           `json:"improvement"`
	TopTopic        string             `json:"top_topic"`
	TotalInterview  int64              `json:"total_interview"`
	ScorePercentage int                `json:"score_percentage"`
	CreatedAt       int64              `json:"created_at"`
}

type OverallFeedbackRepository interface {
	CreateOverallFeedback(c context.Context, overallFeedback OverallFeedback) error
	GetOverallFeedback(c context.Context, userID primitive.ObjectID) ([]OverallFeedback, error)
}

type OverallFeedbackUsecase interface {
	CreateOverallFeedback(c context.Context, overallFeedback OverallFeedback) error
	GetOverallFeedback(c context.Context, userID primitive.ObjectID) ([]OverallFeedback, error)
}
