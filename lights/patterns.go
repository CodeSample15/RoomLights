package lights

type Pattern interface {
	tick()
	reset(color Color)
	get(x int) Color
}

// Generic off pattern
type offPattern struct{}

func (_ *offPattern) tick() {}
func (_ *offPattern) reset(color Color) {}
func (_ *offPattern) get(x int) Color { return Color{} }

// Dim soft lighting for getting up late at night for bathroom
type nightLight struct{}

func (_ *nightLight) tick() {}
func (_ *nightLight) reset(color Color) {}
func (_ *nightLight) get(x int) Color {
	return Color{
		r: 120,
		g: 64, 
		b: 0,
	}
}

