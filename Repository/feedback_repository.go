package repository

import (
	"context"

	domain "github.com/chachidani/interview-coach-backend/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackRepository struct {
	database         *mongo.Database
	collection       string
}

// CreateFeedback implements domain.FeedbackRepository.


// GetFeedback implements domain.FeedbackRepository.
func (f *FeedbackRepository) GetFeedback(c context.Context, roomID string) ([]domain.Feedback, error) {
	objectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}		
	
	filter := bson.M{"room_id": objectID}
	cursor, err := f.database.Collection(f.collection).Find(c, filter)
	if err != nil {
		return nil, err
	}	
	
	var feedbacks []domain.Feedback
	if err := cursor.All(c, &feedbacks); err != nil {
		return nil, err
	}
	
	return feedbacks, nil

}

func NewFeedbackRepository(database mongo.Database, collection string) domain.FeedbackRepository {
	return &FeedbackRepository{
		database:         &database,
		collection:       collection,
	}
}
