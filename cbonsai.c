#include <stdlib.h>
#include <ncurses.h>
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

int lifeStart = 28;
int multiplier = 5;
int live = 0;
double timeStep = 0.03;

int leavesSize = 0;
char* leaves[100];

WINDOW* treeWin;
WINDOW* baseWin;

void finish() {
	clear();
	refresh();
	endwin();

	curs_set(1);	// make cursor visible again

	delwin(baseWin);
	delwin(treeWin);
}

void printHelp() {
	printf("Usage: bonsai [OPTIONS]\n");
	printf("\n");
	printf("bonsai.sh is a beautifully random bonsai tree generator.\n");
	printf("\n");
	printf("optional args:\n");
	printf("  -l, --live             live mode\n");
	printf("  -t, --time TIME        in live mode, minimum time in secs between\n");
	printf("                           steps of growth [default: 0.03]\n");
	printf("  -i, --infinite         infinite mode\n");
	printf("  -w, --wait TIME        in infinite mode, time in secs between\n");
	printf("                           tree generation [default: 4]\n");
	printf("  -n, --neofetch         neofetch mode\n");
	printf("  -m, --message STR      attach message next to the tree\n");
	printf("  -T, --termcolors       use terminal colors\n");
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
	int baseWidth, baseHeight;
	int rows, cols;

	switch(baseType) {
		case 1:
			baseWidth = 30;
			baseHeight = 4;
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
			wprintw(baseWin, "%s", "./~~\\.");
			wattron(baseWin, COLOR_PAIR(2));
			wprintw(baseWin, "%s", "___________");
			wattron(baseWin, COLOR_PAIR(8));
			wprintw(baseWin, "%s", ":");

			mvwprintw(baseWin, 1, 0, "%s", " \\                          / ");
			mvwprintw(baseWin, 2, 0, "%s", "  \\________________________/ ");
			mvwprintw(baseWin, 3, 0, "%s", "  (_)                    (_)");

			wattroff(baseWin, A_BOLD);
			break;
	}
}

// roll (randomize) a given die
void roll(int *dice, int mod) {
	*dice = rand() % mod;
}

int main(int argc, char* argv[]) {
	int rows = 0;
	int cols = 0;
	int y,x;

	int infinite = 0;
	int screensaver = 0;

	int verbosity = 0;
	int termSize = 1;
	int termColors = 0;
	int baseType = 1;
	char *leavesInput = "&";
	char *message;
	char *geometry;


	double timeWait = 4;

	int flag_m = 0;
	int c;

	struct option long_options[] = {
		{"live", no_argument, NULL, 'l'},
		{"time", required_argument, NULL, 't'},
		{"infinite", no_argument, NULL, 'i'},
		{"wait", required_argument, NULL, 'w'},
		{"screensaver", no_argument, NULL, 'S'},
		{"message", required_argument, NULL, 'm'},
		{"termcolors", no_argument, NULL, 'T'},
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

	// seed random number generator using time, and generate tree seed
	srand(time(NULL));
	int seed = rand();

	// parse arguments
	int option_index = 0;
	while ((c = getopt_long(argc, argv, "lt:iw:Sm:Tg:b:c:M:L:s:vh", long_options, &option_index)) != -1) {
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
				message = optarg;
				break;
			case 'T':
				termColors = 1;
				break;
			case 'g':
				geometry = optarg;
				break;
			case 'b':
				baseType = atoi(optarg);
				break;
			case 'c':
				leavesInput = optarg;
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

	initscr();	// init ncurses screen
	savetty();
	noecho();	// don't echo input to screen
	curs_set(0);	// make cursor invisible
	cbreak();	// don't wait for new line to grab user input

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

	/* char* brown =		"[0;33m"; */
	/* char* brownDark =	"[1;33m"; */
	/* char* green =		"[0;32m"; */
	/* char* greenDark =	"[1;32m"; */
	/* char* gray =		"[0;38m"; */
	/* char* reset =		"[0m"; */

	// delimit leaves on "," and add each token to the leaves[] list
	char *token = strtok(leavesInput, ",");
	while (token != NULL) {
		if (leavesSize < 100) leaves[leavesSize] = token;
		printf("%s\n", token);
		token = strtok(NULL, ",");
		leavesSize++;
	}

	// seed tree (no pun intended)
	srand(seed);

	branchesMax = multiplier * 110;
	shootsMax = multiplier;
	shootCounter = rand();

	// create windows and draw base
	drawWins(baseType, &baseWin, &treeWin);

	// grow trunk
	int maxY, maxX;
	getmaxyx(treeWin, maxY, maxX);
	wrefresh(baseWin);
	getch();

	finish();
	return 0;
}
