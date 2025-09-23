package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nileshshrs/infinite-storage/model"
)

// FileRepository handles file storage in MongoDB
type FileRepository interface {
	Insert(file *model.File) error
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
	file.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.Background(), file)
	return err
}
