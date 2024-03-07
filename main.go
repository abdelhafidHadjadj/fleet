package main

import (
	"fleet/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	db := database.ConnectionDB()
	err := database.CreateTables(db)
	if err != nil {
		fmt.Printf("%s", err)
	}
	app.Listen(":3000")

}
