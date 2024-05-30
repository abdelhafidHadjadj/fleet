package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	authHandler "fleet/api/handlers/auth"
	"fleet/api/middlewares"
	router "fleet/api/routers"
	"fleet/config"
	"fleet/database"
	"fleet/mqtt"
)

func main() {
	app := fiber.New()

	// Initialize MQTT client
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
	app.Get("login", authHandler.Login)
	app.Get("/protected", jwt, authHandler.Protected)

	// Start Fiber server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
	/*
		// Setup Socket.IO server
		socketServer, err := socketio.NewServer(nil)
		if err != nil {
			log.Fatal("Error establishing new socketio server")
		}

		socketServer.OnConnect("/", func(so socketio.Conn) error {
			log.Println("Socket.IO connection established")
			so.Join("tracking")
			return nil
		})

		socketServer.OnEvent("/", "track vehicle", func(so socketio.Conn) {
			log.Println("Data received: ", data)
			so.BroadcastToRoom("/", "tracking", "track vehicle", data)
		})

		http.Handle("/socket.io/", socketServer)

		log.Println("Server is running at :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	*/
}
