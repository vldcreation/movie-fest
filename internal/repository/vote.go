package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/vldcreation/movie-fest/db/sqlc"
)

type Vote interface {
	VoteMovie(ctx context.Context, userId uuid.UUID, movieId uuid.UUID) error
}

func (m *Repository) VoteMovie(ctx context.Context, userId uuid.UUID, movieId uuid.UUID) error {
	err := m.execTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}, func(tx *db.Queries) error {
		err := tx.VoteMovie(ctx, db.VoteMovieParams{
			UserID:  userId,
			MovieID: movieId,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
