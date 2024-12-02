package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"streaming-platform/config"
	"streaming-platform/internal/handlers"
	"streaming-platform/internal/storage"
	"streaming-platform/routes"
	"streaming-platform/utils"
)

func main() {
	// Carregar variáveis de ambiente
	config := config.LoadConfig()

	log.Printf("Iniciando servidor com as seguintes configurações básicas...")

	// Criar cliente S3
	s3Client, err := storage.NewS3Client(config.S3Bucket, config.S3Region)
	if err != nil {
		log.Fatalf("Erro ao inicializar cliente S3: %v", err)
	}

	// Configurar handlers
	uploadHandler := handlers.NewUploadHandler()
	processHandler := handlers.NewProcessHandler(s3Client)

	// Configurar rotas
	router := routes.SetupRoutes(uploadHandler, processHandler)

	// Loop de processamento otimizado
	go func() {
		for {
			log.Println("Processando vídeos...")
			if err := utils.ProcessVideos(s3Client, config.Qualities); err != nil {
				log.Printf("Erro no processamento: %v", err)
			}
			time.Sleep(25 * time.Minute) // Intervalo maior para economizar recursos
		}
	}()

	// Configuração da porta pelo Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Padrão
	}
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
