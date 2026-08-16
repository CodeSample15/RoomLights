package lights

import (
	"context"
	"fmt"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Color struct {
	r uint8
	g uint8
	b uint8
}

type colorF struct {
	r float32
	g float32
	b float32
}

type Strip interface {
	SetColor(x int, c Color)
	GetColor(x int) Color
	LedCount() int
	Render()
	Close()
}

type strip struct {
	device *ws2811.WS2811
	ledCount int
}

type LedCommand struct {
	lightPattern Pattern
	transitionSpeed float32
	reset bool
	color Color
}

func NewStrip(pin int, ledCount int, brightness int) Strip {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = pin
	opt.Channels[0].Brightness = brightness
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

func LedService(ctx context.Context, commands chan LedCommand, led Strip) {
	defer led.Close()

	var transition float32
	var transitionSpeed float32
	updateTick := time.Tick(50 * time.Millisecond)

	var targetPattern Pattern = &offPattern{}
	state := make([]colorF, led.LedCount())

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop
		case com:=<-commands:
			targetPattern = com.lightPattern
			transitionSpeed = com.transitionSpeed
			transition = 0
			if com.reset {
				targetPattern.reset(com.color)
			}
		case <-updateTick:
			targetPattern.tick()

			if transition < 1 {
				transition += transitionSpeed
			} else {
				transition = 1
			}

			for c:=range led.LedCount() {
				if transitionSpeed == 0 {
					led.SetColor(c, targetPattern.get(c))
				} else {
					state[c].diff(targetPattern.get(c), transitionSpeed)
					led.SetColor(c, state[c].toColor())
				}
			}

			led.Render()
		}
	}
}

func (led *strip) SetColor(x int, c Color) {
	led.device.Leds(0)[x] = c.toInt()
}

func (led *strip) GetColor(x int) Color {
	res := Color{}
	res.fromInt(led.device.Leds(0)[x])
	return res
}

func (led *strip) LedCount() int {
	return led.ledCount
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

func (color *colorF) diff(other Color, speed float32) {
	color.r += (float32(other.r) - color.r) * speed
	color.g += (float32(other.g) - color.g) * speed
	color.b += (float32(other.b) - color.b) * speed
}

func (color *colorF) toColor() Color {
	return Color{
		r: uint8(color.r),
		g: uint8(color.g),
		b: uint8(color.b),
	}
}