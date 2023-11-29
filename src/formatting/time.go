package formatting

import (
	"fmt"
	"math"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

var (
	day   = 24 * time.Hour
	month = 30 * day
	year  = 12 * month
)

// Converts a timestamp from the past to a high-level human-readable time
// string (e.g. "2 minutes ago", "5 years ago", etc).
func GetRelativeTime(clock Clock, timestamp time.Time) string {
	now := clock.Now()
	if timestamp.After(now) {
		panic(fmt.Sprintf("cannot get relative time for future time %v", timestamp))
	}

	elapsed := now.Sub(timestamp)

	switch {
	case elapsed < 90*time.Second:
		return "a minute ago"
	case elapsed < 50*time.Minute:
		return fmt.Sprintf("%.0f min ago", math.Ceil(elapsed.Minutes()))
	case elapsed < 90*time.Minute:
		return "an hour ago"
	case elapsed < 21*time.Hour:
		return fmt.Sprintf("%.0f hours ago", math.Ceil(elapsed.Hours()))
	case elapsed < 36*time.Hour:
		return "a day ago"
	case elapsed < 25*day:
		return fmt.Sprintf("%.0f days ago", math.Ceil(elapsed.Hours()/24.0))
	case elapsed < 45*day:
		return "a month ago"
	case elapsed < 11*month:
		return fmt.Sprintf("%.0f months ago", math.Ceil(elapsed.Hours()/(24.0*30)))
	case elapsed < 17*month:
		return "a year ago"
	default:
		return fmt.Sprintf("%.0f years ago", math.Ceil(elapsed.Hours()/(24.0*30*12)))
	}
}
