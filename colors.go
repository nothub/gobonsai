package main

import "github.com/fatih/color"

var BOLD_RED = color.New(color.FgRed).Add(color.Bold)

// int init_pair(short pair, short f, short b);
// (set pair 'pair' to fg 'f' and bg 'b')

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
