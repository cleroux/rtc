// +build !windows

package rtc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

type tick struct {
	Time   time.Time
	Delta  time.Duration
	Frame  uint
	Missed uint32
}

type Ticker struct {
	buf   [1]byte
	frame uint
	t     time.Time
	Chan  <-chan tick
	file  *os.File
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
	ch := make(chan tick, 1)
	buf := make([]byte, 4)
	t := &Ticker{
		file:  c.f,
		frame: 0,
		t:     time.Now(),
		Chan:  ch,
	}

	go func() {
		for {
			_, err := syscall.Read(int(c.f.Fd()), buf)
			if err != nil {
				fmt.Printf("got error reading interrupt, breaking loop\n")
				break
			}

			// buf[0] = bit mask encoding the types of interrupt that occurred.
			// buf[1:3] = number of interrupts since last read
			r := binary.LittleEndian.Uint32(buf)
			//irqTypes := r & 0x000000FF
			//fmt.Printf("r: 0x%X, types: 0x%X\n", r, irqTypes)
			cnt := r >> 8

			now := time.Now()
			ch <- tick{
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
		close(ch)
	}()

	return t, nil
}

func (t *Ticker) Stop() {
	t.file.Close()
}
