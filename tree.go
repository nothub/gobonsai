package main

import "time"

type align int

const (
	center align = iota + 1
	left
	right
)

type branch int

const (
	trunk branch = iota
	shootLeft
	shootRight
	// TODO: find better names for these
	dying
	dead
)

type counters struct {
	branches     int
	shoots       int
	shootCounter int
}

// The algorithm for tree generation was ported from the cbonsai codebase and ideally generates identical output.

func drawTree(sc *screen, opts opts) error {
	counters := counters{
		branches:     0,
		shoots:       0,
		shootCounter: rand.Int(),
	}
	life := opts.life
	_, maxY := sc.Size()
	return drawBranch(sc, opts, counters, life, trunk, sc.x, sc.y, maxY-(maxY-sc.y-1))
}

func drawBranch(sc *screen, opts opts, counters counters, life int, kind branch, x int, y int, maxY int) error {
	counters.branches++
	dx := 0
	dy := 0
	age := 0
	shootCooldown := opts.multiplier

	for life > 0 {
		life--
		age = opts.life - life
		dx, dy = deltas(kind, life, age, opts.multiplier)

		// reduce dy if too close to the ground
		if dy > 0 && y > (maxY-2) {
			dy--
		}

		// near-dead branch should branch into a lot of leaves
		if life < 3 {
			err := drawBranch(sc, opts, counters, life, dead, x, y, maxY)
			if err != nil {
				return err
			}

			// dying trunk should branch into a lot of leaves
		} else if kind == trunk && life < (opts.multiplier+2) {
			err := drawBranch(sc, opts, counters, life, dying, x, y, maxY)
			if err != nil {
				return err
			}

			// dying shoot should branch into a lot of leaves
		} else if (kind == shootLeft || kind == shootRight) && life < (opts.multiplier+2) {
			err := drawBranch(sc, opts, counters, life, dying, x, y, maxY)
			if err != nil {
				return err
			}

			// trunks should re-branch either randomly, or upon every <multiplier> steps
		} else if kind == trunk && (((rand.Int() % 3) == 0) || (life%opts.multiplier == 0)) {

			// if trunk is branching and not about to die, create another trunk with random life
			if (rand.Int()%8 == 0) && life > 7 {
				shootCooldown = opts.multiplier * 2 // reset shoot cooldown
				err := drawBranch(sc, opts, counters, life+(rand.Int()%5-2), trunk, x, y, maxY)
				if err != nil {
					return err
				}

				// otherwise create a shoot
			} else if shootCooldown <= 0 {
				shootCooldown = opts.multiplier * 2 // reset shoot cooldown
				shootLife := life + opts.multiplier

				// first shoot is randomly directed
				counters.shoots++
				counters.shootCounter++

				// create shoot
				err := drawBranch(sc, opts, counters, shootLife, branch((counters.shootCounter%2)+1), x, y, maxY)
				if err != nil {
					return err
				}
			}

		}

		shootCooldown--

		// move in x and y directions
		x += dx
		y += dy

		color := chooseColor(kind)

		// choose string to use for this branch
		leaf := chooseLeaf(kind, life, dx, dy, opts)

		sc.x = x
		sc.y = y
		sc.draw(leaf, color)

		if opts.live && active {
			evDrawn(sc)
			// We either await the delay or wait for shutdown.
			select {
			case <-shutdown:
				return nil
			case <-time.After(opts.time):
			}
		}
	}

	return nil
}

// determine change in X and Y coordinates of a given branch
func deltas(branchType branch, life int, age int, multiplier int) (int, int) {
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

			dice := rand.Intn(10)
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
			dice := rand.Intn(10)
			if dice > 2 {
				dy = -1
			} else {
				dy = 0
			}
			dx = (rand.Int() % 3) - 1
		}

	case shootLeft: // trend left and little vertical movement
		dice := rand.Intn(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 7 {
			dy = 0
		} else if dice >= 8 && dice <= 9 {
			dy = 1
		}

		dice = rand.Intn(10)
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
		dice := rand.Intn(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 7 {
			dy = 0
		} else if dice >= 8 && dice <= 9 {
			dy = 1
		}

		dice = rand.Intn(10)
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
		dice := rand.Intn(10)
		if dice >= 0 && dice <= 1 {
			dy = -1
		} else if dice >= 2 && dice <= 8 {
			dy = 0
		} else if dice >= 9 && dice <= 9 {
			dy = 1
		}

		dice = rand.Intn(15)
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
		dice := rand.Intn(10)
		if dice >= 0 && dice <= 2 {
			dy = -1
		} else if dice >= 3 && dice <= 6 {
			dy = 0
		} else if dice >= 7 && dice <= 9 {
			dy = 1
		}
		dx = (rand.Int() % 3) - 1
	}

	return dx, dy
}

func chooseLeaf(branch branch, life int, dx int, dy int, opts opts) string {
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
