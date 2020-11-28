# rtc

Go module for using a hardware Real-Time Clock (RTC) device in Linux.

## Installation

> $ go get github.com/cleroux/rtc

## Usage

`rtc.NewTicker()` and `rtc.NewTimer()` provide high-level interfaces inspired
by Go's `time.NewTicker()` and `time.NewTimer()` with the obvious difference
being that the `rtc` functions trigger directly from the RTC's hardware
interrupts.

The following example creates a ticker that fires at 2 Hz.
```go
t, err := rtc.NewTicker("/dev/rtc", 2)
if err != nil {
  panic(err)
}
defer t.Stop()

for {
  select {
  case tick := <-t.Chan:
    fmt.Printf("Tick.  Frame:%d Time:%v Delta:%v Missed:%d\n", tick.Frame, tick.Time, tick.Delta, tick.Missed)
  }
}
```

The following example sets an alarm for 5 seconds in the future and waits for
the alarm to fire.
```go
t, err := rtc.NewTimer("/dev/rtc", time.Now().Add(time.Second * 5))
if err != nil {
  panic(err)
}
defer t.Stop()

select {
case alarm := <-t.Chan:
  fmt.Printf("Alarm.  Time:%v\n", alarm.Time)
}
```

If more flexible programming of the RTC is needed, `rtc.NewRTC()` instantiates
an object that exposes all RTC functionality. The RTC device file is kept open
until `Close()` is called.

```go
c, err := rtc.NewRTC("/dev/rtc")
if err != nil {
  return err
}
defer c.Close()

t, err := c.Time()
```

Lastly, for convenience, static utility functions are also provided. These
functions are ideal when just a single function or operation is necessary
because they open and close the RTC device. For example, if reading the RTC's
time as in the following example:

```go
t, err := rtc.Time("/dev/rtc")
if err != nil {
  panic(err)
}
fmt.Printf("Current time: %v\n", t)
```

## References

[1] [The Linux kernel user's and administrator's guide: Real Time Clock (RTC) Drivers for Linux](https://www.kernel.org/doc/html/latest/admin-guide/rtc.html)  
[2] [rtc - Linux manual page](https://man7.org/linux/man-pages/man4/rtc.4.html)

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
