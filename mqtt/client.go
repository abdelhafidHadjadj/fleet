// client.go
package mqtt

import (
	"fmt"
	"os"
	"os/signal"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
)

func Client() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	opts := mqttlib.NewClientOptions().AddBroker("mqtt.eclipseprojects.io:1883")
	client := mqttlib.NewClient(opts)
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

	// Subscribe to the same topic as your Python server
	topic := "drone/position/+"

	// Keep the client running
	// select {}
	data, _ := Sub(client, topic)
	return data
}
