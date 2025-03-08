package app

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis/admin"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/internal/apis/user"
	"github.com/vldcreation/movie-fest/internal/handler"
	adminHandler "github.com/vldcreation/movie-fest/internal/handler/admin"
	commonHandler "github.com/vldcreation/movie-fest/internal/handler/common"
	userHandler "github.com/vldcreation/movie-fest/internal/handler/user"
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

	pasetoMaker, err := token.NewPasetoMaker(app.cfg.Token.SecretKey)
	if err != nil {
		panic(err)
	}

	adminGroup := e.Group("")
	userGroup := e.Group("")
	commonGroup := e.Group("")
	adminGroup.Use(AdminMiddleware(pasetoMaker))
	userGroup.Use(UserMiddleware(pasetoMaker))

	var adminServer admin.ServerInterface = newAdminServer(app.cfg, pasetoMaker)
	var userServer user.ServerInterface = newUserServer(app.cfg, pasetoMaker)
	var commonServer common.ServerInterface = newCommonServer(app.cfg, pasetoMaker)

	admin.RegisterHandlers(adminGroup, adminServer)
	user.RegisterHandlers(userGroup, userServer)
	common.RegisterHandlers(commonGroup, commonServer)
	e.Use(middleware.Logger())
	handler.RegisterSwagger(e)
	e.Logger.Fatal(e.Start(":" + app.cfg.APP.Port))
}

func newAdminServer(cfg *config.Config, pasetoMaker token.Maker) *adminHandler.Server {
	repo := repository.NewRepository(cfg)

	return adminHandler.NewServer(cfg, repo, adminHandler.WithTokenMaker(pasetoMaker))
}
func newUserServer(cfg *config.Config, pasetoMaker token.Maker) *userHandler.Server {
	repo := repository.NewRepository(cfg)
	return userHandler.NewServer(cfg, repo, userHandler.WithTokenMaker(pasetoMaker))
}
func newCommonServer(cfg *config.Config, pasetoMaker token.Maker) *commonHandler.Server {
	repo := repository.NewRepository(cfg)
	return commonHandler.NewServer(cfg, repo, commonHandler.WithTokenMaker(pasetoMaker))
}
