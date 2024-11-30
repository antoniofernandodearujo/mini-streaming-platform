package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"streaming-platform/internal/storage"

	"github.com/gorilla/mux"
)

// ListVideosHandler retorna a lista de todos os vídeos disponíveis no bucket
func ListVideosHandler(s3Client *storage.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		prefix := "videos-transcoded/"

		// Listar diretórios (IDs dos vídeos)
		directories, err := s3Client.ListDirectories(ctx, prefix)
		if err != nil {
			http.Error(w, "Erro ao listar vídeos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(directories) == 0 {
			http.Error(w, "Nenhum vídeo encontrado", http.StatusNotFound)
			return
		}

		// Retornar os IDs dos vídeos como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(directories)
	}
}

// ListVideoResolutionsHandler retorna as resoluções disponíveis para um vídeo específico
func ListVideoResolutionsHandler(s3Client *storage.S3Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := context.Background()
        vars := mux.Vars(r)
        videoID := vars["videoKey"]

        // Listar resoluções disponíveis
        resolutions, err := s3Client.ListResolutions(ctx, videoID)
        if err != nil {
            http.Error(w, "Erro ao listar resoluções: "+err.Error(), http.StatusInternalServerError)
            return
        }

        if len(resolutions) == 0 {
            http.Error(w, "Nenhuma resolução encontrada para este vídeo", http.StatusNotFound)
            return
        }

        // Construir a resposta com resoluções e arquivos
        result := make(map[string][]string)

        for _, resolution := range resolutions {
            // Listar os arquivos para a resolução atual
            resolutionPrefix := fmt.Sprintf("videos-transcoded/%s/%s/", videoID, resolution)
            files, err := s3Client.ListFiles(ctx, resolutionPrefix)
            if err != nil {
                http.Error(w, fmt.Sprintf("Erro ao listar arquivos para a resolução %s: %v", resolution, err), http.StatusInternalServerError)
                return
            }

            // Gerar URLs públicas para os arquivos
            var fileURLs []string
            for _, file := range files {
                fileURL, err := s3Client.GetFileURL(file)  // Gerar URL para cada arquivo
                if err != nil {
                    http.Error(w, "Erro ao gerar URL do arquivo: "+err.Error(), http.StatusInternalServerError)
                    return
                }
                fileURLs = append(fileURLs, fileURL)
            }

            // Salvar a resolução e os arquivos com as URLs na resposta
            result[resolution] = fileURLs
        }

        // Retornar a resposta em JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "videoID":     videoID,
            "resolutions": result,
        })
    }
}


func GetVideoByResolutionHandler(s3Client *storage.S3Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		vars := mux.Vars(r)

		// Extrair os parâmetros da URL
		videoID := vars["videoKey"]
		chosenResolution := vars["resolution"]

		// Validar parâmetros
		if videoID == "" || chosenResolution == "" {
			http.Error(w, "Parâmetros 'videoKey' e 'resolution' são obrigatórios. Exemplo: /videos/{videoKey}/{resolution}", http.StatusBadRequest)
			return
		}

		// Construir o caminho completo do arquivo `video.m3u8`
		filePath := "videos-transcoded/" + videoID + "/" + chosenResolution + "/video.m3u8"

		// Usar a função DownloadFile para obter o conteúdo do arquivo
		videoFile, err := s3Client.DownloadFile(ctx, filePath)
		if err != nil {
			http.Error(w, "Erro ao obter o arquivo de vídeo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Configurar cabeçalhos para resposta com o arquivo
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl") // Tipo MIME para HLS
		w.Header().Set("Content-Disposition", "inline; filename=video.m3u8") // Exibir no navegador

		// Enviar o conteúdo do arquivo como resposta
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(videoFile)
		if err != nil {
			http.Error(w, "Erro ao enviar o arquivo de vídeo: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

