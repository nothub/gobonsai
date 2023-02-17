package main

import (
	"errors"
	"fmt"
	"log"
	random "math/rand"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/spf13/pflag"
)

var rand *random.Rand

type branch int

const (
	trunk branch = iota
	shootLeft
	shootRight
	dying
	dead
)

var opts options

type options struct {
	live        *bool
	time        *time.Duration
	infinite    *bool
	wait        *time.Duration
	screensaver *bool
	message     *string
	base        *int
	leaves      []string
	multiplier  *int
	life        *int
	print       *bool
	seed        *int
	save        *string
	load        *string
	verbose     *bool
	help        *bool
	usage       string
}

func flags() *options {
	pflag.CommandLine.SortFlags = false

	var opts options
	opts.live = pflag.BoolP("live", "l", false, "live mode: show each step of growth")
	opts.time = pflag.DurationP("time", "t", 30*time.Millisecond, "in live mode, wait TIME secs between steps of growth (must be larger than 0)")
	opts.infinite = pflag.BoolP("infinite", "i", false, "infinite mode: keep growing trees")
	opts.wait = pflag.DurationP("wait", "w", 4*time.Second, "in infinite mode, wait TIME between each tree generation")
	opts.screensaver = pflag.BoolP("screensaver", "S", false, "screensaver mode equivalent to -liWC and quit on any keypress")
	opts.message = pflag.StringP("message", "m", "", "attach message next to the tree")
	opts.base = pflag.IntP("base", "b", 1, "ascii-art plant base to use, big: 1, small: 2")
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	opts.multiplier = pflag.IntP("multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	opts.life = pflag.IntP("life", "L", 32, "life higher -> more growth (0-200)")
	opts.print = pflag.BoolP("print", "p", false, "print tree to terminal when finished")
	opts.seed = pflag.IntP("seed", "s", 0, "seed random number generator")
	opts.help = pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	opts.leaves = strings.Split(*leavesRaw, ",")
	opts.usage = pflag.CommandLine.FlagUsages()

	return &opts
}

func main() {
	rand = random.New(random.NewSource(time.Now().UnixNano()))

	opts = *flags()
	if *opts.help {
		fmt.Println(opts.usage)
		return
	}

	ui, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ui.Close()

	ui.SetManagerFunc(func(g *gocui.Gui) error {
		width, height := g.Size()
		v, err := g.SetView("hello", 0, 0, width-1, height-1, 0)
		if err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}

			if _, err := g.SetCurrentView("hello"); err != nil {
				return err
			}

			err := v.SetWritePos(16, 8)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(v, "Hello world!")
			if err != nil {
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

	err = ui.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatalln(err.Error())
	}
}

// roll (randomize) a given die
func roll(mod int) int {
	return rand.Int() % mod
}

// determine change in X and Y coordinates of a given branch
func deltas(branchType branch, life int, age int, multiplier int, returnDx *int, returnDy *int) {
	dx := 0
	dy := 0

	switch branchType {
	case trunk:
		if age <= 2 || life < 4 {
			// new or dead trunk
			dy = 0
			dx = (rand.Int() % 3) - 1

		} else if age < (multiplier * 3) {
			// young trunk should grow wide
			// every (multiplier * 0.8) steps, raise tree to next level
			if age%(multiplier/2) == 0 {
				dy = -1
			} else {
				dy = 0
			}

			dice := roll(10)
			if dice >= 0 && dice <= 0 {
				dx = -2
			} else if dice >= 1 && dice <= 3 {
				dx = -1
			} else if dice >= 4 && dice <= 5 {
				dx = 0
			} else if dice >= 6 && dice <= 8 {
				dx = 1
			} else if dice >= 9 && dice <= 9 {
				dx = 2
			}

		} else {
			// middle-aged trunk
			dice := roll(10)
			if dice > 2 {
				dy = -1
			} else {
				dy = 0
			}
			dx = (rand.Int() % 3) - 1
		}

	case shootLeft: // trend left and little vertical movement
		dice := roll(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 7 {
			dy = 0
		} else if dice >= 8 && dice <= 9 {
			dy = 1
		}

		dice = roll(10)
		if dice >= 0 && dice <= 1 {
			dx = -2
		} else if dice >= 2 && dice <= 5 {
			dx = -1
		} else if dice >= 6 && dice <= 8 {
			dx = 0
		} else if dice >= 9 && dice <= 9 {
			dx = 1
		}

	case shootRight: // trend right and little vertical movement
		dice := roll(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 7 {
			dy = 0
		} else if dice >= 8 && dice <= 9 {
			dy = 1
		}

		dice = roll(10)
		if dice >= 0 && dice <= 1 {
			dx = 2
		} else if dice >= 2 && dice <= 5 {
			dx = 1
		} else if dice >= 6 && dice <= 8 {
			dx = 0
		} else if dice >= 9 && dice <= 9 {
			dx = -1
		}

	case dying: // discourage vertical growth(?) trend left/right (-3,3)
		dice := roll(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 8 {
			dy = 0
		} else if dice >= 9 && dice <= 9 {
			dy = 1
		}

		dice = roll(15)
		if dice >= 0 && dice <= 0 {
			dx = -3
		} else if dice >= 1 && dice <= 2 {
			dx = -2
		} else if dice >= 3 && dice <= 5 {
			dx = -1
		} else if dice >= 6 && dice <= 8 {
			dx = 0
		} else if dice >= 9 && dice <= 11 {
			dx = 1
		} else if dice >= 12 && dice <= 13 {
			dx = 2
		} else if dice >= 14 && dice <= 14 {
			dx = 3
		}

	case dead: // fill in surrounding area
		dice := roll(10)
		if dice >= 0 && dice <= 2 {
			dy = -1
		} else if dice >= 3 && dice <= 6 {
			dy = 0
		} else if dice >= 7 && dice <= 9 {
			dy = 1
		}
		dx = (rand.Int() % 3) - 1
	}

	*returnDx = dx
	*returnDy = dy
}

func leaf(branch branch, life int, dx int, dy int) string {
	s := "?"

	if life < 4 {
		branch = dying
	}

	switch branch {
	case trunk:
		if dy == 0 {
			s = "/~"
		} else if dx < 0 {
			s = "\\|"
		} else if dx == 0 {
			s = "/|\\"
		} else if dx > 0 {
			s = "|/"
		}

	case shootLeft:
		if dy > 0 {
			s = "\\"
		} else if dy == 0 {
			s = "\\_"
		} else if dx < 0 {
			s = "\\|"
		} else if dx == 0 {
			s = "/|"
		} else if dx > 0 {
			s = "/"
		}

	case shootRight:
		if dy > 0 {
			s = "/"
		} else if dy == 0 {
			s = "_/"
		} else if dx < 0 {
			s = "\\|"
		} else if dx == 0 {
			s = "/|"
		} else if dx > 0 {
			s = "/"
		}

	case dying, dead:
		s = opts.leaves[rand.Int()%len(opts.leaves)]
	}

	return s
}
