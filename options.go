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
	pot         Pot
	baseX       int
	baseY       int
	align       align
	msg         string // TODO: when not set by flag, try to read from stdin
	msgX        int
	msgY        int
	leaves      []string
	multiplier  int
	life        int
	print       bool
}

func options() opts {
	pflag.CommandLine.SortFlags = false

	var o opts
	pflag.BoolVarP(&o.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&o.time, "time", "t", 33*time.Millisecond, "in live mode, delay between steps of growth")
	pflag.BoolVarP(&o.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&o.wait, "wait", "w", 4*time.Second, "in infinite mode, delay between each tree")
	pflag.BoolVarP(&o.screensaver, "screensaver", "S", false, "screensaver mode: equivalent to -li and quit on any keypress")
	pot := pflag.IntP("base", "b", 1, "base pot: big=1 small=2")
	pflag.IntVarP(&o.baseX, "base-x", "", 0, "column position of upper-left corner of plant base pot")
	pflag.IntVarP(&o.baseY, "base-y", "", 0, "row position of upper-left corner of plant base pot")
	alignRaw := pflag.IntP("align", "a", int(center), "align tree: "+
		"center="+strconv.Itoa(int(center))+" "+
		"left="+strconv.Itoa(int(left))+" "+
		"right="+strconv.Itoa(int(right)))
	pflag.StringVarP(&o.msg, "msg", "m", "", "attach message next to the tree")
	pflag.IntVarP(&o.msgX, "msg-x", "", 4, "column position of upper-left corner of message text")
	pflag.IntVarP(&o.msgY, "msg-y", "", 2, "row position of upper-left corner of message text")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited leaves")
	pflag.IntVarP(&o.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&o.life, "life", "L", 32, "life higher -> more growth (0-127)")
	pflag.BoolVarP(&o.print, "print", "p", false, "print first tree to stdout and exit immediately")
	seed := pflag.Int64P("seed", "s", 0, "seed random number generator (default random)")
	help := pflag.BoolP("help", "h", false, "show help")
	version := pflag.BoolP("version", "V", false, "show version")
	pflag.Parse()

	if *help {
		fmt.Println("A bonsai tree generator")
		fmt.Println("\nUsage:\n" +
			"  gobonsai [flags]")
		fmt.Println("\nExamples:\n" +
			"  gobonsai -p --seed 42\n" +
			"  gobonsai -l -L 48 -M 3\n" +
			"  gobonsai --msg \"hi\" --msg-y 20\n" +
			"  gobonsai -S -c \"&,@,§,$,%,☘️,🌿,🍎,💚,🟢,🟩\"")
		fmt.Printf("\nFlags:\n%s", pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	if *version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	if o.screensaver {
		o.live = true
		o.infinite = true
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

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rand = random.New(random.NewSource(*seed))

	return o
}
