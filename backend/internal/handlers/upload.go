package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"streaming-platform/internal/services"

	"streaming-platform/internal/storage"

	"streaming-platform/utils"
)

type UploadHandler struct {
	S3Client *storage.S3Client
}

func NewUploadHandler(s3Client *storage.S3Client) *UploadHandler {
	return &UploadHandler{S3Client: s3Client}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Salvar o arquivo localmente para transcodificação
	tempFilePath := fmt.Sprintf("/tmp/%s", header.Filename)
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, "Failed to save file locally", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFilePath)
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save file locally", http.StatusInternalServerError)
		return
	}

	videoID := utils.RemoveExtensionID(header.Filename)
	log.Printf("videoID após a remoção da extensão: %s", videoID)
	// Transcodificar e enviar para S3
	qualities := []string{"1080p", "720p", "480p"} // Defina as qualidades desejadas

	err = services.TranscodeVideoToHLS(videoID, tempFilePath, qualities, h.S3Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transcode and upload video: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Video uploaded and transcoded successfully: %s\n", videoID)
}

