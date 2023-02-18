package main

import "github.com/fatih/color"

// TODO: use gocui attribute type instead of fatih/color module

var bold = color.New(color.Bold)

var black = color.New(color.FgBlack)
var blackBold = color.New(color.FgBlack, color.Bold)
var red = color.New(color.FgRed)
var redBold = color.New(color.FgRed, color.Bold)
var green = color.New(color.FgGreen)
var greenBold = color.New(color.FgGreen, color.Bold)
var yellow = color.New(color.FgYellow)
var yellowBold = color.New(color.FgYellow, color.Bold)
var blue = color.New(color.FgBlue)
var blueBold = color.New(color.FgBlue, color.Bold)
var magenta = color.New(color.FgMagenta)
var magentaBold = color.New(color.FgMagenta, color.Bold)
var cyan = color.New(color.FgCyan)
var cyanBold = color.New(color.FgCyan, color.Bold)
var white = color.New(color.FgWhite)
var whiteBold = color.New(color.FgWhite, color.Bold)

var brown = color.New(color.FgRed, color.FgYellow)
var brownBold = color.New(color.FgRed, color.FgYellow, color.Bold)

//	// if terminal has color capabilities, use them
//	if (has_colors()) {
//		start_color();
//
//		// use native background color when possible
//		int bg = COLOR_BLACK;
//		if (use_default_colors() != ERR) bg = -1;
//
//		// define color pairs
//		for(int i=0; i<16; i++){
//			init_pair(i, i, bg);
//		}
//
//		// restrict color pallete in non-256color terminals (e.g. screen or linux)
//		if (COLORS < 256) {
//
//          // ncurses func:
//          // int init_pair(short pair, short f, short b);
//          // (set pair 'pair' to fg 'f' and bg 'b')
//
//			init_pair(8, 7, bg);	// gray will look white
//			init_pair(9, 1, bg);
//			init_pair(10, 2, bg);
//			init_pair(11, 3, bg);
//			init_pair(12, 4, bg);
//			init_pair(13, 5, bg);
//			init_pair(14, 6, bg);
//			init_pair(15, 7, bg);
//		}
//	} else {
//		printf("%s", "Warning: terminal does not have color support.\n");
//	}

// based on type of tree, determine what color a branch should be
func chooseColor(kind branch) *color.Color {
	// TODO: these colors are wrong, we want original cbonsai colors (see above) instead!

	switch kind {
	case dying:
		if rand.Int()%10 == 0 {
			return greenBold
		} else {
			return green
		}

	case dead:
		if rand.Int()%3 == 0 {
			return blackBold
		} else {
			return black
		}

		// trunk | shootLeft | shootRight
	default:
		if rand.Int()%2 == 0 {
			return brownBold
		} else {
			return yellow
		}
	}
}
