package common

import (
	"context"
	"net/http"
	"time"

	goval "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/internal/usecase"
	"github.com/vldcreation/movie-fest/pkg/responsex"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /health)
func (s *Server) GetHealth(ctx echo.Context) error {
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	return ctx.JSON(http.StatusOK, resp)
}

// (POST /auth/register)
func (s *Server) PostAuthRegister(ctx echo.Context) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()
	// Parse request body
	var params common.PostAuthRegisterJSONRequestBody
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.ApiError{Message: err.Error()})
	}

	// Validate request body
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.ApiError{Message: err.Error()})
	}

	// Register User
	res, err := usecase.NewUser(s.cfg, s.tokenMaker, s.repo).Registration(newCtx, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responsex.ApiError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusCreated, res)
}

// (POST /auth/login)
func (s *Server) PostAuthLogin(ctx echo.Context) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()
	// Parse request body
	var params common.PostAuthLoginJSONRequestBody
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.ApiError{Message: err.Error()})
	}

	// Validate request body
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.ApiError{Message: err.Error()})
	}

	// Login User
	res, err := usecase.NewUser(s.cfg, s.tokenMaker, s.repo).Login(newCtx, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responsex.ApiError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, res)
}
