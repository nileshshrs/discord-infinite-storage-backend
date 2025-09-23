package service

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

const chunkSize = 8 * 1024 * 1024 // 8MB

type UploadService struct {
	dg *discordgo.Session
}

func NewUploadService(dg *discordgo.Session) *UploadService {
	return &UploadService{dg: dg}
}

// Trimmed response for DB and API
type UploadedChunk struct {
	MessageID    string `json:"message_id"`
	ChannelID    string `json:"channel_id"`
	AttachmentID string `json:"attachment_id"`
	URL          string `json:"url"`
	Filename     string `json:"filename"`
	Size         int    `json:"size"`
	ContentType  string `json:"content_type"`
	Timestamp    string `json:"timestamp"`
}

// SendFileInOrder uploads file to Discord and returns trimmed responses
func (s *UploadService) SendFileInOrder(channelID, message, filePath string) ([]UploadedChunk, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	var results []UploadedChunk

	// Small file (≤ 8MB)
	if info.Size() <= chunkSize {
		msg, err := s.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content: message,
			Files: []*discordgo.File{
				{
					Name:   info.Name(),
					Reader: mustOpen(filePath),
				},
			},
		})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Uploaded file %s\n", info.Name())
		results = append(results, extractUploadedChunk(msg))
		return results, nil
	}

	// Large file (> 8MB): split into chunks
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	totalSize := info.Size()
	chunks := int(totalSize / chunkSize)
	if totalSize%chunkSize != 0 {
		chunks++
	}

	buffer := make([]byte, chunkSize)
	for i := 0; i < chunks; i++ {
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return nil, err
		}

		// Create temp chunk file
		chunkPath := fmt.Sprintf("%s.part%d", filePath, i+1)
		tmpFile, _ := os.Create(chunkPath)
		tmpFile.Write(buffer[:n])
		tmpFile.Close()

		// Upload chunk
		msg, err := s.dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("%s (part %d/%d)", message, i+1, chunks),
			Files: []*discordgo.File{
				{
					Name:   fmt.Sprintf("%s.part%d", info.Name(), i+1),
					Reader: mustOpen(chunkPath),
				},
			},
		})
		if err != nil {
			os.Remove(chunkPath)
			return nil, err
		}

		fmt.Printf("Uploaded chunk %d/%d of %s\n", i+1, chunks, info.Name())
		results = append(results, extractUploadedChunk(msg))
		os.Remove(chunkPath)
	}

	return results, nil
}

// Extract only useful fields from Discord message
func extractUploadedChunk(msg *discordgo.Message) UploadedChunk {
	att := msg.Attachments[0]
	return UploadedChunk{
		MessageID:    msg.ID,
		ChannelID:    msg.ChannelID,
		AttachmentID: att.ID,
		URL:          att.URL,
		Filename:     att.Filename,
		Size:         att.Size,
		ContentType:  att.ContentType,
		Timestamp:    msg.Timestamp.String(),
	}
}

// helper to open file
func mustOpen(name string) *os.File {
	f, _ := os.Open(name)
	return f
}
