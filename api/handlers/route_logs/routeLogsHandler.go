package routeLogsHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetRoutesLogs(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query(`
	SELECT * FROM ROUTE_LOGS 
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var routesLogs []model.RouteLogs

	for rows.Next() {
		var routeLogs model.RouteLogs
		if err := rows.Scan(&routeLogs.ID, &routeLogs.Lat, &routeLogs.Lng, &routeLogs.DateTime, &routeLogs.VehicleID); err != nil {
			return err
		}
		routesLogs = append(routesLogs, routeLogs)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No routes present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "routes found", "data": routesLogs})
}
