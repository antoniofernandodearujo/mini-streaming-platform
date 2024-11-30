package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"streaming-platform/internal/storage"

	"github.com/gorilla/mux"
)

func GetThumbnailHandler(s3Client *storage.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obter `videoID` dos parâmetros da URL
		vars := mux.Vars(r)
		videoID := vars["videoID"]
		fmt.Println("ID do vídeo recebido:", videoID) // Log para verificar o videoID

		// Montar caminho correto da thumbnail no S3
		ctx := context.Background()
		thumbnailKey := "thumbnails/" + videoID + ".jpg" // Caminho para a pasta thumbnails
		fmt.Println("Procurando arquivo em:", thumbnailKey) // Log para verificar o caminho

		// Baixar thumbnail do S3
		thumbnail, err := s3Client.DownloadFile(ctx, thumbnailKey)
		if err != nil {
			fmt.Println("Erro ao buscar thumbnail:", err.Error()) // Adiciona log de erro completo
			if strings.Contains(err.Error(), "NoSuchKey") {
				http.Error(w, "Thumbnail não encontrada", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar thumbnail: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retornar a thumbnail como resposta
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(thumbnail)
	}
}

