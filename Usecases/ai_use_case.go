package usecases

import (
	"time"

	domain "github.com/chachidani/interview-coach-backend/Domain"
)

type aiUsecase struct {
	aiRepository   domain.GeminiRepository
	ContextTimeout time.Duration
}

// GenerateResponse implements domain.GeminiUsecase.
func (a *aiUsecase) GenerateResponse(request domain.GeminiRequest) (string, error) {
	panic("unimplemented")
}

func NewAiUsecase(aiRepository domain.GeminiRepository, timeout time.Duration) domain.GeminiUsecase {
	return &aiUsecase{
		aiRepository:   aiRepository,
		ContextTimeout: timeout,
	}
}
