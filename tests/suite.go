package tests

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/consts"
	"github.com/vldcreation/movie-fest/pkg/token"
)

type TestSuite struct {
	Cfg        *config.Config
	TokenMaker token.Maker
}

func NewTestSuite() *TestSuite {
	cfg := config.NewConfigFromYaml("../" + consts.ConfigPath)
	tokenMaker, err := token.NewPasetoMaker(cfg.Token.SecretKey)
	if err != nil {
		panic(err)
	}
	return &TestSuite{
		Cfg:        cfg,
		TokenMaker: tokenMaker,
	}
}
