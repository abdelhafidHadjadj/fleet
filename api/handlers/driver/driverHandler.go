package driverHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func GetDrivers(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query("SELECT id, register_number, firstname, lastname,	date_of_birth, phone, email, class, status, created_at, created_by FROM Driver")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var drivers []model.Driver
	for rows.Next() {
		var driver model.Driver
		if err := rows.Scan(&driver.ID, &driver.Firstname, &driver.Lastname, &driver.Phone,
			&driver.Date_of_birth, &driver.Email, &driver.Class, &driver.Status, &driver.CreatedBy, &driver.CreatedAt, &driver.UpdatedAt); err != nil {
			return err
		}
		drivers = append(drivers, driver)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No drivers present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Drivers found", "data": drivers})
}

func GetDriverByID(c *fiber.Ctx) error {
	id := c.Params("driverId")
	db := database.ConnectionDB()
	row, err := db.Query("SELECT id, register_number, firstname, lastname,	date_of_birth, phone, password, email, class, status, created_at, created_by FROM Driver WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var driver model.Driver
	if row.Next() {
		if err := row.Scan(&driver.ID, &driver.Firstname, &driver.Lastname, &driver.Phone,
			&driver.Date_of_birth, &driver.Email, &driver.Password, &driver.Class, &driver.Status, &driver.CreatedBy, &driver.CreatedAt, &driver.UpdatedAt); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": driver})

}

func CreateDriver(c *fiber.Ctx) error {
	var newDriver model.Driver
	if err := c.BodyParser(&newDriver); err != nil {
		return err
	}
	hashedPassword, _ := HashPassword(newDriver.Password)
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO Driver (register_number, firstname, lastname, date_of_birth, phone, email, password, class, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(newDriver.Register_number, newDriver.Firstname, newDriver.Lastname, newDriver.Date_of_birth, newDriver.Email, newDriver.Phone, hashedPassword, newDriver.Class, newDriver.Status)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not create driver"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Driver created successfully"})
}

func UpdateDriver(c *fiber.Ctx) error {
	var driverDetails model.Driver
	id := c.Params("driverId")
	if err := c.BodyParser(&driverDetails); err != nil {
		return err
	}
	db := database.ConnectionDB()
	var updateQuery string
	var params []interface{} /* declaration t3 interface {}tkhli les type des valuers de params flexible*/

	/* manipuler query dynamiquement not best practice*/

	fieldToUpdate := map[string]interface{}{
		"register_number": driverDetails.Register_number,
		"firstname":       driverDetails.Firstname,
		"lastname":        driverDetails.Lastname,
		"date_of_birth":   driverDetails.Date_of_birth,
		"phone":           driverDetails.Phone,
		"email":           driverDetails.Email,
		"password":        driverDetails.Password,
		"class":           driverDetails.Class,
		"status":          driverDetails.Status,
		"createdBy":       driverDetails.CreatedBy,
	}

	for field, value := range fieldToUpdate {
		if value != "" && field != "phone" || (field == "phone" && value != 0) {
			updateQuery += field + " = ?, "
			params = append(params, value)
		}
	}
	if len(params) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "No fields to update"})
	}

	/* Pour supprimer virgule et l'espace à la fin de string */
	updateQuery = strings.TrimSuffix(updateQuery, ", ")

	stmt, err := db.Prepare("UPDATE User SET " + updateQuery + " WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "User Not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User updated successfully"})
}

func DeleteDriver(c *fiber.Ctx) error {
	id := c.Params("driverId")
	db := database.ConnectionDB()
	stmt, err := db.Prepare("DELETE FROM Driver WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Driver not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Driver Deleted successfully"})
}

/************************* Some Feature ******************************/

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
