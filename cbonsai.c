#include <stdlib.h>
#include <ncurses.h>
#include <unistd.h>

void finish() {
	clear();
	refresh();
	endwin();
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
	int row = 0;
	int col = 0;
	int centerx, ch;

	int live = 0;
	int infinite = 0;

	int termSize = 1;
	int termColors = 0;
	char leafStrs[] = "&";
	int baseType = 1;
	char message[] = "";

	double multiplier = 5;
	int lifeStart = 28;

	double timeStep = 0.03;
	double timeWait = 4;

	int flag_m = 0;

	int c;

	// parse arguments
	while ((c = getopt(argc, argv, "hlt:w:ig:c:Tm:b:M:L:s:vn")) != EOF) {
		switch (c) {
			case 'h':
				printHelp();
				return 0;
				break;
			case 'l':
				live = 1;
				break;
			case 't':
				timeStep = atof(optarg);
				break;
		}
	}

	initscr();	// init ncurses screen
	savetty();
	noecho();	// don't echo input to screen
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
