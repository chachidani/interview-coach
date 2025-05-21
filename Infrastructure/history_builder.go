package infrastructure

import (
	"fmt"
	"strings"

	domain "github.com/chachidani/interview-coach-backend/Domain"
)

// BuildMessageHistory builds a formatted message history for a single room
func BuildMessageHistory(room domain.Room) string {
	var messageHistory strings.Builder
	for _, msg := range room.Messages {
		if msg.Sender == "user" {
			messageHistory.WriteString(fmt.Sprintf("User: %s\n", msg.Text))
		} else {
			messageHistory.WriteString(fmt.Sprintf("AI: %s\n", msg.Text))
		}
	}
	return messageHistory.String()
}

// BuildRoomsMessageHistory builds a formatted message history for multiple rooms
func BuildRoomsMessageHistory(rooms []domain.Room) string {
	var allHistory strings.Builder
	for _, room := range rooms {
		allHistory.WriteString(BuildMessageHistory(room))
	}
	return allHistory.String()
}

// BuildRoomMessageHistoryWithTimestamps builds a formatted message history with timestamps for a single room
func BuildRoomMessageHistoryWithTimestamps(room domain.Room) string {
	var messageHistory strings.Builder

	for _, msg := range room.Messages {
		timestamp := fmt.Sprintf("[%d]", msg.Timestamp)
		if msg.Sender == "user" {
			messageHistory.WriteString(fmt.Sprintf("%s User: %s\n", timestamp, msg.Text))
		} else {
			messageHistory.WriteString(fmt.Sprintf("%s AI: %s\n", timestamp, msg.Text))
		}
	}

	return messageHistory.String()
}

// BuildRoomsMessageHistoryWithTimestamps builds a formatted message history with timestamps for multiple rooms
func BuildRoomsMessageHistoryWithTimestamps(rooms []domain.Room) string {
	var allHistory strings.Builder

	for i, room := range rooms {
		allHistory.WriteString(fmt.Sprintf("=== Room %d ===\n", i+1))
		allHistory.WriteString(fmt.Sprintf("Topic: %s\n", room.Topic))
		allHistory.WriteString(fmt.Sprintf("Role: %s\n", room.Role))
		allHistory.WriteString("Messages:\n")
		allHistory.WriteString(BuildRoomMessageHistoryWithTimestamps(room))
		allHistory.WriteString("\n")
	}

	return allHistory.String()
}
