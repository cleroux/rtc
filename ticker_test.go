package rtc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func ExampleNewTicker() {
	ticker, err := NewTicker("/dev/rtc", 2)
	if err != nil {
		panic(err)
	}
	defer ticker.Stop()

	for {
		select {
		case tick := <-ticker.Chan:
			fmt.Printf("Tick.  Frame:%d Time:%v Delta:%v Over:%d\n", tick.Frame, tick.Time, tick.Delta, tick.Missed)
		}
	}
}

func TestTicker(t *testing.T) {

	var freq uint = 2
	interval := time.Duration(time.Second.Nanoseconds() / int64(freq))

	var cnt uint
	var prevTick tick

	done := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second + interval/2)
		done <- true
	}()

	ticker, err := NewTicker(devRtc, freq)
	require.NoError(t, err)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-done:
			break loop
		case tick := <-ticker.Chan:
			assert.Equal(t, uint32(0), tick.Missed)
			assert.Equal(t, cnt, tick.Frame)
			//assert.WithinDuration(t, interval, tick.Delta, time.Millisecond) // TODO
			if !prevTick.Time.Equal(time.Time{}) {
				assert.WithinDuration(t, prevTick.Time.Add(interval), tick.Time, time.Millisecond)
			}
			cnt++
		case <-time.After(time.Second):
			break
		}
	}

	assert.Equal(t, freq, cnt)
}
