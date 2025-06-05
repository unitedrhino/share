package def

import "time"

var TimerPriority = map[int64]string{
	6: "critical",
	3: "default",
	1: "low",
}

const (
	TimedUnitedRhinoQueueGroupCode = "unitedRhino-queue"
)

type TimeUnit string

const (
	TimeUnitA = "a" // (毫秒,默认),
	TimeUnitD = "d" // (天),
	TimeUnitH = "h" // (小时),
	TimeUnitM = "m" // (分钟),
	TimeUnitN = "n" // (月),
	TimeUnitS = "s" // (秒),
	TimeUnitU = "u" // (微秒),
	TimeUnitW = "w" // (周),
	TimeUnitY = "y" // (年)
)

func (t TimeUnit) String() string {
	return string(t)
}
func (t TimeUnit) ToPgStr() string {
	switch t {
	case TimeUnitD:
		return "day"
	case TimeUnitH:
		return "hour"
	case TimeUnitM:
		return "minute"
	case TimeUnitN:
		return "month"
	case TimeUnitW:
		return "week"
	case TimeUnitY:
		return "year"
	default:
		return "second"
	}
}
func (t TimeUnit) Truncate(ts time.Time, tim int64) time.Time {
	if tim <= 1 {
		switch t {
		case TimeUnitY:
			ts = time.Date(ts.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		case TimeUnitN:
			ts = time.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0, time.Local)
		case TimeUnitD, TimeUnitW:
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.Local)
		case TimeUnitH:
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), 0, 0, 0, time.Local)
		case TimeUnitM:
			ts = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), 0, 0, time.Local)
		default:
			goto DefaultRun
		}
		return ts
	}
DefaultRun:
	return ts.Truncate(t.ToDuration(tim))
}

func (t TimeUnit) ToDuration(in int64) time.Duration {
	switch t {
	case TimeUnitD:
		return time.Hour * 24 * time.Duration(in)
	case TimeUnitH:
		return time.Hour * time.Duration(in)
	case TimeUnitM:
		return time.Minute * time.Duration(in)
	case TimeUnitN:
		return time.Hour * 24 * 30 * time.Duration(in)
	case TimeUnitS:
		return time.Second * time.Duration(in)
	case TimeUnitU:
		return time.Millisecond * time.Duration(in)
	case TimeUnitW:
		return time.Hour * 24 * 7 * time.Duration(in)
	case TimeUnitY:
		return time.Hour * 24 * 365 * time.Duration(in)
	default:
		return time.Millisecond * time.Duration(in)
	}
}
