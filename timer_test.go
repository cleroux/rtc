package rtc

import (
	"testing"
	"time"
)

func TestNewTimerAt(t *testing.T) {
	timer, err := NewTimerAt("/dev/rtc", time.Now().UTC().Add(time.Second))
	if err != nil {
		t.Error("failed to start timer", err)
		t.FailNow()
	}
	defer timer.Stop()

	select {
	case <-timer.Chan:
	case <-time.After(3 * time.Second):
		t.Error("alarm did not trigger in time")
	}
}

func TestNewTimer(t *testing.T) {
	timer, err := NewTimer("/dev/rtc", time.Second)
	if err != nil {
		t.Error("failed to start timer", err)
		t.FailNow()
	}
	defer timer.Stop()

	select {
	case <-timer.Chan:
	case <-time.After(3 * time.Second):
		t.Error("alarm did not trigger in time")
	}
}
