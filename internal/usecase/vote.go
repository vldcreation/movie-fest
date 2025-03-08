package usecase

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/repository"
	"github.com/vldcreation/movie-fest/pkg/token"
)

type Vote struct {
	cfg        *config.Config
	tokenMaker token.Maker
	repo       repository.RepositoryInterface
}

func NewVote(cfg *config.Config, tokenMaker token.Maker, repo repository.RepositoryInterface) *User {
	return &User{
		cfg:        cfg,
		repo:       repo,
		tokenMaker: tokenMaker,
	}
}
