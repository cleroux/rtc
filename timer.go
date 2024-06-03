//go:build !windows
// +build !windows

package rtc

import (
	"fmt"
	"sync/atomic"
	"syscall"
	"time"
)

type Alarm struct {
	Time time.Time
}

type Timer struct {
	done  chan struct{}
	rtc   *RTC
	fired atomic.Bool
	C     <-chan Alarm
}

// NewTimerAt creates a new Timer that will send an Alarm on its channel after the given time.
func NewTimerAt(dev string, t time.Time) (*Timer, error) {
	c, err := NewRTC(dev)
	if err != nil {
		return nil, err
	}

	if err := c.SetAlarm(t); err != nil {
		_ = c.Close()
		return nil, err
	}

	if err := c.SetAlarmInterrupt(true); err != nil {
		_ = c.Close()
		return nil, err
	}

	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// on the floor until the client catches up.
	ch := make(chan Alarm, 1)
	timer := &Timer{
		done: make(chan struct{}),
		rtc:  c,
		C:    ch,
	}

	go func() {
		buf := make([]byte, 4)
		_, err := syscall.Read(c.fd, buf)
		if err != nil {
			fmt.Printf("got error reading interrupt, returning\n")
			return
		}

		select {
		case <-timer.done:
		// Don't send alarm if Stop() has been called
		default:
			timer.fired.Store(true)
		}

		// buf[0] = bit mask encoding the types of interrupt that occurred.
		// buf[1:3] = number of interrupts since last read
		//r := binary.LittleEndian.Uint32(buf)
		//irqTypes := r & 0x000000FF
		//fmt.Printf("r: 0x%X, types: 0x%X\n", r, irqTypes)
		//cnt := r >> 8

		ch <- Alarm{
			Time: time.Now(),
		}
	}()

	return timer, nil
}

// NewTimer creates a new Timer that will send an Alarm with the current time on its channel after at least duration d.
func NewTimer(dev string, d time.Duration) (*Timer, error) {
	c, err := NewRTC(dev)
	if err != nil {
		return nil, err
	}

	t, err := c.GetTime()
	if err != nil {
		return nil, err
	}

	if err := c.SetAlarm(t.Add(d)); err != nil {
		_ = c.Close()
		return nil, err
	}

	if err := c.SetAlarmInterrupt(true); err != nil {
		_ = c.Close()
		return nil, err
	}

	ch := make(chan Alarm, 1)
	buf := make([]byte, 4)
	timer := &Timer{
		done: make(chan struct{}),
		rtc:  c,
		C:    ch,
	}

	go func() {
		_, err := syscall.Read(c.fd, buf)
		if err != nil {
			fmt.Printf("got error reading interrupt, returning: %v\n", err)
			return
		}

		select {
		case <-timer.done:
		// Don't send alarm if Stop() has been called
		default:
		}

		ch <- Alarm{
			Time: time.Now(),
		}
	}()

	return timer, nil
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
//
// To ensure the channel is empty after a call to Stop, check the
// return value and drain the channel.
// For example, assuming the program has not received from t.C already:
//
//	if !t.Stop() {
//		<-t.C
//	}
//
// This cannot be done concurrent to other receives from the Timer's
// channel or other calls to the Timer's Stop method.
func (t *Timer) Stop() bool {
	close(t.done)
	_ = t.rtc.Close()
	return t.fired.Load()
}
