package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func TranscodeVideoToHLS(videoID, inputPath string, qualities []string) error {
	tempDir := filepath.Join(os.TempDir(), "videos", videoID)
	log.Printf("Diretório temporário para transcodificação: %s", tempDir)

	// Criar diretório base
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	masterPlaylistPath := filepath.Join(tempDir, "master.m3u8")
	masterFile, err := os.Create(masterPlaylistPath)
	if err != nil {
		return fmt.Errorf("erro ao criar master playlist: %v", err)
	}
	defer masterFile.Close()
	masterFile.WriteString("#EXTM3U\n")

	for _, quality := range qualities {
		qualityDir := filepath.Join(tempDir, quality)
		if err := os.MkdirAll(qualityDir, os.ModePerm); err != nil {
			return fmt.Errorf("erro ao criar diretório da qualidade %s: %v", quality, err)
		}

		outputPath := filepath.Join(qualityDir, "playlist.m3u8")
		cmd := exec.Command("ffmpeg",
			"-i", inputPath,
			"-preset", "veryfast",
			"-b:v", fmt.Sprintf("%dk", getBandwidth(quality)/1000),
			"-s", getResolutionString(quality),
			"-c:v", "libx264",
			"-hls_time", "10",
			"-hls_playlist_type", "vod",
			outputPath,
		)
		log.Printf("Executando FFmpeg para qualidade %s", quality)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao transcodificar %s: %v", quality, err)
		}
		masterFile.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s\n",
			getBandwidth(quality), getResolutionString(quality), filepath.Join(quality, "playlist.m3u8")))
	}

	log.Printf("Transcodificação para HLS concluída: %s", masterPlaylistPath)
	return nil
}

func getBandwidth(quality string) int {
	switch quality {
	case "1080p":
		return 3000000
	case "720p":
		return 1500000
	case "480p":
		return 800000
	default:
		return 400000
	}
}

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
