package main

import (
	router "fleet/api/routers"
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

	router.SetupRoutes(app)

	err = app.Listen(":8080")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
