package vehicleHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetVehicles(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query(`
	SELECT VEHICLE.id, VEHICLE.register_number, VEHICLE.name, VEHICLE.model, VEHICLE.type, VEHICLE.type_charge, VEHICLE.current_charge, VEHICLE.charge_capacity,
	VEHICLE.current_distance ,VEHICLE.current_position , VEHICLE.status, VEHICLE.connection_key,VEHICLE.created_at, VEHICLE.created_by ,USER.firstname, USER.lastname 
	FROM VEHICLE 
	INNER JOIN USER ON VEHICLE.created_by = USER.id`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var vehicles []model.Vehicle
	var userFirstname, userLastname string
	for rows.Next() {
		var vehicle model.Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Register_number, &vehicle.Name, &vehicle.Model,
			&vehicle.Type, &vehicle.Type_charge, &vehicle.Current_charge, &vehicle.Charge_capacity, &vehicle.Current_distance, &vehicle.Current_position, &vehicle.Status, &vehicle.Connection_key, &vehicle.CreatedAt, &vehicle.CreatedBy, &userFirstname, &userLastname); err != nil {
			return err
		}
		vehicle.CreatedBy = userFirstname + " " + userLastname
		vehicles = append(vehicles, vehicle)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No vehicles present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Vehicles found", "data": vehicles})
}

func GetVehicleByID(c *fiber.Ctx) error {
	id := c.Params("vehicleId")
	db := database.ConnectionDB()
	row, err := db.Query(`
	SELECT VEHICLE.id, VEHICLE.register_number, VEHICLE.name, VEHICLE.model, VEHICLE.type, VEHICLE.type_charge, VEHICLE.current_charge, VEHICLE.charge_capacity,
	VEHICLE.current_distance ,VEHICLE.current_position , VEHICLE.status, VEHICLE.connection_key,VEHICLE.created_at, VEHICLE.created_by ,USER.firstname, USER.lastname 
	FROM VEHICLE 
	INNER JOIN USER ON VEHICLE.created_by = USER.id
	WHERE VEHICLE.id = ?
	`, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var vehicle model.Vehicle
	var userFirstname, userLastname string
	if row.Next() {
		if err := row.Scan(&vehicle.ID, &vehicle.Register_number, &vehicle.Name, &vehicle.Model,
			&vehicle.Type, &vehicle.Type_charge, &vehicle.Current_charge, &vehicle.Charge_capacity, &vehicle.Current_distance, &vehicle.Current_position, &vehicle.Status, &vehicle.Connection_key, &vehicle.CreatedAt, &vehicle.CreatedBy, &userFirstname, &userLastname); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})

	}
	vehicle.CreatedBy = userFirstname + " " + userLastname
	return c.JSON(fiber.Map{"status": "success", "message": "Vehicle found", "data": vehicle})
}

func CreateVehicle(c *fiber.Ctx) error {
	var newVehicle model.Vehicle
	if err := c.BodyParser(&newVehicle); err != nil {
		return err
	}
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO VEHICLE (register_number, name, model, type, type_charge, current_charge, charge_capacity, current_distance, current_position, status, connection_key, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(newVehicle.Register_number, newVehicle.Name, newVehicle.Model, newVehicle.Type, newVehicle.Type_charge, newVehicle.Current_charge, newVehicle.Charge_capacity, newVehicle.Current_distance, newVehicle.Current_position, newVehicle.Status, newVehicle.Connection_key, newVehicle.CreatedBy)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not create vehicle"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Vehicle created successfully"})

}

func UpdateVehicle(c *fiber.Ctx) error {
	var vehicleDetails model.Vehicle
	id := c.Params("vehicleId")
	if err := c.BodyParser(&vehicleDetails); err != nil {
		return err
	}
	db := database.ConnectionDB()
	var updateQuery string
	var params []interface{} /* declaration t3 interface {}tkhli les type des valuers de params flexible*/

	/* manipuler query dynamiquement not best practice*/

	if vehicleDetails.Register_number != "" {
		updateQuery += "register_number = ?, "
		params = append(params, vehicleDetails.Register_number)
	}
	if vehicleDetails.Name != "" {
		updateQuery += "name = ?, "
		params = append(params, vehicleDetails.Name)
	}
	if vehicleDetails.Model != "" {
		updateQuery += "model = ?, "
		params = append(params, vehicleDetails.Model)
	}
	if vehicleDetails.Type != "" {
		updateQuery += "type = ?, "
		params = append(params, vehicleDetails.Type)
	}
	if vehicleDetails.Type_charge != "" {
		updateQuery += "type_charge = ?, "
		params = append(params, vehicleDetails.Type_charge)
	}
	if vehicleDetails.Current_charge != "" {
		updateQuery += "current_charge = ?, "
		params = append(params, vehicleDetails.Current_charge)
	}
	if vehicleDetails.Charge_capacity != "" {
		updateQuery += "charge_capacity = ?, "
		params = append(params, vehicleDetails.Charge_capacity)
	}
	if vehicleDetails.Current_distance != "" {
		updateQuery += "current_distance = ?, "
		params = append(params, vehicleDetails.Current_distance)
	}
	if vehicleDetails.Current_position != "" {
		updateQuery += "current_position = ?, "
		params = append(params, vehicleDetails.Current_position)
	}
	if vehicleDetails.Status != "" {
		updateQuery += "status = ?, "
		params = append(params, vehicleDetails.Status)
	}
	if vehicleDetails.Connection_key != "" {
		updateQuery += "connection_key = ?, "
		params = append(params, vehicleDetails.Connection_key)
	}

	/* Pour supprimer virgule et l'espace à la fin de string */
	updateQuery = strings.TrimSuffix(updateQuery, ", ")

	stmt, err := db.Prepare("UPDATE VEHICLE SET " + updateQuery + " WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Vehicle Not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Vehicle updated successfully"})
}

func DeleteVehicle(c *fiber.Ctx) error {
	id := c.Params("vehicleId")
	db := database.ConnectionDB()
	stmt, err := db.Prepare("DELETE FROM VEHICLE WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Vehicle not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Vehicle Deleted successfully"})
}
