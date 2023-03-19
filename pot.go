package main

import (
	"github.com/gdamore/tcell/v2"
)

type Pot struct {
	w int
	h int
	d func(sc *screen, px int, py int)
}

// draw pot and set cursor to tree start pos
func (p Pot) draw(sc *screen) {
	px, py := potPos(sc, p.w, p.h)

	sc.x = px
	sc.y = py

	p.d(sc, px, py)

	// tree grows from here upwards
	sc.x = px + (p.w / 2)
	sc.y = py - 1
}

var bigPot = Pot{
	w: 31,
	h: 4,
	d: func(sc *screen, px int, py int) {
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
	},
}

var smallPot = Pot{
	w: 15,
	h: 3,
	d: func(sc *screen, px int, py int) {
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
	},
}

func potPos(sc *screen, pw int, ph int) (x int, y int) {
	vw, vh := sc.Size()
	x = (vw / 2) - (pw / 2)
	y = vh - ph
	return x, y
}
