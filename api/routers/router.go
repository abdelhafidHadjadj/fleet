package router

import (
	userRoute "fleet/api/routers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	userRoute.SetupUserRoute(api)
}
