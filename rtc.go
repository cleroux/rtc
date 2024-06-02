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
	"errors"
	"fmt"
	"os"
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
	f *os.File
}

// NewRTC opens a real-time clock device.
func NewRTC(dev string) (*RTC, error) {
	f, err := os.Open(dev)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open rtc: %s", err.Error()))
	}
	return &RTC{
		f: f,
	}, nil
}

// Close closes a real-time clock device.
func (c *RTC) Close() (err error) {
	err = c.f.Close()
	c.f = nil
	return err
}

// Epoch returns the real-time clock's epoch.
func (c *RTC) Epoch() (epoch uint, err error) {
	e := new(uint32)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_EPOCH_READ, uintptr(unsafe.Pointer(e))); errno != 0 {
		return 0, errors.New(fmt.Sprintf("failed to read real-time clock epoch: %s", errno.Error()))
	}
	return uint(*e), nil
}

// SetEpoch sets the real-time clock's epoch.
func (c *RTC) SetEpoch(epoch uint) (err error) {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_EPOCH_SET, uintptr(epoch)); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock epoch: %s", errno.Error()))
	}
	return nil
}

// Time returns the specified real-time clock device time.
func (c *RTC) Time() (t time.Time, err error) {
	tm := new(rtcTime)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_RD_TIME, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return time.Time{}, errors.New(fmt.Sprintf("failed to read real-time clock time: %s", errno.Error()))
	}
	return tm.time(), nil
}

// SetTime sets the time for the specified real-time clock device.
func (c *RTC) SetTime(t time.Time) (err error) {
	tm := timeRtc{Time: t}.rtcTime()
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_SET_TIME, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock time: %s", errno.Error()))
	}
	return nil
}

// Frequency returns the periodic interrupt frequency.
func (c *RTC) Frequency() (frequency uint, err error) {
	f := new(uint)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_IRQP_READ, uintptr(unsafe.Pointer(f))); errno != 0 {
		return 0, errors.New(fmt.Sprintf("failed to read real-time clock frequency: %s", errno.Error()))
	}
	return *f, nil
}

// SetFrequency sets the frequency of the real-time clock's periodic interrupt.
func (c *RTC) SetFrequency(frequency uint) (err error) {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_IRQP_SET, uintptr(frequency)); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock frequency: %s", errno.Error()))
	}
	return nil
}

// SetPeriodicInterrupt enables or disables the real-time clock's periodic interrupts.
func (c *RTC) SetPeriodicInterrupt(enable bool) (err error) {
	op := unix.RTC_PIE_ON
	if !enable {
		op = unix.RTC_PIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), uintptr(op), 0); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock interrupts: %s", errno.Error()))
	}
	return nil
}

// SetAlarmInterrupt enables or disables the real-time clock's alarm interrupt.
func (c *RTC) SetAlarmInterrupt(enable bool) (err error) {
	op := unix.RTC_AIE_ON
	if !enable {
		op = unix.RTC_AIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), uintptr(op), 0); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock alarm interrupt: %s", errno.Error()))
	}
	return nil
}

// SetUpdateInterrupt enables or disables the real-time clock's update interrupt.
func (c *RTC) SetUpdateInterrupt(enable bool) (err error) {
	op := unix.RTC_UIE_ON
	if !enable {
		op = unix.RTC_UIE_OFF
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), uintptr(op), 0); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock update interrupt: %s", errno.Error()))
	}
	return nil
}

// Alarm returns the real-time clock's alarm time.
func (c *RTC) Alarm() (t time.Time, err error) {
	tm := new(rtcTime)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_ALM_READ, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return time.Time{}, errors.New(fmt.Sprintf("failed to read real-time clock alarm: %s", errno.Error()))
	}
	return tm.time(), nil
}

// SetAlarm sets the real-time clock's alarm time.
func (c *RTC) SetAlarm(t time.Time) (err error) {
	tm := timeRtc{Time: t}.rtcTime()
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_ALM_SET, uintptr(unsafe.Pointer(tm))); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock alarm: %s", errno.Error()))
	}
	return nil
}

// WakeAlarm returns the real-time clock's wake alarm time.
func (c *RTC) WakeAlarm() (enabled bool, pending bool, t time.Time, err error) {
	a := new(unix.RTCWkAlrm)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_ALM_READ, uintptr(unsafe.Pointer(a))); errno != 0 {
		return false, false, time.Time{}, errors.New(fmt.Sprintf("failed to read real-time clock wake alarm: %s", errno.Error()))
	}
	return a.Enabled == 1, a.Pending == 1, rtcTime{a.Time}.time(), nil
}

// SetWakeAlarm sets the real-time clock's wake alarm time.
func (c *RTC) SetWakeAlarm(t time.Time) (err error) {
	a := &unix.RTCWkAlrm{
		Enabled: 1,
		Time:    *timeRtc{Time: t}.rtcTime(),
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_WKALM_SET, uintptr(unsafe.Pointer(a))); errno != 0 {
		return errors.New(fmt.Sprintf("failed to set real-time clock wake alarm: %s", errno.Error()))
	}
	return nil
}

// CancelWakeAlarm cancels the real-time clock's wake alarm.
func (c *RTC) CancelWakeAlarm() (err error) {
	a := &unix.RTCWkAlrm{
		Enabled: 0,
		Time:    *timeRtc{Time: time.Time{}}.rtcTime(),
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, c.f.Fd(), unix.RTC_WKALM_SET, uintptr(unsafe.Pointer(a))); errno != 0 {
		return errors.New(fmt.Sprintf("failed to cancel real-time clock wake alarm: %s", errno.Error()))
	}
	return nil
}
