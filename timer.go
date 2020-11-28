package rtc

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type alarm struct {
	Time time.Time
}

type Timer struct {
	Chan <-chan alarm
	file *os.File
}

// NewTimerAt creates a timer that will trigger an alarm at the given time.
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
	ch := make(chan alarm, 1)
	buf := make([]byte, 4)
	timer := &Timer{
		file: c.f,
		Chan: ch,
	}

	go func() {
		_, err := syscall.Read(int(c.f.Fd()), buf)
		if err != nil {
			fmt.Printf("got error reading interrupt, returning\n")
			return
		}

		// TODO: Clean up these comments?
		// buf[0] = bit mask encoding the types of interrupt that occurred.
		// buf[1:3] = number of interrupts since last read
		//r := binary.LittleEndian.Uint32(buf)
		//irqTypes := r & 0x000000FF
		//fmt.Printf("r: 0x%X, types: 0x%X\n", r, irqTypes)
		//cnt := r >> 8

		now := time.Now()
		ch <- alarm{
			Time: now,
		}
		close(ch)
	}()

	return timer, nil
}

// TODO: Timer resolution limited to 1 second
// TODO: What to do if d < 1 second?
// TODO: Consider mimicking the time.After() patterns
func NewTimer(dev string, d time.Duration) (*Timer, error) {

	c, err := NewRTC(dev)
	if err != nil {
		return nil, err
	}

	t, err := c.Time()
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
	return nil, nil

	// TODO: Finish this function
	// copy other timer function or create a reusable function since code is the same?
}

func (t *Timer) Stop() {
	//close(t.Chan) // TODO?
	t.file.Close()
}
