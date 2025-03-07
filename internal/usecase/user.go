package usecase

import (
	"context"
	"time"

	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/internal/repository"
	"github.com/vldcreation/movie-fest/pkg/token"
)

type User struct {
	cfg        *config.Config
	tokenMaker token.Maker
	repo       repository.RepositoryInterface
}

func NewUser(cfg *config.Config, tokenMaker token.Maker, repo repository.RepositoryInterface) *User {
	return &User{
		cfg:        cfg,
		repo:       repo,
		tokenMaker: tokenMaker,
	}
}

func (u *User) Registration(ctx context.Context, arg common.UserRegistration) (common.UserResponse, error) {
	var res common.UserResponse
	user, roles, err := u.repo.Registration(ctx, common.UserRegistration{
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
	})

	if err != nil {
		return res, err
	}
	res = common.UserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		Roles: common.Roles{
			Id:   roles.ID,
			Name: roles.Name,
		},
	}

	return res, nil
}

func (u *User) Login(ctx context.Context, arg common.UserLogin) (common.LoginResponse, error) {
	var res common.LoginResponse
	user, role, err := u.repo.Login(ctx, common.UserLogin{
		Email:    arg.Email,
		Password: arg.Password,
	})
	if err != nil {
		return res, err
	}

	metadata := map[string]any{
		"role_id": role.ID,
		"role":    role.Name,
		"user_id": user.ID,
	}

	accessToken, err := u.tokenMaker.CreateToken(user.Username, u.cfg.Token.AccessTokenDuration, metadata)
	if err != nil {
		return res, err
	}
	res = common.LoginResponse{
		Token:     accessToken,
		ExpiresAt: time.Now().Add(u.cfg.Token.AccessTokenDuration),
	}

	return res, nil
}
