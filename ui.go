package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type screen struct {
	tcell.Screen
	x, y int // cursor
}

func (sc *screen) draw(text string, style tcell.Style) {
	w, h := sc.Size()

	for _, r := range []rune(text) {
		rw := runewidth.RuneWidth(r)

		// skip oob writes
		if sc.x+rw >= w || sc.y >= h {
			// TODO: error handling required?
			continue
		}

		sc.put(r, style)
	}
}

func (sc *screen) put(r rune, style tcell.Style) {
	if !active {
		// stop drawing while shutdown to
		// avoid a data race with screen.Fini()
		return
	}

	w := runewidth.RuneWidth(r)

	// skip all non-printable
	if w < 1 {
		return
	}

	sc.SetContent(sc.x, sc.y, r, nil, style)
	sc.x = sc.x + w
}

func newScreen() (sc *screen) {
	tsc, err := tcell.NewScreen()
	if err != nil {
		log.Panicln(err.Error())
	}

	err = tsc.Init()
	if err != nil {
		log.Panicln(err.Error())
	}

	tsc.DisablePaste()
	tsc.DisableMouse()

	tsc.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))

	tsc.Clear()

	return &screen{Screen: tsc}
}
