package main

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
