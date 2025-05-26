package usecases

import (
	"context"
	"time"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OverallFeedbackUsecase struct {
	overallFeedbackRepository domain.OverallFeedbackRepository
	ContextTimeout time.Duration
}


func NewOverallFeedbackUsecase(overallFeedbackRepository domain.OverallFeedbackRepository, timeout time.Duration) *OverallFeedbackUsecase {
	return &OverallFeedbackUsecase{
		overallFeedbackRepository: overallFeedbackRepository,
		ContextTimeout: timeout,
	}
}

func (u *OverallFeedbackUsecase) CreateOverallFeedback(ctx context.Context, overallFeedback domain.OverallFeedback) error {
	return u.overallFeedbackRepository.CreateOverallFeedback(ctx, overallFeedback)
}

func (u *OverallFeedbackUsecase) GetOverallFeedback(ctx context.Context, userID primitive.ObjectID) ([]domain.OverallFeedback, error) {
	return u.overallFeedbackRepository.GetOverallFeedback(ctx, userID)
}
