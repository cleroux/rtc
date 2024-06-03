//go:build !windows
// +build !windows

package rtc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"syscall"
	"time"
)

type Tick struct {
	Time   time.Time
	Delta  time.Duration
	Frame  uint
	Missed uint32
}

type Ticker struct {
	done  chan struct{}
	frame uint
	rtc   *RTC
	t     time.Time
	wait  sync.WaitGroup
	C     <-chan Tick
}

func NewTicker(dev string, frequency uint) (*Ticker, error) {
	if frequency == 0 {
		return nil, errors.New("zero frequency for NewTicker")
	}

	c, err := NewRTC(dev)
	if err != nil {
		return nil, err
	}

	if err := c.SetFrequency(frequency); err != nil {
		_ = c.Close()
		return nil, err
	}

	if err := c.SetPeriodicInterrupt(true); err != nil {
		_ = c.Close()
		return nil, err
	}

	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// until the client catches up.
	ch := make(chan Tick, 1)
	buf := make([]byte, 4)
	t := &Ticker{
		done:  make(chan struct{}),
		rtc:   c,
		frame: 0,
		t:     time.Now(),
		C:     ch,
	}

	t.wait.Add(1)
	go func() {
		defer t.wait.Done()
	loop:
		for {
			select {
			case <-t.done:
				break loop
			default:
			}

			_, err := syscall.Read(c.fd, buf)
			if err != nil {
				fmt.Printf("got error reading interrupt, breaking loop: %v\n", err)
				break
			}

			// buf[0] = bit mask encoding the types of interrupt that occurred.
			// buf[1:3] = number of interrupts since last read
			r := binary.LittleEndian.Uint32(buf)
			//irqTypes := r & 0x000000FF
			//fmt.Printf("r: 0x%X, types: 0x%X\n", r, irqTypes)
			cnt := r >> 8

			now := time.Now()
			ch <- Tick{
				Time:   now,
				Delta:  now.Sub(t.t),
				Frame:  t.frame,
				Missed: cnt - 1,
			}

			// Save current time
			t.t = now

			// Increment frame count
			t.frame = t.frame + 1
			if t.frame >= frequency {
				t.frame = 0
			}
		}

		// Disable interrupts and close RTC device
		_ = c.SetPeriodicInterrupt(false)
		_ = c.Close()
	}()

	return t, nil
}

func (t *Ticker) Stop() {
	close(t.done)
	t.wait.Wait()
}
