package rtc

import (
	"path/filepath"
	"time"
)

// Clocks returns a list of real-time clocks in the system.
func Clocks() (devices []string, err error) {
	return filepath.Glob("/dev/rtc*")
}

// Epoch reads the epoch from the specified real-time clock device.
func Epoch(dev string) (epoch uint, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.Epoch()
}

// Epoch sets the epoch on the specified real-time clock device.
func SetEpoch(dev string, epoch uint) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetEpoch(epoch)
}

// Time reads the time from the specified real-time clock device.
func Time(dev string) (t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return time.Time{}, err
	}
	defer c.Close()
	return c.Time()
}

// SetTime sets the time for the specified real-time clock device.
func SetTime(dev string, t time.Time) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetTime(t)
}

// Frequency returns the frequency of the specified real-time clock device.
func Frequency(dev string) (frequency uint, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.Frequency()
}

// SetFrequency sets the periodic interrupt frequency of the specified real-time clock device.
func SetFrequency(dev string, frequency uint) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetFrequency(frequency)
}

// SetPeriodicInterrupt enables or disables periodic interrupts for the specified real-time clock device.
func SetPeriodicInterrupt(dev string, enable bool) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return nil
	}
	defer c.Close()
	return c.SetPeriodicInterrupt(enable)
}

// SetAlarmInterrupt enables or disables the alarm interrupt for the specified real-time clock device.
func SetAlarmInterrupt(dev string, enable bool) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return nil
	}
	defer c.Close()
	return c.SetAlarmInterrupt(enable)
}

// SetUpdateInterrupt enables or disables the update interrupt for the specified real-time clock device.
func SetUpdateInterrupt(dev string, enable bool) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return nil
	}
	defer c.Close()
	return c.SetUpdateInterrupt(enable)
}

// Alarm returns the alarm time for the specified real-time clock device.
func Alarm(dev string) (t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return time.Time{}, err
	}
	defer c.Close()
	return c.Alarm()
}

// SetAlarm sets the alarm time for the specified real-time clock device.
func SetAlarm(dev string, t time.Time) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetAlarm(t)
}

// WakeAlarm returns the current state of the wake alarm for the specified real-time clock device.
func WakeAlarm(dev string) (enabled bool, pending bool, t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return false, false, time.Time{}, err
	}
	defer c.Close()
	return c.WakeAlarm()
}

// SetWakeAlarm sets the wake alarm time for the specified real-time clock device.
func SetWakeAlarm(dev string, t time.Time) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return
	}
	defer c.Close()
	return c.SetWakeAlarm(t)
}

// CancelWakeAlarm cancels the wake alarm for the specified real-time clock device.
func CancelWakeAlarm(dev string) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return
	}
	defer c.Close()
	return c.CancelWakeAlarm()
}
