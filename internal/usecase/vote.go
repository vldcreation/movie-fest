package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/consts"
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

func (u *User) VoteMovie(ctx context.Context, id uuid.UUID) error {
	user, ok := ctx.Value(consts.AuthKey).(*token.Payload)
	if !ok {
		return errors.New("user not found")
	}

	userId, ok := user.GetCustomClaims("user_id")
	if !ok {
		return errors.New("user not found")
	}

	parseUserId, err := uuid.Parse(userId.(string))
	if err != nil {
		return err
	}

	return u.repo.VoteMovie(ctx, id, parseUserId)
}
