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
