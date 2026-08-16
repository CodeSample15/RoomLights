package lights

type Pattern interface {
	tick()
	reset(color Color)
	get(x int) Color
}

type pattern struct{
	state []Color
}

var (
	offPattern Pattern
)