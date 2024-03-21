package routeHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetRoutes(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query("SELECT * FROM route")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var routes []model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.Status, &route.Departue_date, &route.Arrival_date,
			&route.DriverID, &route.VehicleID, &route.CreatedAt); err != nil {
			return err
		}
		routes = append(routes, route)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No routes present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "routes found", "data": routes})
}

func GetRouteByID(c *fiber.Ctx) error {
	id := c.Params("routeId")
	db := database.ConnectionDB()
	row, err := db.Query("SELECT * FROM Route WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var route model.Route
	if row.Next() {
		if err := row.Scan(&route.ID, &route.Status, &route.Departue_date,
			&route.Arrival_date, &route.DriverID, &route.VehicleID, &route.CreatedAt); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Route found", "data": route})

}

func CreateRoute(c *fiber.Ctx) error {
	var newRoute model.Route
	if err := c.BodyParser(&newRoute); err != nil {
		return err
	}
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO Route (status, departure_date, arrival_date, driver_id, vehicle_id, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(newRoute.Status, newRoute.Departue_date, newRoute.Arrival_date, newRoute.DriverID, newRoute.VehicleID, newRoute.CreatedAt, newRoute.CreatedBy)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not create route"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Route created successfully"})

}

func UpdateRoute(c *fiber.Ctx) error {
	var routeDetails model.Route
	id := c.Params("routeId")
	if err := c.BodyParser(&routeDetails); err != nil {
		return err
	}
	db := database.ConnectionDB()
	var updateQuery string
	var params []interface{} /* declaration t3 interface {}tkhli les type des valuers de params flexible*/

	/* manipuler query dynamiquement not best practice*/

	if routeDetails.Status != "" {
		updateQuery += "status = ?, "
		params = append(params, routeDetails.Status)
	}
	if routeDetails.Departue_date != "" {
		updateQuery += "departure_date = ?, "
		params = append(params, routeDetails.Departue_date)
	}
	if routeDetails.Arrival_date != "" {
		updateQuery += "arrival_date = ?, "
		params = append(params, routeDetails.Arrival_date)
	}
	if routeDetails.DriverID != "" {
		updateQuery += "driver_id = ?, "
		params = append(params, routeDetails.DriverID)
	}
	if routeDetails.VehicleID != "" {
		updateQuery += "vehicle_id = ?, "
		params = append(params, routeDetails.VehicleID)
	}

	/* Pour supprimer virgule et l'espace à la fin de string */
	updateQuery = strings.TrimSuffix(updateQuery, ", ")

	stmt, err := db.Prepare("UPDATE Route SET " + updateQuery + " WHERE id = ?")
	params = append(params, id) /* hna nzid appender id */
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(params...) /* destruct params*/
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Error executing statement", "data": nil})
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Route Not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Route updated successfully"})
}

func DeleteRoute(c *fiber.Ctx) error {
	id := c.Params("routeId")
	db := database.ConnectionDB()
	stmt, err := db.Prepare("DELETE FROM Route WHERE id = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Database error"})
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Route not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Route Deleted successfully"})
}
