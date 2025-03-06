package main

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/consts"
	"github.com/vldcreation/movie-fest/internal/app"
)

func main() {
	cfg := config.NewConfigFromYaml(consts.ConfigPath)
	app.NewApp(cfg).Run()
}
