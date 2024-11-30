package routes

import (
	"net/http"
	"streaming-platform/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// SetupRoutes configura todas as rotas da aplicação.
func SetupRoutes(uploadHandler *handlers.UploadHandler, processHandler *handlers.ProcessHandler) http.Handler {
	router := mux.NewRouter()

	// Rota para listar todos os vídeos
	router.HandleFunc("/videos", handlers.ListVideosHandler(processHandler.S3Client)).Methods("GET")
	// Rota para listar resoluções de um vídeo
	router.HandleFunc("/videos/{videoKey}", handlers.ListVideoResolutionsHandler(processHandler.S3Client)).Methods("GET")

	// Configurar CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler

	return corsHandler(router) // Retorna como http.Handler
}
