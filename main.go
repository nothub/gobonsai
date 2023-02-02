package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/gdamore/tcell/v2"
	nc "github.com/rthornton128/goncurses"
	"github.com/spf13/pflag"
)

type baseType int

const (
	bigPot baseType = iota
	smallPot
)

type branchType int

const (
	trunk branchType = iota
	shootLeft
	shootRight
	dying
	dead
)

type counters struct {
	branches     int
	shoots       int
	shootCounter int
}

var opts options
var tui tuiElements

type options struct {
	live        *bool
	time        *time.Duration
	infinite    *bool
	wait        *time.Duration
	screensaver *bool
	message     *string
	base        *baseType
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

type tuiElements struct {
	baseWin          *nc.Window
	treeWin          *nc.Window
	messageBorderWin *nc.Window
	messageWin       *nc.Window

	basePanel          *nc.Panel
	treePanel          *nc.Panel
	messageBorderPanel *nc.Panel
	messagePanel       *nc.Panel
}

func (e *tuiElements) cleanup() {
	windows := []*nc.Window{e.baseWin, e.treeWin, e.messageBorderWin, e.messageWin}
	for i := range windows {
		err := windows[i].Delete()
		if err != nil {
			log.Println(err)
			//log.Fatal(err)
		}
	}

	panels := []*nc.Panel{e.basePanel, e.treePanel, e.messageBorderPanel, e.messagePanel}
	for i := range panels {
		err := panels[i].Delete()
		if err != nil {
			log.Println(err)
			//log.Fatal(err)
		}
	}
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
	opts.base = (*baseType)(pflag.IntP("base", "b", 1, "ascii-art plant base to use, big: 1, small: 2"))
	leavesRaw := pflag.StringP("leaves", "c", "&", "list of comma-delimited strings randomly chosen for leaves")
	opts.multiplier = pflag.IntP("multiplier", "M", 5, "branch multiplier higher -> more branching (0-20)")
	opts.life = pflag.IntP("life", "L", 32, "life higher -> more growth (0-200)")
	opts.print = pflag.BoolP("print", "p", false, "print tree to terminal when finished")
	opts.seed = pflag.IntP("seed", "s", 0, "seed random number generator")
	opts.save = pflag.StringP("save", "W", "~/.cache/cbonsai", "save progress to file")
	opts.load = pflag.StringP("load", "C", "~/.cache/cbonsai", "load progress from file")
	opts.verbose = pflag.BoolP("verbose", "v", false, "increase output verbosity")
	opts.help = pflag.BoolP("help", "h", false, "show help")

	pflag.Parse()

	opts.leaves = strings.Split(*leavesRaw, ",")
	opts.usage = pflag.CommandLine.FlagUsages()

	return &opts
}

func main() {
	// init
	rand.Seed(time.Now().UnixNano())

	// read options
	opts = *flags()
	if *opts.help {
		fmt.Println(opts.usage)
		return
	}

	// init curses window
	stdscr, err := nc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.End()

	nc.Echo(false)
	nc.CBreak(false)
	// TODO: pr to goncurses: NoDelay(stdscr, bool)
	err = nc.Cursor(0)
	if err != nil {
		log.Fatal(err)
	}

	if nc.HasColors() {
		err := nc.StartColor()
		if err != nil {
			log.Fatal(err)
		}

		// use native background color when possible
		bg := nc.C_BLACK
		if nc.UseDefaultColors() == nil {
			bg = -1
		}

		// define color pairs
		var i int16
		for i = 1; i <= 16; i++ {
			initPair(i, i, bg)
		}

		// restrict color pallete in non-256color terminals (e.g. screen or linux)
		if nc.Colors() < 256 {
			initPair(8, 7, bg) // gray will look white
			initPair(9, 1, bg)
			initPair(10, 2, bg)
			initPair(11, 3, bg)
			initPair(12, 4, bg)
			initPair(13, 5, bg)
			initPair(14, 6, bg)
			initPair(15, 7, bg)
		}
	} else {
		log.Println("Warning: terminal does not have color support.")
	}

	// base pot
	var baseWidth int
	var baseHeight int
	switch *opts.base {
	case bigPot:
		baseWidth = 31
		baseHeight = 4
	case smallPot:
		baseWidth = 15
		baseHeight = 3
	default:
		baseWidth = 0
		baseHeight = 0
	}

	// base position
	rows, cols := stdscr.MaxYX()
	baseOriginY := rows - baseHeight
	baseOriginX := (cols / 2) - (baseWidth / 2)

	tui.cleanup()

	// create windows
	tui.baseWin, err = nc.NewWindow(baseHeight, baseWidth, baseOriginY, baseOriginX)
	if err != nil {
		log.Fatal(err)
	}
	tui.treeWin, err = nc.NewWindow(rows-baseHeight, cols, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// create or replace panels
	if tui.basePanel == nil {
		tui.basePanel = nc.NewPanel(tui.baseWin)
	} else {
		err = tui.basePanel.Replace(tui.baseWin)
		if err != nil {
			log.Fatal(err)
		}
	}
	if tui.treePanel == nil {
		tui.treePanel = nc.NewPanel(tui.treeWin)
	} else {
		err = tui.treePanel.Replace(tui.treeWin)
		if err != nil {
			log.Fatal(err)
		}
	}

	drawBase(tui.baseWin, opts.base)

	stdscr.Refresh()
	stdscr.GetChar()
}

func initPair(pair int16, fg int16, bg int16) {
	err := nc.InitPair(pair, fg, bg)
	if err != nil {
		log.Fatal(err)
	}
}

func drawBase(window *nc.Window, base *baseType) {
	switch *base {
	case bigPot:
		window.AttrOn(nc.A_BOLD | nc.ColorPair(8))
		window.Print(":")
		window.AttrOn(nc.ColorPair(2))
		window.Print("___________")
		window.AttrOn(nc.ColorPair(11))
		window.Print("./~~~\\.")
		window.AttrOn(nc.ColorPair(2))
		window.Print("___________")
		window.AttrOn(nc.ColorPair(8))
		window.Print(":")
		window.MovePrint(1, 0, " \\                           / ")
		window.MovePrint(2, 0, "  \\_________________________/ ")
		window.MovePrint(3, 0, "  (_)                     (_)")
		window.AttrOff(nc.A_BOLD)
	case smallPot:
		window.AttrOn(nc.ColorPair(8))
		window.Print("(")
		window.AttrOn(nc.ColorPair(2))
		window.Print("---")
		window.AttrOn(nc.ColorPair(11))
		window.Print("./~~~\\.")
		window.AttrOn(nc.ColorPair(2))
		window.Print("---")
		window.AttrOn(nc.ColorPair(8))
		window.Print(")")
		window.MovePrint(1, 0, " (           ) ")
		window.MovePrint(2, 0, "  (_________)  ")
	}
}

// display changes
func updateScreen() {
	nc.UpdatePanels()
	err := nc.Update()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(*opts.time)
}

// based on type of tree, determine what color a branch should be
func chooseColor(branchType branchType, treeWin *nc.Window) {
	switch branchType {
	case trunk, shootLeft, shootRight:
		if rand.Int()%2 == 0 {
			treeWin.AttrOn(nc.A_BOLD | nc.ColorPair(11))
		} else {
			treeWin.AttrOn(nc.ColorPair(3))
		}
	case dying:
		if rand.Int()%10 == 0 {
			treeWin.AttrOn(nc.A_BOLD | nc.ColorPair(2))
		} else {
			treeWin.AttrOn(nc.ColorPair(2))
		}
	case dead:
		if rand.Int()%3 == 0 {
			treeWin.AttrOn(nc.A_BOLD | nc.ColorPair(10))
		} else {
			treeWin.AttrOn(nc.ColorPair(10))
		}
	}
}

// roll (randomize) a given die
func roll(mod int) int {
	return rand.Int() % mod
}

// determine change in X and Y coordinates of a given branch
func setDeltas(branchType branchType, life int, age int, multiplier int, returnDx *int, returnDy *int) {
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

func chooseString(branchType branchType, life int, dx int, dy int) string {
	branchStr := "?"

	if life < 4 {
		branchType = dying
	}

	switch branchType {
	case trunk:
		if dy == 0 {
			branchStr = "/~"
		} else if dx < 0 {
			branchStr = "\\|"
		} else if dx == 0 {
			branchStr = "/|\\"
		} else if dx > 0 {
			branchStr = "|/"
		}

	case shootLeft:
		if dy > 0 {
			branchStr = "\\"
		} else if dy == 0 {
			branchStr = "\\_"
		} else if dx < 0 {
			branchStr = "\\|"
		} else if dx == 0 {
			branchStr = "/|"
		} else if dx > 0 {
			branchStr = "/"
		}

	case shootRight:
		if dy > 0 {
			branchStr = "/"
		} else if dy == 0 {
			branchStr = "_/"
		} else if dx < 0 {
			branchStr = "\\|"
		} else if dx == 0 {
			branchStr = "/|"
		} else if dx > 0 {
			branchStr = "/"
		}

	case dying, dead:
		branchStr = opts.leaves[rand.Int()%len(opts.leaves)]
	}

	return branchStr
}
