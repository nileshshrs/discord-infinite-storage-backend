package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nileshshrs/infinite-storage/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileRepository handles file storage in MongoDB
type FileRepository interface {
	Insert(file *model.File) error
	FindByUser(userID primitive.ObjectID) ([]*model.File, error)
	FindAll() ([]*model.File, error)
}

// fileRepository implements FileRepository
type fileRepository struct {
	collection *mongo.Collection
}

// NewFileRepository returns a new file repository
func NewFileRepository(col *mongo.Collection) FileRepository {
	return &fileRepository{collection: col}
}

// Insert saves a file document into MongoDB
func (r *fileRepository) Insert(file *model.File) error {
	file.CreatedAt = file.CreatedAt.UTC()
	_, err := r.collection.InsertOne(context.Background(), file)
	return err
}

// FindByUser returns all files for a given user
func (r *fileRepository) FindByUser(userID primitive.ObjectID) ([]*model.File, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var files []*model.File
	for cursor.Next(context.Background()) {
		var f model.File
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}

// FindAll returns all files in the collection
func (r *fileRepository) FindAll() ([]*model.File, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var files []*model.File
	for cursor.Next(context.Background()) {
		var f model.File
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}
