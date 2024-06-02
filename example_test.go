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

	for tick := range ticker.Chan {
		fmt.Printf("Tick.  Frame:%d Time:%v Delta:%v Missed:%d\n", tick.Frame, tick.Time, tick.Delta, tick.Missed)
	}
}

func ExampleNewTimer() {
	timer, err := rtc.NewTimer("/dev/rtc", time.Minute)
	if err != nil {
		panic(err)
	}
	defer timer.Stop()

	alarm := <-timer.Chan
	fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
}

func ExampleClocks() {
	clocks, err := rtc.Clocks()
	if err != nil {
		panic(err)
	}
	for _, clock := range clocks {
		fmt.Printf("Clock found: %s\n", clock)
	}
}

func ExampleTime() {
	t, err := rtc.Time("/dev/rtc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current time: %v\n", t)
}

func ExampleSetTime() {
	if err := rtc.SetTime("/dev/rtc", time.Now()); err != nil {
		panic(err)
	}
}

func ExampleEpoch() {
	epoch, err := rtc.Epoch("/dev/rtc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current epoch: %d\n", epoch)
}

func ExampleSetEpoch() {
	if err := rtc.SetEpoch("/dev/rtc", 99); err != nil {
		panic(err)
	}
}

func ExampleAlarm() {
	t, err := rtc.Alarm("/dev/rtc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current alarm time: %v\n", t)
}

func ExampleSetAlarm() {
	if err := rtc.SetAlarm("/dev/rtc", time.Now().Add(time.Minute)); err != nil {
		panic(err)
	}
}

func ExampleSetAlarmInterrupt() {
	if err := rtc.SetAlarmInterrupt("/dev/rtc", true); err != nil {
		panic(err)
	}
}

func ExampleFrequency() {
	frequency, err := rtc.Frequency("/dev/rtc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Frequency: %d\n", frequency)
}

func ExampleSetFrequency() {
	if err := rtc.SetFrequency("/dev/rtc", 64); err != nil {
		panic(err)
	}
}
