package main

import (
	"github.com/spf13/pflag"
	"log"
	"strconv"
	"strings"
	"time"
)

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

	/* cbonsai flags:

	   Usage: cbonsai [OPTION]...

	   cbonsai is a beautifully random bonsai tree generator.

	   Options:
	     -l, --live             Live mode: show each step of growth
	     -t, --time=TIME        In live mode, wait TIME secs between steps of growth (must be larger than 0) [default: 0.03]
	     -i, --infinite         Infinite mode: keep growing trees
	     -w, --wait=TIME        In infinite mode, wait TIME between each tree generation [default: 4.00]
	     -S, --screensaver      Screensaver mode; equivalent to -li and quit on any keypress
	     -m, --message=STR      Attach message next to the tree
	     -T, --textOrigin=y,x   Display text from STDIN at row Y, column X
	     -b, --base=INT         Ascii-art plant base to use, 0 is none
	     -y, --baseY=INT        Row of the upper-left corner of the plant base
	     -x, --baseX=INT        Column of the upper-left corner of the plant base
	     -c, --leaf=LIST        List of comma-delimited strings randomly chosen for leaves
	     -M, --multiplier=INT   Branch multiplier; higher -> more branching (0-20) [default: 5]
	     -L, --life=INT         Life; higher -> more growth (0-200) [default: 32]
	     -p, --print            Print tree to terminal when finished
	     -s, --seed=INT         Seed random number generator
	     -W, --save=FILE        Save progress to file [default: $XDG_CACHE_HOME/cbonsai or $HOME/.cache/cbonsai]
	     -C, --load=FILE        Load progress from file [default: $XDG_CACHE_HOME/cbonsai]
	     -v, --verbose          Increase output verbosity
	     -h, --help             Show help
	*/

	var opts options
	pot := pflag.IntP("pot", "P", 1, "plant pot to use, big: 1, small: 2")
	pflag.Int64VarP(&opts.seed, "seed", "s", time.Now().UnixNano(), "seed random number generator")
	pflag.BoolVarP(&opts.help, "help", "h", false, "show help")
	pflag.IntVarP(&opts.multiplier, "multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	pflag.IntVarP(&opts.life, "life", "L", 32, "life higher -> more growth (0-200)")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	// TODO:
	pflag.BoolVarP(&opts.live, "live", "l", false, "live mode: show each step of growth")
	pflag.DurationVarP(&opts.time, "time", "t", 50*time.Millisecond, "in live mode, wait TIME secs between steps of growth (must be larger than 0)")
	pflag.BoolVarP(&opts.infinite, "infinite", "i", false, "infinite mode: keep growing trees")
	pflag.DurationVarP(&opts.wait, "wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	pflag.BoolVarP(&opts.screensaver, "screensaver", "S", false, "screensaver mode equivalent to -liWC and quit on any keypress")
	pflag.StringVarP(&opts.message, "message", "m", "", "attach message next to the tree")
	pflag.BoolVarP(&opts.print, "print", "p", false, "print tree to terminal when finished")
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
