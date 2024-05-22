package rtc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cleroux/rtc"
)

func ExampleNewTimer() {
	t, err := rtc.NewTimer("/dev/rtc", time.Minute)
	if err != nil {
		panic(err)
	}
	defer t.Stop()

	select {
	case alarm := <-t.Chan:
		fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
	}
}

func TestNewTimerAt(t *testing.T) {
	timer, err := rtc.NewTimerAt("/dev/rtc", time.Now().UTC().Add(time.Second))
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
	timer, err := rtc.NewTimer("/dev/rtc", time.Second)
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
