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
				evQuit(sc)
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
