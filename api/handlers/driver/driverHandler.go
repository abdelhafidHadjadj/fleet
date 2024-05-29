package driverHandler

import (
	"fleet/database"
	model "fleet/models"
	"fleet/utils"
	"log"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func GetDrivers(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query(`
	SELECT DRIVER.id, DRIVER.register_number, DRIVER.firstname, DRIVER.lastname, DRIVER.date_of_birth, DRIVER.phone, DRIVER.email, DRIVER.class, DRIVER.status, DRIVER.created_at, DRIVER.updated_at, DRIVER.created_by, USER.firstname, USER.lastname 
	FROM DRIVER
	INNER JOIN USER ON DRIVER.created_by = USER.id
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var drivers []model.Driver
	var userFirstname, userLastname string
	for rows.Next() {
		var driver model.Driver
		if err := rows.Scan(&driver.ID, &driver.Register_number, &driver.Firstname, &driver.Lastname,
			&driver.Date_of_birth, &driver.Phone, &driver.Email, &driver.Class, &driver.Status, &driver.CreatedAt, &driver.UpdatedAt, &driver.CreatedBy, &userFirstname, &userLastname); err != nil {
			return err
		}
		driver.CreatedByUserName = userFirstname + " " + userLastname
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
	row, err := db.Query(`
	SELECT DRIVER.id, DRIVER.register_number, DRIVER.firstname, DRIVER.lastname, DRIVER.date_of_birth, DRIVER.phone, DRIVER.email, DRIVER.class, DRIVER.status, DRIVER.created_at, DRIVER.updated_at ,DRIVER.created_by ,USER.firstname, USER.lastname 
	FROM DRIVER
	INNER JOIN USER ON DRIVER.created_by = USER.id
	WHERE DRIVER.id = ?
	`, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var driver model.Driver
	var userFirstname, userLastname string

	if row.Next() {
		if err := row.Scan(&driver.ID, &driver.Register_number, &driver.Firstname, &driver.Lastname,
			&driver.Date_of_birth, &driver.Phone, &driver.Email, &driver.Class, &driver.Status, &driver.CreatedAt, &driver.UpdatedAt, &driver.CreatedBy, &userFirstname, &userLastname); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})
	}
	driver.CreatedByUserName = userFirstname + " " + userLastname
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": driver})

}
func CreateDriver(c *fiber.Ctx) error {
	var newDriver model.Driver
	if err := c.BodyParser(&newDriver); err != nil {
		return err
	}

	// Parse created_by from form data
	createdBy := c.FormValue("created_by")
	if createdBy == "" {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "created_by is required"})
	}

	// Convert created_by to an integer
	createdByInt, err := strconv.Atoi(createdBy)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "created_by must be a valid integer"})
	}
	newDriver.CreatedBy = createdByInt

	// Handle avatar upload
	avatarPath, err := utils.SaveFile(c, "avatar", "./uploads/avatars")
	if err != nil && err.Error() != "http: no such file" {
		log.Printf("Error saving file: %v", err) // Log error
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not upload avatar"})
	} else if err != nil && err.Error() == "http: no such file" {
		avatarPath = "" // No file uploaded
	} else {
		log.Printf("File saved at: %s", avatarPath) // Log success
	}
	newDriver.Avatar = avatarPath

	// Hash the password
	hashedPassword, err := HashPassword(newDriver.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err) // Log error
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not hash password"})
	}

	// Get database connection
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO DRIVER (register_number, firstname, lastname, date_of_birth, avatar, phone, email, password, class, status, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not prepare SQL statement"})
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(newDriver.Register_number, newDriver.Firstname, newDriver.Lastname, newDriver.Date_of_birth, newDriver.Avatar, newDriver.Phone, newDriver.Email, hashedPassword, newDriver.Class, newDriver.Status, newDriver.CreatedBy)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not create driver", "error": err.Error()})
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
	}

	for field, value := range fieldToUpdate {
		if value != "" {
			updateQuery += field + " = ?, "
			params = append(params, value)
		}
	}
	if len(params) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "failed", "message": "No fields to update"})
	}

	/* Pour supprimer virgule et l'espace à la fin de string */
	updateQuery = strings.TrimSuffix(updateQuery, ", ")

	stmt, err := db.Prepare("UPDATE DRIVER SET " + updateQuery + " WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Driver Not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Driver updated successfully"})
}

func DeleteDriver(c *fiber.Ctx) error {
	id := c.Params("driverId")
	db := database.ConnectionDB()
	stmt, err := db.Prepare("DELETE FROM DRIVER WHERE id = ?")
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
