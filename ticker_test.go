package rtc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

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
