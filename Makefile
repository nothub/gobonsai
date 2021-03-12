.POSIX:
CC	= cc
CFLAGS	= -Wall -Wextra -pedantic
PKG_CONFIG	?= pkg-config
LDLIBS	= $(shell $(PKG_CONFIG) --libs ncurses panel || echo "-lncurses -ltinfo -lpanel")
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
