package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nileshshrs/infinite-storage/handler"
	"github.com/nileshshrs/infinite-storage/repository"
	"github.com/nileshshrs/infinite-storage/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func loadRoutes(mongoCollection *mongo.Collection) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// separate collections
	userCollection := mongoCollection.Database().Collection("users")
	sessionCollection := mongoCollection.Database().Collection("sessions")

	// repositories
	userRepo := repository.NewUserRepository(userCollection)
	sessionRepo := repository.NewSessionRepository(sessionCollection)

	// service & handler
	authService := service.NewAuthService(userRepo, sessionRepo)
	authHandler := handler.NewAuthHandler(authService)

	// routes
	router.Route("/api/v1", func(api chi.Router) {
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/sign-up", authHandler.Register)
			auth.Post("/sign-in", authHandler.Login)
			auth.Post("/refresh", authHandler.RefreshToken)
		})
	})

	return router
}
