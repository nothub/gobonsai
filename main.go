package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	random "math/rand"
	"strings"
	"time"
)

var rand *random.Rand
var shutdown = make(chan bool)
var active = true // TODO: get rid of this to avoid race conditions

func main() {
	opts := options()
	sc := newScreen()

	var out string
	quit := func() {
		p := recover()
		sc.Fini()
		if p != nil {
			panic(p)
		}
		if out != "" {
			fmt.Println(out)
		}
	}
	defer quit()

	go func() {
		for {
			if !active {
				<-shutdown
				break
			}

			sc.Clear()

			// default center position
			px, py := opts.pot.ulPos(sc)

			// align by moving 1/4 screen size
			switch opts.align {
			case left:
				sw, _ := sc.Size()
				px = px - (sw / 4)
			case right:
				sw, _ := sc.Size()
				px = px + (sw / 4)
			}

			// user defined position overrides
			if opts.baseX != 0 {
				px = opts.baseX
			}
			if opts.baseY != 0 {
				py = opts.baseY
			}

			// draw pot
			opts.pot.draw(sc, px, py)

			// draw from pot upwards
			err := drawTree(sc, opts)
			if err != nil {
				log.Panicln(err.Error())
			}

			// draw message box
			if opts.msg != "" {
				sc.drawMessage(opts.msg, opts.msgX, opts.msgY)
			}

			// emit drawn event to trigger screen refresh
			evDrawn(sc)

			if opts.print {
				// Store the tree for printing later,
				// when screen cleanup is finished.
				var tree []string
				w, h := sc.Size()
				for y := 0; y < h; y++ {
					var sb strings.Builder
					for x := 0; x < w; x++ {
						r, _, _, _ := sc.GetContent(x, y)
						sb.WriteRune(r)
					}
					s := sb.String()
					// Ignore empty space above tree.
					if strings.TrimSpace(s) != "" {
						tree = append(tree, s)
					}
				}
				out = strings.Join(tree, "\n")

				// Send the quit event,
				evQuit(sc)
				// then wait for shutdown.
				<-shutdown
				break
			}

			if opts.infinite {
				// We either await the delay or wait for shutdown.
				select {
				case <-shutdown:
					return
				case <-time.After(opts.wait):
				}
			} else {
				// When not in infinite mode, we just
				// draw 1 tree and wait for shutdown.
				<-shutdown
			}
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
			active = false

			// signal shutdown to main loop
			shutdown <- true

			// we can just exit here, the shutdown hook will clean up the terminal
			return
		}
	}
}
