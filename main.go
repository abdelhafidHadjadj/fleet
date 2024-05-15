package main

import (
	router "fleet/api/routers"
	"fleet/database"
	"fleet/mqtt"
	"fmt"
	//"log"
	//"net/http"
	"github.com/gofiber/fiber/v2"
	//socketio "github.com/googollee/go-socket.io"
)

func main() {
	app := fiber.New()
	data := mqtt.Client()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	db := database.ConnectionDB()
	err := database.CreateTables(db)
	if err != nil {
		fmt.Printf("%s", err)
	}
	router.SetupRoutes(app)
	app.Listen(":8080")
	go func() {
		if err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()
	socketServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("Error establishing new socketio server")
	}
	socketServer.On("connection", func(so socketio.Socket) {
		log.Println("Socket.IO connection established")
		so.Join("tracking")
		so.On("track vehicle", func() {
			log.Println("Data received: " + data)
			so.BroadcastTo("tracking", "track vehicle", data)
		})
	})
	http.Handle("/socket.io/", socketServer)

}
