package admin

import (
	"context"
	"log"
	"net/http"
	"time"

	goval "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vldcreation/movie-fest/internal/apis/admin"
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

// (POST /movies)
func (s *Server) PostAdminMovies(ctx echo.Context) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	// Parse request body
	var params admin.MovieCreateRequest
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Validate request body
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Post Movie
	res, err := usecase.NewMovie(s.cfg, s.repo).CreateMovie(newCtx, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, admin.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, res)
}

// (PUT /movies/:id)
func (s *Server) PutAdminMoviesId(ctx echo.Context, id string) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	// Parse id
	parseID, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Parse request body
	var params admin.MovieUpdateRequest

	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Validate request body
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}

	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Update Movie
	res, err := usecase.NewMovie(s.cfg, s.repo).UpdateMovie(newCtx, parseID, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, admin.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

// (GET /movies/most-viewed)
func (s *Server) GetAdminMoviesMostViewed(ctx echo.Context, params admin.GetAdminMoviesMostViewedParams) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()
	// Parse request query
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Validate request query
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Get Most Viewed Movie
	res, err := usecase.NewMovie(s.cfg, s.repo).GetAdminMoviesMostViewed(newCtx, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, admin.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

// (GET /movies/most-viewed/genre)
func (s *Server) GetAdminMoviesMostViewedGenres(ctx echo.Context, params admin.GetAdminMoviesMostViewedGenresParams) error {
	newCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	scopes, _ := ctx.Get(admin.BearerAuthScopes).([]string)
	log.Printf("scopes: %v", scopes)
	// Parse request query
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}
	// Validate request query
	ctx.Echo().Validator = &Validator{validator: goval.New(goval.WithRequiredStructEnabled())}
	if err := ctx.Validate(params); err != nil {
		return ctx.JSON(http.StatusBadRequest, admin.ErrorResponse{Message: err.Error()})
	}

	// Get Most Viewed Movie Genre
	res, err := usecase.NewMovie(s.cfg, s.repo).GetAdminMoviesMostViewedGenre(newCtx, params)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, admin.ErrorResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, res)
}
