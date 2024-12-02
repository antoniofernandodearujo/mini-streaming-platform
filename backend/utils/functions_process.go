package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"streaming-platform/internal/services"
	"streaming-platform/internal/storage"
)

func ProcessVideos(s3Client *storage.S3Client, qualities []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listar todos os vídeos no bucket na pasta 'videos/'
	videoKeys, err := s3Client.ListFiles(ctx, "videos/")
	if err != nil {
		return fmt.Errorf("erro ao listar vídeos: %v", err)
	}

	// Canal para distribuir trabalhos
	videoQueue := make(chan string, len(videoKeys))

	// Canal para sinalizar erros
	errorChannel := make(chan error, len(videoKeys))

	// Iniciar workers
	const workerCount = 5 // Número de goroutines
	for i := 0; i < workerCount; i++ {
		go videoWorker(ctx, videoQueue, errorChannel, s3Client, qualities)
	}

	// Enviar vídeos para a fila
	for _, videoKey := range videoKeys {
		videoQueue <- videoKey
	}
	close(videoQueue)

	// Monitorar erros
	var hasError bool
	for i := 0; i < len(videoKeys); i++ {
		if err := <-errorChannel; err != nil {
			hasError = true
			fmt.Printf("Erro no processamento: %v\n", err)
		}
	}

	if hasError {
		return fmt.Errorf("houve erros no processamento de alguns vídeos")
	}

	return nil
}

// Worker que processa vídeos de forma concorrente
func videoWorker(
	ctx context.Context,
	videoQueue <-chan string,
	errorChannel chan<- error,
	s3Client *storage.S3Client,
	qualities []string,
) {
	for videoKey := range videoQueue {
		select {
		case <-ctx.Done():
			// Interromper se o contexto for cancelado
			errorChannel <- fmt.Errorf("processamento cancelado: %s", ctx.Err())
			return
		default:
			err := processSingleVideo(ctx, videoKey, s3Client, qualities)
			errorChannel <- err
		}
	}
}

// Processar um único vídeo
func processSingleVideo(ctx context.Context, videoKey string, s3Client *storage.S3Client, qualities []string) error {
	fmt.Printf("Processando vídeo: %s\n", videoKey)

	// Baixar o vídeo para processamento local
	videoData, err := s3Client.DownloadFile(ctx, videoKey)
	if err != nil {
		return fmt.Errorf("erro ao baixar vídeo %s: %v", videoKey, err)
	}

	tempFile := filepath.Join(os.TempDir(), filepath.Base(videoKey))
	err = os.WriteFile(tempFile, videoData, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar vídeo localmente %s: %v", videoKey, err)
	}
	defer os.Remove(tempFile)

	videoID := filepath.Base(videoKey)
	err = services.TranscodeVideoToHLS(videoID, tempFile, qualities)
	if err != nil {
		return fmt.Errorf("erro ao transcodificar vídeo %s: %v", videoKey, err)
	}

	for _, quality := range qualities {
		playlistPath := filepath.Join(os.Getenv("STORAGE_PATH"), "hls", videoID, quality, "video.m3u8")
		err := s3Client.UploadFileFromPath(ctx, fmt.Sprintf("videos-transcoded/%s/%s/video.m3u8", videoID, quality), playlistPath)
		if err != nil {
			return fmt.Errorf("erro ao fazer upload da qualidade %s para %s: %v", quality, videoKey, err)
		}
	}

	thumbnailPath := filepath.Join(os.TempDir(), "thumbnail.jpg")
	err = services.GenerateThumbnail(tempFile, thumbnailPath)
	if err != nil {
		return fmt.Errorf("erro ao gerar miniatura para vídeo %s: %v", videoKey, err)
	}
	defer os.Remove(thumbnailPath)

	err = s3Client.UploadFileFromPath(ctx, fmt.Sprintf("thumbnails/%s.jpg", videoID), thumbnailPath)
	if err != nil {
		return fmt.Errorf("erro ao fazer upload da miniatura do vídeo %s: %v", videoKey, err)
	}

	return nil
}