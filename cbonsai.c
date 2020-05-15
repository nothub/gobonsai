#include <ncurses.h>

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

// attron(A_BOLD);
// attrset(A_NORMAL); || standend()

int main() {
	int row, col, centerx, ch;

	initscr();	// init curses mode
	savetty();
	noecho();	// don't echo input to screen
	cbreak();	// don't wait for new line to grab user input

	if (has_colors()) {
		start_color();	// allow us to use color features

		// define color pairs
		init_pair(COLOR_BLACK, -1, -1);
		init_pair(COLOR_RED, COLOR_RED, -1);
		init_pair(COLOR_GREEN, COLOR_GREEN, -1);
		init_pair(COLOR_YELLOW, COLOR_YELLOW, -1);
		init_pair(COLOR_BLUE, COLOR_BLUE, -1);
		init_pair(COLOR_MAGENTA, COLOR_MAGENTA, -1);
		init_pair(COLOR_CYAN, COLOR_CYAN, -1);
		init_pair(COLOR_WHITE, COLOR_WHITE, -1);
	} else {
		finish();
		return 1;
	}

	getmaxyx(stdscr, row, col);
	centerx = col / 2;

	attron(COLOR_PAIR(COLOR_RED));
	printw("Hello World!"); // print to stdscr

	ch = getch();

	mvaddch(5, centerx, ch);
	refresh();		// update screen
	getch();

	finish();
}
