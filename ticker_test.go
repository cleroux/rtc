package rtc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TestTicker checks that the ticker fires the expected number of times in one second.
func TestTicker(t *testing.T) {
	const frequencyHz uint = 2

	// Calculate the expected time interval between ticks
	interval := time.Duration(time.Second.Nanoseconds() / int64(frequencyHz))

	var tickCount uint
	var prevTick tick

	// Sleep for 1 second + margin for timing error, then signal to end the test.
	done := make(chan bool)
	go func() {
		time.Sleep(time.Second + interval/2)
		close(done)
	}()

	ticker, err := NewTicker(devRtc, frequencyHz)
	require.NoError(t, err)
	defer ticker.Stop()

	// Count ticks
loop:
	for {
		select {
		case <-done:
			break loop
		case tick := <-ticker.Chan:
			// Expect we have not missed any ticks
			assert.Equal(t, uint32(0), tick.Missed)

			// Expect Frame to be the same as the tick count since the test is only one second and
			// the frame counter will not roll over.
			assert.Equal(t, tickCount, tick.Frame)

			if !prevTick.Time.Equal(time.Time{}) {
				assert.WithinDuration(t, prevTick.Time.Add(interval), tick.Time, time.Millisecond)
			}
			tickCount++
		}
	}

	// Expect the tick count to equal the ticker's frequency.
	assert.Equal(t, frequencyHz, tickCount)
}
