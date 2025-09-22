package application

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router http.Handler
}

func New(mongoCollection *mongo.Collection) *App {
	return &App{
		router: loadRoutes(mongoCollection),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":6278",
		Handler: a.router,
	}
	return server.ListenAndServe()
}
