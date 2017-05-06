package timeHelper

import (
	"time"

	"github.com/ewhal/nyaa/util/log"
)

func FewDaysLater(day int) time.Time {
	return FewDurationLater(time.Duration(day) * 24 * time.Hour)
}

func TwentyFourHoursLater() time.Time {
	return FewDurationLater(time.Duration(24) * time.Hour)
}

func SixHoursLater() time.Time {
	return FewDurationLater(time.Duration(6) * time.Hour)
}

func InTimeSpan(start, end, check time.Time) bool {
	log.Debugf("check after before: %s %t %t\n", check, check.After(start), check.Before(end))
	return check.After(start) && check.Before(end)
}

func InTimeSpanNow(start, end time.Time) bool {
	now := time.Now()
	return InTimeSpan(start, end, now)
}

func FewDurationLater(duration time.Duration) time.Time {
	// When Save time should considering UTC
	// baseTime := time.Now()
	// log.Debugf("basetime : %s", baseTime)
	fewDurationLater := time.Now().Add(duration)
	log.Debugf("time : %s", fewDurationLater)
	return fewDurationLater
}

func FewDurationLaterMillisecond(duration time.Duration) int64 {
	return FewDurationLater(duration).UnixNano() / int64(time.Millisecond)
}

func IsExpired(expirationTime time.Time) bool {
	// baseTime := time.Now()
	// log.Debugf("basetime : %s", baseTime)
	log.Debugf("expirationTime : %s", expirationTime)
	// elapsed := time.Since(expirationTime)
	// log.Debugf("elapsed : %s", elapsed)
	after := time.Now().After(expirationTime)
	log.Debugf("after : %t", after)
	return after
}
