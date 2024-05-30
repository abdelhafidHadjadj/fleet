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
	SELECT ROUTE_LOGS.id ,ROUTE_LOGS.lat, ROUTE_LOGS.lng, ROUTE_LOGS.datetime, ROUTE_LOGS.vehicle_id, VEHICLE.register_number, VEHICLE.name, VEHICLE.model
	FROM ROUTE_LOGS
	INNER JOIN VEHICLE ON ROUTE_LOGS.vehicle_id = VEHICLE.id 
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var routesLogs []model.RouteLogs

	for rows.Next() {
		var routeLogs model.RouteLogs
		var vehicle model.Vehicle
		if err := rows.Scan(&routeLogs.ID, &routeLogs.Lat, &routeLogs.Lng, &routeLogs.DateTime, &routeLogs.VehicleID, &vehicle.Register_number, &vehicle.Name, &vehicle.Model); err != nil {
			return err
		}
		routesLogs = append(routesLogs, routeLogs)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No routes present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "routes found", "data": routesLogs})
}

func GetRoutesLogsByVehID(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	id := c.Params("vehicle_id")
	rows, err := db.Query(`
	SELECT ROUTE_LOGS.id ROUTE_LOGS.lat, ROUTE_LOGS.lng, ROUTE_LOGS.datetime, ROUTE_LOGS.vehicle_id, VEHICLE.register_number, VEHICLE.name, VEHICLE.model
	FROM ROUTE_LOGS
	INNER JOIN VEHICLE ON ROUTE_LOGS.vehicle_id = VEHICLE.id 
	WHERE ROUTE_LOGS.vehicle_id = ?
	`, id)
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
