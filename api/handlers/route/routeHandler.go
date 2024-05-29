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
	rows, err := db.Query(`
	SELECT ROUTE.id, ROUTE.status, ROUTE.departure_date, ROUTE.arrival_date, 
	ROUTE.lat_start, ROUTE.lng_start, ROUTE.lat_end, ROUTE.lng_end, ROUTE.departure_city, ROUTE.arrival_city,
	ROUTE.driver_id ,DRIVER.firstname, DRIVER.lastname, ROUTE.vehicle_id, VEHICLE.register_number, VEHICLE.name, VEHICLE.model, ROUTE.created_at, USER.firstname, USER.lastname 	
	FROM ROUTE
	LEFT JOIN 
    DRIVER ON ROUTE.driver_id = DRIVER.id
	LEFT JOIN 
    VEHICLE ON ROUTE.vehicle_id = VEHICLE.id
	LEFT JOIN 
	USER ON ROUTE.created_by = USER.id
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var routes []model.Route
	var driverFirstName, driverLastName, vehicleName, vehicleModel, userFirstName, userLastName string

	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.Status, &route.Departue_date, &route.Arrival_date, &route.Lat_start, &route.Lng_start, &route.Lat_end, &route.Lng_end, &route.Departure_city, &route.Arrival_city,
			&route.DriverID, &driverFirstName, &driverLastName, &route.VehicleID, &route.VehicleRegisterNumber, &vehicleName, &vehicleModel, &route.CreatedAt, &userFirstName, &userLastName); err != nil {
			return err
		}
		route.DriverName = driverFirstName + " " + driverLastName
		route.VehicleName = vehicleName + " (" + vehicleModel + ")"
		route.CreatedBy = userFirstName + " " + userLastName
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
	row, err := db.Query(`
	SELECT ROUTE.id, ROUTE.status, ROUTE.departure_date, ROUTE.arrival_date, 
	ROUTE.lat_start, ROUTE.lng_start, ROUTE.lat_end, ROUTE.lng_end, ROUTE.departure_city, ROUTE.arrival_city,
	ROUTE.driver_id ,DRIVER.firstname, DRIVER.lastname, ROUTE.vehicle_id ,VEHICLE.register_number, VEHICLE.name, VEHICLE.model, ROUTE.created_at, USER.firstname, USER.lastname 	
	FROM ROUTE
	LEFT JOIN 
    DRIVER ON ROUTE.driver_id = DRIVER.id
	LEFT JOIN 
    VEHICLE ON ROUTE.vehicle_id = VEHICLE.id
	LEFT JOIN 
	USER ON ROUTE.created_by = USER.id
	WHERE ROUTE.id = ?
	`, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var route model.Route
	var driverFirstName, driverLastName, vehicleName, vehicleModel, userFirstName, userLastName string

	if row.Next() {
		if err := row.Scan(&route.ID, &route.Status, &route.Departue_date, &route.Arrival_date, &route.Lat_start, &route.Lng_start, &route.Lat_end, &route.Lng_end, &route.Departure_city, &route.Arrival_city,
			&route.DriverID, &driverFirstName, &driverLastName, &route.VehicleID, &route.VehicleRegisterNumber, &vehicleName, &vehicleModel, &route.CreatedAt, &userFirstName, &userLastName); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})

	}
	route.DriverName = driverFirstName + " " + driverLastName
	route.VehicleName = vehicleName + " (" + vehicleModel + ")"
	route.CreatedBy = userFirstName + " " + userLastName

	return c.JSON(fiber.Map{"status": "success", "message": "Route found", "data": route})

}

func CreateRoute(c *fiber.Ctx) error {
	var newRoute model.Route
	if err := c.BodyParser(&newRoute); err != nil {
		return err
	}
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO ROUTE (status, departure_date, arrival_date, lat_start, lng_start, lat_end, lng_end, departure_city, arrival_city, driver_id, vehicle_id, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(newRoute.Status, newRoute.Departue_date, newRoute.Arrival_date, newRoute.Lat_start, newRoute.Lng_start, newRoute.Lat_end, newRoute.Lng_end, newRoute.Departure_city, newRoute.Arrival_city, newRoute.DriverID, newRoute.VehicleID, newRoute.CreatedBy)
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
	if routeDetails.Departure_city != "" {
		updateQuery += "arrival_city = ?, "
		params = append(params, routeDetails.Departure_city)
	}
	if routeDetails.Arrival_city != "" {
		updateQuery += "arrival_city = ?, "
		params = append(params, routeDetails.Arrival_city)
	}
	if routeDetails.Lat_start != "" {
		updateQuery += "lat_start = ?, "
		params = append(params, routeDetails.Lat_start)
	}
	if routeDetails.Lng_start != "" {
		updateQuery += "lng_start = ?, "
		params = append(params, routeDetails.Lng_start)
	}
	if routeDetails.Lat_end != "" {
		updateQuery += "lat_end = ?, "
		params = append(params, routeDetails.Lat_end)
	}
	if routeDetails.Lng_end != "" {
		updateQuery += "lng_end = ?, "
		params = append(params, routeDetails.Lng_end)
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

	stmt, err := db.Prepare("UPDATE ROUTE SET " + updateQuery + " WHERE id = ?")
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
