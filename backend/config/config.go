package config

import (
	"log"
	"os"
	"path/filepath"
	"streaming-platform/utils"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	StoragePath  string
	VideoBaseDir string
	HLSBaseDir   string
	Qualities    []string
	S3Bucket     string
	S3Region     string
}

func LoadConfig() Config {
	// Tenta carregar o .env apenas se existir
	_ = godotenv.Load(".env")

	// Pega as variáveis de ambiente ou usa valores padrão
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "/app/videos" // Ajustado para funcionar no Railway
	}

	// Criar diretórios necessários
	utils.EnsureDirectoryExists(filepath.Join(storagePath, "hls"))
	utils.EnsureDirectoryExists(filepath.Join(storagePath, "videos"))

	videoBaseDir := os.Getenv("VIDEO_BASE_DIR")
	if videoBaseDir == "" {
		videoBaseDir = "videos"
	}

	hlsBaseDir := os.Getenv("HLS_BASE_DIR")
	if hlsBaseDir == "" {
		hlsBaseDir = "hls"
	}

	qualities := os.Getenv("VIDEO_QUALITIES")
	if qualities == "" {
		qualities = "1080p,720p,480p"
	}

	s3Bucket := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("AWS_REGION")

	log.Printf("Configuração carregada: StoragePath=%s, VideoBaseDir=%s, HLSBaseDir=%s, Qualities=%s",
		storagePath, videoBaseDir, hlsBaseDir, qualities)

	return Config{
		StoragePath:  storagePath,
		VideoBaseDir: videoBaseDir,
		HLSBaseDir:   hlsBaseDir,
		Qualities:    strings.Split(qualities, ","),
		S3Bucket:     s3Bucket,
		S3Region:     s3Region,
	}
}
