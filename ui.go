package main

import (
	"log"
	"strings"

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
			continue
		}

		sc.put(r, style)
	}
}

func (sc *screen) drawMessage(msg string, x int, y int) {
	// upper border
	sc.x = x
	sc.y = y
	sc.draw("+"+strings.Repeat("-", len(msg)+2)+"+", styleGray)
	// center with message and front- and back-border
	sc.x = x
	sc.y = y + 1
	sc.draw("| ", styleGray)
	sc.draw(msg, styleDefault)
	sc.draw(" |", styleGray)
	// lower border
	sc.x = x
	sc.y = y + 2
	sc.draw("+"+strings.Repeat("-", len(msg)+2)+"+", styleGray)
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
