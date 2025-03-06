package usecase

import (
	"context"
	"time"

	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis"
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

func (u *User) Registration(ctx context.Context, arg apis.UserRegistration) (apis.UserResponse, error) {
	var res apis.UserResponse
	user, roles, err := u.repo.Registration(ctx, apis.UserRegistration{
		Username: arg.Username,
		Email:    arg.Email,
		Password: arg.Password,
	})

	if err != nil {
		return res, err
	}
	res = apis.UserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		Roles: apis.Roles{
			Id:   roles.ID,
			Name: roles.Name,
		},
	}

	return res, nil
}

func (u *User) Login(ctx context.Context, arg apis.UserLogin) (apis.LoginResponse, error) {
	var res apis.LoginResponse
	user, _, err := u.repo.Login(ctx, apis.UserLogin{
		Email:    arg.Email,
		Password: arg.Password,
	})
	if err != nil {
		return res, err
	}

	accessToken, err := u.tokenMaker.CreateToken(user.Username, u.cfg.Token.AccessTokenDuration)
	if err != nil {
		return res, err
	}
	res = apis.LoginResponse{
		Token:     accessToken,
		ExpiresAt: time.Now().Add(u.cfg.Token.AccessTokenDuration),
	}

	return res, nil
}
