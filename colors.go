package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"strconv"
)

var (
	styleDefault   = tcell.StyleDefault
	styleBrown     = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.Color94)
	styleBrownBold = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.Color94).Bold(true)
	styleGreen     = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	styleGreenBold = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen).Bold(true)
	styleWhiteBold = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite).Bold(true)
)

// TODO: check supported number of colors
// TODO: color schemes? (8-col, original cbonsai, monochrome)

// based on type of tree, determine what color a branch should be
func chooseColor(kind branch) tcell.Style {
	switch kind {
	case dying:
		if rand.Int()%10 == 0 {
			return styleGreenBold
		} else {
			return styleGreen
		}

	case dead:
		if rand.Int()%3 == 0 {
			return styleGreenBold
		} else {
			return styleGreen
		}

		// trunk | shootLeft | shootRight
	default:
		if rand.Int()%2 == 0 {
			return styleBrownBold
		} else {
			return styleBrown
		}
	}
}

func listColors(sc *screen) {
	// num := sc.Colors()

	sc.x = 0
	sc.y = 0
	c := 0

	for col := tcell.ColorValid; col < tcell.ColorYellowGreen; col++ {
		style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(col).Bold(true)
		sc.draw(strconv.FormatUint(uint64(col), 10), style)
		sc.x++
		c++
		if c%8 == 0 {
			sc.x = 0
			sc.y++
		}
	}

	err := sc.PostEvent(EvDrawn())
	if err != nil {
		log.Panicln(err.Error())
	}
}
