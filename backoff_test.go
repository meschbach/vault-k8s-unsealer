package main

import (
	"testing"
	"time"
)

func assertDurationEqual(t *testing.T, expected time.Duration, got time.Duration)  {
	if expected != got {
		t.Errorf("Expected %s got %s", humanizeDuration(expected), humanizeDuration(got))
	}
}

func TestPowerDuration_1(t *testing.T)  {
	assertDurationEqual(t, 10 * time.Millisecond, powerDuration(time.Millisecond, 10, 1))
}

func TestPowerDuration_2(t *testing.T)  {
	assertDurationEqual(t, 100 * time.Millisecond, powerDuration(time.Millisecond, 10,2))
}

func TestPowerDuration_3(t *testing.T)  {
	assertDurationEqual(t, 1 * time.Second, powerDuration(time.Millisecond, 10,3))
}
