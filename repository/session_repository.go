package repository

import (
	"context"
	"errors"
	"github.com/nileshshrs/infinite-storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	result, err := r.collection.InsertOne(context.Background(), session)
	if err != nil {
		return err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to get inserted ID")
	}
	session.ID = oid
	return nil
}

// FindByUserID retrieves the session for a given user
func (r *SessionRepo) FindByUserID(userID interface{}) (*model.Session, error) {
	var session model.Session
	err := r.collection.FindOne(context.Background(), bson.M{"userID": userID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByID retrieves a session by its ID
func (r *SessionRepo) FindByID(sessionID primitive.ObjectID) (*model.Session, error) {
	var session model.Session
	err := r.collection.FindOne(context.Background(), bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Update updates fields of a session document
func (r *SessionRepo) Update(sessionID primitive.ObjectID, updates map[string]interface{}) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": sessionID},
		bson.M{"$set": updates},
	)
	return err
}

// Optional: Delete session by ID (for logout)
func (r *SessionRepo) DeleteByID(sessionID interface{}) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": sessionID})
	return err
}
