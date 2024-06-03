//go:build !windows
// +build !windows

// Package rtc facilitates working with real-time clocks (RTCs).
// High level functions such as NewTicker and NewTimer encapsulate the details
// of working with the RTC while providing interfaces that are similar to Go's
// time.NewTicker and time.NewTimer respectively.
// If more flexible programming of the RTC is needed, the NewRTC function
// returns an rtc object that exposes all RTC functionality. When this object
// is instantiated, the RTC device file is kept open until the Close function
// is called.
// For convenience, static utility functions are also provided to open and
// close the RTC when only one function is needed. For example, reading the
// clock once is possible simply by calling rtc.Time().
// Note that when working with the RTC, the highest resolution for time values
// is one second as defined in unix.RTCTime.
// https://www.kernel.org/doc/html/latest/admin-guide/rtc.html
// https://blog.cloudflare.com/its-go-time-on-linux/
// https://man7.org/linux/man-pages/man4/rtc.4.html
// https://code.woboq.org/linux/linux/drivers/char/rtc.c.html
package rtc

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type rtcTime struct {
	unix.RTCTime
}

func (r rtcTime) time() time.Time {
	return time.Date(int(r.Year+1900), time.Month(r.Mon+1), int(r.Mday), int(r.Hour), int(r.Min), int(r.Sec), 0, time.UTC)
}

type timeRtc struct {
	time.Time
}

func (t timeRtc) rtcTime() *unix.RTCTime {
	return &unix.RTCTime{
		Sec:  int32(t.Second()),
		Min:  int32(t.Minute()),
		Hour: int32(t.Hour()),
		Mday: int32(t.Day()),
		Mon:  int32(t.Month() - 1),
		Year: int32(t.Year() - 1900),
	}
}

type RTC struct {
	fd int
}

// NewRTC opens a real-time clock device.
func NewRTC(dev string) (*RTC, error) {
	fd, err := syscall.Open(dev, syscall.O_RDWR, uint32(0600))
	if err != nil {
		return nil, fmt.Errorf("failed to open rtc: %w", err)
	}
	return &RTC{
		fd: fd,
	}, nil
}

// Close closes a real-time clock device.
func (c *RTC) Close() (err error) {
	err = syscall.Close(c.fd)
	c.fd = 0
	return err
}

// GetEpoch returns the real-time clock's epoch.
func (c *RTC) GetEpoch() (epoch uint, err error) {
	e := new(uint32)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_EPOCH_READ, uintptr(unsafe.Pointer(e))); errno != 0 {
		return 0, fmt.Errorf("failed to read real-time clock epoch: %w", errno)
	}
	return uint(*e), nil
}

// SetEpoch sets the real-time clock's epoch.
func (c *RTC) SetEpoch(epoch uint) (err error) {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_EPOCH_SET, uintptr(epoch)); errno != 0 {
		return fmt.Errorf("failed to set real-time clock epoch: %w", errno)
	}
	return nil
}

// GetTime returns the specified real-time clock device time.
func (c *RTC) GetTime() (t time.Time, err error) {
	tm := new(rtcTime)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_RD_TIME, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return time.Time{}, fmt.Errorf("failed to read real-time clock time: %w", errno)
	}
	return tm.time(), nil
}

// SetTime sets the time for the specified real-time clock device.
func (c *RTC) SetTime(t time.Time) (err error) {
	tm := timeRtc{Time: t}.rtcTime()
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_SET_TIME, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return fmt.Errorf("failed to set real-time clock time: %w", errno)
	}
	return nil
}

// GetFrequency returns the periodic interrupt frequency.
func (c *RTC) GetFrequency() (frequency uint, err error) {
	f := new(uint)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_IRQP_READ, uintptr(unsafe.Pointer(f))); errno != 0 {
		return 0, fmt.Errorf("failed to read real-time clock frequency: %w", errno)
	}
	return *f, nil
}

// SetFrequency sets the frequency of the real-time clock's periodic interrupt.
func (c *RTC) SetFrequency(frequency uint) (err error) {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_IRQP_SET, uintptr(frequency)); errno != 0 {
		return fmt.Errorf("failed to set real-time clock frequency: %w", errno)
	}
	return nil
}

// SetPeriodicInterrupt enables or disables the real-time clock's periodic interrupts.
func (c *RTC) SetPeriodicInterrupt(enable bool) (err error) {
	op := unix.RTC_PIE_ON
	if !enable {
		op = unix.RTC_PIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), uintptr(op), 0); errno != 0 {
		return fmt.Errorf("failed to set real-time clock interrupts: %w", errno)
	}
	return nil
}

// SetAlarmInterrupt enables or disables the real-time clock's alarm interrupt.
func (c *RTC) SetAlarmInterrupt(enable bool) (err error) {
	op := unix.RTC_AIE_ON
	if !enable {
		op = unix.RTC_AIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), uintptr(op), 0); errno != 0 {
		return fmt.Errorf("failed to set real-time clock alarm interrupt: %w", errno)
	}
	return nil
}

// SetUpdateInterrupt enables or disables the real-time clock's update interrupt.
func (c *RTC) SetUpdateInterrupt(enable bool) (err error) {
	op := unix.RTC_UIE_ON
	if !enable {
		op = unix.RTC_UIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), uintptr(op), 0); errno != 0 {
		return fmt.Errorf("failed to set real-time clock update interrupt: %w", errno)
	}
	return nil
}

// GetAlarm returns the real-time clock's alarm time.
func (c *RTC) GetAlarm() (t time.Time, err error) {
	tm := new(rtcTime)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_ALM_READ, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return time.Time{}, fmt.Errorf("failed to read real-time clock alarm: %w", errno)
	}
	return tm.time(), nil
}

// SetAlarm sets the real-time clock's alarm time.
func (c *RTC) SetAlarm(t time.Time) (err error) {
	tm := timeRtc{Time: t}.rtcTime()
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_ALM_SET, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return fmt.Errorf("failed to set real-time clock alarm: %w", errno)
	}
	return nil
}

// GetWakeAlarm returns the real-time clock's wake alarm time.
func (c *RTC) GetWakeAlarm() (enabled bool, pending bool, t time.Time, err error) {
	a := new(unix.RTCWkAlrm)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_ALM_READ, uintptr(unsafe.Pointer(a))); errno != 0 {
		return false, false, time.Time{}, fmt.Errorf("failed to read real-time clock wake alarm: %w", errno)
	}
	return a.Enabled == 1, a.Pending == 1, rtcTime{a.Time}.time(), nil
}

// SetWakeAlarm sets the real-time clock's wake alarm time.
func (c *RTC) SetWakeAlarm(t time.Time) (err error) {
	a := &unix.RTCWkAlrm{
		Enabled: 1,
		Time:    *timeRtc{Time: t}.rtcTime(),
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_WKALM_SET, uintptr(unsafe.Pointer(a))); errno != 0 {
		return fmt.Errorf("failed to set real-time clock wake alarm: %w", errno)
	}
	return nil
}

// CancelWakeAlarm cancels the real-time clock's wake alarm.
func (c *RTC) CancelWakeAlarm() (err error) {
	a := &unix.RTCWkAlrm{
		Enabled: 0,
		Time:    *timeRtc{Time: time.Time{}}.rtcTime(),
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(c.fd), unix.RTC_WKALM_SET, uintptr(unsafe.Pointer(a))); errno != 0 {
		return fmt.Errorf("failed to cancel real-time clock wake alarm: %w", errno)
	}
	return nil
}
