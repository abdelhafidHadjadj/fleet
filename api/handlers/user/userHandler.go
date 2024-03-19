package userHandler

import (
	"fleet/database"
	model "fleet/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.ConnectionDB()
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname,
			&user.Email, &user.Password, &user.Phone, &user.Role, &user.Status, &user.CreatedAt); err != nil {
			return err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Users found", "data": users})
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("userId")
	db := database.ConnectionDB()
	row, err := db.Query("SELECT * FROM User WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var user model.User
	for row.Next() {
		if err := row.Scan(&user.ID, &user.Firstname, &user.Lastname,
			&user.Email, &user.Password, &user.Phone, &user.Role, &user.Status, &user.CreatedAt); err != nil {
			return err
		}
	}
	if err = row.Err(); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})

}

func CreateUser(c *fiber.Ctx) error {

}

/*
func UpdateUser(c *fiber.Ctx) {
	fmt.Print("update user")
}
func DeleteUser(c *fiber.Ctx) {
	fmt.Print("Delete user")
}
*/
