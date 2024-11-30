package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"streaming-platform/internal/services"
	"streaming-platform/internal/storage"
	"streaming-platform/utils"
)

type ProcessHandler struct {
	S3Client *storage.S3Client
}


func NewProcessHandler(s3Client *storage.S3Client) *ProcessHandler {
	return &ProcessHandler{
		S3Client: s3Client,
	}
}

func (h *ProcessHandler) HandleProcess(w http.ResponseWriter, r *http.Request) {
	videoKey := r.URL.Query().Get("videoKey")
	if videoKey == "" {
		http.Error(w, "Missing video key", http.StatusBadRequest)
		return
	}

	// Remover a extensão .mp4 se presente
	videoID := utils.RemoveExtensionID(videoKey)
	log.Printf("videoID para processamento: %s", videoID) // Log para verificar

	ctx := context.Background()

	// Baixar vídeo do S3
	videoData, err := h.S3Client.DownloadFile(ctx, "videos/"+videoKey)
	if err != nil {
		http.Error(w, "Failed to download video from S3", http.StatusInternalServerError)
		return
	}

	// Salvar vídeo localmente
	tempFile := filepath.Join(os.TempDir(), videoKey)
	err = os.WriteFile(tempFile, videoData, 0644)
	if err != nil {
		http.Error(w, "Failed to save temporary video file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile)

	// Transcodificar o vídeo
	qualities := []string{"1080p", "720p", "480p"}
	err = services.TranscodeVideoToHLS(videoID, tempFile, qualities, h.S3Client)
	if err != nil {
		http.Error(w, "Failed to transcode video", http.StatusInternalServerError)
		return
	}

	// Upload das resoluções transcodificadas para o S3
	for _, quality := range qualities {
		playlistPath := filepath.Join("storage/hls", videoID, quality, "playlist.m3u8")
		err := h.S3Client.UploadFileFromPath(ctx, "videos-transcoded/"+videoID+"/"+quality+"/playlist.m3u8", playlistPath)
		if err != nil {
			http.Error(w, "Failed to upload transcoded video", http.StatusInternalServerError)
			return
		}
	}

	// Geração de thumbnail
	thumbnailPath := filepath.Join(os.TempDir(), "thumbnail.jpg")
	err = services.GenerateThumbnail(tempFile, thumbnailPath)
	if err != nil {
		http.Error(w, "Failed to generate thumbnail", http.StatusInternalServerError)
		return
	}
	defer os.Remove(thumbnailPath)

	// Upload da thumbnail para o S3
	err = h.S3Client.UploadFileFromPath(ctx, "thumbnails/"+videoID+".jpg", thumbnailPath)
	if err != nil {
		http.Error(w, "Failed to upload thumbnail", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Video %s processed and uploaded successfully\n", videoID)
}




