package main

import (
	"errors"
	"fmt"
	"log"
	random "math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/spf13/pflag"
)

var rand *random.Rand

var treeRunes = []string{"🌳", "🌲", "🌴", "🎄", "🎋", "🥦", "🌱"}

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
	seed        int
	help        bool
	usage       string
}

func flags() options {
	pflag.CommandLine.SortFlags = false

	var opts options
	pot := pflag.IntP("pot", "P", 1, "plant pot to use, big: 1, small: 2")
	pflag.IntVarP(&opts.seed, "seed", "s", 0, "seed random number generator")
	pflag.BoolVarP(&opts.help, "help", "h", false, "show help")
	// TODO:
	pflag.BoolVarP(&opts.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&opts.time, "time", "t", 50*time.Millisecond, "in live mode, wait TIME secs between steps of growth (must be larger than 0)")
	pflag.BoolVarP(&opts.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&opts.wait, "wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	pflag.BoolVarP(&opts.screensaver, "screensaver", "S", false, "screensaver mode equivalent to -liWC and quit on any keypress")
	pflag.StringVarP(&opts.message, "message", "m", "", "attach message next to the tree")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	pflag.IntVarP(&opts.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&opts.life, "life", "L", 32, "life higher -> more growth (0-200)")
	pflag.BoolVarP(&opts.print, "print", "p", false, "print tree to terminal when finished")
	pflag.Parse()

	switch *pot {
	case 1:
		opts.pot = bigPot
	case 2:
		opts.pot = smallPot
	default:
		log.Fatalln("unknown pot type", strconv.Itoa(*pot))
	}

	opts.leaves = strings.Split(*leavesRaw, ",")
	opts.usage = pflag.CommandLine.FlagUsages()

	return opts
}

func roll(mod int) int {
	return rand.Int() % mod
}

func main() {
	opts = flags()

	if opts.help {
		fmt.Println(opts.usage)
		return
	}

	if opts.seed <= 0 {
		rand = random.New(random.NewSource(time.Now().Unix()))
	} else {
		rand = random.New(random.NewSource(int64(opts.seed)))
	}

	drawPot := opts.pot

	ui, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ui.Close()

	ui.SetManagerFunc(func(g *gocui.Gui) error {
		w, h := g.Size()
		_, err := g.SetView("main", 0, 0, w-1, h-1, 0)
		if err != nil {
			if errors.Is(err, gocui.ErrUnknownView) {
				_, err := g.SetCurrentView("main")
				if err != nil {
					return err
				}
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
		log.Fatalln(err.Error())
	}

	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			ui.Update(func(g *gocui.Gui) error {
				v, err := g.View("main")
				if err != nil {
					return err
				}

				v.Clear()

				err = v.SetWritePos(8, 4)
				if err != nil {
					return err
				}

				_, err = fmt.Fprintln(v, "Hello, World!", strconv.Itoa(rand.Int()))
				if err != nil {
					return err
				}

				err = drawPot(v)
				if err != nil {
					return err
				}

				// TODO: tree
				fmt.Fprintf(v, treeRunes[rand.Intn(len(treeRunes))])

				return nil
			})
			<-t.C
		}
	}()

	err = ui.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatalln(err.Error())
	}
}
