package repository

import (
	"context"
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
	collection := r.database.Collection(r.collection)
	_, err := collection.DeleteOne(c, bson.M{"_id": roomID})
	if err != nil {
		return err
	}
	return nil
}

// GetRoom implements domain.RoomRepository.
func (r *roomRepository) GetRoom(c context.Context, roomID string) (domain.Room, error) {
	collection := r.database.Collection(r.collection)
	var room domain.Room
	err := collection.FindOne(c, bson.M{"_id": roomID}).Decode(&room)
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
	collection := r.database.Collection(r.collection)
	_, err := collection.UpdateOne(c, bson.M{"_id": roomID}, bson.M{"$set": room})
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
		return "", err
	}

	// Add initial message to room
	room.Messages = []domain.Message{
		{
			Sender:    "ai",
			Text:      initialMessage,
			Timestamp: time.Now().Unix(),
		},
	}
	room.CreatedAt = time.Now().Unix()

	// Create a composite ID using roomID and topic
	roomID := primitive.NewObjectID()
	compositeID := fmt.Sprintf("%s_%s", roomID.Hex(), room.Topic)
	room.ID = compositeID

	// Save room to database
	collection := r.database.Collection(r.collection)
	_, err = collection.InsertOne(c, room)
	if err != nil {
		return "", err
	}

	// Update user's rooms array with the composite ID
	userCollection := r.database.Collection("users")
	_, err = userCollection.UpdateOne(
		c,
		bson.M{"_id": room.UserID},
		bson.M{"$push": bson.M{"rooms": compositeID}},
	)
	if err != nil {
		// If updating user fails, we should delete the room to maintain consistency
		_, _ = collection.DeleteOne(c, bson.M{"_id": compositeID})
		return "", fmt.Errorf("failed to update user with room ID: %v", err)
	}

	return "Room created successfully", nil
}

// AddMessageToRoom implements domain.RoomRepository.
func (r *roomRepository) AddMessageToRoom(c context.Context, roomID string, message domain.Message) (domain.Room, error) {
	// First get the current room
	collection := r.database.Collection(r.collection)
	var room domain.Room
	err := collection.FindOne(c, bson.M{"_id": roomID}).Decode(&room)
	if err != nil {
		return domain.Room{}, err
	}

	// Add user's message to the room
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
		Sender:    "ai",
		Text:      aiResponse,
		Timestamp: time.Now().Unix(),
	}
	room.Messages = append(room.Messages, aiMessage)

	// Update room in database
	_, err = collection.UpdateOne(
		c,
		bson.M{"_id": roomID},
		bson.M{"$set": bson.M{"messages": room.Messages}},
	)
	if err != nil {
		return domain.Room{}, err
	}

	return room, nil
}
