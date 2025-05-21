package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	domain "github.com/chachidani/interview-coach-backend/Domain"
)

type geminiRepository struct{}

func NewGeminiRepository() domain.GeminiRepository {
	return &geminiRepository{}
}

func (g *geminiRepository) GenerateResponse(request domain.GeminiRequest) (string, error) {
	apiKey := bootstrap.NewEnv().GeminiAPIKey
	if apiKey == "" {
		apiKey = "AIzaSyCNv7mKF-hXXjyuKlWBSDi9CojC1S7ca3M"
		fmt.Println("GEMINI_API_KEY environment variable is not set")
		// return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": request.Contents[0].Parts[0].Text}}},
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is not OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var geminiResp domain.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(body))
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response: %s", string(body))
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in candidate: %s", string(body))
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
