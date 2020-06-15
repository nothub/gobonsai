#include <stdlib.h>
#include <ncurses.h>
#include <panel.h>
#include <unistd.h>
#include <getopt.h>
#include <time.h>
#include <string.h>
#include <ctype.h>

// global variables
int branches = 0;
int shoots = 0;
int shootsMax = 0;
int shootCounter;

int seed = 0;
int baseType = 1;
int lifeStart = 32;
int multiplier = 5;
int live = 0;
double timeStep = 0.03;
double timeWait = 4;
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
	printf("                           steps of growth [default: %f]\n", timeStep);
	printf("  -i, --infinite         infinite mode\n");
	printf("  -w, --wait TIME        in infinite mode, time in secs between\n");
	printf("                           tree generation [default: %f]\n", timeWait);
	printf("  -m, --message STR      attach message next to the tree\n");
	printf("  -g, --geometry X,Y     set custom geometry\n");
	printf("  -b, --base INT         ascii-art plant base to use, 0 is none\n");
	printf("  -c, --leaf STR1,STR2,STR3...   list of strings randomly chosen for leaves\n");
	printf("  -M, --multiplier INT   branch multiplier; higher -> more\n");
	printf("                           branching (0-20) [default: %i]\n", lifeStart);
	printf("  -L, --life INT         life; higher -> more growth (0-200) [default: %i]\n", lifeStart);
	printf("  -s, --seed INT         seed random number generator\n");
	printf("  -v, --verbose          increase output verbosity\n");
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

	// add windows to array of panels
	myPanels[0] = new_panel(baseWin);
	myPanels[1] = new_panel(treeWin);

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

void branch(int y, int x, int type, int life) {
	branches++;
	int dx = 0;
	int dy = 0;
	int dice = 0;
	int age = 0;

	while (life > 0) {
		life--;		// decrement remaining life counter
		age = lifeStart - life;

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

					roll(&dice, 10);
					if (dice >= 0 && dice <=0) dx = -2;
					else if (dice >= 1 && dice <= 3) dx = -1;
					else if (dice >= 4 && dice <= 5) dx = 0;
					else if (dice >= 6 && dice <= 8) dx = 1;
					else if (dice >= 9 && dice <= 9) dx = 2;
				}
				// middle-aged trunk
				else {
					roll(&dice, 10);
					if (dice > 2) dy = -1;
					else dy = 0;
					dx = (rand() % 3) - 1;
				}
				break;

			case 1: // left shoot: trend left and little vertical movement
				roll(&dice, 10);
				if (dice >= 0 && dice <= 1) dy = -1;
				else if (dice >= 2 && dice <= 7) dy = 0;
				else if (dice >= 8 && dice <= 9) dy = 1;

				roll(&dice, 10);
				if (dice >= 0 && dice <=1) dx = -2;
				else if (dice >= 2 && dice <= 5) dx = -1;
				else if (dice >= 6 && dice <= 8) dx = 0;
				else if (dice >= 9 && dice <= 9) dx = 1;
				break;

			case 2: // right shoot: trend right and little vertical movement
				roll(&dice, 10);
				if (dice >= 0 && dice <= 1) dy = -1;
				else if (dice >= 2 && dice <= 7) dy = 0;
				else if (dice >= 8 && dice <= 9) dy = 1;

				roll(&dice, 10);
				if (dice >= 0 && dice <=1) dx = 2;
				else if (dice >= 2 && dice <= 5) dx = 1;
				else if (dice >= 6 && dice <= 8) dx = 0;
				else if (dice >= 9 && dice <= 9) dx = -1;
				break;

			case 3: // dying: discourage vertical growth(?); trend left/right (-3,3)
				roll(&dice, 10);
				if (dice >= 0 && dice <=1) dy = -1;
				else if (dice >= 2 && dice <=8) dy = 0;
				else if (dice >= 9 && dice <=9) dy = 1;

				roll(&dice, 15);
				if (dice >= 0 && dice <=0) dx = -3;
				else if (dice >= 1 && dice <= 2) dx = -2;
				else if (dice >= 3 && dice <= 5) dx = -1;
				else if (dice >= 6 && dice <= 8) dx = 0;
				else if (dice >= 9 && dice <= 11) dx = 1;
				else if (dice >= 12 && dice <= 13) dx = 2;
				else if (dice >= 14 && dice <= 14) dx = 3;
				break;

			case 4: // dead: fill in surrounding area
				roll(&dice, 10);
				if (dice >= 0 && dice <= 2) dy = -1;
				else if (dice >= 3 && dice <= 6) dy = 0;
				else if (dice >= 7 && dice <= 9) dy = 1;
				dx = (rand() % 3) - 1;
				break;
		}

		int maxY, maxX;
		getmaxyx(treeWin, maxY, maxX);
		if (dy > 0 && y > (maxY - 2)) dy--; // reduce dy if too close to the ground

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
					else if (dx < 0) strcpy(branchChar, "\\|");
					else if (dx == 0) strcpy(branchChar, "/|\\");
					else if (dx > 0) strcpy(branchChar, "|/");
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
	if (message == NULL) return 1;

	// determine dimensions of window box
	int maxY, maxX;
	getmaxyx(stdscr, maxY, maxX);
	int boxWidth = 0;
	int boxHeight = 0;
	if (strlen(message) + 3 <= (0.25 * maxX)) {
		boxWidth = strlen(message) + 1;
		boxHeight = 1;
	} else {
		boxWidth = 0.25 * maxX;
		boxHeight = (strlen(message) / boxWidth) + (strlen(message) / boxWidth);
	}
	if (verbosity) mvwprintw(treeWin, 8, 5, "boxWidth: %0d", boxWidth);

	// create separate box for message border
	messageBorderWin = newwin(boxHeight + 2, boxWidth + 4, (maxY * 0.7) - 1, (maxX * 0.7) - 2);
	messageWin = newwin(boxHeight, boxWidth + 1, maxY * 0.7, maxX * 0.7);

	// draw box
	wattron(messageBorderWin, COLOR_PAIR(8));
	box(messageBorderWin, 0, 0);

	// assign new windows to array of panels
	myPanels[2] = new_panel(messageBorderWin);
	myPanels[3] = new_panel(messageWin);

	// word wrap message as it is written
	unsigned int i = 0;
	int linePosition = 0;
	int wordLength = 0;
	char wordBuffer[500];
	wordBuffer[0] = '\0';
	char thisChar;
	int messageBoxWidth = boxWidth - 1;
	while (true) {
		thisChar = message[i];
		if (verbosity) {
			mvwprintw(treeWin, 9, 5, "index: %03d", i);
			mvwprintw(treeWin, 10, 5, "linePosition: %02d", linePosition);
		}

		// if char is not a space or null char
		if (!(isspace(thisChar) || thisChar == '\0')) {
			strncat(wordBuffer, &thisChar, 1); // append thisChar to wordBuffer
			wordLength++;
			linePosition++;
		}

		// if char is space or null char
		else if (isspace(thisChar) || thisChar == '\0') {

			// if current line can fit word, add word to current line
			if (linePosition <= messageBoxWidth) {
				wprintw(messageWin, "%s", wordBuffer);	// print word
				wordLength = 0;		// reset word length
				wordBuffer[0] = '\0';	// clear word buffer

				void addSpaces(int count) {
					// add spaces if there's enough space
					if (linePosition < (messageBoxWidth - count)) {
						if (verbosity) mvwprintw(treeWin, 12, 5, "inserting a space: linePosition: %02d, wordLength: %02d", linePosition, wordLength);

						// add spaces up to width
						for (int j = 0; j < count; j++) {
							wprintw(messageWin, "%s", " ");
							linePosition++;
						}
					}
				}

				switch (thisChar) {
					case ' ':
						addSpaces(1);
						break;
					case '\t':
						addSpaces(1);
						break;
					case '\n':
						waddch(messageWin, thisChar);
						linePosition = 0;
						break;
				}

			}

			// if word can't fit within a single line, just print it
			else if (wordLength > messageBoxWidth) {
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

		if (verbosity >= 2) {
			update_panels();
			doupdate();
			usleep(100000);
			mvwprintw(treeWin, 11, 5, "word buffer: |% 15s|", wordBuffer);
		}
		if (thisChar == '\0') break;	// quit when we reach the end of the message
		i++;
	}
	return 0;
}

void init() {
	savetty();	// save terminal settings
	initscr();	// init ncurses screen
	noecho();	// don't echo input to screen
	curs_set(0);	// make cursor invisible
	cbreak();	// don't wait for new line to grab user input

	// if terminal has color capabilities, use them
	if (has_colors()) {
		start_color();

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
	}

	// define and draw windows, then create panels
	drawWins(baseType, &baseWin, &treeWin);
	drawMessage();
}

void growTree() {
	int maxY, maxX;
	getmaxyx(treeWin, maxY, maxX);

	branches = 0;
	shoots = 0;
	shootCounter = rand();

	if (verbosity > 0) {
		mvwprintw(treeWin, 2, 5, "maxX: %03d, maxY: %03d", maxX, maxY);
		mvwprintw(treeWin, 3, 5, "seed: %d", seed);
	}
	branch(maxY - 1, (maxX / 2), 0, lifeStart);	// grow tree trunk

	// display changes
	update_panels();
	doupdate();
}

int main(int argc, char* argv[]) {
	int infinite = 0;
	int screensaver = 0;

	int termSize = 1;
	char *leavesInput = "&";
	char *geometry;

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

	// delimit leaves on "," and add each token to the leaves[] list
	char *token = strtok(leavesInput, ",");
	while (token != NULL) {
		if (leavesSize < 100) leaves[leavesSize] = token;
		token = strtok(NULL, ",");
		leavesSize++;
	}

	branchesMax = multiplier * 110;
	shootsMax = multiplier;

	// seed random number generator
	if (seed == 0) seed = time(NULL);
	srand(seed);

	do {
		init();
		growTree();
		if (infinite) {
			sleep(timeWait);
			clear();
			refresh();

			// seed random number generator
			seed = time(NULL);
			srand(seed);
		}
	} while (infinite);

	wgetch(treeWin);	// quit upon any input
	finish();
	return 0;
}
