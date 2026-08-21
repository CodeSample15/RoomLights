package lights

import (
	"context"
	"fmt"
	"math"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type color struct {
	r uint8
	g uint8
	b uint8
}

type colorF struct {
	r float64
	g float64
	b float64
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

	LED_COUNT = 300
)

var patterns []Pattern = []Pattern{
	&sinWavePattern{
		rate:  5,
		state: make([]colorF, LED_COUNT),
	},
	&stars{
		createTime:      50 * time.Millisecond,
		decayRate:       0.01,
		backgroundColor: color{0, 20, 70},
		state:           make([]colorF, LED_COUNT),
	},
}

func NewStrip(pin int, brightness int) Strip {
	opt := ws2811.DefaultOptions
	opt.Channels[0].GpioPin = pin
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = LED_COUNT

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("Unable to create strip object")
		return &strip{}
	}

	dev.Init()
	return &strip{
		device:   dev,
		ledCount: LED_COUNT,
	}
}

func LedService(ctx context.Context, commands chan LedCommand, led Strip) {
	defer led.Close()

	var transitionSpeed float64
	updateTick := time.Tick(10 * time.Millisecond)

	var targetPattern Pattern = &offPattern{}
	patternIndex := 0
	lightsOn := false
	state := make([]colorF, led.LedCount())

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop
		case com := <-commands:
			switch com {
			case LedCommand_on:
				targetPattern = patterns[patternIndex]
				transitionSpeed = 0.1
				lightsOn = true

			case LedCommand_off:
				targetPattern = &offPattern{}

				transitionSpeed = 0.1
				lightsOn = false

			case LedCommand_up:
				patternIndex += 1
				if patternIndex >= len(patterns) {
					patternIndex = 0
				}

				if lightsOn {
					targetPattern = patterns[patternIndex]
					transitionSpeed = 0.1
				}

			case LedCommand_down:
				patternIndex -= 1
				if patternIndex < 0 {
					patternIndex = len(patterns) - 1
				}

				if lightsOn {
					targetPattern = patterns[patternIndex]
					transitionSpeed = 0.1
				}

			case LedCommand_middle:
				break

			case LedCommand_middleNight:
				targetPattern = &nightLight{
					nextPixelTime: 50 * time.Millisecond,
				}
				targetPattern.reset(color{})
				transitionSpeed = 0.01
			}
		case <-updateTick:
			targetPattern.tick(led.LedCount())

			for c := range led.LedCount() {
				target := targetPattern.get(c)

				if transitionSpeed == 0 {
					state[c] = colorF{float64(target.r), float64(target.g), float64(target.b)}
				} else {
					state[c].diff(target, transitionSpeed)
				}

				led.SetColor(c, state[c].toColor())
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

func (col *colorF) diff(other color, speed float64) {
	col.r += (float64(other.r) - col.r) * speed
	col.g += (float64(other.g) - col.g) * speed
	col.b += (float64(other.b) - col.b) * speed
}

func (col *colorF) toColor() color {
	return color{
		r: uint8(math.Max(0, math.Min(col.r, 255))),
		g: uint8(math.Max(0, math.Min(col.g, 255))),
		b: uint8(math.Max(0, math.Min(col.b, 255))),
	}
}
