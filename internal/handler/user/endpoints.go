package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
