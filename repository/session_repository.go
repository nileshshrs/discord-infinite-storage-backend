package repository

import (
	"context"

	"github.com/nileshshrs/infinite-storage/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepo struct {
	collection *mongo.Collection
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(col *mongo.Collection) *SessionRepo {
	return &SessionRepo{collection: col}
}

// Create inserts a new session document into MongoDB
func (r *SessionRepo) Create(session *model.Session) error {
	_, err := r.collection.InsertOne(context.Background(), session)
	return err
}

// FindByUserID retrieves the session for a given user
func (r *SessionRepo) FindByUserID(userID interface{}) (*model.Session, error) {
	var session model.Session
	err := r.collection.FindOne(context.Background(), map[string]interface{}{
		"userID": userID,
	}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Optional: Delete session by ID (for logout)
func (r *SessionRepo) DeleteByID(sessionID interface{}) error {
	_, err := r.collection.DeleteOne(context.Background(), map[string]interface{}{
		"_id": sessionID,
	})
	return err
}
