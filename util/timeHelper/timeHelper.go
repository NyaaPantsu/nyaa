package timeHelper

import (
	"time"

	"github.com/NyaaPantsu/nyaa/util/log"
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
	fewDurationLater := time.Now().Add(duration)
	log.Debugf("time : %s", fewDurationLater)
	return fewDurationLater
}

func FewDurationLaterMillisecond(duration time.Duration) int64 {
	return FewDurationLater(duration).UnixNano() / int64(time.Millisecond)
}

func IsExpired(expirationTime time.Time) bool {
	log.Debugf("expirationTime : %s", expirationTime)
	after := time.Now().After(expirationTime)
	log.Debugf("after : %t", after)
	return after
}
