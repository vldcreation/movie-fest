package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/vldcreation/movie-fest/db/sqlc"
	"github.com/vldcreation/movie-fest/internal/apis/admin"
	"github.com/vldcreation/movie-fest/internal/apis/user"
	"github.com/vldcreation/movie-fest/pkg/util"
)

type CreateMovieTxResult struct {
	Movie  db.Movies
	Genres []db.Genres
	Artits []db.Artists
	Votes  []db.Votes
	Views  []db.Views
}

type UpdateMovieTxResult struct {
	Movie  db.Movies
	Genres []db.Genres
	Artits []db.Artists
	Votes  []db.Votes
	Views  []db.Views
}

type Movie interface {
	CreateMovie(ctx context.Context, arg admin.MovieCreateRequest) (CreateMovieTxResult, error)
	UpdateMovie(ctx context.Context, id uuid.UUID, arg admin.MovieUpdateRequest) (UpdateMovieTxResult, error)
	GetMostViewedMovie(ctx context.Context, arg admin.GetAdminMoviesMostViewedParams) ([]db.GetMostViewedMoviesRow, error)
	GetMostViewedMovieGenre(ctx context.Context, arg admin.GetAdminMoviesMostViewedGenresParams) ([]db.GetMostViewedGenresRow, error)
	GetMovies(ctx context.Context, arg user.GetMoviesParams) ([]db.GetMoviesRow, error)
}

func (m *Repository) CreateMovie(ctx context.Context, arg admin.MovieCreateRequest) (CreateMovieTxResult, error) {
	var result CreateMovieTxResult

	err := m.execTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}, func(tx *db.Queries) error {
		movie, err := tx.CreateMovie(ctx, db.CreateMovieParams{
			Title:       arg.Title,
			Description: util.ToSQLNullableString(arg.Description),
			Duration:    int32(arg.Duration),
			WatchUrl:    arg.WatchUrl,
		})

		if err != nil {
			return err
		}

		result.Movie = movie

		for _, genre := range arg.Genres {
			resultGenre := db.Genres{}
			genreID, err := uuid.Parse(genre)
			if err != nil || genreID == uuid.Nil {
				genre, err := tx.CreateGenre(ctx, genre)

				if err != nil {
					return err
				}

				resultGenre = genre
			} else {
				genre, err := tx.GetGenreByID(ctx, genreID)

				if err != nil {
					return err
				}

				resultGenre = genre
			}

			// add genre to movie
			err = tx.AddGenreToMovie(ctx, db.AddGenreToMovieParams{
				MovieID: movie.ID,
				GenreID: resultGenre.ID,
			})

			if err != nil {
				return err
			}

			result.Genres = append(result.Genres, resultGenre)
		}

		for _, artist := range arg.Artists {
			resultArtist := db.Artists{}
			artistID, err := uuid.Parse(artist)

			if err != nil || artistID == uuid.Nil {
				artist, err := tx.CreateArtist(ctx, artist)

				if err != nil {
					return err
				}

				resultArtist = artist
			} else {
				artist, err := tx.GetArtistByID(ctx, artistID)

				if err != nil {
					return err
				}

				resultArtist = artist
			}

			// add artist to movie
			err = tx.AddArtistToMovie(ctx, db.AddArtistToMovieParams{
				MovieID:  movie.ID,
				ArtistID: resultArtist.ID,
			})

			if err != nil {
				return err
			}

			result.Artits = append(result.Artits, resultArtist)
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *Repository) UpdateMovie(ctx context.Context, id uuid.UUID, arg admin.MovieUpdateRequest) (UpdateMovieTxResult, error) {
	var result UpdateMovieTxResult

	// check if movie exists
	existingMovie, err := m.querier.GetMovieByID(ctx, id)
	if err != nil {
		return result, err
	}
	if existingMovie.ID == uuid.Nil {
		return result, errors.New("movie not found")
	}

	err = m.execTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}, func(tx *db.Queries) error {
		// Prepare update params with existing values
		updateParams := db.UpdateMovieParams{
			ID:          id,
			Title:       existingMovie.Title,
			Description: existingMovie.Description,
			Duration:    existingMovie.Duration,
			WatchUrl:    existingMovie.WatchUrl,
		}

		// Only update fields that are provided
		if arg.Title != "" {
			updateParams.Title = arg.Title
		}
		if arg.Description != "" {
			updateParams.Description = util.ToSQLNullableString(arg.Description)
		}
		if arg.Duration > 0 {
			updateParams.Duration = int32(arg.Duration)
		}
		if arg.WatchUrl != "" {
			updateParams.WatchUrl = arg.WatchUrl
		}

		movie, err := tx.UpdateMovie(ctx, updateParams)
		if err != nil {
			return err
		}
		result.Movie = movie

		// Only update genres if provided
		if len(arg.Genres) > 0 {
			// First, remove existing genres
			err = tx.RemoveAllGenresFromMovie(ctx, movie.ID)
			if err != nil {
				return err
			}

			// Then add new genres
			for _, genre := range arg.Genres {
				resultGenre := db.Genres{}
				genreID, err := uuid.Parse(genre)
				if err != nil || genreID == uuid.Nil {
					genre, err := tx.UpsertGenre(ctx, genre)
					if err != nil {
						return err
					}
					resultGenre = genre
				} else {
					genre, err := tx.GetGenreByID(ctx, genreID)
					if err != nil {
						return err
					}
					resultGenre = genre
				}

				err = tx.AddGenreToMovie(ctx, db.AddGenreToMovieParams{
					MovieID: movie.ID,
					GenreID: resultGenre.ID,
				})
				if err != nil {
					return err
				}
				result.Genres = append(result.Genres, resultGenre)
			}
		}

		// Only update artists if provided
		if len(arg.Artists) > 0 {
			// First, remove existing artists
			err = tx.RemoveAllArtistsFromMovie(ctx, movie.ID)
			if err != nil {
				return err
			}

			// Then add new artists
			for _, artist := range arg.Artists {
				resultArtist := db.Artists{}
				artistID, err := uuid.Parse(artist)
				if err != nil || artistID == uuid.Nil {
					artist, err := tx.UpsertArtist(ctx, artist)
					if err != nil {
						return err
					}
					resultArtist = artist
				} else {
					artist, err := tx.GetArtistByID(ctx, artistID)
					if err != nil {
						return err
					}
					resultArtist = artist
				}

				err = tx.AddArtistToMovie(ctx, db.AddArtistToMovieParams{
					MovieID:  movie.ID,
					ArtistID: resultArtist.ID,
				})
				if err != nil {
					return err
				}
				result.Artits = append(result.Artits, resultArtist)
			}
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *Repository) GetMostViewedMovie(ctx context.Context, arg admin.GetAdminMoviesMostViewedParams) ([]db.GetMostViewedMoviesRow, error) {
	offset := (arg.Page - 1) * arg.Limit
	movies, err := m.querier.GetMostViewedMovies(ctx, db.GetMostViewedMoviesParams{
		Limit:  arg.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *Repository) GetMostViewedMovieGenre(ctx context.Context, arg admin.GetAdminMoviesMostViewedGenresParams) ([]db.GetMostViewedGenresRow, error) {
	offset := (arg.Page - 1) * arg.Limit
	movies, err := m.querier.GetMostViewedGenres(ctx, db.GetMostViewedGenresParams{
		Limit:  int32(arg.Limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *Repository) GetMovies(ctx context.Context, arg user.GetMoviesParams) ([]db.GetMoviesRow, error) {
	offset := (arg.Page - 1) * arg.Limit
	movies, err := m.querier.GetMovies(ctx, db.GetMoviesParams{
		Limit:   arg.Limit,
		Offset:  offset,
		Column1: arg.Search,
	})
	if err != nil {
		return nil, err
	}
	return movies, nil
}
