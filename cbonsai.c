#include <stdlib.h>
#include <ncurses.h>
#include <panel.h>
#include <unistd.h>
#include <getopt.h>
#include <time.h>
#include <string.h>

// global variables
int branches = 0;
int shoots = 0;
int branchesMax = 0;
int shootsMax = 0;
int shootCounter;

int baseType = 1;
int lifeStart = 32;
int multiplier = 5;
int live = 0;
double timeStep = 0.03;
int verbosity = 0;
char *message = NULL;
int leavesSize = 0;
char* leaves[100];

WINDOW *baseWin, *treeWin, *messageBorderWin, *messageWin;
PANEL *myPanels[4];

void finish() {
	delwin(baseWin);
	delwin(treeWin);
	delwin(messageBorderWin);
	delwin(messageWin);

	clear();
	refresh();
	endwin();	// delete ncurses screen
	curs_set(1);
}

void printHelp() {
	printf("Usage: cbonsai [OPTIONS]\n");
	printf("\n");
	printf("cbonsai is a beautifully random bonsai tree generator.\n");
	printf("\n");
	printf("optional args:\n");
	printf("  -l, --live             live mode\n");
	printf("  -t, --time TIME        in live mode, minimum time in secs between\n");
	printf("                           steps of growth [default: 0.03]\n");
	printf("  -i, --infinite         infinite mode\n");
	printf("  -w, --wait TIME        in infinite mode, time in secs between\n");
	printf("                           tree generation [default: 4]\n");
	printf("  -m, --message STR      attach message next to the tree\n");
	printf("  -g, --geometry X,Y     set custom geometry\n");
	printf("  -b, --base INT         ascii-art plant base to use, 0 is none\n");
	printf("  -c, --leaf STR1,STR2,STR3...   list of strings randomly chosen for leaves\n");
	printf("  -M, --multiplier INT   branch multiplier; higher -> more\n");
	printf("                           branching (0-20) [default: 5]\n");
	printf("  -L, --life INT         life; higher -> more growth (0-200) [default: 28]\n");
	printf("  -s, --seed INT         seed random number generator (0-32767)\n");
	printf("  -v, --verbose          print information each step of generation\n");
	printf("  -h, --help             show help	\n");
}

void drawWins(int baseType, WINDOW* *baseWinPtr, WINDOW* *treeWinPtr) {
	int baseWidth = 0;
	int baseHeight = 0;
	int rows, cols;

	switch(baseType) {
		case 1:
			baseWidth = 31;
			baseHeight = 4;
			break;
		case 2:
			baseWidth = 15;
			baseHeight = 3;
			break;
	}

	// calculate where base should go
	getmaxyx(stdscr, rows, cols);
	int baseOriginY = (rows - baseHeight);
	int baseOriginX = (cols / 2) - (baseWidth / 2);

	// create windows
	*baseWinPtr = newwin(baseHeight, baseWidth, baseOriginY, baseOriginX);
	*treeWinPtr = newwin(rows - baseHeight, cols, 0, 0);

	WINDOW *baseWin = *baseWinPtr;
	WINDOW *treeWin = *treeWinPtr;

	// draw art
	switch(baseType) {
		case 1:
			wattron(baseWin, A_BOLD | COLOR_PAIR(8));
			wprintw(baseWin, "%s", ":");
			wattron(baseWin, COLOR_PAIR(2));
			wprintw(baseWin, "%s", "___________");
			wattron(baseWin, COLOR_PAIR(11));
			wprintw(baseWin, "%s", "./~~~\\.");
			wattron(baseWin, COLOR_PAIR(2));
			wprintw(baseWin, "%s", "___________");
			wattron(baseWin, COLOR_PAIR(8));
			wprintw(baseWin, "%s", ":");

			mvwprintw(baseWin, 1, 0, "%s", " \\                           / ");
			mvwprintw(baseWin, 2, 0, "%s", "  \\_________________________/ ");
			mvwprintw(baseWin, 3, 0, "%s", "  (_)                     (_)");

			wattroff(baseWin, A_BOLD);
			break;
		case 2:
			wattron(baseWin, COLOR_PAIR(8));
			wprintw(baseWin, "%s", "(");
			wattron(baseWin, COLOR_PAIR(2));
			wprintw(baseWin, "%s", "---");
			wattron(baseWin, COLOR_PAIR(11));
			wprintw(baseWin, "%s", "./~~~\\.");
			wattron(baseWin, COLOR_PAIR(2));
			wprintw(baseWin, "%s", "---");
			wattron(baseWin, COLOR_PAIR(8));
			wprintw(baseWin, "%s", ")");

			mvwprintw(baseWin, 1, 0, "%s", " (           ) ");
			mvwprintw(baseWin, 2, 0, "%s", "  (_________)  ");
			break;
	}
}

// roll (randomize) a given die
void roll(int *dice, int mod) { *dice = rand() % mod; }

int branch(int y, int x, int type, int life) {
	branches++;
	int dx, dy;
	int d10;
	int age;

	while (life > 0) {
		life--;		// decrement remaining life counter
		age = lifeStart - life;
		roll(&d10, 10); // randomize growth

		switch (type) {
			case 0: // trunk

				// new or dead trunk
				if (age <= 2 || life < 4) {
					dy = 0;
					dx = (rand() % 3) - 1;
				}
				// young trunk should grow wide
				else if (age < (multiplier * 3)) {

					// every (multiplier * 0.8) steps, raise tree to next level
					if (age % (int) (multiplier * 0.5) == 0) dy = -1;
					else dy = 0;
					dx = (rand() % 5) - 2;
				}
				// middle-aged trunk
				else {
					if (d10 > 2) dy = -1;
					else dy = 0;
					dx = (rand() % 3) - 1;
				}
				break;

			case 1: // left shoot: trend left and little vertical movement
				if (d10 >= 0 && d10 <= 1) dy = -1;
				else if (d10 >= 2 && d10 <= 7) dy = 0;
				else if (d10 >= 8 && d10 <= 9) dy = 1;

				roll(&d10, 10);
				if (d10 >= 0 && d10 <=1) dx = -2;
				else if (d10 >= 2 && d10 <= 5) dx = -1;
				else if (d10 >= 6 && d10 <= 8) dx = 0;
				else if (d10 >= 9 && d10 <= 9) dx = 1;
				break;

			case 2: // right shoot: trend right and little vertical movement
				if (d10 >= 0 && d10 <= 1) dy = -1;
				else if (d10 >= 2 && d10 <= 7) dy = 0;
				else if (d10 >= 8 && d10 <= 9) dy = 1;

				roll(&d10, 10);
				if (d10 >= 0 && d10 <=1) dx = 2;
				else if (d10 >= 2 && d10 <= 5) dx = 1;
				else if (d10 >= 6 && d10 <= 8) dx = 0;
				else if (d10 >= 9 && d10 <= 9) dx = -1;
				break;

			case 3: // dying: discourage vertical growth(?); trend left/right (-3,3)
				if (d10 >= 0 && d10 <=1) dy = -1;
				else if (d10 >= 2 && d10 <=8) dy = 0;
				else if (d10 >= 9 && d10 <=9) dy = 1;
				dx = (rand() % 7) - 3;
				break;

			case 4: // dead: fill in surrounding area
				if (d10 >= 0 && d10 <= 2) dy = -1;
				else if (d10 >= 3 && d10 <= 6) dy = 0;
				else if (d10 >= 7 && d10 <= 9) dy = 1;
				dx = (rand() % 3) - 1;
				break;
		}

		int maxY, maxX;
		getmaxyx(treeWin, maxY, maxX);
		if (dy > 0 && y > (maxY - 2)) dy--; // reduce dy if too close to the ground

		// re-branch upon certain conditions
		if (branches < branchesMax) {

			// near-dead branch should branch into a lot of leaves
			if (life < 3) branch(y, x, 4, life);

			// dying trunk should branch into a lot of leaves
			else if (type == 0 && life < (multiplier + 2)) branch(y, x, 3, life);

			// dying shoot should branch into a lot of leaves
			else if ((type == 1 || type == 2) && life < (multiplier + 2)) branch(y, x, 3, life);

			// trunks should re-branch if not close to ground AND either randomly, or upon every <multiplier> steps
			else if (type == 0 && y < (maxY - multiplier + 1) && ( \
					(rand() % (16 - multiplier)) == 0 || \
					(life > multiplier && life % multiplier == 0)
					) ) {

				// if trunk is branching and not about to die, create another trunk
				if ((rand() % 3 == 0) && life > 7) branch(y, x, 0, life);

				// otherwise create a shoot
				else if (shoots < shootsMax) {
					int shootLife = (life + multiplier);
					if (shootLife < 0) shootLife = 0;

					// first shoot is randomly directed
					shoots++;
					shootCounter++;
					if (verbosity) mvwprintw(treeWin, 4, 5, "shoots: %02d", shoots);

					// create shoot
					branch(y, x, (shootCounter % 2) + 1, shootLife);
				}
			}
		} else {
			if (verbosity) mvwprintw(treeWin, 2, 5, "%s", "Max branches hit!");
		}

		// move in x and y directions
		if (verbosity > 0) {
			mvwprintw(treeWin, 5, 5, "dx: %02d", dx);
			mvwprintw(treeWin, 6, 5, "dy: %02d", dy);
		}
		x = (x + dx);
		y = (y + dy);

		// choose color
		switch(type) {
			case 0:
			case 1:
			case 2: // trunk or shoot
				if (rand() % 2 == 0) wattron(treeWin, A_BOLD | COLOR_PAIR(11));
				else wattron(treeWin, COLOR_PAIR(3));
				break;

			case 3: // dying
				if (rand() % 10 == 0) wattron(treeWin, A_BOLD | COLOR_PAIR(2));
				else wattron(treeWin, COLOR_PAIR(2));
				break;

			case 4: // dead
				if (rand() % 3 == 0) wattron(treeWin, A_BOLD | COLOR_PAIR(10));
				else wattron(treeWin, COLOR_PAIR(10));
				break;
		}

		// choose string to use for this branch
		char *branchChar = malloc(100);
		strcpy(branchChar, "?");	// fallback character
		if (life < 4 || type >= 3) strcpy(branchChar, leaves[rand() % leavesSize]);
		else {
			switch(type) {
				case 0: // trunk
					if (dy == 0) strcpy(branchChar, "/~");
					else if (dx < 0) strcpy(branchChar, "\\");
					else if (dx == 0) strcpy(branchChar, "/|");
					else if (dx > 0) strcpy(branchChar, "/");
					break;
				case 1: // left shoot
					if (dy > 0) strcpy(branchChar, "\\");
					else if (dy == 0) strcpy(branchChar, "\\_");
					else if (dx < 0) strcpy(branchChar, "\\|");
					else if (dx == 0) strcpy(branchChar, "/|");
					else if (dx > 0) strcpy(branchChar, "/");
					break;
				case 2: // right shoot
					if (dy > 0) strcpy(branchChar, "/");
					else if (dy == 0) strcpy(branchChar, "_/");
					else if (dx < 0) strcpy(branchChar, "\\|");
					else if (dx == 0) strcpy(branchChar, "/|");
					else if (dx > 0) strcpy(branchChar, "/");
					break;
			}
		}

		mvwprintw(treeWin, y, x, "%s", branchChar);
		free(branchChar);
		wattroff(treeWin, A_BOLD);

		// if live, show progress
		if (live) {
			struct timespec tm1, tm2;
			tm1.tv_sec = 0;
			tm1.tv_nsec = (timeStep * 100);

			// display changes
			update_panels();
			doupdate();

			/* nanosleep(&tm1, &tm2); */
			usleep(timeStep * 1000000);
		}

	}
}

int drawMessage() {
	// draw message
	int maxY, maxX;
	getmaxyx(stdscr, maxY, maxX);
	if (message != NULL) {

		// determine dimensions of window box
		int boxWidth = 0;
		int boxHeight = 0;
		if (strlen(message) <= (0.25 * maxX)) {
			boxWidth = strlen(message);
			boxHeight = 1;
		} else {
			boxWidth = 0.25 * maxX;
			boxHeight = (strlen(message) / boxWidth) + (strlen(message) / boxWidth * 0.9);
		}
		if (verbosity) mvwprintw(treeWin, 5, 5, "boxWidth: %0d", boxWidth);

		// create separate box for message border
		messageBorderWin = newwin(boxHeight + 2, boxWidth + 4, (maxY * 0.7) - 1, (maxX * 0.7) - 2);
		messageWin = newwin(boxHeight, boxWidth + 1, maxY * 0.7, maxX * 0.7);

		// draw box
		wattron(messageBorderWin, COLOR_PAIR(8));
		box(messageBorderWin, 0, 0);

		// word wrap message as it is written
		unsigned int i = 0;
		int linePosition = 1;
		int wordLength = 0;
		char wordBuffer[500];
		wordBuffer[0] = '\0';
		while (true) {
			if (verbosity) {
				mvwprintw(treeWin, 9, 5, "index: %03d", i);
				mvwprintw(treeWin, 10, 5, "linePosition: %02d", linePosition);
			}

			// if char is not a space or null char
			if (message[i] != ' ' && message[i] != '\0' && i < sizeof(wordBuffer)) {
				strncat(wordBuffer, &message[i], 1); // append message[i] to wordBuffer
				wordLength++;
				linePosition++;
			}

			// if char is space or null char
			else if (message[i] == ' ' || message[i] == '\0') {

				// if current line can fit word, add word to current line
				if (linePosition - 1 <= boxWidth) {
					wprintw(messageWin, "%s", wordBuffer);	// print word
					wordLength = 0;		// reset word length
					wordBuffer[0] = '\0';	// clear word buffer
					if (linePosition < (boxWidth - 2)) {
						if (verbosity) mvwprintw(treeWin, 12, 5, "inserting a space: linePosition: %02d, wordLength: %02d", linePosition, wordLength);
						wprintw(messageWin, "%s", " ");
						linePosition++;
					}
				}

				// if word can't fit within a single line, just print it
				else if (wordLength > boxWidth) {
					wprintw(messageWin, "%s ", wordBuffer);	// print word
					wordLength = 0;		// reset word length
					wordBuffer[0] = '\0';	// clear word buffer

					// our line position on this new line is the x coordinate
					int y;
					getyx(messageWin, y, linePosition);
				}

				// if current line can't fit word, go to next line
				else {
					if (verbosity) mvwprintw(treeWin, (i / 24) + 28, 5, "couldn't fit word. linePosition: %02d, wordLength: %02d", linePosition, wordLength);
					wprintw(messageWin, "\n%s ", wordBuffer); // print newline, then word
					linePosition = wordLength;	// reset line position
					wordLength = 0;		// reset word length
					wordBuffer[0] = '\0';	// clear word buffer
				}
			}
			else {
				printf("%s", "Error while parsing message");
				return 1;
			}

			if (verbosity) mvwprintw(treeWin, 11, 5, "word buffer: |% 12s|", wordBuffer);
			if (message[i] == '\0' || i > strlen(message)) break;	// quit when we reach the end of the message
			i++;
		}
	}
}

void init() {
	savetty();	// save terminal settings
	initscr();	// init ncurses screen
	noecho();	// don't echo input to screen
	curs_set(0);	// make cursor invisible
	cbreak();	// don't wait for new line to grab user input

	// define and draw windows
	drawWins(baseType, &baseWin, &treeWin);
	drawMessage();

	// create panels
	myPanels[0] = new_panel(baseWin);
	myPanels[1] = new_panel(treeWin);
	myPanels[2] = new_panel(messageBorderWin);
	myPanels[3] = new_panel(messageWin);
}

int main(int argc, char* argv[]) {
	int infinite = 0;
	int screensaver = 0;

	int termSize = 1;
	char *leavesInput = "&";
	char *geometry;

	int seed = 0;
	double timeWait = 4;

	struct option long_options[] = {
		{"live", no_argument, NULL, 'l'},
		{"time", required_argument, NULL, 't'},
		{"infinite", no_argument, NULL, 'i'},
		{"wait", required_argument, NULL, 'w'},
		{"screensaver", no_argument, NULL, 'S'},
		{"message", required_argument, NULL, 'm'},
		{"geometry", required_argument, NULL, 'g'},
		{"base", required_argument, NULL, 'b'},
		{"leaf", required_argument, NULL, 'c'},
		{"multiplier", required_argument, NULL, 'M'},
		{"life", required_argument, NULL, 'L'},
		{"seed", required_argument, NULL, 's'},
		{"verbose", no_argument, NULL, 'v'},
		{"help", no_argument, NULL, 'h'},
		{0, 0, 0, 0}
	};

	// parse arguments
	int option_index = 0;
	int c;
	while ((c = getopt_long(argc, argv, "lt:iw:Sm:g:b:c:M:L:s:vh", long_options, &option_index)) != -1) {
		switch (c) {
			case 'l':
				live = 1;
				break;
			case 't':
				timeStep = atof(optarg);
				break;
			case 'i':
				infinite = 1;
				break;
			case 'w':
				timeWait = atof(optarg);
				break;
			case 'S':
				screensaver = 1;
				break;
			case 'm':
				message = strdup(optarg);
				break;
			case 'g':
				geometry = optarg;
				break;
			case 'b':
				baseType = atoi(optarg);
				break;
			case 'c':
				leavesInput = strdup(optarg);
				break;
			case 'M':
				multiplier = atoi(optarg);
				break;
			case 'L':
				lifeStart = atoi(optarg);
				break;
			case 's':
				seed = atoi(optarg);
				break;
			case 'v':
				verbosity++;
				break;

			// '?' represents unknown option. Treat it like --help.
			case '?':
			case 'h':
				printHelp();
				return 0;
				break;
		}
	}

	init();

	// if terminal has color capabilities, use them
	if (has_colors()) {
		start_color();	// allow us to use color capabilities

		// define color pairs
		init_pair(0, 0, COLOR_BLACK);
		init_pair(1, 1, COLOR_BLACK);
		init_pair(2, 2, COLOR_BLACK);
		init_pair(3, 3, COLOR_BLACK);
		init_pair(4, 4, COLOR_BLACK);
		init_pair(5, 5, COLOR_BLACK);
		init_pair(6, 6, COLOR_BLACK);
		init_pair(7, 7, COLOR_BLACK);
		init_pair(8, 8, COLOR_BLACK);
		init_pair(9, 9, COLOR_BLACK);
		init_pair(10, 10, COLOR_BLACK);
		init_pair(11, 11, COLOR_BLACK);
		init_pair(12, 12, COLOR_BLACK);
		init_pair(13, 13, COLOR_BLACK);
		init_pair(14, 14, COLOR_BLACK);
		init_pair(15, 15, COLOR_BLACK);
	} else {
		printf("Exiting: terminal does not support colors\n");
		finish();
		return 1;
	}

	// delimit leaves on "," and add each token to the leaves[] list
	char *token = strtok(leavesInput, ",");
	while (token != NULL) {
		if (leavesSize < 100) leaves[leavesSize] = token;
		printf("%s\n", token);
		token = strtok(NULL, ",");
		leavesSize++;
	}

	branchesMax = multiplier * 110;
	shootsMax = multiplier;
	shootCounter = rand();

	// seed random number generator
	if (seed == 0) seed = time(NULL);
	srand(seed);

	void growTree() {
		int maxY, maxX;
		getmaxyx(treeWin, maxY, maxX);

		branches = 0;
		shoots = 0;

		if (verbosity > 0) {
			mvwprintw(treeWin, 2, 5, "maxX: %03d, maxY: %03d", maxX, maxY);
			mvwprintw(treeWin, 3, 5, "seed: %d", seed);
		}
		branch(maxY - 1, (maxX / 2), 0, lifeStart);	// grow tree trunk

		// display changes
		update_panels();
		doupdate();

		seed = time(NULL);
		srand(time(NULL));	// re-seed tree
	}

	do {
		init();
		growTree();
		if (infinite) {
			sleep(timeWait);
			clear();
			refresh();
		}
	} while (infinite);

	wgetch(treeWin);	// quit upon any input
	finish();
	return 0;
}
