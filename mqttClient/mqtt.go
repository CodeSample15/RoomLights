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
    var received message
    err := json.Unmarshal(msg.Payload(), &received)
    if err == nil {
        fmt.Printf("Received pico command: %s\n", received.Type)

        switch received.Type {
        case "on":
            commandChan<-lights.LedCommand_on
        case "off":
            commandChan<-lights.LedCommand_off
        case "raise":
            commandChan<-lights.LedCommand_up
        case "lower":
            commandChan<-lights.LedCommand_down
        case "middle":
            commandChan<-lights.LedCommand_middle
        case "middle_night":
            commandChan<-lights.LedCommand_middleNight
        }
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