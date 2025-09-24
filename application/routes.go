package application

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nileshshrs/infinite-storage/config"
	"github.com/nileshshrs/infinite-storage/handler"
	"github.com/nileshshrs/infinite-storage/middlewares"
	"github.com/nileshshrs/infinite-storage/repository"
	"github.com/nileshshrs/infinite-storage/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func loadRoutes(mongoCollection *mongo.Collection, dg *discordgo.Session, cfg *config.Config) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Health check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Collections
	userCollection := mongoCollection.Database().Collection("users")
	sessionCollection := mongoCollection.Database().Collection("sessions")
	fileCollection := mongoCollection.Database().Collection("files")

	// Repositories
	userRepo := repository.NewUserRepository(userCollection)
	sessionRepo := repository.NewSessionRepository(sessionCollection)
	fileRepo := repository.NewFileRepository(fileCollection)

	// Services
	authService := service.NewAuthService(userRepo, sessionRepo)
	userService := service.NewUserService(userRepo)
	uploadService := service.NewUploadService(dg)
	fileService := service.NewFileService(fileRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	uploadHandler := handler.NewUploadHandler(uploadService, fileService, cfg)
	fileHandler := handler.NewFileHandler(fileService)

	// API routes
	router.Route("/api/v1", func(api chi.Router) {
		// Auth routes
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/sign-up", authHandler.Register)
			auth.Post("/sign-in", authHandler.Login)
			auth.Post("/refresh", authHandler.RefreshToken)
		})

		// Protected user routes
		api.Route("/users", func(users chi.Router) {
			users.Use(middlewares.Authenticate)
			users.Get("/", userHandler.GetAllUsers)
		})

		// File routes
		api.Route("/files", func(files chi.Router) {
			files.Use(middlewares.Authenticate) // protect all file routes
			files.Get("/", fileHandler.GetUserFiles)  // get all files for authenticated user
			files.Post("/upload", uploadHandler.HandleUpload) // upload files
		})
	})

	return router
}
