// This file contains the repository implementation layer.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vldcreation/movie-fest/config"
	db "github.com/vldcreation/movie-fest/db/sqlc"
)

type RepositoryInterface interface {
	Movie
	User
	Vote
}

type Repository struct {
	db      *sql.DB
	querier db.Querier
}

func NewRepository(cfg *config.Config) RepositoryInterface {
	dbSource := "postgresql://" + cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Host + ":" + cfg.DB.Port + "/" + cfg.DB.Database + "?sslmode=disable"
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		panic(err)
	}

	return &Repository{
		db:      conn,
		querier: db.New(conn),
	}
}

func (m *Repository) execTx(ctx context.Context, opts *sql.TxOptions, fn func(*db.Queries) error) error {
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
