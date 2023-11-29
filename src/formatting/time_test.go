package formatting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FakeClock struct {
	Current time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.Current
}

func TestRelativeTime(t *testing.T) {
	now := time.Unix(10000000, 0)
	clock := FakeClock{now}

	fiveSecondsAgo := now.Add(-5 * time.Second)
	fortyFiveSecondsAgo := now.Add(-45 * time.Second)
	eightySecondsAgo := now.Add(-80 * time.Second)
	fortyFiveMinutesAgo := now.Add(-45 * time.Minute)
	fiftyFiveMinutesAgo := now.Add(-55 * time.Minute)
	eightyFiveMinutesAgo := now.Add(-85 * time.Minute)
	twoHoursAgo := now.Add(-2 * time.Hour)
	twentyTwoHoursAgo := now.Add(-22 * time.Hour)
	thirtyHoursAgo := now.Add(-30 * time.Hour)
	sixDaysAgo := now.Add(-6 * day)
	twentyOneDaysAgo := now.Add(-21 * day)
	twentySixDaysAgo := now.Add(-26 * day)
	fortyDaysAgo := now.Add(-40 * day)
	tenMonthsAgo := now.Add(-10 * month)
	twelveMonthsAgo := now.Add(-12 * month)
	sixteenMonthsAgo := now.Add(-16 * month)
	eighteenMonthsAgo := now.Add(-18 * month)
	fourYearsAgo := now.Add(-4 * year)

	assert.Equal(t, "a minute ago", GetRelativeTime(&clock, fiveSecondsAgo))
	assert.Equal(t, "a minute ago", GetRelativeTime(&clock, fortyFiveSecondsAgo))
	assert.Equal(t, "a minute ago", GetRelativeTime(&clock, eightySecondsAgo))
	assert.Equal(t, "45 min ago", GetRelativeTime(&clock, fortyFiveMinutesAgo))
	assert.Equal(t, "an hour ago", GetRelativeTime(&clock, fiftyFiveMinutesAgo))
	assert.Equal(t, "an hour ago", GetRelativeTime(&clock, eightyFiveMinutesAgo))
	assert.Equal(t, "2 hours ago", GetRelativeTime(&clock, twoHoursAgo))
	assert.Equal(t, "a day ago", GetRelativeTime(&clock, twentyTwoHoursAgo))
	assert.Equal(t, "a day ago", GetRelativeTime(&clock, thirtyHoursAgo))
	assert.Equal(t, "6 days ago", GetRelativeTime(&clock, sixDaysAgo))
	assert.Equal(t, "21 days ago", GetRelativeTime(&clock, twentyOneDaysAgo))
	assert.Equal(t, "a month ago", GetRelativeTime(&clock, twentySixDaysAgo))
	assert.Equal(t, "a month ago", GetRelativeTime(&clock, fortyDaysAgo))
	assert.Equal(t, "10 months ago", GetRelativeTime(&clock, tenMonthsAgo))
	assert.Equal(t, "a year ago", GetRelativeTime(&clock, twelveMonthsAgo))
	assert.Equal(t, "a year ago", GetRelativeTime(&clock, sixteenMonthsAgo))
	assert.Equal(t, "2 years ago", GetRelativeTime(&clock, eighteenMonthsAgo))
	assert.Equal(t, "4 years ago", GetRelativeTime(&clock, fourYearsAgo))
}

func TestRelativeTimePanicsWithFutureValues(t *testing.T) {
	now := time.Unix(10000000, 0)
	clock := FakeClock{now}

	assert.Panics(t, func() {
		oneMonthFromNow := now.Add(month)
		GetRelativeTime(&clock, oneMonthFromNow)
	}, "cannot get relative time for future time")
}
