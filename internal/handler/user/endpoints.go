package user

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vldcreation/movie-fest/consts"
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

func (s *Server) PostMoviesIdVote(ctx echo.Context, id string) error {
	newCtx := context.WithValue(ctx.Request().Context(), consts.AuthKey, ctx.Get(consts.AuthKey))
	newCtx, cancel := context.WithTimeout(newCtx, 30*time.Second)
	defer cancel()

	parseID, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.NewApiError(http.StatusBadRequest, responsex.WithErrorMessage("invalid id format")))
	}

	err = usecase.NewVote(s.cfg, s.tokenMaker, s.repo).VoteMovie(newCtx, parseID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, responsex.NewApiError(http.StatusBadRequest, responsex.WithErrorMessage(err.Error())))
	}

	return ctx.JSON(http.StatusOK, responsex.NewApiSuccess(http.StatusOK, responsex.WithSuccessMessage("success")))
}
