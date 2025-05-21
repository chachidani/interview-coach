package domain

type OverallFeedback struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Strength  []string `json:"strength"`
	Improvement []string `json:"improvement"`
	TopTopic string `json:"top_topic"`
	TotalInterview int64 `json:"total_interview"`
	ScorePercentage int      `json:"score_percentage"`
	CreatedAt   int64    `json:"created_at"`
}

type OverallFeedbackRepository interface {
	CreateOverallFeedback(overallFeedback OverallFeedback) error
	GetOverallFeedback(userID string) ([]OverallFeedback, error)
}

type OverallFeedbackUsecase interface {
	CreateOverallFeedback(overallFeedback OverallFeedback) error
	GetOverallFeedback(userID string) ([]OverallFeedback, error)
}
