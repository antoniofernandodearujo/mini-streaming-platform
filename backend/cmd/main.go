package main

import (
	"fmt"
	"log"
	"net/http"
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

	// Criar o cliente S3
	s3Client, err := storage.NewS3Client(config.S3Bucket, config.S3Region)
	if err != nil {
		log.Fatalf("Erro ao inicializar cliente S3: %v", err)
	}

	// Configurar handlers específicos
	uploadHandler := handlers.NewUploadHandler(s3Client) // Exemplo de criação de handler
	processHandler := handlers.NewProcessHandler(s3Client)

	// Configurar rotas
	router := routes.SetupRoutes(uploadHandler, processHandler)

	// Iniciar o loop de processamento de vídeos em goroutine separada
	go func() {
		for {
			fmt.Println("Iniciando processamento de vídeos...")
			err := utils.ProcessVideos(s3Client, config.Qualities)
			if err != nil {
				fmt.Printf("Erro durante processamento: %v\n", err)
			}

			fmt.Println("Aguardando para reiniciar processamento...")
			time.Sleep(10 * time.Minute)
		}
	}()

	// Iniciar o servidor HTTP
	port := "8080"
	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

