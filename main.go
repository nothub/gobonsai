package main

import (
	"log"
	random "math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var rand *random.Rand

func main() {
	opts := options()

	sc, sh := newScreen()
	defer sh()

	go func() {
		t := time.NewTicker(opts.wait)
		for {
			sc.Clear()

			// draw pot and set cursor to tree start pos
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
				evQuit(sc)
				break
			}

			// wait for tick
			<-t.C
		}
	}()

	for {
		switch ev := sc.PollEvent().(type) {
		case *tcell.EventResize:
			// resize event will be emitted once initially
			sc.Sync()

		case *tcell.EventKey:
			if opts.screensaver ||
				ev.Key() == tcell.KeyEscape ||
				ev.Key() == tcell.KeyCtrlC ||
				ev.Key() == tcell.KeyCtrlD {
				evQuit(sc)
			}

		case *eventDrawn:
			// finished drawing, show changes
			sc.Show()

		case *eventQuit:
			if opts.print {
				// TODO: print buffer to stdout
				// -p, --print  Print tree to terminal when finished
			}
			// we can just exit here, the shutdown hook will clean up the terminal
			return
		}
	}
}
