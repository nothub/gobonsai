package main

import (
	"errors"
	"fmt"
	"log"
	random "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/pflag"
)

var rand *random.Rand

var opts options

type options struct {
	live        bool
	time        time.Duration
	infinite    bool
	wait        time.Duration
	screensaver bool
	message     string
	pot         Pot
	leaves      []string
	multiplier  int
	life        int
	print       bool
	seed        int64
	help        bool
	usage       string
}

func flags() options {
	pflag.CommandLine.SortFlags = false

	var opts options
	pot := pflag.IntP("pot", "P", 1, "plant pot to use, big: 1, small: 2")
	pflag.Int64VarP(&opts.seed, "seed", "s", time.Now().UnixNano(), "seed random number generator")
	pflag.BoolVarP(&opts.help, "help", "h", false, "show help")
	pflag.BoolVarP(&opts.print, "print", "p", false, "print tree to terminal when finished")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	pflag.IntVarP(&opts.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&opts.life, "life", "L", 32, "life higher -> more growth (0-200)")
	// TODO:
	pflag.BoolVarP(&opts.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&opts.time, "time", "t", 50*time.Millisecond, "in live mode, wait TIME secs between steps of growth (must be larger than 0)")
	pflag.BoolVarP(&opts.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&opts.wait, "wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	pflag.BoolVarP(&opts.screensaver, "screensaver", "S", false, "screensaver mode equivalent to -liWC and quit on any keypress")
	pflag.StringVarP(&opts.message, "message", "m", "", "attach message next to the tree")
	pflag.Parse()

	switch *pot {
	case 1:
		opts.pot = bigPot
	case 2:
		opts.pot = smallPot
	default:
		log.Panicln("unknown pot type", strconv.Itoa(*pot))
	}

	opts.leaves = strings.Split(*leavesRaw, ",")
	opts.usage = pflag.CommandLine.FlagUsages()

	return opts
}

func main() {
	opts = flags()

	if opts.help {
		fmt.Println(opts.usage)
		return
	}

	rand = random.New(random.NewSource(opts.seed))

	sc, sh := newScreen()
	defer sh()

	// TODO: color schemes?
	colors := sc.Colors()
	listColors(sc, colors)

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
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}

	os.Exit(0)

	drawPot := opts.pot // TODO generalize the pot func and use type for sizes struct

	ui, err := gocui.NewGui(gocui.Output256, true)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer ui.Close()

	ui.SetManagerFunc(func(g *gocui.Gui) error {
		w, h := g.Size()
		_, err := g.SetView("main", 0, 0, w-1, h-1, 0)
		if err != nil {
			if errors.Is(err, gocui.ErrUnknownView) {
				v, err := g.SetCurrentView("main")
				if err != nil {
					return err
				}
				v.Frame = false
			} else {
				return err
			}
		}
		return nil
	})

	err = ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		log.Panicln(err.Error())
	}

	var out string

	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			ui.Update(func(g *gocui.Gui) error {
				v, err := g.View("main")
				if err != nil {
					return err
				}

				v.Clear()

				/* print seed:
				   				_, y := v.Size()
				   				err = v.SetWritePos(2, y-2)
				   				if err != nil {
				   				    return err
				   				}
				   				_, err = fmt.Fprintf(v, "Seed: %s", strconv.Itoa(int(opts.seed)))
				   				if err != nil {
				   				    return err
				                   }
				*/

				err = drawPot(v)
				if err != nil {
					return err
				}

				potHeight := 3

				err = drawTree(v, opts, potHeight)
				if err != nil {
					return err
				}

				if opts.print {
					out = fmt.Sprint(v.Buffer())
					g.Close()
					for _, l := range strings.Split(out, "\n") {
						if len(strings.TrimSpace(l)) > 0 {
							fmt.Println(l)
						}
					}
					return gocui.ErrQuit
				}

				return nil
			})
			<-t.C
		}
	}()

	err = ui.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err.Error())
	}
}
