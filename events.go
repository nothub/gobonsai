package main

import (
	"time"
)

type EventDrawn struct {
	t time.Time
}

func (ev *EventDrawn) When() time.Time {
	return ev.t
}

func EvDrawn() *EventDrawn {
	return &EventDrawn{t: time.Now()}
}
