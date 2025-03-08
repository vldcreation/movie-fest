package usecase

import (
	"context"
	"time"

	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/internal/apis/user"
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

func (m *Movie) GetMovies(ctx context.Context, params user.GetMoviesParams) (*user.PaginatedMoviesResponse, error) {
	movies, err := m.repo.GetMovies(ctx, params)
	if err != nil {
		return nil, err
	}
	response := &user.PaginatedMoviesResponse{
		Page:  int(params.Page),
		Limit: int(params.Limit),
		Data:  make([]user.MovieResponse, 0),
	}
	for _, movie := range movies {
		movieResp := user.MovieResponse{
			Id:          movie.ID.String(),
			Title:       movie.Title,
			Description: movie.Description.String,
			Duration:    int(movie.Duration),
			WatchUrl:    movie.WatchUrl,
			Votes:       int(movie.VoteCount),
			Views:       int(movie.ViewCount),
			CreatedAt:   &movie.CreatedAt.Time,
			UpdatedAt:   &movie.UpdatedAt.Time,
		}
		response.Data = append(response.Data, movieResp)
	}
	response.Total = len(movies)
	response.TotalPages = (response.Total + int(params.Limit) - 1) / int(params.Limit)
	return response, nil
}
