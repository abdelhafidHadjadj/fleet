package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	authHandler "fleet/api/handlers/auth"
	"fleet/api/middlewares"
	router "fleet/api/routers"
	"fleet/config"
	"fleet/database"
	"fleet/mqtt"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8001", // Replace with the actual URL where your Svelte app is hosted
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	})) // Initialize MQTT client
	mqtt.Client()
	// Basic route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Database setup
	db := database.ConnectionDB()
	err := database.CreateTables(db)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Setup routes
	router.SetupRoutes(app)
	jwt := middlewares.NewAuthMiddleware(config.Secret)
	app.Post("/login", authHandler.Login)
	app.Get("/protected", jwt, authHandler.Protected)

	// Start Fiber + websockets server
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		go func() {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					message := []byte("Message from server")
					if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
						log.Println("Write error:", err)
						return
					}
				}
			}
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s\n", msg)
		}
	}))

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}

}
