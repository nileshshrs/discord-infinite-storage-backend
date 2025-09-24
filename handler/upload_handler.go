package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nileshshrs/infinite-storage/config"
	"github.com/nileshshrs/infinite-storage/middlewares"
	"github.com/nileshshrs/infinite-storage/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadHandler struct {
	uploadService *service.UploadService
	fileService   *service.FileService
	cfg           *config.Config
}

func NewUploadHandler(uploadService *service.UploadService, fileService *service.FileService, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		fileService:   fileService,
		cfg:           cfg,
	}
}

// HandleUpload handles POST /api/v1/files/upload
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 100MB)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse multipart form"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file not found"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file to a temp directory
	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, header.Filename)
	out, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, `{"error":"could not save temp file"}`, http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}

	// Upload to Discord (wait for all chunks)
	uploadedChunks, err := h.uploadService.SendFileInOrder(h.cfg.DiscordChannelID, "New file uploaded", tempPath)
	if err != nil {
		log.Printf("Failed to send file: %v", err)
		http.Error(w, `{"error":"failed to send file to Discord"}`, http.StatusInternalServerError)
		os.Remove(tempPath)
		return
	}

	// Get file size
	fileInfo, err := os.Stat(tempPath)
	if err != nil {
		log.Printf("Failed to get file info: %v", err)
		http.Error(w, `{"error":"failed to get file info"}`, http.StatusInternalServerError)
		os.Remove(tempPath)
		return
	}

	// Get user ID from middleware context
	var userID *primitive.ObjectID
	if userIDHex, ok := r.Context().Value(middlewares.UserIDKey).(string); ok && userIDHex != "" {
		if oid, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			userID = &oid
		}
	}

	// Save metadata to MongoDB
	fileDoc, err := h.fileService.SaveUploadedFile(
		header.Filename,        // original filename
		fileInfo.Size(),        // file size
		h.cfg.DiscordChannelID, // channel ID
		userID,                 // authenticated user
		uploadedChunks,         // uploaded chunks
	)
	if err != nil {
		log.Printf("Failed to save file metadata: %v", err)
		http.Error(w, `{"error":"failed to save file metadata"}`, http.StatusInternalServerError)
		os.Remove(tempPath)
		return
	}

	// Clean up temp file
	os.Remove(tempPath)

	// Return success + trimmed responses
	resp := map[string]interface{}{
		"message": "File sent to Discord and saved successfully",
		"file":    fileDoc,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
