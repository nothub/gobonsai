package main

import (
	"fmt"
	"log"
	random "math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var rand *random.Rand

func main() {
	opts = flags()

	if opts.help {
		fmt.Println(opts.usage)
		return
	}

	rand = random.New(random.NewSource(opts.seed))

	sc, sh := newScreen()
	defer sh()

	go func() {
		t := time.NewTicker(1 * time.Second)
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

			if opts.print {
				// TODO: print buffer and exit
			}

			// refresh screen
			err = sc.PostEvent(EvDrawn())
			if err != nil {
				log.Panicln(err.Error())
			}

			// wait for tick
			<-t.C
		}
	}()

	for {
		ev := sc.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			// resize event will be emitted once initially
			sc.Sync()

		case *EventDrawn:
			// finished drawing, show changes
			sc.Show()

		case *tcell.EventKey:
			// TODO: handle relevant unix signals
			// ev.Key() == tcell.KeyCtrlC
			if ev.Key() == tcell.KeyEscape {
				return
			}
		}
	}
}
