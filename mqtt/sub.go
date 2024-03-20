package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
