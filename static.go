package rtc

import (
	"path/filepath"
	"time"
)

// GetClocks returns a list of real-time clocks in the system.
func GetClocks() (devices []string, err error) {
	return filepath.Glob("/dev/rtc*")
}

// GetEpoch reads the epoch from the specified real-time clock device.
func GetEpoch(dev string) (epoch uint, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.GetEpoch()
}

// SetEpoch sets the epoch on the specified real-time clock device.
func SetEpoch(dev string, epoch uint) (err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetEpoch(epoch)
}

// GetTime reads the time from the specified real-time clock device.
func GetTime(dev string) (t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return time.Time{}, err
	}
	defer c.Close()
	return c.GetTime()
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

// GetFrequency returns the frequency of the specified real-time clock device.
func GetFrequency(dev string) (frequency uint, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return 0, err
	}
	defer c.Close()
	return c.GetFrequency()
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

// GetAlarm returns the alarm time for the specified real-time clock device.
func GetAlarm(dev string) (t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return time.Time{}, err
	}
	defer c.Close()
	return c.GetAlarm()
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

// GetWakeAlarm returns the current state of the wake alarm for the specified real-time clock device.
func GetWakeAlarm(dev string) (enabled bool, pending bool, t time.Time, err error) {
	c, err := NewRTC(dev)
	if err != nil {
		return false, false, time.Time{}, err
	}
	defer c.Close()
	return c.GetWakeAlarm()
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
