package mqtt

import (
	"encoding/json"
	"fleet/database"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PositionLogs struct {
	Drone_id string  `json: "drone_id"`
	Lat      float64 `json: "lat"`
	Long     float64 `json: "long"`
	Time     string  `json: "time"`
}

func Sub(client mqtt.Client, topic string) (string, error) {
	var position PositionLogs
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
		err := json.Unmarshal(message.Payload(), &position)
		if err != nil {
			fmt.Println("Error parsing message payload : ", err)
			return
		}
		logToDb(position.Drone_id, position.Lat, position.Long, position.Time)
	}
	token := client.Subscribe(topic, 1, messageHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
	jsonData, err := json.Marshal(position)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func logToDb(vehicle_id string, lat float64, long float64, time string) {
	fmt.Printf("id: %s", vehicle_id)
	fmt.Printf("lat: %f", lat)
	fmt.Printf("long: %f", long)
	fmt.Printf("time: %s", time)
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO Route_logs (vehicle_id, lat, lng, datetime) VALUES (?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(vehicle_id, lat, long, time)
	if err != nil {
		log.Fatal(err)
	}
}
