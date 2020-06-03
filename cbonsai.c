#include <stdlib.h>
#include <ncurses.h>
#include <unistd.h>
#include <getopt.h>
#include <time.h>
#include <string.h>

void finish() {
	clear();
	refresh();
	endwin();

	curs_set(1);	// make cursor visible again
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

// attrset(A_NORMAL); || standend()

int main(int argc, char* argv[]) {
	int rows = 0;
	int cols = 0;
	int y,x;

	int live = 0;
	int infinite = 0;
	int screensaver = 0;

	int verbosity = 0;
	int termSize = 1;
	int termColors = 0;
	int baseType = 1;
	char *leafStrs = "&";
	char *message;
	char *geometry;

	double multiplier = 5;
	int lifeStart = 28;

	double timeStep = 0.03;
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
				leafStrs = optarg;
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
		init_pair(COLOR_BLACK, COLOR_WHITE, COLOR_BLACK);
		init_pair(COLOR_RED, COLOR_RED, COLOR_BLACK);
		init_pair(COLOR_GREEN, COLOR_GREEN, COLOR_BLACK);
		init_pair(COLOR_YELLOW, COLOR_YELLOW, COLOR_BLACK);
		init_pair(COLOR_BLUE, COLOR_BLUE, COLOR_BLACK);
		init_pair(COLOR_MAGENTA, COLOR_MAGENTA, COLOR_BLACK);
		init_pair(COLOR_CYAN, COLOR_CYAN, COLOR_BLACK);
		init_pair(COLOR_WHITE, COLOR_WHITE, COLOR_BLACK);
	} else {
		printf("Exiting: terminal does not support colors\n");
		finish();
		return 1;
	}

	getmaxyx(stdscr, row, col);
	centerx = col / 2;

	printw("Hello World!"); // print to stdscr

	ch = getch();

	mvaddch(5, centerx, ch);
	refresh();		// update screen
	getch();

	finish();
}
