package app

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vldcreation/movie-fest/consts"
	"github.com/vldcreation/movie-fest/pkg/responsex"
	"github.com/vldcreation/movie-fest/pkg/token"
)

var (
	ErrUnauthorized = responsex.NewApiError(http.StatusUnauthorized)
)

func AdminMiddleware(token token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authorizationHeader := ctx.Request().Header.Get(consts.AuthHeader)
			if len(authorizationHeader) == 0 {
				err := ErrUnauthorized.WithMessage("authorization header is not provided")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := ErrUnauthorized.WithMessage("invalid authorization header format")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			schema := fields[0]
			if schema != consts.BearerSchema {
				err := ErrUnauthorized.WithMessage("unsupported authorization schema")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			accessToken := fields[1]
			payload, err := token.VerifyToken(accessToken)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			// validate scopes
			err = mustAuthenticatedAdmin(*payload)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			ctx.Set(consts.AuthKey, payload)
			return next(ctx)
		}
	}
}

func UserMiddleware(token token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authorizationHeader := ctx.Request().Header.Get(consts.AuthHeader)
			if len(authorizationHeader) == 0 {
				err := ErrUnauthorized.WithMessage("authorization header is not provided")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := ErrUnauthorized.WithMessage("invalid authorization header format")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			schema := fields[0]
			if schema != consts.BearerSchema {
				err := ErrUnauthorized.WithMessage("unsupported authorization schema")
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			accessToken := fields[1]
			payload, err := token.VerifyToken(accessToken)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			if err := payload.Valid(); err != nil {
				ctx.JSON(http.StatusUnauthorized, responsex.NewApiError(http.StatusUnauthorized, responsex.WithErrorMessage(err.Error())))
				return err
			}

			// validate scopes
			err = mustAuthenticatedUser(*payload)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return err
			}

			ctx.Set(consts.AuthKey, payload)
			return next(ctx)
		}
	}
}

func mustAuthenticatedAdmin(authPayload token.Payload) error {
	role, ok := authPayload.GetCustomClaims("role")
	if !ok {
		return ErrUnauthorized.WithMessage("role not found")
	}
	if role != "admin" {
		return ErrUnauthorized.WithMessage("unauthorized")
	}
	return nil
}

func mustAuthenticatedUser(authPayload token.Payload) error {
	role, ok := authPayload.GetCustomClaims("role")
	if !ok {
		return ErrUnauthorized.WithMessage("role not found")
	}
	if role != "admin" && role != "user" {
		return ErrUnauthorized.WithMessage("unauthorized")
	}
	return nil
}
