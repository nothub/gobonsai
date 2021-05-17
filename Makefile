.POSIX:
CC	= cc
PKG_CONFIG	?= pkg-config
CFLAGS	+= -Wall -Wextra -Wshadow -Wpointer-arith -Wcast-qual -pedantic $(shell $(PKG_CONFIG) --cflags ncursesw panelw)
LDLIBS	= $(shell $(PKG_CONFIG) --libs ncursesw panelw || echo "-lncursesw -ltinfo -lpanelw")
PREFIX	= /usr/local

cbonsai: cbonsai.c

install: cbonsai
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -m 0755 cbonsai $(DESTDIR)$(PREFIX)/bin/cbonsai

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/cbonsai

clean:
	rm -f cbonsai

.PHONY: install uninstall clean
