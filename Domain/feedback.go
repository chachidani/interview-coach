package domain

type Feedback struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	RoomID      string   `json:"room_id"`
	Question    string   `json:"question"`
	Answer      string   `json:"answer"`
	Strength    []string `json:"strength"`
	ToImprove   []string `json:"to_improve"`
	ScorePercentage int      `json:"score_percentage"`
	CreatedAt   int64    `json:"created_at"`
}

type FeedbackRepository interface {
	CreateFeedback(feedback Feedback) error
	GetFeedback(roomID string) ([]Feedback, error)
}

type FeedbackUsecase interface {
	CreateFeedback(feedback Feedback) error
	GetFeedback(roomID string) ([]Feedback, error)
}
