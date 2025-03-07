package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	db "github.com/vldcreation/movie-fest/db/sqlc"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/pkg/util"
)

type User interface {
	Registration(ctx context.Context, arg common.UserRegistration) (db.Users, db.Roles, error)
	Login(ctx context.Context, arg common.UserLogin) (db.Users, db.Roles, error)
}

func (m *Repository) Registration(ctx context.Context, arg common.UserRegistration) (db.Users, db.Roles, error) {
	var result db.Users
	passwordHash, err := util.HashPassword(arg.Password)
	if err != nil {
		return result, db.Roles{}, err
	}

	email, err := arg.Email.MarshalJSON()
	if err != nil {
		return result, db.Roles{}, err
	}
	var emailString string
	err = json.Unmarshal(email, &emailString)
	if err != nil {
		return result, db.Roles{}, err
	}

	err = m.execTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	}, func(tx *db.Queries) error {
		user, err := tx.CreateUser(ctx, db.CreateUserParams{
			Username:     arg.Username,
			Email:        emailString,
			PasswordHash: passwordHash,
		})
		if err != nil {
			return err
		}

		// assign role
		err = tx.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
			UserID: user.ID,
			RoleID: 2, // user
		})
		if err != nil {
			return err
		}
		result = user
		return nil
	})
	if err != nil {
		return result, db.Roles{}, err
	}

	return result, db.Roles{
		ID:   2,
		Name: "user",
	}, nil
}

func (m *Repository) Login(ctx context.Context, arg common.UserLogin) (db.Users, db.Roles, error) {
	var result db.Users

	email, err := arg.Email.MarshalJSON()
	if err != nil {
		return result, db.Roles{}, err
	}
	var emailString string
	err = json.Unmarshal(email, &emailString)
	if err != nil {
		return result, db.Roles{}, err
	}

	user, err := m.querier.GetUserByEmail(ctx, emailString)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, db.Roles{}, errors.New("user not found")
		}
		return result, db.Roles{}, err
	}

	err = util.ComparePassword(user.PasswordHash, arg.Password)
	if err != nil {
		return result, db.Roles{}, err
	}
	role, err := m.querier.GetUserRoles(ctx, user.ID)
	if err != nil {
		return result, db.Roles{}, err
	}
	return user, role, nil
}
