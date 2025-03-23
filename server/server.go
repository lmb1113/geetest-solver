package server

import (
	"fmt"
	"github.com/brianxor/geetest-solver/server/middlewares"
	"github.com/brianxor/geetest-solver/server/routes"
	"github.com/labstack/echo/v4"
)

func Start(serverHost string, serverPort string) error {
	app := echo.New()

	middlewares.SetupLoggerMiddleware(app)
	middlewares.SetupTimeoutMiddleware(app)

	routes.SetupRoutes(app)

	serverAddress := fmt.Sprintf("%s:%s", serverHost, serverPort)

	if err := app.Start(serverAddress); err != nil {
		return err
	}

	return nil
}
