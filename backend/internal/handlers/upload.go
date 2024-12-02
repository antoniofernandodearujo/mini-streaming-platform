package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"streaming-platform/internal/services"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempDir := filepath.Join(os.TempDir(), "uploads")
	os.MkdirAll(tempDir, os.ModePerm)

	filePath := filepath.Join(tempDir, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	videoID := header.Filename
	qualities := []string{"720p", "480p"} // Redução de qualidades para economizar

	if err := services.TranscodeVideoToHLS(videoID, filePath, qualities); err != nil {
		http.Error(w, fmt.Sprintf("Erro na transcodificação: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vídeo %s transcodificado com sucesso", videoID)
}
