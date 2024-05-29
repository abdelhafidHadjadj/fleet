package userHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query(`
	SELECT id, firstname, lastname, phone, email, role, status, created_at 
	FROM USER
	`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Phone,
			&user.Email, &user.Role, &user.Status, &user.CreatedAt); err != nil {
			return err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "No users present", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("userId")
	db := database.ConnectionDB()
	row, err := db.Query("SELECT id, firstname, lastname, phone, email, role, status FROM USER WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var user model.User
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Firstname, &user.Lastname,
			&user.Phone, &user.Email, &user.Role, &user.Status); err != nil {
			return err
		}
	} else {
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "Not found", "data": nil})
	}
	if err = row.Err(); err != nil {
		log.Println("Row error:", err)
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Row error"})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})

}

func CreateUser(c *fiber.Ctx) error {
	var newUser model.User
	if err := c.BodyParser(&newUser); err != nil {
		return err
	}
	hashedPassword, _ := HashPassword(newUser.Password)
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO USER (firstname, lastname, email, password, phone, role, status) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(newUser.Firstname, newUser.Lastname, newUser.Email, hashedPassword, newUser.Phone, newUser.Role, newUser.Status)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).JSON(fiber.Map{"status": "failed", "message": "Could not create user"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User created successfully"})
}

func UpdateUser(c *fiber.Ctx) error {
	var userDetails model.User
	id := c.Params("userId")
	if err := c.BodyParser(&userDetails); err != nil {
		return err
	}
	db := database.ConnectionDB()
	var updateQuery string
	var params []interface{} /* declaration t3 interface {}tkhli les type des valuers de params flexible*/

	/* manipuler query dynamiquement not best practice*/

	if userDetails.Firstname != "" {
		updateQuery += "firstname = ?, "
		params = append(params, userDetails.Firstname)
	}
	if userDetails.Lastname != "" {
		updateQuery += "lastname = ?, "
		params = append(params, userDetails.Lastname)
	}
	if userDetails.Email != "" {
		updateQuery += "email = ?, "
		params = append(params, userDetails.Email)
	}
	if userDetails.Password != "" {
		updateQuery += "password = ?, "
		params = append(params, userDetails.Password)
	}
	if userDetails.Phone != "" {
		updateQuery += "phone = ?, "
		params = append(params, userDetails.Phone)
	}
	if userDetails.Role != "" {
		updateQuery += "role = ?, "
		params = append(params, userDetails.Role)
	}
	if userDetails.Status != "" {
		updateQuery += "status = ?, "
		params = append(params, userDetails.Status)
	}

	/* Pour supprimer virgule et l'espace à la fin de string */
	updateQuery = strings.TrimSuffix(updateQuery, ", ")

	stmt, err := db.Prepare("UPDATE USER SET " + updateQuery + " WHERE id = ?")
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

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("userId")
	db := database.ConnectionDB()
	stmt, err := db.Prepare("DELETE FROM USER WHERE id = ?")
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
		return c.Status(404).JSON(fiber.Map{"status": "failed", "message": "User not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User Deleted successfully"})
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
