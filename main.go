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
				sc.drawMessage(opts.msg, opts.msgX, opts.msgY)
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

				// Store the tree for printing later,
				// when screen cleanup is finished.
				out = strings.Join(trimmed, "\n")

				// Send the quit event,
				evQuit(sc)
				// then wait for shutdown.
				<-shutdown
				break
			}

			if opts.infinite {
				// When in infinite mode, we either
				// await the delay or wait for shutdown.
				select {
				case <-shutdown:
					break
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
