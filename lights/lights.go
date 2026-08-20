package lights

import (
	"context"
	"fmt"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type color struct {
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
	SetColor(x int, c color)
	GetColor(x int) color
	LedCount() int
	Render()
	Close()
}

type strip struct {
	device   *ws2811.WS2811
	ledCount int
}

type LedCommand uint16

const (
	LedCommand_on          LedCommand = 0
	LedCommand_off         LedCommand = 1
	LedCommand_up          LedCommand = 2
	LedCommand_down        LedCommand = 3
	LedCommand_middle      LedCommand = 4
	LedCommand_middleNight LedCommand = 5
)

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
		device:   dev,
		ledCount: ledCount,
	}
}

func LedService(ctx context.Context, commands chan LedCommand, led Strip) {
	defer led.Close()

	var transition float32
	var transitionSpeed float32
	updateTick := time.Tick(10 * time.Millisecond)

	var targetPattern Pattern = &offPattern{}
	state := make([]colorF, led.LedCount())

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop
		case com := <-commands:
			transition = 0

			switch com {
			case LedCommand_on:
				break
			case LedCommand_off:
				targetPattern = &offPattern{}
				transitionSpeed = 0.1
			case LedCommand_up:
				break
			case LedCommand_down:
				break
			case LedCommand_middle:
				break
			case LedCommand_middleNight:
				targetPattern = &nightLight{}
				transitionSpeed = 0.05
			}
		case <-updateTick:
			targetPattern.tick()

			if transition < 1 {
				transition += transitionSpeed
			} else {
				transition = 1
			}

			for c := range led.LedCount() {
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

	fmt.Println("LED service shutdown")
}

func (led *strip) SetColor(x int, c color) {
	led.device.Leds(0)[x] = c.toInt()
}

func (led *strip) GetColor(x int) color {
	res := color{}
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

func (col *color) fromInt(c uint32) {
	col.b = uint8(c & 0xFF)
	c >>= 8
	col.g = uint8(c & 0xFF)
	c >>= 8
	col.r = uint8(c & 0xFF)
}

func (col *color) toInt() uint32 {
	return uint32(col.r)<<16 | uint32(col.g)<<8 | uint32(col.b)
}

func (col *colorF) diff(other color, speed float32) {
	col.r += (float32(other.r) - col.r) * speed
	col.g += (float32(other.g) - col.g) * speed
	col.b += (float32(other.b) - col.b) * speed
}

func (col *colorF) toColor() color {
	return color{
		r: uint8(col.r),
		g: uint8(col.g),
		b: uint8(col.b),
	}
}
