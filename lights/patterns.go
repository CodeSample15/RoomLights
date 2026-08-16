package lights

type Pattern interface {
	tick()
	reset(col color)
	get(x int) color
}

// Generic off pattern
type offPattern struct{}

func (_ *offPattern) tick() {}
func (_ *offPattern) reset(col color) {}
func (_ *offPattern) get(x int) color { return color{} }

// Dim soft lighting for getting up late at night for bathroom
type nightLight struct{}

func (_ *nightLight) tick() {}
func (_ *nightLight) reset(col color) {}
func (_ *nightLight) get(x int) color {
	return color{
		r: 120,
		g: 64, 
		b: 0,
	}
}

