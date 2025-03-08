package handler

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vldcreation/movie-fest/internal/apis/admin"
	"github.com/vldcreation/movie-fest/internal/apis/common"
	"github.com/vldcreation/movie-fest/internal/apis/user"
	// Import your generated API packages
)

func RegisterSwagger(e *echo.Echo) {
	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.DeepLinking(true),
		echoSwagger.DocExpansion("list"),
		echoSwagger.URL("/swagger/admin/doc.json"),
		echoSwagger.URL("/swagger/user/doc.json"),
		echoSwagger.URL("/swagger/common/doc.json"),
	))
	// Serve OpenAPI specs for each module
	e.GET("/swagger/admin/doc.json", func(c echo.Context) error {
		swagger, err := admin.GetSwagger()
		if err != nil {
			return err
		}
		swagger.OpenAPI = "3.0.0"
		return c.JSON(200, swagger)
	})

	e.GET("/swagger/user/doc.json", func(c echo.Context) error {
		swagger, err := user.GetSwagger()
		if err != nil {
			return err
		}
		swagger.OpenAPI = "3.0.0"
		return c.JSON(200, swagger)
	})

	e.GET("/swagger/common/doc.json", func(c echo.Context) error {
		swagger, err := common.GetSwagger()
		if err != nil {
			return err
		}
		swagger.OpenAPI = "3.0.0"
		return c.JSON(200, swagger)
	})
}
