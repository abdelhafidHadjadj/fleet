package routeRouter

import (
	routeHandler "fleet/api/handlers/route"

	"github.com/gofiber/fiber/v2"
)

func SetupRouteRouter(router fiber.Router) {
	route := router.Group("/route")
	route.Get("/getall", routeHandler.GetRoutes)
	route.Get("/get/:routeId", routeHandler.GetRouteByID)
	route.Post("/add", routeHandler.CreateRoute)
	route.Patch("/update/:routeId", routeHandler.UpdateRoute)
	route.Delete("/delete/:routeId", routeHandler.DeleteRoute)
	route.Get("/route_number", routeHandler.GetRouteNumber)
	route.Get("/route_number_status/:status", routeHandler.GetRouteNumberByStatus)
	route.Get("/my_routes/:driverId", routeHandler.GetMyRoutes)
}
