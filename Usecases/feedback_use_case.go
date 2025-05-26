package usecases

import (
	"context"



	domain "github.com/chachidani/interview-coach-backend/Domain"
)

type FeedbackUsecase struct {
	feedbackRepository domain.FeedbackRepository
}

func NewFeedbackUsecase(feedbackRepository domain.FeedbackRepository) *FeedbackUsecase {
	return &FeedbackUsecase{
		feedbackRepository: feedbackRepository,
	}
}

func (u *FeedbackUsecase) GetFeedback(ctx context.Context, roomID string) ([]domain.Feedback, error) {
	return u.feedbackRepository.GetFeedback(ctx, roomID)
}
		