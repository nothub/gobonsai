package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	random "math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type opts struct {
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
	align       align
	print       bool
	help        bool
}

func options() opts {
	pflag.CommandLine.SortFlags = false
	// TODO: sort flags

	var o opts
	pot := pflag.IntP("base", "b", 1, "base pot: big=1 small=2")
	seed := pflag.Int64P("seed", "s", time.Now().UnixNano(), "seed random number generator")
	pflag.BoolVarP(&o.help, "help", "h", false, "show help")
	pflag.IntVarP(&o.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&o.life, "life", "L", 32, "life higher -> more growth (0-200)")
	pflag.BoolVarP(&o.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&o.wait, "wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	pflag.BoolVarP(&o.screensaver, "screensaver", "S", false, "screensaver mode: equivalent to -li and quit on any keypress")
	/* TODO:
	   -T, --textOrigin=y,x   Display text from STDIN at row Y, column X
	   -y, --baseY=INT        Row of the upper-left corner of the plant base
	   -x, --baseX=INT        Column of the upper-left corner of the plant base
	   -p, --print            Print tree to terminal when finished
	*/
	pflag.BoolVarP(&o.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&o.time, "time", "t", 50*time.Millisecond, "in live mode, wait TIME secs between steps of growth (must be larger than 0)")
	pflag.StringVarP(&o.message, "message", "m", "", "attach message next to the tree")
	pflag.BoolVarP(&o.print, "print", "p", false, "print tree to terminal when finished")
	alignRaw := pflag.IntP("align", "a", 0, "align tree: center=0 left=1 right=2")
	pflag.Parse()

	if o.help {
		fmt.Println(pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	switch *pot {
	case 1:
		o.pot = bigPot
	case 2:
		o.pot = smallPot
	default:
		log.Panicln("unknown pot type", strconv.Itoa(*pot))
	}

	o.align = align(*alignRaw)

	o.leaves = strings.Split(*leavesRaw, ",")

	if o.screensaver {
		o.live = true
		o.infinite = true
	}

	rand = random.New(random.NewSource(*seed))

	return o
}
