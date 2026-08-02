package lights

type Pattern interface {
	tick()
	reset()
}
