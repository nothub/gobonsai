package main

import (
	"fmt"
	"log"
	random "math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type opts struct {
	live        bool
	time        time.Duration
	infinite    bool
	wait        time.Duration
	screensaver bool
	message     string // TODO: when not set by flag, try to read from stdin
	pot         Pot
	leaves      []string
	multiplier  int
	life        int
	align       align
	baseX       uint
	baseY       uint
	msgX        uint
	msgY        uint
	print       bool
	noColor     bool
	help        bool
}

func options() opts {
	pflag.CommandLine.SortFlags = false
	// TODO: sort flags

	var o opts
	pflag.BoolVarP(&o.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&o.time, "time", "t", 33*time.Millisecond, "in live mode, wait TIME between steps of growth")
	pflag.BoolVarP(&o.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&o.wait, "wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	pflag.BoolVarP(&o.screensaver, "screensaver", "S", false, "screensaver mode: equivalent to -li and quit on any keypress")
	pot := pflag.IntP("base", "b", 1, "base pot: big=1 small=2")
	pflag.UintVarP(&o.baseX, "base-x", "", 0, "column position of upper-left corner of plant base pot")
	pflag.UintVarP(&o.baseY, "base-y", "", 0, "row position of upper-left corner of plant base pot")
	alignRaw := pflag.IntP("align", "a", int(center), "align tree: center=0 left=1 right=2")
	pflag.StringVarP(&o.message, "message", "m", "", "attach message next to the tree")
	pflag.UintVarP(&o.msgX, "message-x", "", 0, "column position of upper-left corner of message text")
	pflag.UintVarP(&o.msgY, "message-y", "", 0, "row position of upper-left corner of message text")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	pflag.IntVarP(&o.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&o.life, "life", "L", 32, "life higher -> more growth (0-200)")
	pflag.BoolVarP(&o.print, "print", "p", false, "print first tree to stdout and exit immediately")
	pflag.BoolVarP(&o.noColor, "no-color", "n", false, "disable all colors")
	seed := pflag.Int64P("seed", "s", 42, "seed random number generator")
	pflag.BoolVarP(&o.help, "help", "h", false, "show help")
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

	// TODO: use align to set base-x and base-y values relative to window-size
	o.align = align(*alignRaw)

	o.leaves = strings.Split(*leavesRaw, ",")

	if o.screensaver {
		o.live = true
		o.infinite = true
	}

	rand = random.New(random.NewSource(*seed))

	return o
}
