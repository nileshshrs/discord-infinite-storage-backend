package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nileshshrs/infinite-storage/handler"
	"github.com/nileshshrs/infinite-storage/middlewares"
	"github.com/nileshshrs/infinite-storage/repository"
	"github.com/nileshshrs/infinite-storage/service"

	"go.mongodb.org/mongo-driver/mongo"
)

// loadRoutes sets up all routes for the application
func loadRoutes(mongoCollection *mongo.Collection) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// health check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// separate collections
	userCollection := mongoCollection.Database().Collection("users")
	sessionCollection := mongoCollection.Database().Collection("sessions")

	// repositories
	userRepo := repository.NewUserRepository(userCollection)
	sessionRepo := repository.NewSessionRepository(sessionCollection)

	// services
	authService := service.NewAuthService(userRepo, sessionRepo)
	userService := service.NewUserService(userRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

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
			users.Use(middlewares.Authenticate) // protect these routes
			users.Get("/", userHandler.GetAllUsers)
		})
	})

	return router
}
