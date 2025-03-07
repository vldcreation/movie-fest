package user

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vldcreation/movie-fest/internal/usecase"
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
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	parseID, err := uuid.Parse(id)
	log.Println("id", id)
	log.Println("err", err)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = usecase.NewVote(s.cfg, s.tokenMaker, s.repo).VoteMovie(newCtx, parseID)
	if err != nil {
		log.Println("err", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
