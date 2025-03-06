package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetHello(t *testing.T) {
	s := NewServer(nil, nil)

	t.Run("GetHello", func(t *testing.T) {
		// Create a new Echo instance
		e := echo.New()

		// Create a new request with the desired parameters
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Call the handler function
		err := s.GetHealth(ctx)

		// Assertions
		assert.NoError(t, err)                   // Check that there was no error
		assert.Equal(t, http.StatusOK, rec.Code) // Check that the status code is 200

		// Check the response body
		assert.NoError(t, err)
	})
}
