package router

import (
	driverRouter "fleet/api/routers/driver"
	routeRouter "fleet/api/routers/route"
	userRoute "fleet/api/routers/user"
	vehicleRouter "fleet/api/routers/vehicle"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	userRoute.SetupUserRoute(api)
	vehicleRouter.SetupVehicleRoute(api)
	routeRouter.SetupRouteRouter(api)
	driverRouter.SetupDriverRoute(api)
}
