package rtc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTimerAt(t *testing.T) {
	timer, err := NewTimerAt("/dev/rtc", time.Now().UTC().Add(time.Second))
	require.NoError(t, err)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-time.After(3 * time.Second):
		t.Error("alarm did not trigger in time")
	}
}

func TestNewTimer(t *testing.T) {
	timer, err := NewTimer("/dev/rtc", time.Second)
	require.NoError(t, err)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-time.After(3 * time.Second):
		t.Error("alarm did not trigger in time")
	}
}
