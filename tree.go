package main

type branch int

const (
	trunk branch = iota
	shootLeft
	shootRight
	dying
	dead
)

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

//void branch(struct config *conf, struct ncursesObjects *objects, struct counters *myCounters, int y, int x, enum branchType type, int life) {
//	myCounters->branches++;
//	int dx = 0;
//	int dy = 0;
//	int age = 0;
//	int shootCooldown = conf->multiplier;
//
//	while (life > 0) {
//		if (checkKeyPress(conf, myCounters) == 1)
//			quit(conf, objects, 0);
//
//		life--;		// decrement remaining life counter
//		age = conf->lifeStart - life;
//
//		setDeltas(type, life, age, conf->multiplier, &dx, &dy);
//
//		int maxY = getmaxy(objects->treeWin);
//		if (dy > 0 && y > (maxY - 2)) dy--; // reduce dy if too close to the ground
//
//		// near-dead branch should branch into a lot of leaves
//		if (life < 3)
//			branch(conf, objects, myCounters, y, x, dead, life);
//
//		// dying trunk should branch into a lot of leaves
//		else if (type == 0 && life < (conf->multiplier + 2))
//			branch(conf, objects, myCounters, y, x, dying, life);
//
//		// dying shoot should branch into a lot of leaves
//		else if ((type == shootLeft || type == shootRight) && life < (conf->multiplier + 2))
//			branch(conf, objects, myCounters, y, x, dying, life);
//
//		// trunks should re-branch if not close to ground AND either randomly, or upon every <multiplier> steps
//		/* else if (type == 0 && ( \ */
//		/* 		(rand() % (conf.multiplier)) == 0 || \ */
//		/* 		(life > conf.multiplier && life % conf.multiplier == 0) */
//		/* 		) ) { */
//		else if (type == trunk && (((rand() % 3) == 0) || (life % conf->multiplier == 0))) {
//
//			// if trunk is branching and not about to die, create another trunk with random life
//			if ((rand() % 8 == 0) && life > 7) {
//				shootCooldown = conf->multiplier * 2;	// reset shoot cooldown
//				branch(conf, objects, myCounters, y, x, trunk, life + (rand() % 5 - 2));
//			}
//
//			// otherwise create a shoot
//			else if (shootCooldown <= 0) {
//				shootCooldown = conf->multiplier * 2;	// reset shoot cooldown
//
//				int shootLife = (life + conf->multiplier);
//
//				// first shoot is randomly directed
//				myCounters->shoots++;
//				myCounters->shootCounter++;
//				if (conf->verbosity) mvwprintw(objects->treeWin, 4, 5, "shoots: %02d", myCounters->shoots);
//
//				// create shoot
//				branch(conf, objects, myCounters, y, x, (myCounters->shootCounter % 2) + 1, shootLife);
//			}
//		}
//		shootCooldown--;
//
//		if (conf->verbosity > 0) {
//			mvwprintw(objects->treeWin, 5, 5, "dx: %02d", dx);
//			mvwprintw(objects->treeWin, 6, 5, "dy: %02d", dy);
//			mvwprintw(objects->treeWin, 7, 5, "type: %d", type);
//			mvwprintw(objects->treeWin, 8, 5, "shootCooldown: % 3d", shootCooldown);
//		}
//
//		// move in x and y directions
//		x += dx;
//		y += dy;
//
//		chooseColor(type, objects->treeWin);
//
//		// choose string to use for this branch
//		char *branchStr = chooseString(conf, type, life, dx, dy);
//
//		// grab wide character from branchStr
//		wchar_t wc = 0;
//		mbstate_t *ps = 0;
//		mbrtowc(&wc, branchStr, 32, ps);
//
//		// print, but ensure wide characters don't overlap
//		if(x % wcwidth(wc) == 0)
//			mvwprintw(objects->treeWin, y, x, "%s", branchStr);
//
//		wattroff(objects->treeWin, A_BOLD);
//		free(branchStr);
//
//		// if live, update screen
//		// skip updating if we're still loading from file
//		if (conf->live && !(conf->load && myCounters->branches < conf->targetBranchCount))
//			updateScreen(conf->timeStep);
//	}
//}
