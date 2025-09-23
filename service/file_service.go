package service

import (
	"time"

	"github.com/nileshshrs/infinite-storage/model"
	"github.com/nileshshrs/infinite-storage/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileService handles saving file metadata to MongoDB
type FileService struct {
	repo repository.FileRepository
}

// NewFileService returns a new FileService
func NewFileService(repo repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

// SaveUploadedFile creates a File model from Discord chunks and inserts it into MongoDB
func (s *FileService) SaveUploadedFile(
	originalName string,
	size int64,
	channelID string,
	userID *primitive.ObjectID, // optional
	uploadedChunks []UploadedChunk,
) (*model.File, error) {

	// Map UploadedChunk -> model.Chunk
	var chunks []model.Chunk
	for i, c := range uploadedChunks {
		chunks = append(chunks, model.Chunk{
			Index:     i + 1,
			MessageID: c.MessageID,
			URL:       c.URL,
			Filename:  c.Filename,
			Size:      int64(c.Size),
		})
	}

	file := &model.File{
		Name:      originalName,
		Size:      size,
		ChannelID: channelID,
		UserID:    userID,
		Chunks:    chunks,
		CreatedAt: time.Now(),
	}

	// Insert into MongoDB
	if err := s.repo.Insert(file); err != nil {
		return nil, err
	}

	return file, nil
}
