package main

import (
	"main/lights"
	mqttclient "main/mqttClient"
	"time"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	ledGpioPin = 12
	ledCount   = 300
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting MQTT Client")
	comChan := mqttclient.StartMQTTClient(ctx)

	fmt.Println("Starting LED service")
	strip := lights.NewStrip(ledGpioPin, 255)
	if strip.LedCount() == 0 {
		return
	}

	go lights.LedService(ctx, comChan, strip)

	fmt.Println("Running")

	// Wait for interrupt signal to gracefully shutdown the subscriber
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down")

	time.Sleep(500 * time.Millisecond)
}
