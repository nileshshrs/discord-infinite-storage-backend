package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nileshshrs/infinite-storage/config"
	"github.com/nileshshrs/infinite-storage/service"
)

type UploadHandler struct {
	uploadService *service.UploadService
	cfg           *config.Config
}

func NewUploadHandler(uploadService *service.UploadService, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		cfg:           cfg,
	}
}

// HandleUpload handles POST /api/v1/upload
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
	results, err := h.uploadService.SendFileInOrder(h.cfg.DiscordChannelID, "New file uploaded", tempPath)
	if err != nil {
		log.Printf("Failed to send file: %v", err)
		http.Error(w, `{"error":"failed to send file to Discord"}`, http.StatusInternalServerError)
		os.Remove(tempPath)
		return
	}

	// Clean up temp file
	os.Remove(tempPath)

	// Return success + trimmed responses
	resp := map[string]interface{}{
		"message":  "File sent to Discord successfully",
		"uploads":  results,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
