// main.go
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"

	authHandler "fleet/api/handlers/auth"
	"fleet/api/middlewares"
	router "fleet/api/routers"
	"fleet/config"
	"fleet/database"
	"fleet/mqtt"
)

var (
	clients    = make(map[*websocket.Conn]bool)
	broadcast  = make(chan []byte)
	clientLock = sync.Mutex{}
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8001", // Replace with the actual URL where your Svelte app is hosted
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Initialize MQTT client with message handler
	mqtt.SetupClient(func(topic string, payload []byte) {
		log.Printf("MQTT message received: %s - %s", topic, payload)
		broadcast <- payload
	})

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

	// WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clientLock.Lock()
		clients[c] = true
		clientLock.Unlock()
		log.Println("WebSocket connection opened")

		defer func() {
			clientLock.Lock()
			delete(clients, c)
			clientLock.Unlock()
			log.Println("WebSocket connection closed")
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
		}
	}))

	// Broadcast MQTT messages to WebSocket clients
	go func() {
		for {
			msg := <-broadcast
			clientLock.Lock()
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
			clientLock.Unlock()
		}
	}()

	// Start Fiber server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
