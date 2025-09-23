package application

import (
	"context"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/nileshshrs/infinite-storage/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router http.Handler
	cfg    *config.Config
}

func New(mongoCollection *mongo.Collection, dg *discordgo.Session, cfg *config.Config) *App {
	return &App{
		router: loadRoutes(mongoCollection, dg, cfg), // pass dg and cfg
		cfg:    cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + a.cfg.Port, // use port from config instead of hardcoding
		Handler: a.router,
	}
	return server.ListenAndServe()
}
