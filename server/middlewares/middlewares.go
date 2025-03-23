package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func SetupLoggerMiddleware(app *echo.Echo) {
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))
}

func SetupTimeoutMiddleware(app *echo.Echo) {
	app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Second * 15,
	}))
}
