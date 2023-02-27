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

			if opts.msg != "" {
				// upper border
				sc.x = opts.msgX
				sc.y = opts.msgY
				sc.draw("+"+strings.Repeat("-", len(opts.msg)+2)+"+", styleGray)
				// center with message and front- and back-border
				sc.x = opts.msgX
				sc.y = opts.msgY + 1
				sc.draw("| ", styleGray)
				sc.draw(opts.msg, styleDefault)
				sc.draw(" |", styleGray)
				// lower border
				sc.x = opts.msgX
				sc.y = opts.msgY + 2
				sc.draw("+"+strings.Repeat("-", len(opts.msg)+2)+"+", styleGray)
			}

			// refresh screen
			evDrawn(sc)

			if opts.print {
				// TODO: retain colors when printing

				// convert screen content to string
				var sb strings.Builder
				w, h := sc.Size()
				for y := 0; y < h; y++ {
					for x := 0; x < w; x++ {
						r, _, _, _ := sc.GetContent(x, y)
						sb.WriteRune(r)
					}
				}

				// trim empty space above tree
				// TODO: fix tree space trimming on print
				var trimmed []string
				split := strings.Split(sb.String(), "\n")
				for _, s := range split {
					t := strings.TrimSpace(s)
					if t != "" {
						trimmed = append(trimmed, s)
					}
				}

				// store output for printing later when screen cleanup is done
				out = strings.Join(trimmed, "\n")

				// send quit event
				evQuit(sc)
				// wait for shutdown signal
				<-shutdown
				break
			}

			if !opts.infinite {
				// not in infinite (regrowing trees) mode
				// so we just block here until shutdown
				<-shutdown
				break
			}

			switch {
			case <-shutdown:
				// break loop when supposed to shutdown
				break
			default:
				// chill out a bit
				time.Sleep(opts.wait)
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
