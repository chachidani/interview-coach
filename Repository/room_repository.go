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

type roomRepository struct {
	database         mongo.Database
	collection       string
	geminiRepository domain.GeminiRepository
}

// DeleteRoom implements domain.RoomRepository.
func (r *roomRepository) DeleteRoom(c context.Context, roomID string) error {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return fmt.Errorf("invalid room ID format: %v", err)
	}

	collection := r.database.Collection(r.collection)
	_, err = collection.DeleteOne(c, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return nil
}

// GetRoom implements domain.RoomRepository.
func (r *roomRepository) GetRoom(c context.Context, roomID string) (domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return domain.Room{}, fmt.Errorf("invalid room ID format: %v", err)
	}

	collection := r.database.Collection(r.collection)
	var room domain.Room
	err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Room{}, fmt.Errorf("room not found")
		}
		return domain.Room{}, err
	}
	return room, nil
}

// GetRoomsWithUserID implements domain.RoomRepository.
func (r *roomRepository) GetRoomsWithUserID(c context.Context, userID primitive.ObjectID) ([]domain.Room, error) {
	collection := r.database.Collection(r.collection)
	var rooms []domain.Room
	cursor, err := collection.Find(c, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(c, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

// UpdateRoom implements domain.RoomRepository.
func (r *roomRepository) UpdateRoom(c context.Context, roomID string, room domain.Room) (domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return domain.Room{}, fmt.Errorf("invalid room ID format: %v", err)
	}

	collection := r.database.Collection(r.collection)
	_, err = collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": room})
	if err != nil {
		return domain.Room{}, err
	}
	return room, nil
}

func NewRoomRepository(database mongo.Database, collection string, geminiRepository domain.GeminiRepository) domain.RoomRepository {
	return &roomRepository{
		database:         database,
		collection:       collection,
		geminiRepository: geminiRepository,
	}
}

// CreateRoom implements domain.RoomRepository.
func (r *roomRepository) CreateRoom(c context.Context, room domain.Room) (string, error) {
	// Generate initial message using Gemini
	prompt := fmt.Sprintf("You are an AI interviewer. The role is %s and the topic is %s. Please start the interview with a greeting and your first question.", room.Role, room.Topic)

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

	initialMessage, err := r.geminiRepository.GenerateResponse(geminiRequest)
	if err != nil {
		return "", fmt.Errorf("failed to generate initial message: %v", err)
	}

	// Ensure room has a valid ID
	room.ID = primitive.NewObjectID()

	// Add initial message to room
	room.Messages = []domain.Message{
		{
			ID:        primitive.NewObjectID(),
			Sender:    "ai",
			Text:      initialMessage,
			Timestamp: time.Now().Unix(),
		},
	}
	room.CreatedAt = time.Now().Unix()
	room.Status = "active"

	// Save room to database
	collection := r.database.Collection(r.collection)
	_, err = collection.InsertOne(c, room)
	if err != nil {
		return "", fmt.Errorf("failed to create room: %v", err)
	}

	// Update user's rooms array with the room ID
	userCollection := r.database.Collection("users")
	_, err = userCollection.UpdateOne(
		c,
		bson.M{"_id": room.UserID},
		bson.M{"$push": bson.M{"rooms": room.ID.Hex()}},
	)
	if err != nil {
		// If updating user fails, we should delete the room to maintain consistency
		_, _ = collection.DeleteOne(c, bson.M{"_id": room.ID})
		return "", fmt.Errorf("failed to update user with room ID: %v", err)
	}

	return "Room created successfully", nil
}

// AddMessageToRoom implements domain.RoomRepository.
func (r *roomRepository) AddMessageToRoom(c context.Context, roomID string, message domain.Message) (domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return domain.Room{}, fmt.Errorf("invalid room ID format: %v", err)
	}

	// First get the current room
	collection := r.database.Collection(r.collection)
	var room domain.Room
	err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		return domain.Room{}, err
	}

	// Add user's message to the room
	message.ID = primitive.NewObjectID() // Ensure message has an ID
	room.Messages = append(room.Messages, message)

	// Format message history for prompt using the reusable builder
	messageHistory := infrastructure.BuildMessageHistory(room)

	// Create prompt for Gemini including context of the interview
	prompt := fmt.Sprintf(`You are an AI interviewer. The role is %s and the topic is %s.

Previous conversation:
%s

Please provide a relevant follow-up question or response based on the conversation above.`,
		room.Role,
		room.Topic,
		messageHistory)

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

	// Get response from Gemini
	aiResponse, err := r.geminiRepository.GenerateResponse(geminiRequest)
	if err != nil {
		return domain.Room{}, err
	}

	// Add AI's response to the room
	aiMessage := domain.Message{
		ID:        primitive.NewObjectID(),
		Sender:    "ai",
		Text:      aiResponse,
		Timestamp: time.Now().Unix(),
	}
	room.Messages = append(room.Messages, aiMessage)

	// Update room in database
	_, err = collection.UpdateOne(
		c,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"messages": room.Messages}},
	)
	if err != nil {
		return domain.Room{}, err
	}

	return room, nil
}

// CompletedRoom implements domain.RoomRepository.
func (r *roomRepository) CompletedRoom(c context.Context, userID primitive.ObjectID, roomID string) (domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return domain.Room{}, fmt.Errorf("invalid room ID format: %v", err)
	}

	collection := r.database.Collection(r.collection)
	var room domain.Room
	err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		return domain.Room{}, err
	}

	room.Status = "completed"

	// Calculate performance percentage based on feedback
	room.PerformancePercentage = 0
	room.Feedback = []domain.Feedback{}

	// Process messages in pairs (AI question + User answer)
	for i := 0; i < len(room.Messages)-1; i += 2 {
		if i+1 >= len(room.Messages) {
			break
		}

		aiMessage := room.Messages[i]
		userMessage := room.Messages[i+1]

		// Skip if not a valid AI-User pair
		if aiMessage.Sender != "ai" || userMessage.Sender != "user" {
			continue
		}

		// Create feedback for this message pair
		feedback := domain.Feedback{
			ID:              primitive.NewObjectID(),
			UserID:          userID,
			RoomID:          room.ID,
			MessageID:       userMessage.ID,
			Question:        aiMessage.Text,
			Answer:          userMessage.Text,
			Strength:        []string{}, // Will be populated by AI
			ToImprove:       []string{}, // Will be populated by AI
			ScorePercentage: 0,          // Will be calculated by AI
			CreatedAt:       time.Now().Unix(),
		}

		// Generate feedback using Gemini
		prompt := fmt.Sprintf(`You are an AI interviewer providing feedback. The role is %s and the topic is %s.

Question: %s
Answer: %s

Please provide:
1. List of strengths in the answer
2. Areas for improvement
3. Score percentage (0-100)

Format your response as JSON:
{
    "strength": ["strength1", "strength2"],
    "to_improve": ["improvement1", "improvement2"],
    "score_percentage": 85
}`, room.Role, room.Topic, aiMessage.Text, userMessage.Text)

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
			return domain.Room{}, fmt.Errorf("failed to generate feedback: %v", err)
		}

		// Parse the JSON response
		var feedbackData struct {
			Strength        []string `json:"strength"`
			ToImprove       []string `json:"to_improve"`
			ScorePercentage int      `json:"score_percentage"`
		}

		if err := json.Unmarshal([]byte(feedbackResponse), &feedbackData); err != nil {
			return domain.Room{}, fmt.Errorf("failed to parse feedback response: %v", err)
		}

		// Update feedback with AI-generated data
		feedback.Strength = feedbackData.Strength
		feedback.ToImprove = feedbackData.ToImprove
		feedback.ScorePercentage = feedbackData.ScorePercentage

		// Add feedback to room
		room.Feedback = append(room.Feedback, feedback)

		// Update performance percentage
		room.PerformancePercentage += int64(feedback.ScorePercentage)
	}

	// Calculate average performance percentage
	if len(room.Feedback) > 0 {
		room.PerformancePercentage = room.PerformancePercentage / int64(len(room.Feedback))
	}

	// Save feedbacks to feedback collection
	feedbackCollection := r.database.Collection(domain.FeedbackCollection)
	for _, feedback := range room.Feedback {
		_, err := feedbackCollection.InsertOne(c, feedback)
		if err != nil {
			return domain.Room{}, fmt.Errorf("failed to save feedback: %v", err)
		}
	}

	// Update room in database
	_, err = collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": room})
	if err != nil {
		return domain.Room{}, err
	}

	return room, nil
}
