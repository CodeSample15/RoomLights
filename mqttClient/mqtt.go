package mqttclient

import (
	"context"
	"encoding/json"
	"fmt"
	"main/lights"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
    // hardcoded cuz it's a tailscale IP lol
	brokerAddress = "100.93.66.64:1883"
	clientID   = "RoomLightsNode"
	topic      = "pico_commands"
)

type message struct {
    Type string
}

var commandChan chan lights.LedCommand

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Recieved message: %s", string(msg.Payload()))

    var received message
    err := json.Unmarshal(msg.Payload(), received)
    if err == nil {
        fmt.Printf("Unmarshaled to: %s", received.Type)
    }
}

func StartMQTTClient(ctx context.Context) chan lights.LedCommand {
    commandChan = make(chan lights.LedCommand)

    opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

    // Subscribe to the topic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)

    go mqttClientShutdown(ctx, client)

    return commandChan
}

func mqttClientShutdown(ctx context.Context, client mqtt.Client) {
    <-ctx.Done()

    // Unsubscribe and disconnect
	fmt.Println("Unsubscribing and disconnecting...")
	client.Unsubscribe(topic)
	client.Disconnect(250)
}