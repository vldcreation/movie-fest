package app

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis"
	"github.com/vldcreation/movie-fest/internal/handler"
	"github.com/vldcreation/movie-fest/internal/repository"
	"github.com/vldcreation/movie-fest/pkg/token"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (app *App) Run() {
	e := echo.New()

	var server apis.ServerInterface = newServer(app.cfg)

	apis.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":" + app.cfg.APP.Port))
}

func newServer(cfg *config.Config) *handler.Server {
	repo := repository.NewRepository(cfg)
	pasetoMaker, err := token.NewPasetoMaker(cfg.Token.SecretKey)
	if err != nil {
		panic(err)
	}

	return handler.NewServer(cfg, repo, handler.WithTokenMaker(pasetoMaker))
}
