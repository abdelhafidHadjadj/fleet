// client.go
package mqtt

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
)

var client mqttlib.Client

func SetupClient(messageHandler func(topic string, payload []byte)) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	opts := mqttlib.NewClientOptions().AddBroker("tcp://mqtt.eclipseprojects.io:1883").SetClientID("fleet_client")
	client = mqttlib.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		for range c {
			fmt.Println("Exiting...")
			client.Disconnect(250)
			os.Exit(0)
		}
	}()

	// Subscribe to the topic
	topic := "vehicle/position/+"
	client.Subscribe(topic, 1, func(client mqttlib.Client, msg mqttlib.Message) {
		messageHandler(msg.Topic(), msg.Payload())

		type Position struct {
			VehicleID string `json:"vehicle_id"`
			Lat       string `json:"lat"`
			Lng       string `json:"lng"`
		}

		// Parse the JSON payload
		var position Position
		err := json.Unmarshal(msg.Payload(), &position)
		if err != nil {
			fmt.Printf("Error parsing JSON payload: %v\n", err)
			return
		}

		// Call logToDb
		logToDb(position.VehicleID, position.Lat, position.Lng)

	})
}

func Client() mqttlib.Client {
	return client
}
