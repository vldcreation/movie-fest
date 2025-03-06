package usecase

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/vldcreation/movie-fest/config"
	"github.com/vldcreation/movie-fest/internal/apis"
	"github.com/vldcreation/movie-fest/internal/repository"
)

type Movie struct {
	cfg  *config.Config
	repo repository.RepositoryInterface
}

func NewMovie(cfg *config.Config, repo repository.RepositoryInterface) *Movie {
	return &Movie{
		cfg:  cfg,
		repo: repo,
	}
}

func (m *Movie) CreateMovie(ctx context.Context, param apis.MovieCreateRequest) (*apis.MovieResponse, error) {
	var res apis.MovieResponse
	result, err := m.repo.CreateMovie(ctx, param)
	if err != nil {
		return nil, err
	}

	log.Printf("result: %+v\n", result)
	if result.Movie.ID != uuid.Nil {
		res.Id = result.Movie.ID.String()
		res.Title = result.Movie.Title
		res.Description = result.Movie.Description.String
		res.Duration = int(result.Movie.Duration)
		res.WatchUrl = result.Movie.WatchUrl
		res.CreatedAt = &result.Movie.CreatedAt.Time
		res.UpdatedAt = &result.Movie.UpdatedAt.Time

		if result.Genres != nil {
			for _, genre := range result.Genres {
				res.Genres = append(res.Genres, genre.Name)
			}
		}

		if result.Artits != nil {
			for _, artist := range result.Artits {
				res.Artists = append(res.Artists, artist.Name)
			}
		}

		res.Views = len(result.Views)
		res.Votes = len(result.Votes)
	}

	return &res, nil

}

func (m *Movie) UpdateMovie(ctx context.Context, id uuid.UUID, param apis.MovieUpdateRequest) (*apis.MovieResponse, error) {
	var res apis.MovieResponse
	result, err := m.repo.UpdateMovie(ctx, id, param)
	if err != nil {
		return nil, err
	}

	log.Printf("result: %+v\n", result)
	if result.Movie.ID != uuid.Nil {
		res.Id = result.Movie.ID.String()
		res.Title = result.Movie.Title
		res.Description = result.Movie.Description.String
		res.Duration = int(result.Movie.Duration)
		res.WatchUrl = result.Movie.WatchUrl
		res.CreatedAt = &result.Movie.CreatedAt.Time
		res.UpdatedAt = &result.Movie.UpdatedAt.Time

		if result.Genres != nil {
			for _, genre := range result.Genres {
				res.Genres = append(res.Genres, genre.Name)
			}
		}

		if result.Artits != nil {
			for _, artist := range result.Artits {
				res.Artists = append(res.Artists, artist.Name)
			}
		}

		res.Views = len(result.Views)
		res.Votes = len(result.Votes)
	}

	return &res, nil
}

func (m *Movie) GetAdminMoviesMostViewed(ctx context.Context, params apis.GetAdminMoviesMostViewedParams) (*apis.PaginatedMovieViewsResponse, error) {
	movies, err := m.repo.GetMostViewedMovie(ctx, params)
	if err != nil {
		return nil, err
	}

	response := &apis.PaginatedMovieViewsResponse{
		Page:  int(params.Page),
		Limit: int(params.Limit),
		Data:  make([]apis.MovieViewsResponse, 0),
	}

	for _, movie := range movies {
		movieResp := apis.MovieViewsResponse{
			Movie: apis.MovieResponse{
				Id:          movie.ID.String(),
				Title:       movie.Title,
				Description: movie.Description.String,
				Duration:    int(movie.Duration),
				WatchUrl:    movie.WatchUrl,
				CreatedAt:   &movie.CreatedAt.Time,
				UpdatedAt:   &movie.UpdatedAt.Time,
			},
			Views: movie.ViewCount,
		}
		response.Data = append(response.Data, movieResp)
	}

	response.Total = len(movies)
	response.TotalPages = (response.Total + int(params.Limit) - 1) / int(params.Limit)

	return response, nil
}

func (m *Movie) GetAdminMoviesMostViewedGenre(ctx context.Context, params apis.GetAdminMoviesMostViewedGenresParams) (*apis.PaginatedGenreViewResponse, error) {
	genres, err := m.repo.GetMostViewedMovieGenre(ctx, params)
	if err != nil {
		return nil, err
	}

	response := &apis.PaginatedGenreViewResponse{
		Page:  int(params.Page),
		Limit: int(params.Limit),
		Data:  make([]apis.GenreViewsResponse, 0),
	}

	for _, genre := range genres {
		genreResp := apis.GenreViewsResponse{
			Genre: genre.Name,
			Views: genre.ViewCount,
		}
		response.Data = append(response.Data, genreResp)
	}

	response.Total = len(genres)
	response.TotalPages = (response.Total + int(params.Limit) - 1) / int(params.Limit)

	return response, nil
}
