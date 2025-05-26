package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	infrastructure "github.com/chachidani/interview-coach-backend/Infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OverallFeedbackRepository struct {
	database         *mongo.Database
	collection       string
	geminiRepository domain.GeminiRepository
	roomRepository   domain.RoomRepository
}

func NewOverallFeedbackRepository(database *mongo.Database, collection string, geminiRepository domain.GeminiRepository, roomRepository domain.RoomRepository) *OverallFeedbackRepository {
	return &OverallFeedbackRepository{
		database:         database,
		collection:       collection,
		geminiRepository: geminiRepository,
		roomRepository:   roomRepository,
	}
}

func (r *OverallFeedbackRepository) CreateOverallFeedback(c context.Context, overallFeedback domain.OverallFeedback) error {
	// Get all completed rooms for the user
	rooms, err := r.roomRepository.GetRoomsWithUserID(c, overallFeedback.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user rooms: %v", err)
	}

	// Filter only completed rooms
	var completedRooms []domain.Room
	for _, room := range rooms {
		if room.Status == "completed" {
			completedRooms = append(completedRooms, room)
		}
	}

	if len(completedRooms) == 0 {
		return fmt.Errorf("no completed rooms found for user")
	}

	// Build message history for all completed rooms
	messageHistory := infrastructure.BuildRoomsMessageHistory(completedRooms)

	// Generate overall feedback using Gemini
	prompt := fmt.Sprintf(`You are an AI interviewer providing overall feedback. Analyze the following interview sessions and provide comprehensive feedback.

Interview History:
%s

Please provide:
1. List of overall strengths across all interviews
2. Areas that need improvement
3. Top performing topic/area
4. Overall score percentage (0-100)

Format your response as JSON:
{
    "strength": ["strength1", "strength2"],
    "improvement": ["improvement1", "improvement2"],
    "top_topic": "topic name",
    "score_percentage": 85
}`, messageHistory)

	geminiRequest := domain.GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	feedbackResponse, err := r.geminiRepository.GenerateResponse(geminiRequest)
	if err != nil {
		return fmt.Errorf("failed to generate overall feedback: %v", err)
	}

	// Parse the JSON response
	var feedbackData struct {
		Strength        []string `json:"strength"`
		Improvement     []string `json:"improvement"`
		TopTopic        string   `json:"top_topic"`
		ScorePercentage int      `json:"score_percentage"`
	}

	if err := json.Unmarshal([]byte(feedbackResponse), &feedbackData); err != nil {
		return fmt.Errorf("failed to parse feedback response: %v", err)
	}

	// Update overall feedback with AI-generated data
	overallFeedback.ID = primitive.NewObjectID()
	overallFeedback.Strength = feedbackData.Strength
	overallFeedback.Improvement = feedbackData.Improvement
	overallFeedback.TopTopic = feedbackData.TopTopic
	overallFeedback.ScorePercentage = feedbackData.ScorePercentage
	overallFeedback.TotalInterview = int64(len(completedRooms))
	overallFeedback.CreatedAt = time.Now().Unix()

	// Save to database
	collection := r.database.Collection(r.collection)
	_, err = collection.InsertOne(c, overallFeedback)
	if err != nil {
		return fmt.Errorf("failed to save overall feedback: %v", err)
	}

	return nil
}

func (r *OverallFeedbackRepository) GetOverallFeedback(c context.Context, userID primitive.ObjectID) ([]domain.OverallFeedback, error) {
	collection := r.database.Collection(r.collection)
	var feedbacks []domain.OverallFeedback

	cursor, err := collection.Find(c, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find overall feedbacks: %v", err)
	}
	defer cursor.Close(c)

	if err := cursor.All(c, &feedbacks); err != nil {
		return nil, fmt.Errorf("failed to decode overall feedbacks: %v", err)
	}

	return feedbacks, nil
}
