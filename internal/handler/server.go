package handler

import (
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/repository"
	"github.com/vldcreation/movie-fest/pkg/token"
)

type Server struct {
	cfg        *config.Config
	repo       repository.RepositoryInterface
	tokenMaker token.Maker
}

type ServerOpts func(s *Server)

func WithConfig(cfg *config.Config) ServerOpts {
	return func(s *Server) {
		s.cfg = cfg
	}
}

func WithRepository(repo repository.RepositoryInterface) ServerOpts {
	return func(s *Server) {
		s.repo = repo
	}
}

func WithTokenMaker(tokenMaker token.Maker) ServerOpts {
	return func(s *Server) {
		s.tokenMaker = tokenMaker
	}
}

func NewServer(cfg *config.Config, repo repository.RepositoryInterface, opts ...ServerOpts) *Server {
	s := &Server{
		cfg:  cfg,
		repo: repo,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}
