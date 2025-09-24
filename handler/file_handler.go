package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nileshshrs/infinite-storage/middlewares"
	"github.com/nileshshrs/infinite-storage/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler returns a new FileHandler
func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// GetUserFiles handles GET /api/v1/files
func (h *FileHandler) GetUserFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDHex, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userIDHex == "" {
		http.Error(w, `{"error":"user not authenticated"}`, http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusInternalServerError)
		return
	}

	files, err := h.fileService.GetFilesByUser(userID)
	if err != nil {
		http.Error(w, `{"error":"failed to retrieve files"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"files": files,
	}
	json.NewEncoder(w).Encode(resp)
}
