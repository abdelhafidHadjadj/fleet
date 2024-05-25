package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"github.com/pion/webrtc/v3"

	router "fleet/api/routers"
	"fleet/database"
	"fleet/mqtt"
)

var peerConnection *webrtc.PeerConnection

func main() {
	app := fiber.New()

	// Initialize MQTT client
	data := mqtt.Client()

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

	// Start Fiber server
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

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

	// WebRTC signaling
	socketServer.OnEvent("/", "offer", func(so socketio.Conn, msg string) {
		var offer webrtc.SessionDescription
		if err := json.Unmarshal([]byte(msg), &offer); err != nil {
			log.Println("Error unmarshaling offer:", err)
			return
		}

		pc, err := setupWebRTC()
		if err != nil {
			log.Println("Error setting up WebRTC:", err)
			return
		}

		if err := pc.SetRemoteDescription(offer); err != nil {
			log.Println("Error setting remote description:", err)
			return
		}

		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			log.Println("Error creating answer:", err)
			return
		}

		if err := pc.SetLocalDescription(answer); err != nil {
			log.Println("Error setting local description:", err)
			return
		}

		answerJSON, err := json.Marshal(answer)
		if err != nil {
			log.Println("Error marshaling answer:", err)
			return
		}

		so.Emit("answer", string(answerJSON))
	})

	socketServer.OnEvent("/", "candidate", func(so socketio.Conn, msg string) {
		var candidate webrtc.ICECandidateInit
		if err := json.Unmarshal([]byte(msg), &candidate); err != nil {
			log.Println("Error unmarshaling candidate:", err)
			return
		}

		if peerConnection != nil {
			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println("Error adding ICE candidate:", err)
			}
		}
	})

	http.Handle("/socket.io/", socketServer)

	log.Println("Server is running at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupWebRTC() (*webrtc.PeerConnection, error) {
	mediaEngine := webrtc.MediaEngine{}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}

	pc, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	peerConnection = pc

	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo)
	if err != nil {
		return nil, err
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			candidateJSON, err := json.Marshal(candidate.ToJSON())
			if err == nil {
				// Send ICE candidate to the client
				// Assuming you have access to the current socket
				// Replace `so.Emit` with the appropriate socket emit function
				// so.Emit("candidate", string(candidateJSON))
			}
		}
	})

	return pc, nil
}
