package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"time"
)

func post(sc *screen, ev tcell.Event) {
	err := sc.PostEvent(ev)
	if err != nil {
		log.Panicln(err.Error())
	}
}

type eventDrawn struct {
	t time.Time
}

func (ev *eventDrawn) When() time.Time {
	return ev.t
}

func evDrawn(sc *screen) {
	post(sc, &eventDrawn{t: time.Now()})
}

type eventQuit struct {
	t time.Time
}

func (ev *eventQuit) When() time.Time {
	return ev.t
}

func evQuit(sc *screen) {
	post(sc, &eventQuit{t: time.Now()})
}
