package lights

import (
	"math"
	"math/rand"
	"time"
)

type Pattern interface {
	tick(ledCount int)
	reset(col color)
	get(x int) color
}

// Generic off pattern
type offPattern struct{}

func (_ *offPattern) tick(ledCount int) {}
func (_ *offPattern) reset(col color)   {}
func (_ *offPattern) get(x int) color   { return color{} }

// Dim soft lighting for getting up late at night for bathroom
type nightLight struct{}

func (_ *nightLight) tick(ledCount int) {}
func (_ *nightLight) reset(col color)   {}
func (_ *nightLight) get(x int) color {
	return color{
		r: 120,
		g: 64,
		b: 0,
	}
}

// Cool sin wave thing that I came up with and reuse in every light strip project I do
type sinWavePattern struct {
	offset float64
	rate   float64
	state  []colorF
}

func (pat *sinWavePattern) tick(ledCount int) {
	for i := range ledCount {
		red := (math.Cos(float64(i)/pat.rate+pat.offset) + 1) / 2 * 255
		green := (math.Sin((float64(ledCount-i)/pat.rate + pat.offset)) + 1) / 2 * 255
		pat.state[i] = colorF{red, green, 30}
	}
	pat.offset += 0.04
}

func (_ *sinWavePattern) reset(col color)   {}
func (pat *sinWavePattern) get(x int) color { return pat.state[x].toColor() }

// Twinkling star pattern
type stars struct {
	lastCreate time.Time

	createTime      time.Duration
	decayRate       float32
	backgroundColor color
	state           []colorF
}

func (pat *stars) tick(ledCount int) {
	if time.Since(pat.lastCreate) > pat.createTime {
		pat.lastCreate = time.Now()

		pos := rand.Intn(ledCount)
		pat.state[pos].g = 255
		pat.state[pos].b = 220
		pat.state[pos].r = 30
	}

	for i := range ledCount {
		pat.state[i].diff(pat.backgroundColor, float64(pat.decayRate))
	}
}

func (_ *stars) reset(col color)   {}
func (pat *stars) get(x int) color { return pat.state[x].toColor() }
