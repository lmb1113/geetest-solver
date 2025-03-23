package routes

import (
	"github.com/brianxor/geetest-solver/server/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {
	geetest := app.Group("/geetest")
	geetestV4 := geetest.Group("/v4")
	geetestV4Puzzle := geetestV4.Group("/puzzle")

	geetestV4Puzzle.POST("/solve", handlers.V4PuzzleSolveHandler)
}
