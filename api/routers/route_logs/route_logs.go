package routeLogsRouter

import (
	routeLogsHandler "fleet/api/handlers/route_logs"

	"github.com/gofiber/fiber/v2"
)

func SetupRouteLogsRouter(router fiber.Router) {
	route := router.Group("/route_logs")
	route.Get("/getall", routeLogsHandler.GetRoutesLogs)

}
