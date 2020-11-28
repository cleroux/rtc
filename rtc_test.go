package rtc

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const devRtc = "/dev/rtc"

func ExampleRTC() {
	c, err := NewRTC(devRtc)
	if err != nil {
		return
	}
	defer c.Close()

	if t, err := c.Time(); err == nil {
		fmt.Printf("Current RTC Time: %v\n", t)
	}
}

func TestRtcEpoch(t *testing.T) {
	c, err := NewRTC(devRtc)
	require.NoError(t, err)
	defer c.Close()

	// Read the current epoch value
	curEpoch, err := c.Epoch()
	if strings.Contains(err.Error(), "inappropriate ioctl for device") {
		// Epoch not supported by this hardware
		t.SkipNow()
	}
	require.NoError(t, err)
	require.NotZero(t, curEpoch)

	// Set a random epoch value
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	newEpoch := uint(rng.Uint32())
	require.NoError(t, c.SetEpoch(newEpoch))

	// Read the new epoch value
	readEpoch, err := c.Epoch()
	require.NoError(t, err)
	assert.Equal(t, newEpoch, readEpoch)

	// Set the original epoch value
	assert.NoError(t, c.SetEpoch(curEpoch))
}

func TestRtcTime(t *testing.T) {
	c, err := NewRTC(devRtc)
	require.NoError(t, err)
	defer c.Close()

	curTime, err := c.Time()
	require.NoError(t, err)

	// Change the time
	//rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	//newTime := curTime.Add(time.Duration(rng.Int()))
	newTime := curTime.Add(time.Minute * 10)
	require.NoError(t, c.SetTime(newTime))

	// Read the time
	readTime, err := c.Time()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, readTime.UnixNano(), newTime.UnixNano())

	// TODO: Use /sys/class/rtc/rtc0 files to validate?
	// TODO: Use /proc/driver/rtc?

	// Set the original time value
	assert.NoError(t, c.SetTime(curTime))
}

func TestRtcFrequency(t *testing.T) {
	c, err := NewRTC(devRtc)
	require.NoError(t, err)
	defer c.Close()

	// Read the current frequency
	curFreq, err := c.Frequency()
	require.NoError(t, err)

	// Set a new frequency
	var newFreq uint = 8
	require.NoError(t, c.SetFrequency(newFreq))

	// Read the frequency
	readFreq, err := c.Frequency()
	require.NoError(t, err)
	assert.Equal(t, newFreq, readFreq)

	// Restore the original frequency value
	assert.NoError(t, c.SetFrequency(curFreq))
}
