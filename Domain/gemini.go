package domain

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiUsecase interface {
	GenerateResponse(request GeminiRequest) (string, error)
}

type GeminiRepository interface {
	GenerateResponse(request GeminiRequest) (string, error)
}




