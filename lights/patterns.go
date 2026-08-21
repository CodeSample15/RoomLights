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
type nightLight struct {
	nextPixelTime time.Duration

	lastPixel time.Time
	currPixel int
}

func (pat *nightLight) tick(ledCount int) {
	if time.Since(pat.lastPixel) > pat.nextPixelTime {
		pat.lastPixel = time.Now()
		pat.currPixel++
	}
}

func (pat *nightLight) reset(col color) {
	pat.lastPixel = time.Now()
	pat.currPixel = 0
}

func (pat *nightLight) get(x int) color {
	if x < pat.currPixel {
		return color{
			r: 85,
			g: 70,
			b: 0,
		}
	}

	return color{}
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
	pat.offset += 0.02
}

func (_ *sinWavePattern) reset(col color) {}

func (pat *sinWavePattern) get(x int) color {
	return pat.state[x].toColor()
}

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
		pat.state[pos].g = float64(randRange(100, 255))
		pat.state[pos].b = float64(randRange(50, 250))
		pat.state[pos].r = float64(randRange(0, 90))
	}

	for i := range ledCount {
		pat.state[i].diff(pat.backgroundColor, float64(pat.decayRate))
	}
}

func (_ *stars) reset(col color) {}

func (pat *stars) get(x int) color {
	return pat.state[x].toColor()
}

// Rainbow (many colors but not really a rainbow). Scrolls constantly
type rainbow struct {
	offset float64
}

func (pat *rainbow) tick(ledCount int) {

}

func (_ *rainbow) reset(col color) {}

func (pat *rainbow) get(x int) color {
	sinR := math.Sin(float64(x) / 50)
	sinR = (sinR + 1) / 2

	sinG := math.Sin((float64(x) + math.Pi) / 50)
	sinG = (sinG + 1) / 2

	sinB := math.Sin((float64(x) + math.Pi*2) / 50)
	sinB = (sinB + 1) / 2

	c := colorF{
		r: 255 * sinR,
		g: 255 * sinG,
		b: 255 * sinB,
	}

	return c.toColor()
}

// HELPER FUNCTIONS
func randRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
