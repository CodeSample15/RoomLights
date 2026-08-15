package lights

type Pattern interface {
	tick()
	reset(color Color)
	render(led *strip)
}

type pattern struct{
	state []Color
}

type NightLight pattern