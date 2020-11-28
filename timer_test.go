package rtc

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNewTimer() {
	t, err := NewTimer("/dev/rtc", time.Minute)
	if err != nil {
		panic(err)
	}
	defer t.Stop()

	select {
	case alarm := <-t.Chan:
		fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
	}
}

func DisabledTestNewTimerAt(t *testing.T) {
	// TODO: Fix this test
	timer, err := NewTimerAt(devRtc, time.Now())
	if err != nil {
		t.Error("failed to start timer", err)
		t.FailNow()
	}
	//defer timer.Stop()
	_ = timer
}

func TestNewTimer(t *testing.T) {
	// TODO
}
