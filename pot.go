package main

import (
	"github.com/gdamore/tcell/v2"
)

// TODO: generalize the pot func and use type for sizes struct
type Pot func(sc *screen) error

var bigPot = func(sc *screen) error {
	pw, ph := 31, 4
	px, py := potPos(sc, pw, ph)

	sc.x = px
	sc.y = py

	sc.draw(":", styleWhiteBold)
	sc.draw("___________", styleGreenBold)
	sc.draw("./~~~\\.", styleBrown)
	sc.draw("___________", styleGreenBold)
	sc.draw(":", styleWhiteBold)
	sc.x = px
	sc.y = py + 1
	sc.draw(" \\                           / ", styleDefault)
	sc.x = px
	sc.y = py + 2
	sc.draw("  \\_________________________/ ", styleDefault)
	sc.x = px
	sc.y = py + 3
	sc.draw("  (_)                     (_)", styleDefault)

	// tree grows from here upwards
	sc.x = px + (pw / 2)
	sc.y = py - 1

	return nil
}

var smallPot = func(sc *screen) error {
	pw, ph := 15, 3
	px, py := potPos(sc, pw, ph)

	sc.x = px
	sc.y = py

	sc.draw("(", styleWhiteBold)
	sc.draw("---", styleGreenBold)
	sc.draw("./~~~\\.", styleBrown)
	sc.draw("---", styleGreenBold)
	sc.draw(")", styleWhiteBold)
	sc.x = px
	sc.y = py + 1
	sc.draw(" (           ) ", tcell.StyleDefault)
	sc.x = px
	sc.y = py + 2
	sc.draw("  (_________)  ", tcell.StyleDefault)

	// tree grows from here upwards
	sc.x = px + (pw / 2)
	sc.y = py - 1

	return nil
}

func potPos(sc *screen, pw int, ph int) (x int, y int) {
	vw, vh := sc.Size()
	x = (vw / 2) - (pw / 2)
	y = vh - ph
	return x, y
}
