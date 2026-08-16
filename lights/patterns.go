package lights

type Pattern interface {
	tick()
	reset(color Color)
	get(x int) Color
}

var (
	offPattern Pattern
)