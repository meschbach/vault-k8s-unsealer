package main

import (
	"log"
	"time"
)

//todo: convert to time.Backoff
type exponetialBackoff struct {
	unit time.Duration
	base int
	limit time.Duration
	currentIncrement int
	state int
}

func (e *exponetialBackoff) nextBackoff(state int) time.Duration {
	if e.state != state {
		e.state = state
		e.currentIncrement = 0
	} else {
		e.currentIncrement++
	}
	target := powerDuration(e.unit,e.base,e.currentIncrement)
	if target > e.limit {
		target = e.limit
	}
	return target
}

func (e *exponetialBackoff) performBackoff(state int)  {
	sleepFor := e.nextBackoff(state)
	log.Printf("Sleeping for %s...", humanizeDuration(sleepFor))
	time.Sleep(sleepFor)
}
