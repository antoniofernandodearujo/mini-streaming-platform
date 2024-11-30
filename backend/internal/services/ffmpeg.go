package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"streaming-platform/internal/storage"
	"strings"
)

func TranscodeVideoToHLS(videoID, inputPath string, qualities []string, s3Client *storage.S3Client) error {
	baseStoragePath := os.Getenv("STORAGE_PATH")
	if baseStoragePath == "" {
		baseStoragePath = "/app/videos/hls"
	}

	// Diretório base para o vídeo
	masterDir := filepath.Join(baseStoragePath, "hls", videoID) // Sem .mp4
	log.Printf("Criando masterDir: %s", masterDir) // Log para verificar

	err := os.MkdirAll(masterDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório da master playlist: %v", err)
	}

	// Criação da master playlist
	masterPlaylistPath := filepath.Join(masterDir, "master.m3u8")
	masterPlaylistFile, err := os.Create(masterPlaylistPath)
	if err != nil {
		return fmt.Errorf("erro ao criar master playlist: %v", err)
	}
	defer masterPlaylistFile.Close()

	_, _ = masterPlaylistFile.WriteString("#EXTM3U\n")

	for _, quality := range qualities {
		err := TranscodeVideo(videoID, inputPath, quality)
		if err != nil {
			return fmt.Errorf("erro ao transcodificar para qualidade %s: %v", quality, err)
		}

		playlistPath := filepath.Join(quality, "playlist.m3u8")
		_, _ = masterPlaylistFile.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s\n",
			getBandwidth(quality), getResolutionString(quality), playlistPath))
	}

	// Fazer upload dos arquivos HLS para S3
	err = uploadHLSFilesToS3(masterDir, videoID, s3Client)
	if err != nil {
		return fmt.Errorf("erro ao fazer upload dos arquivos HLS para o S3: %v", err)
	}

	fmt.Printf("Master playlist criada e arquivos enviados para S3: %s\n", masterPlaylistPath)
	return nil
}



// TranscodeVideo transcodifica o vídeo para a qualidade especificada utilizando o FFmpeg.
func TranscodeVideo(videoID, inputPath, quality string) error {
	// Caminho de saída do arquivo transcodificado
	outputPath := filepath.Join(os.Getenv("STORAGE_PATH"), "hls", videoID, quality, "video.m3u8")
	log.Printf("OutputPath para FFmpeg: %s", outputPath) // Log para verificar

	// Garantir que o diretório de saída existe
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório para transcodificação: %v", err)
	}

	// Construir o comando FFmpeg
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-preset", "fast",
		"-b:v", fmt.Sprintf("%dk", getBandwidth(quality)/1000),
		"-s", getResolutionString(quality),
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "experimental",
		outputPath,
	)

	// Executar o comando e capturar a saída
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao transcodificar o vídeo com FFmpeg: %v, output: %s", err, string(output))
	}

	fmt.Printf("Vídeo transcodificado para: %s\n", outputPath)
	return nil
}


// getBandwidth retorna a largura de banda aproximada com base na qualidade.
func getBandwidth(quality string) int {
	switch quality {
	case "1080p":
		return 5000000
	case "720p":
		return 2500000
	case "480p":
		return 1000000
	default:
		return 500000
	}
}

// getResolutionString retorna a resolução em formato string para o manifesto.
func getResolutionString(quality string) string {
	switch quality {
	case "1080p":
		return "1920x1080"
	case "720p":
		return "1280x720"
	case "480p":
		return "854x480"
	default:
		return "640x360"
	}
}

func uploadHLSFilesToS3(localDir, videoID string, s3Client *storage.S3Client) error {
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Cria o caminho no bucket (exemplo: `videos-transcoded/{videoID}/1080p/video0.ts`)
			relativePath := strings.TrimPrefix(path, localDir+"/")
			s3Key := fmt.Sprintf("videos-transcoded/%s/%s", videoID, relativePath)
			log.Printf("Fazendo upload para o S3 com a chave: %s", s3Key) // Log para verificar

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Faz o upload do arquivo para o S3
			ctx := context.Background()
			err = s3Client.UploadFile(ctx, file, s3Key)
			if err != nil {
				return fmt.Errorf("erro ao fazer upload do arquivo %s: %v", path, err)
			}
			fmt.Printf("Arquivo enviado para o S3: %s\n", s3Key)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro ao iterar pelos arquivos HLS: %v", err)
	}
	return nil
}


