package lights

type Pattern interface {
	tick()
	reset(color Color)
	get(x int) Color
}

type offPattern struct{}

func (_ *offPattern) tick() {}
func (_ *offPattern) reset(color Color) {}
func (_ *offPattern) get(x int) Color { return Color{} }