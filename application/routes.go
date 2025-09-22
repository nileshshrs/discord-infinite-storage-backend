package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nileshshrs/infinite-storage/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api/v1", func(api chi.Router) {
		api.Route("/auth", loadAuthRoutes)
	})
	return router
}

func loadAuthRoutes(router chi.Router) {
	orderHandler := &handler.User{}

	router.Post("/sign-up", orderHandler.Register)
	router.Post("/sign-in", orderHandler.Login)
	router.Patch("/update-user/{id}", orderHandler.Update)

}
