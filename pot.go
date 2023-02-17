package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type Pot func(v *gocui.View) error

var bigPot = func(v *gocui.View) error {
	pw, ph := 31, 4
	px, py := potPos(v, pw, ph)

	// TODO: attributes
	// TODO: errorhandling

	//AttrOn(A_BOLD | ColorPair(8))
	v.SetWritePos(px, py)
	fmt.Fprintf(v, ":")
	//AttrOn(ColorPair(2))
	fmt.Fprintf(v, "___________")
	//AttrOn(ColorPair(11))
	fmt.Fprintf(v, "./~~~\\.")
	//AttrOn(ColorPair(2))
	fmt.Fprintf(v, "___________")
	//AttrOn(ColorPair(8))
	fmt.Fprintf(v, ":")
	v.SetWritePos(px, py+1)
	fmt.Fprintf(v, " \\                           / ")
	v.SetWritePos(px, py+2)
	fmt.Fprintf(v, "  \\_________________________/ ")
	v.SetWritePos(px, py+3)
	fmt.Fprintf(v, "  (_)                     (_)")
	//AttrOff(A_BOLD)

	// tree grows from here upwards
	x, y := px+(pw/2), py-1
	v.SetWritePos(x, y)

	return nil
}

var smallPot = func(v *gocui.View) error {
	pw, ph := 15, 3
	px, py := potPos(v, pw, ph)

	//AttrOn(ColorPair(8))
	v.SetWritePos(px, py)
	fmt.Fprintf(v, "(")
	//AttrOn(ColorPair(2))
	fmt.Fprintf(v, "---")
	//AttrOn(ColorPair(11))
	fmt.Fprintf(v, "./~~~\\.")
	//AttrOn(ColorPair(2))
	fmt.Fprintf(v, "---")
	//AttrOn(ColorPair(8))
	fmt.Fprintf(v, ")")
	v.SetWritePos(px, py+1)
	fmt.Fprintf(v, " (           ) ")
	v.SetWritePos(px, py+2)
	fmt.Fprintf(v, "  (_________)  ")
	//AttrOff(A_BOLD)

	// tree grows from here upwards
	x, y := px+(pw/2), py-1
	v.SetWritePos(x, y)

	return nil
}

func potPos(v *gocui.View, pw int, ph int) (x int, y int) {
	vw, vh := v.Size()
	x = (vw / 2) - (pw / 2)
	y = vh - ph
	return x, y
}
