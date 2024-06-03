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

	for tick := range ticker.C {
		fmt.Printf("Tick.  Frame:%d Time:%v Delta:%v Missed:%d\n", tick.Frame, tick.Time, tick.Delta, tick.Missed)
	}
}

func ExampleNewTimer() {
	timer, err := rtc.NewTimer("/dev/rtc", time.Minute)
	if err != nil {
		panic(err)
	}
	defer timer.Stop()

	alarm := <-timer.C
	fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
}

func ExampleRTC() {
	clock, err := rtc.NewRTC("/dev/rtc")
	if err != nil {
		panic(err)
	}
	defer clock.Close()

	t, err := clock.GetTime()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current RTC Time: %v\n", t)
}

func ExampleGetClocks() {
	clocks, err := rtc.GetClocks()
	if err != nil {
		panic(err)
	}
	for _, clock := range clocks {
		fmt.Printf("Clock found: %s\n", clock)
	}
}

func ExampleGetTime() {
	t, err := rtc.GetTime("/dev/rtc")
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

func ExampleGetEpoch() {
	epoch, err := rtc.GetEpoch("/dev/rtc")
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

func ExampleGetAlarm() {
	t, err := rtc.GetAlarm("/dev/rtc")
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

func ExampleGetFrequency() {
	frequency, err := rtc.GetFrequency("/dev/rtc")
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
