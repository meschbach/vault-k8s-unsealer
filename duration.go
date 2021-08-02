package main

import (
	"fmt"
	"strings"
	"time"
)

func powerDuration(unit time.Duration, base int, power int ) time.Duration  {
	multiplier :=  1
	for i := 0; i < power ; i++ {
		multiplier = multiplier * base
	}
	return time.Duration(multiplier) * unit
}

func remainder(from time.Duration, unit time.Duration) (wholes time.Duration, remainder time.Duration)  {
	wholes = from / unit
	remainder = from - (wholes * unit)
	return wholes, remainder
}

func humanizeDuration(what time.Duration) string {
	hours, partHours := remainder(what, time.Hour )
	minutes, partMinutes := remainder(partHours, time.Minute)
	seconds, partSeconds := remainder(partMinutes, time.Second)
	ms, _ := remainder(partSeconds, time.Millisecond)

	var out []string
	if hours > 0 {
		out = append(out, fmt.Sprintf("%d hrs", hours) )
	}
	if minutes > 0 {
		out = append(out, fmt.Sprintf("%d mins", minutes) )
	}
	if seconds > 0 {
		out = append(out, fmt.Sprintf("%d secs", seconds) )
	}
	if ms > 0 {
		out = append(out, fmt.Sprintf("%d ms", ms) )
	}

	if len(out) == 0 {
		return "a blank of am eye"
	} else {
		return strings.Join(out," ")
	}
}
