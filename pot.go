package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type Pot func(v *gocui.View) error

var bigPot = func(v *gocui.View) error {
	pw, ph := 31, 4
	px, py := potPos(v, pw, ph)

	v.SetWritePos(px, py)
	fmt.Fprintf(v, whiteBold.Sprint(":"))
	fmt.Fprintf(v, greenBold.Sprint("___________"))
	fmt.Fprintf(v, yellowBold.Sprint("./~~~\\."))
	fmt.Fprintf(v, greenBold.Sprint("___________"))
	fmt.Fprintf(v, whiteBold.Sprint(":"))
	v.SetWritePos(px, py+1)
	fmt.Fprintf(v, " \\                           / ")
	v.SetWritePos(px, py+2)
	fmt.Fprintf(v, "  \\_________________________/ ")
	v.SetWritePos(px, py+3)
	fmt.Fprintf(v, "  (_)                     (_)")

	// tree grows from here upwards
	x, y := px+(pw/2), py-1
	err := v.SetWritePos(x, y)
	if err != nil {
		return err
	}

	return nil
}

var smallPot = func(v *gocui.View) error {
	pw, ph := 15, 3
	px, py := potPos(v, pw, ph)

	v.SetWritePos(px, py)
	fmt.Fprintf(v, whiteBold.Sprint("("))
	fmt.Fprintf(v, greenBold.Sprint("---"))
	fmt.Fprintf(v, yellowBold.Sprint("./~~~\\."))
	fmt.Fprintf(v, greenBold.Sprint("---"))
	fmt.Fprintf(v, whiteBold.Sprint(")"))
	v.SetWritePos(px, py+1)
	fmt.Fprintf(v, " (           ) ")
	v.SetWritePos(px, py+2)
	fmt.Fprintf(v, "  (_________)  ")

	// tree grows from here upwards
	x, y := px+(pw/2), py-1
	err := v.SetWritePos(x, y)
	if err != nil {
		return err
	}

	return nil
}

func potPos(v *gocui.View, pw int, ph int) (x int, y int) {
	vw, vh := v.Size()
	x = (vw / 2) - (pw / 2)
	y = vh - ph
	return x, y
}
