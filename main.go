package main

import (
	"fmt"
	"log"
	random "math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

var rand *random.Rand
var shutdown = make(chan bool)
var shouldPrint = true

func main() {
	opts := options()
	sc := newScreen()

	quit := func() {
		var out strings.Builder
		if opts.print {
			// TODO: retain colors when printing
			w, h := sc.Size()
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					r, _, _, _ := sc.GetContent(x, y)
					out.WriteRune(r)
				}
			}
		}
		p := recover()
		sc.Fini()
		if p != nil {
			panic(p)
		}
		fmt.Println(out.String())
	}
	defer quit()

	go func() {
		for {
			sc.Clear()

			// draw pot and set cursor to tree start pos
			// TODO: separate base position search and pot drawing
			// -a --align=INT  Align tree: center=0 left=1 right=2
			err := opts.pot(sc)
			if err != nil {
				log.Panicln(err.Error())
			}

			// draw from pot upwards
			err = drawTree(sc, opts)
			if err != nil {
				log.Panicln(err.Error())
			}

			// TODO: draw message
			// -m, --message=STR  Attach message next to the tree

			// refresh screen
			evDrawn(sc)

			if !opts.infinite {
				// not in infinite (regrowing trees) mode
				// so we just wait here for the shutdown
				<-shutdown
				break
			}

			// chill out a bit
			time.Sleep(opts.wait)
		}
	}()

	for {
		switch ev := sc.PollEvent().(type) {
		case *tcell.EventResize:
			// resize event will be emitted once initially
			sc.Sync()

		case *tcell.EventKey:
			if opts.screensaver {
				evQuit(sc)
				break
			}
			switch ev.Key() {
			case tcell.KeyEscape:
				evQuit(sc)
			case tcell.KeyCtrlC:
				// SIGINT
				evQuit(sc)
			case tcell.KeyCtrlD:
				// SIGQUIT
				evQuit(sc)
			}

		case *eventDrawn:
			// finished drawing, show changes
			sc.Show()

		case *eventQuit:
			// stop drawing while shutdown to
			// avoid a data race with screen.Fini()
			shouldPrint = false

			// signal shutdown to main loop
			shutdown <- true

			// we can just exit here, the shutdown hook will clean up the terminal
			return
		}
	}
}
