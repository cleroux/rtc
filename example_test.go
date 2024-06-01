package rtc_test

import (
	"fmt"
	"time"

	"github.com/cleroux/rtc"
)

func ExampleNewTicker() {
	ticker, err := rtc.NewTicker("/dev/rtc", 2)
	if err != nil {
		panic(err)
	}
	defer ticker.Stop()

	for {
		select {
		case tick := <-ticker.Chan:
			fmt.Printf("Tick.  Frame:%d Time:%v Delta:%v Missed:%d\n", tick.Frame, tick.Time, tick.Delta, tick.Missed)
		}
	}
}

func ExampleNewTimer() {
	timer, err := rtc.NewTimer("/dev/rtc", time.Minute)
	if err != nil {
		panic(err)
	}
	defer timer.Stop()

	select {
	case alarm := <-timer.Chan:
		fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
	}
}
