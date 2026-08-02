package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	mqttBroker = "100.93.66.64:1883"
	clientID   = "RoomLightsNode"
	topic      = "pico_commands"

	ledCounts = 300
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Recieved message: %s", string(msg.Payload()))
}

func main() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = 12
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = ledCounts

	fmt.Println("Creating device...")
	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("Unable to create strip object")
		return
	}

	dev.Init()
	defer dev.Fini()

	color := uint32(0x000000)
	dev.Leds(0)[0] = color
	dev.Render()
	fmt.Println("done")

	//mqtt test code
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
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

	// Wait for interrupt signal to gracefully shutdown the subscriber
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Unsubscribe and disconnect
	fmt.Println("Unsubscribing and disconnecting...")
	client.Unsubscribe(topic)
	client.Disconnect(250)
}
