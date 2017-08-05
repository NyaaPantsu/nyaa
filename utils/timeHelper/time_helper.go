package timeHelper

import (
	"time"

	"github.com/NyaaPantsu/nyaa/utils/log"
)

// FewDaysLater : Give time now + some days
func FewDaysLater(day int) time.Time {
	return FewDurationLater(time.Duration(day) * 24 * time.Hour)
}

// TwentyFourHoursLater : Give time now + 24 hours
func TwentyFourHoursLater() time.Time {
	return FewDurationLater(time.Duration(24) * time.Hour)
}

// SixHoursLater : Give time now + 6 hours
func SixHoursLater() time.Time {
	return FewDurationLater(time.Duration(6) * time.Hour)
}

// InTimeSpan : check if time given is in the given time encapsulation
func InTimeSpan(start, end, check time.Time) bool {
	log.Debugf("check after before: %s %t %t\n", check, check.After(start), check.Before(end))
	return check.After(start) && check.Before(end)
}

// InTimeSpanNow : check if time now is in the given time encapsulation
func InTimeSpanNow(start, end time.Time) bool {
	now := time.Now()
	return InTimeSpan(start, end, now)
}

// FewDurationLater : Give time now + some time duration
func FewDurationLater(duration time.Duration) time.Time {
	// When Save time should considering UTC
	fewDurationLater := time.Now().Add(duration)
	log.Debugf("time : %s", fewDurationLater)
	return fewDurationLater
}

// FewDurationLaterMillisecond : Give time now + some millisecond
func FewDurationLaterMillisecond(duration time.Duration) int64 {
	return FewDurationLater(duration).UnixNano() / int64(time.Millisecond)
}

// IsExpired : check if time given is expired
func IsExpired(expirationTime time.Time) bool {
	log.Debugf("expirationTime : %s", expirationTime)
	after := time.Now().After(expirationTime)
	log.Debugf("after : %t", after)
	return after
}
