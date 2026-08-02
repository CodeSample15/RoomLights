package main

import (
	"fmt"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	ledCounts = 300
)

func main() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("Unable to create strip object")
		return
	}

	defer dev.Fini()

	color := uint32(0x0000ff)
	dev.Leds(0)[0] = color
	dev.Render()
}
