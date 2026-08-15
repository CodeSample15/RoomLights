package lights

import (
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Color struct {
	r uint8
	g uint8
	b uint8
}

type Strip interface {
	SetColor(x uint32, c Color)
	GetColor(x uint32) Color
	Render()
	Close()
}

type strip struct {
	device *ws2811.WS2811
	ledCount int
}

func NewStrip(pin int, ledCount int, brightness int) Strip {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = 12
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = ledCount

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("Unable to create strip object")
		return &strip{}
	}

	dev.Init()
	return &strip{
		device: dev,
		ledCount: ledCount,
	}
}

func (led *strip) SetColor(x uint32, c Color) {
	led.device.Leds(0)[0] = c.toInt()
}

func (led *strip) GetColor(x uint32) Color {
	res := Color{}
	res.fromInt(led.device.Leds(0)[x])
	return res
}

func (led *strip) Render() {
	led.device.Render()
}

func (led *strip) Close() {
	led.device.Fini()
}

func (color *Color) fromInt(c uint32) {
	color.b = uint8(c & 0xFF)
	c >>= 8
	color.g = uint8(c & 0xFF)
	c >>= 8
	color.r = uint8(c & 0xFF)
}

func (color *Color) toInt() uint32 {
	return uint32(color.r) << 16 | uint32(color.g) << 8 | uint32(color.b)
}