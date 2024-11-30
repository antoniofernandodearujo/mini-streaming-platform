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

// Função para garantir que um diretório existe
func ensureDirectoryExists(path string) {
	// Diagnóstico: Log do diretório a ser verificado/criado
	log.Printf("Verificando/criando diretório: %s", utils.RemoveExtensionID(path))

	// Verificar se o diretório já existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Tentativa de criar o diretório
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Erro ao criar o diretório '%s': %v", path, err)
		}
		log.Printf("Diretório criado: %s", path)
	} else {
		log.Printf("O diretório '%s' já existe", path)
	}
}

func LoadConfig() Config {
	// Carrega o arquivo .env se existir
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Erro ao carregar o arquivo .env: ", err)
	}

	// Obtém o diretório atual de trabalho
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter o diretório atual: %v", err)
	}
	log.Printf("Diretório atual de trabalho: %s", currentDir)

	// Caminho de armazenamento principal
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = filepath.Join(currentDir, "videos", "hls") // Caminho relativo ao diretório atual
	}

	// Garantir que o diretório principal exista
	ensureDirectoryExists(storagePath)

	// Diretório para os vídeos
	videoBaseDir := os.Getenv("VIDEO_BASE_DIR")
	if videoBaseDir == "" {
		videoBaseDir = "videos"
	}
	videoFullPath := filepath.Join(storagePath, videoBaseDir)
	ensureDirectoryExists(videoFullPath)

	// Diretório para HLS
	hlsBaseDir := os.Getenv("HLS_BASE_DIR")
	if hlsBaseDir == "" {
		hlsBaseDir = "hls"
	}
	hlsFullPath := filepath.Join(storagePath, hlsBaseDir)
	ensureDirectoryExists(hlsFullPath)

	// Qualidades de vídeo
	qualities := os.Getenv("VIDEO_QUALITIES")
	if qualities == "" {
		qualities = "1080p,720p,480p" // Qualidades padrão
	}

	// Bucket e região S3
	s3Bucket := os.Getenv("S3_BUCKET_NAME")
	s3Region := os.Getenv("AWS_REGION")

	log.Printf("Configuração carregada: StoragePath=%s, VideoBaseDir=%s, HLSBaseDir=%s, Qualities=%s",
		storagePath, videoBaseDir, hlsBaseDir, qualities)

	return Config{
		StoragePath:  storagePath,
		VideoBaseDir: videoFullPath,
		HLSBaseDir:   hlsFullPath,
		Qualities:    strings.Split(qualities, ","),
		S3Bucket:     s3Bucket,
		S3Region:     s3Region,
	}
}
