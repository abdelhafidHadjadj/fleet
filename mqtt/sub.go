package mqtt

import (
	"encoding/json"
	"fleet/database"
	model "fleet/models"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Sub(client mqtt.Client, topic string) (string, error) {
	var position model.RouteLogs
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
		err := json.Unmarshal(message.Payload(), &position)
		if err != nil {
			fmt.Println("Error parsing message payload : ", err)
			return
		}
		logToDb(position.VehicleID, position.Lat, position.Lng)
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

func logToDb(vehicle_id string, lat string, lng string) {
	fmt.Printf("id: %s", vehicle_id)
	fmt.Printf("lat: %s", lat)
	fmt.Printf("lng: %s", lng)
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO ROUTE_LOGS (vehicle_id, lat, lng) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(vehicle_id, lat, lng)
	if err != nil {
		log.Fatal(err)
	}
}
