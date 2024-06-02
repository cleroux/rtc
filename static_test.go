package rtc

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// procDriverRtc reads /proc/driver/rtc and returns a map of the values it contains.
func procDriverRtc(t *testing.T) (values map[string]string) {
	t.Helper()

	b, err := os.ReadFile("/proc/driver/rtc")
	require.NoError(t, err, "Unable to read /proc/driver/rtc")

	lines := strings.Split(string(b), "\n")
	values = make(map[string]string)
	for _, l := range lines {
		fields := strings.SplitN(l, ":", 2)
		if len(fields) < 2 {
			continue
		}
		values[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
	}

	return values
}

// sysClassRtc reads /sys/class/rtc and returns the value contained in the specified file.
func sysClassRtc(t *testing.T) (value string) {
	t.Helper()

	b, err := os.ReadFile("/sys/class/rtc/rtc0")
	require.NoError(t, err)
	return string(b)
}

func TestTime(t *testing.T) {
	pdr := procDriverRtc(t)
	pdrTime, ok := pdr["rtc_time"]
	require.True(t, ok, "/proc/driver/rtc did not report rtc_time")
	pdrDate, ok := pdr["rtc_date"]
	require.True(t, ok, "/proc/driver/rtc did not report rtc_date")

	tm, err := Time(devRtc)
	require.NoError(t, err)

	assert.Equal(t, pdrTime, tm.Format("15:04:05"), "time read from RTC did not match the value reported by /proc/driver/rtc")
	assert.Equal(t, pdrDate, tm.Format("2006-01-02"), "date read from RTC did not match the value reported by /proc/driver/rtc")
}

func TestAlarm(t *testing.T) {
	pdr := procDriverRtc(t)
	pdrTime, ok := pdr["alrm_time"]
	require.True(t, ok, "/proc/driver/rtc did not report alrm_time")
	pdrDate, ok := pdr["alrm_date"]
	require.True(t, ok, "/proc/driver/rtc did not report alrm_time")

	tm, err := Alarm(devRtc)
	require.NoError(t, err)

	assert.Equal(t, pdrTime, tm.Format("15:04:05"), "alarm time read from RTC did not match the value reported by /proc/driver/rtc")
	assert.Equal(t, pdrDate, tm.Format("2006-01-02"), "alarm date read from RTC did not match the value reported by /proc/driver/rtc")
}

func TestFrequency(t *testing.T) {
	pdr := procDriverRtc(t)
	pdrFreqStr, ok := pdr["periodic IRQ frequency"]
	require.True(t, ok, "/proc/driver/rtc did not report periodic IRQ frequency")
	pdrFreq, err := strconv.ParseUint(pdrFreqStr, 10, 32)
	require.NoError(t, err)

	// According to `man rtc`, frequency can be in the range of 2 Hz to 8192 Hz
	assert.GreaterOrEqual(t, pdrFreq, uint64(2))
	assert.LessOrEqual(t, pdrFreq, uint64(8192))

	freq, err := Frequency(devRtc)
	require.NoError(t, err)
	assert.Equal(t, uint(pdrFreq), freq, "frequency read from RTC did not match the value reported by /proc/driver/rtc")
}

func TestSetFrequency(t *testing.T) {
	require.NoError(t, SetFrequency("/dev/rtc", 64))
	f, err := Frequency(devRtc)
	require.NoError(t, err)
	require.Equal(t, f, uint(64))

	require.NoError(t, SetFrequency("/dev/rtc", 32))
	f, err = Frequency(devRtc)
	require.NoError(t, err)
	require.Equal(t, f, uint(32))
}
